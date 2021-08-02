# Inkle

![master](https://github.com/abrampers/inkle/workflows/master/badge.svg?event=push)
[![Codecov](https://img.shields.io/codecov/c/github/abrampers/inkle)](https://codecov.io/gh/abrampers/inkle)
[![Docker Image Version (latest semver)](https://img.shields.io/docker/v/abrampers/inkle?color=blue&label=docker&sort=semver)](https://hub.docker.com/repository/docker/abrampers/inkle/tags)
[![Helm Chart](https://img.shields.io/badge/helm-chart-%2306227f)](https://github.com/abrampers/inkle/tree/master/helm-chart)

Log Based gRPC Tracing System 𝌘

This is my bachelor thesis project to obtain Computer Science bachelor degree from Institut Teknologi Bandung.

This [paper](https://ieeexplore.ieee.org/abstract/document/9429054) was published on [ICAICTA 2020](http://icaicta.cs.tut.ac.jp/2020/).

## Log format
```
grpc_service_name,grpc_method_name,src_ip,src_tcp,dst_ip,dst_tcp,grpc_status_code,duration, info

e.g:
helloworld.Greeter,SayHello,::1,53412,::1,8000,0,161626,Request - Response
datetime.Datetime,GetDatetime,::1,53413,::1,9000,0,10120,Request - Response
```

## Installation

### Kubernetes Environment

Make sure `kubectl` is properly configured to an active Kubernetes cluster.

#### Using Helm

This requires [Helm](https://helm.sh/docs/intro/install/) to be installed. For more information about installing Inkle using Helm, see the [Inkle Helm Chart](https://github.com/abrampers/inkle/tree/master/helm-chart).

```sh
$ helm repo add abrampers https://abram.id/helm/abrampers
$ helm install inkle abrampers/inkle
```

#### Apply manifest manually

```sh
$ git clone https://github.com/abrampers/inkle
$ kubectl apply -f manifests/kubernetes.yaml
```

### Local environment

#### Docker

```sh
$ docker run --name inkle --network=host --privileged --rm -it abrampers/inkle:v0.1.0 [PARAMS]
```

#### Build from source

```sh
$ go get -u github.com/abrampers/inkle
$ cd $GOPATH/src/github.com/abrampers/inkle
$ go build -o .
```

## Usage
```sh
$ ./inkle -stdout
```
Available flags:
| Flag | Type | Default | Description |
| ---- | ---- | ------- | ----------- |
| `-device=cni0` | string | `eth0` | Network Device to be intercepted. |
| `-stdout` | bool | `false` | Write logs to stdout. |
| `-output=/var/log` | string | `.` | Write log file to specified directory (ignored if `-stdout` is set). |
| `-timeout=200ms` | time.Duration | `800ms` | Set request timeout. |
| `-filter-by-host-cidr` | bool | `false` | If this flag is set, Inkle will get the valid IP range of the network device specified in `-device` and will only print logs with source IP addres within that range. |
| `-h` | n/a | n/a | Print out help message. |

## Roadmap
- [ ] Repo description.
- [ ] Repo architecture.
- [x] HTTP/2 frame classification.
- [x] State management to support gRPC connection reuse.
- [x] Supports source IP address filtering by host CIDR.
- [ ] Log rotation.
- [ ] HTTPS support
- [ ] Ensure correctness while ignoring unsupported streams.
- [ ] Support for gRPC streams.
