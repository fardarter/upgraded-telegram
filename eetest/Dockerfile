FROM golang:1.15.7-buster AS base

LABEL module="eetest"
LABEL project="ee"

WORKDIR /workspace/eetest

COPY Makefile .
COPY vendor vendor
COPY pkg pkg
COPY go.mod go.mod
COPY cmd cmd

ARG GIT_HASH=unknown
ENV HASH_FROM_DOCKER=${GIT_HASH}
RUN make linux

FROM debian:buster-slim

RUN addgroup --gid 2000 eetest && \
    adduser --system --uid 2000 --ingroup eetest eetest
# This will download all certificates (ca-certificates) and builds it in a
# single file under /etc/ssl/certs/ca-certificates.crt (update-ca-certificates)
RUN apt-get update && \
    apt-get install -y ca-certificates openssl && \
    update-ca-certificates

USER eetest:eetest

WORKDIR /usr/bin/
COPY --from=base /workspace/eetest/eetest .

EXPOSE 8080

ENTRYPOINT [ "eetest" ]
