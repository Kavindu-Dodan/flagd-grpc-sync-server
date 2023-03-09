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

### Generate certificates ? 

Given below are some commands you can use to generate CA cert and Server cert to used with `localhost`


#### Generate CA cert

- CA Private Key: `openssl ecparam -name prime256v1 -genkey -noout -out ca.key`
- CA Certificate: `openssl req -new -x509 -sha256 -key ca.key -out ca.cert`

#### Generate Server certificate

- Server private key:  `openssl ecparam -name prime256v1 -genkey -noout -out server.key`
- Server signing request:  `openssl req -new -sha256 -addext "subjectAltName=DNS:localhost" -key server.key -out server.csr`
- Server cert:  `openssl x509 -req -in server.csr -CA ca.cert -CAkey ca.key  -out server.crt -days 1000 -sha256 -extfile opnessl.conf`

Where the file `opnessl.conf` contains following,

`subjectAltName = DNS:localhost`

#### Running grpc server with certificates

`go run main.go -s=true -certPath=server.crt -keyPath=server.key`

#### Running flagd with certificates

`go run main.go start --sources='[{"uri":"grpcs://localhost:9090","provider":"grpc", "certPath":"ca.cert"}]'`