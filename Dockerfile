FROM debian:stable-slim

WORKDIR /app

COPY dist/go-ebpf-logger_linux_amd64_v1/go-ebpf-logger /app/go-ebpf-logger

COPY monitor.bpf.o /app/monitor.bpf.o

RUN chmod +x /app/go-ebpf-logger

ENTRYPOINT ["/app/go-ebpf-logger"]