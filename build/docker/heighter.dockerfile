FROM golang:alpine as builder

WORKDIR /app

RUN apk add --no-cache ca-certificates
RUN apk add --no-cache make
RUN apk add --no-cache protobuf protobuf-dev
RUN apk add --no-cache protoc
RUN apk add --no-cache git
RUN apk add --no-cache tzdata

COPY go.mod go.sum .
RUN go mod download
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

COPY . /app

RUN make protos
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/heighter

FROM scratch

WORKDIR /app

COPY --from=builder /app/app /app/app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT [ "./app" ]
