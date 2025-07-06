echo "Enter version of file monitor to use with docker (e.g. 1.0.0):"

if [ -z "$1" ]; then
  read -r version
else
  version="$1"
fi

echo "Building Docker image with version: $version"

if ! command -v gh &> /dev/null; then
    echo "gh CLI is not installed. Please install it first and proceed with authentication."
    exit 1
fi

if ! gh auth status &> /dev/null; then
    echo "You are not authenticated with GitHub CLI. Please authenticate first."
    exit 1
fi

gh release download v${version} --repo SentinalFS/file-monitor --clobber --pattern "monitor.bpf.o" 

if [ $? -ne 0 ]; then
    echo "Failed to download the file monitor binary. Please check the version and try again."
    exit 1
fi

echo "Running go releaser"

if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go and try again."
    exit 1
fi

if ! command -v goreleaser &> /dev/null; then
    echo "Goreleaser is not installed. Please install Goreleaser and try again."
    exit 1
fi

goreleaser release --snapshot --skip=publish --clean 


if docker --version &> /dev/null; then
    echo "Docker is installed. Proceeding with build."
else
    echo "Docker is not installed. Please install Docker and try again."
    exit 1
fi

arch=$(uname -m)
if [ "$arch" = "aarch64" ] || [ "$arch" = "arm64" ]; then
    echo "Detected ARM architecture. Building Docker image for ARM."
    docker build -t go-logger-arm:latest -f Dockerfile.amd64 --build-arg TARGETARCH=arm64 .
    if [ $? -ne 0 ]; then
        echo "Docker build for ARM failed."
        exit 1
    fi
elif [ "$arch" = "x86_64" ]; then
    echo "Detected x86_64 architecture. Building Docker image for x86_64."
    docker build -t go-logger-amd64:latest -f Dockerfile.amd64 --build-arg TARGETARCH=amd64 .
    if [ $? -ne 0 ]; then
        echo "Docker build for x86_64 failed."
        exit 1
    fi
else
    echo "Non-ARM architecture detected. No additional ARM build required."
fi