# Inkle

![master](https://github.com/abrampers/inkle/workflows/master/badge.svg?event=push)
[![coverage](https://codecov.io/gh/abrampers/inkle/branch/master/graph/badge.svg)](https://codecov.io/gh/abrampers/inkle)

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

# Available flags
# -stdout         (bool)      : Write logs to stdout.
# -output=.       (string)    : Write log file to specified directory (ignored if -stdout is set).
# -timeout=200ms  (duration)  : Set request timeout.
# -h                          : Help.
```

## Roadmap
- [ ] Repo description.
- [ ] Repo architecture.
- [x] HTTP/2 frame classification.
- [x] State management to support gRPC connection reuse.
- [ ] Ensure correctness while ignoring unsupported streams.
- [ ] Support for gRPC streams.
