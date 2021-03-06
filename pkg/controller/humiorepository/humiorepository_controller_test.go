package humiorepository

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"testing"

	humioapi "github.com/humio/cli/api"
	corev1alpha1 "github.com/humio/humio-operator/pkg/apis/core/v1alpha1"
	"github.com/humio/humio-operator/pkg/humio"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// TODO: Add tests for updating repository

func TestReconcileHumioRepository_Reconcile(t *testing.T) {
	tests := []struct {
		name            string
		humioRepository *corev1alpha1.HumioRepository
		humioClient     *humio.MockClientConfig
	}{
		{
			"test simple repository reconciliation",
			&corev1alpha1.HumioRepository{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "humiorepository",
					Namespace: "logging",
				},
				Spec: corev1alpha1.HumioRepositorySpec{
					ManagedClusterName: "example-humiocluster",
					Name:               "example-repository",
					Description:        "important description",
					Retention: corev1alpha1.HumioRetention{
						TimeInDays:      30,
						IngestSizeInGB:  5,
						StorageSizeInGB: 1,
					},
				},
			},
			humio.NewMocklient(humioapi.Cluster{}, nil, nil, nil, ""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, req := reconcileInitWithHumioClient(tt.humioRepository, tt.humioClient)
			defer r.logger.Sync()

			_, err := r.Reconcile(req)
			if err != nil {
				t.Errorf("reconcile: (%v)", err)
			}

			updatedRepository, err := r.humioClient.GetRepository(tt.humioRepository)
			if err != nil {
				t.Errorf("get HumioRepository: (%v)", err)
			}

			expectedRepository := humioapi.Repository{
				Name:                   tt.humioRepository.Spec.Name,
				Description:            tt.humioRepository.Spec.Description,
				RetentionDays:          float64(tt.humioRepository.Spec.Retention.TimeInDays),
				IngestRetentionSizeGB:  float64(tt.humioRepository.Spec.Retention.IngestSizeInGB),
				StorageRetentionSizeGB: float64(tt.humioRepository.Spec.Retention.StorageSizeInGB),
			}

			if !reflect.DeepEqual(*updatedRepository, expectedRepository) {
				t.Errorf("repository %#v, does not match expected %#v", *updatedRepository, expectedRepository)
			}
		})
	}
}

func reconcileInitWithHumioClient(humioRepository *corev1alpha1.HumioRepository, humioClient *humio.MockClientConfig) (*ReconcileHumioRepository, reconcile.Request) {
	r, req := reconcileInit(humioRepository)
	r.humioClient = humioClient
	return r, req
}

func reconcileInit(humioRepository *corev1alpha1.HumioRepository) (*ReconcileHumioRepository, reconcile.Request) {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar().With("Request.Namespace", humioRepository.Namespace, "Request.Name", humioRepository.Name)

	humioCluster := &corev1alpha1.HumioCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      humioRepository.Spec.ManagedClusterName,
			Namespace: humioRepository.Namespace,
		},
	}

	apiTokenSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-admin-token", humioRepository.Spec.ManagedClusterName),
			Namespace: humioRepository.Namespace,
		},
		StringData: map[string]string{
			"token": "secret-api-token",
		},
	}

	// Objects to track in the fake client.
	objs := []runtime.Object{
		humioCluster,
		apiTokenSecret,
		humioRepository,
	}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(corev1alpha1.SchemeGroupVersion, humioRepository)
	s.AddKnownTypes(corev1alpha1.SchemeGroupVersion, humioCluster)

	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)

	// Create a ReconcileHumioRepository object with the scheme and fake client.
	r := &ReconcileHumioRepository{
		client: cl,
		scheme: s,
		logger: sugar,
	}

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      humioRepository.Name,
			Namespace: humioRepository.Namespace,
		},
	}
	return r, req
}
