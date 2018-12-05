# Builder container
FROM golang:1.11 AS builder
WORKDIR /go/src/app

# We could use dep to manage dependencies here
# for now we're using the simplest possible example.
# curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
# COPY Gopkg.lock Gopkg.toml ./
# RUN dep ensure

# Copy source and compile
COPY *.go ./
RUN go get
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo -ldflags="-w -s" -o main

# Executable container
FROM scratch
WORKDIR /app

EXPOSE 8080

# Copy executable and run
COPY --from=builder /go/src/app/main ./
ENTRYPOINT [ "./main" ]