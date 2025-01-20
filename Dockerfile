FROM golang:1.23

WORKDIR /app

COPY . /app

ENV GO111MODULE=on \
CGO_ENABLED=0 \
GOPROXY=https://proxy.golang.org,direct

RUN apt-get update && apt-get install -y \
    vim \
    git \
    curl \
    && rm -rf /var/lib/apt/lists/*

RUN go install golang.org/x/tools/cmd/goimports@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \

EXPOSE 8083

CMD ["bash"]