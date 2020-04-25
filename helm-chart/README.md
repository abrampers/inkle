# Inkle Helm Chart

This helm chart is an easier way to configure and run Inkle.

## Notice

[0.1.0](https://github.com/abrampers/inkle/releases/tag/v0.1.0) release is tested on microk8s clusters.

## Requirements

* [Helm](https://helm.sh/).
* Kubernetes 1.16 - 1.8
* Minimum cluster requirements include the following to run this chart with default settings (Elastic stack enabled). All of these settings are configurable.
  * Two Kubernetes nodes with master (2 vCPU, 3GiB), worker (1 vCPU, 2GiB).

## Usage notes and getting started

* This chart deploys Inkle as a DaemonSet and will write the logs to path specified by `logPath` on each node.
* This chart also deploys Filebeat as a DaemonSet, Logstash, Kibana, and Elasticsearch cluster with one node to investigate logs produced by Inkle.

## Installing

### Using Helm repository

* Add the elastic helm charts repo
  ```
  helm repo add abrampers https://abram.id/helm/abrampers
  ```
* Install it
  ```
  helm install inkle abrampers/inkle
  ```

### Using master branch

* Clone the git repo
  ```
  git clone https://github.com/abrampers/inkle.git
  ```
* Install it
  ```
  helm install inkle ./helm-charts/inkle
  ```

## Configuration

| Parameter         | Description | Default |
| ----------------- | ----------- | ------- |
| `device`          | Network interface name where Inkle will intercept packets. | `cni0` |
| `image.repository`| The Inkle docker image | `abrampers/inkle` |
| `imagePullPolicy` | The Kubernetes [imagePullPolicy](https://kubernetes.io/docs/concepts/containers/images/#updating-images) value                                                                                                                                                                                                  | `IfNotPresent` |
| `elastic.enabled` | Deploys ELKB stack to do process logs produced by Inkle. | `true` |
| `logPath`         | The path where Inkle will write the logs.  | `/var/log` |
| `filterByHost`    | Filters logs originated from other host. If this parameter is set to true, it is guaranteed that every log sent to Elasticsearch is unique.  | `true` |
| `resources`       | Allows you to set the [resources](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/) for the DaemonSet                                                                                                                                                                          | `requests.cpu: 100m`<br>`requests.memory: 200Mi`<br>`limits.cpu: 100m`<br>`limits.memory: 200Mi`|
