# go-ebpf-logger

The golang code that actually runs the file monitor ebpf code

## Pre-requisite

Install golang, visit [link](https://go.dev/doc/install)

Install goreleaser, visit [link](https://goreleaser.com/install/#aur)

Install gh, visit [link](https://cli.github.com/)

## Screenshots

It works!

![Screenshot showing go-ebpf-logger in action](./docs/assets/working.png)


## Run it on Local

Get file monitor binary from the repo

Put the version accordingly here in the below command at `vX.Y.Z`

```sh
gh release download vX.Y.Z --repo SentinalFS/file-monitor --pattern "monitor.bpf.o"
```

Run it

```sh
sudo go run main.go
```

## Run it on docker

Get file monitor binary from the repo

Put the version accordingly here in the below command at `vX.Y.Z`

```sh
gh release download vX.Y.Z --repo SentinalFS/file-monitor --pattern "monitor.bpf.o"
```

Run go releaser on local

```sh
goreleaser release --snapshot --skip=publish --clean
```

Build it

```sh
docker build --build-arg TARGETARCH=amd64 -t go-ebpf-logger -f Dockerfile.amd64 .
```

Run it

```sh
sudo docker run --rm -it --privileged  -v /sys/fs/bpf:/sys/fs/bpf:rw go-ebpf-logger
```