# Inkle

![master](https://github.com/abrampers/inkle/workflows/master/badge.svg?event=push)
[![Codecov](https://img.shields.io/codecov/c/github/abrampers/inkle)](https://codecov.io/gh/abrampers/inkle)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/abrampers/inkle?color=%2329beb0)
[![Docker Image Version (latest semver)](https://img.shields.io/docker/v/abrampers/inkle?color=blue&label=docker&sort=semver)](https://hub.docker.com/repository/docker/abrampers/inkle/tags)


Log Based gRPC Tracing System 𝌘

## Log format
```
grpc_service_name,grpc_method_name,src_ip,src_tcp,dst_ip,dst_tcp,grpc_status_code,duration, info

e.g:
helloworld.Greeter,SayHello,::1,53412,::1,8000,0,161626,Request - Response
datetime.Datetime,GetDatetime,::1,53413,::1,9000,0,10120,Request - Response
```

## TODO: Installation

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
| --- | --- | --- | --- |
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
