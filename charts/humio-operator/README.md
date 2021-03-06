# humio-operator

[humio-operator](https://github.com/humio/humio-operator) Kubernetes Operator for running Humio on top of Kubernetes.

## Introduction

This chart bootstraps a humio-operator deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- [Kubernetes](https://kubernetes.io) 1.16+
- [cert-manager](https://cert-manager.io) v0.16+ (by default, but can be disabled with `certmanager` set to `false`)
- [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx) controller v0.34.1 (only required if configuring HumioCluster CR's with `ingress.controller` set to `nginx`)

## Installing the CRD's

```bash
kubectl apply -f https://raw.githubusercontent.com/humio/humio-operator/humio-operator-0.0.12/deploy/crds/core.humio.com_humioclusters_crd.yaml
kubectl apply -f https://raw.githubusercontent.com/humio/humio-operator/humio-operator-0.0.12/deploy/crds/core.humio.com_humioexternalclusters_crd.yaml
kubectl apply -f https://raw.githubusercontent.com/humio/humio-operator/humio-operator-0.0.12/deploy/crds/core.humio.com_humioingesttokens_crd.yaml
kubectl apply -f https://raw.githubusercontent.com/humio/humio-operator/humio-operator-0.0.12/deploy/crds/core.humio.com_humioparsers_crd.yaml
kubectl apply -f https://raw.githubusercontent.com/humio/humio-operator/humio-operator-0.0.12/deploy/crds/core.humio.com_humiorepositories_crd.yaml
```

## Installing the Chart

To install the chart with the release name `humio-operator`:

```bash
# Helm v3+
helm install humio-operator humio-operator/humio-operator --namespace humio-operator -f values.yaml

# Helm v2
helm install humio-operator/humio-operator --name humio-operator --namespace humio-operator -f values.yaml
```

> **Note**: By default, we expect cert-manager to be installed in order to configure TLS. If you do not have cert-manager installed, or if you know you do not want TLS, see the [configuration](#configuration) section for how to disable this.

> **Note**: By default, we expect a non-OpenShift installation, see the [configuration](#configuration) section for how to enable OpenShift specific functionality.

The command deploys humio-operator on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `humio-operator` deployment:

```bash
helm delete humio-operator --namespace humio-operator
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the ingress-nginx chart and their default values.

Parameter | Description | Default
--- | --- | ---
`operator.image.repository` | operator container image repository | `humio/humio-operator`
`operator.image.tag` | operator container image tag | `0.0.12`
`operator.rbac.create` | automatically create operator RBAC resources | `true`
`operator.watchNamespaces` | list of namespaces the operator will watch for resources (if empty, it watches all namespaces) | `[]`
`installCRDs` | automatically install CRDs. NB: if this is set to true, custom resources will be removed if the Helm chart is uninstalled | `false`
`openshift` | install additional RBAC resources specific to OpenShift | `false`
`certmanager` | whether cert-manager is present on the cluster, which will be used for TLS functionality | `true`

These parameters can be passed via Helm's `--set` option

```bash
# Helm v3+
helm install humio-operator humio-operator/humio-operator \
  --set operator.image.tag=0.0.12

# Helm v2
helm install humio-operator --name humio-operator \
  --set operator.image.tag=0.0.12
```

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example,

```bash
# Helm v3+
helm install humio-operator humio-operator/humio-operator --namespace humio-operator -f values.yaml

# Helm v2
helm install humio-operator/humio-helm-charts --name humio-operator --namespace humio-operator -f values.yaml
```
