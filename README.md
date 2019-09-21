# S3 Bucket Permission Exporter

A Prometheus Exporter for gathering status from TrustedAdvisor's S3 Bucket Permissions check.

## Usage

```shell
usage: s3-bucket-permissions-exporter [<flags>]

Flags:
  -h, --help              Show context-sensitive help (also try --help-long and --help-man).
  -l, --web.listen-address=":9199"
                          Address to listen on for web interface and telemetry.
      --web.telemetry-path="/metrics"
                          Path under which to expose metrics.
      --collector.ignored-buckets="^$"
                          Regexp of buckets to ignore from collecting.
      --log.level="info"  Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]
      --log.format="logger:stderr"
                          Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or "logger:stdout?json=true"
      --version           Show application version.
```

## Run with Docker

```shell
docker run -p 9199:9199 simenduev/s3-bucket-permissions-exporter
```

## AWS IAM Policy

The following IAM policy required for the exporter to be able to scrape from TrustedAdvisor API

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "support:*"
            ],
            "Resource": [
                "*"
            ]
        }
    ]
}
```

## License

Licensed under the MIT License. See the `LICENSE` file for details.
