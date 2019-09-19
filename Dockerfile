FROM golang:1.13-alpine AS builder

RUN apk --no-cache add git make

WORKDIR /src/s3-bucket-permissions-exporter
COPY . .
RUN make build

FROM busybox
LABEL maintainer "Alex Simenduev <shamil.si@gmail.com>"

COPY --from=builder /src/s3-bucket-permissions-exporter/s3-bucket-permissions-exporter /usr/local/bin/
ENTRYPOINT ["s3-bucket-permissions-exporter"]
