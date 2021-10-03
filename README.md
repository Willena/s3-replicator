# S3 Replicator

A simple service that is used to replicate and maintain a backup up to date using event from S3 bucket.

## Feature

- Replicate every action to another S3
- Select event that should not be applied
- Replicate to another S3 types (glacier)
- Read from multiple message broker types
    - Kafka
    - AMQP
    - ...

## Configure

```yml
s3:
  source:
    url: ""
    bucket: ""
  destination:
    url: ""
    bucket: ""
amqp:
  url: "<endpoint>"
  exchange: "<string>"
  exchange_type: "<string>"
  routing_key: "<string>"
  mandatory: "<string>"
  durable: "<string>"
  no_wait: "<string>"
  internal: "<string>"
  auto_deleted: "<string>"
  delivery_mode: "<string>"
  queue_dir: "<string>"
  queue_limit: "<string>"
kafka:
  broker-list: "abc:555,abc:66"
  topic:
  consumer:
    props: "ddd"
```

