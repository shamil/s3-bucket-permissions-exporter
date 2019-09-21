FROM golang:1.13-alpine AS builder

RUN apk --no-cache add git make

WORKDIR /src/s3-bucket-permissions-exporter
COPY . .
RUN make build

FROM busybox
LABEL maintainer "Alex Simenduev <shamil.si@gmail.com>"

ENTRYPOINT ["s3-bucket-permissions-exporter"]

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /src/s3-bucket-permissions-exporter/s3-bucket-permissions-exporter /usr/local/bin/
