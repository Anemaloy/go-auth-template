FROM golang:1.19.3-alpine3.15 as builder

RUN apk add --update --no-cache build-base curl git && \
  rm -rf /var/cache/apk/*

ENV \
  GOLANG_PROTOBUF_VERSION=1.31.0 \
  GOLANG_PROTOBUF_GRPC_VERSION=1.2.0
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v${GOLANG_PROTOBUF_VERSION} && \
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${GOLANG_PROTOBUF_GRPC_VERSION} && \
  mv /go/bin/protoc-gen-go* /usr/local/bin/

ENV GEN_VALIDATE_VERSION=0.9.0
RUN git clone --depth 1 -b v${GEN_VALIDATE_VERSION} https://github.com/envoyproxy/protoc-gen-validate.git /tmp/protoc-gen-validate && \
  cd /tmp/protoc-gen-validate && go build && \
  mv /tmp/protoc-gen-validate/protoc-gen-validate /usr/local/bin/protoc-gen-validate && \
  mv /tmp/protoc-gen-validate/validate /usr/local/include

FROM bufbuild/buf:1.9.0

WORKDIR /workspace

COPY --from=builder /usr/local/bin /usr/local/bin
