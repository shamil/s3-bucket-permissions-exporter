FROM golang:1.13-alpine AS builder

ARG promu_version=0.5.0

RUN apk --no-cache add git \
    && wget -qO- https://github.com/prometheus/promu/releases/download/v${promu_version}/promu-${promu_version}.linux-amd64.tar.gz | tar -xzvf - -C /opt \
    && ln -s /opt/promu-${promu_version}.linux-amd64/promu /usr/local/bin/promu

WORKDIR /src/s3-bucket-permissions-exporter
COPY . .
RUN promu build

FROM busybox
LABEL maintainer "Alex Simenduev <shamil.si@gmail.com>"

COPY --from=builder /src/s3-bucket-permissions-exporter/s3-bucket-permissions-exporter /usr/local/bin/
ENTRYPOINT ["s3-bucket-permissions-exporter"]
