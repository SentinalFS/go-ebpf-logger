FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Download the precompiled BPF object: 
# TODO: Figure out wget for object file
# RUN wget <url> -O monitor.bpf.o

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /listener ./cmd/listener

ENTRYPOINT ["/listener", "monitor.bpf.o"]
