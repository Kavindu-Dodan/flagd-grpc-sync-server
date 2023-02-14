## Flagd GRPC sync provider

A simple POC for [flagd](https://github.com/open-feature/flagd) which provides GRPC flag configuration syncs.

Utilize buf definitions at https://buf.build/kavindudodan/flagd 

### How to run ?

```shell
go run main.go <OPTIONS>
```

Following options are available,

```text
 -certPath string
        certificate path for tls connection
  -h string
        hostDefault of the server (default "localhost")
  -keyPath string
        certificate key for tls connection
  -p string
        portDefault of the server (default "9090")
  -s    enable tls
```

For example, to start with TLS certs,

```shell
go run main.go -s=true -certPath=server.crt -keyPath=server.key
```

Then start your GRPC sync enabled flagd.

Related to - https://github.com/open-feature/flagd/pull/297