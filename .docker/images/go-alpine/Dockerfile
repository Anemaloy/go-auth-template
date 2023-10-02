FROM alpine:3

RUN apk update && \
  apk add \
    zlib \
    libcrypto1.1 \
    libssl1.1 && \
    apk add --upgrade libcrypto3 &&\
    apk add --upgrade libssl3