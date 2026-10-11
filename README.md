# Fintech Platform

A set of backend microservices for a small fintech system, written in **Go** with the
**Kratos** framework. Each service exposes the same API over **gRPC** and **HTTP**, defined
first in **Protocol Buffers**.

## Services

| Service | HTTP | gRPC | What it manages |
|---|---|---|---|
| `customer` | 8009 | 9009 | Customers |
| `account` | 8007 | 9007 | Accounts |
| `card` | 8001 | 9001 | Cards |
| `payments` | 8003 | 9003 | Payments (publishes and consumes NATS events) |
| `transaction` | 8004 | 9004 | Transactions |

Other folders:

- `proto/`: the shared Protobuf API definitions for all services (managed with Buf).
- `nats/`: a small publisher and subscriber used to try out the messaging flow.

## How a service is built

Every service follows the Kratos layout:

```
cmd/            entry point and Wire dependency injection
api/            generated gRPC, HTTP and OpenAPI code
internal/
  service/      transport layer: implements the Protobuf service
  biz/          business logic and use cases
  data/         repositories (GORM on PostgreSQL)
  server/       gRPC and HTTP server setup
  conf/         typed configuration
configs/        config.yaml
```

## Tech

Go · Kratos · gRPC · Protocol Buffers (Buf) · Google Wire · GORM · PostgreSQL · NATS · Snowflake IDs · Docker

## Run a service

You need Go, PostgreSQL and (for `payments`) a NATS server running locally.

```sh
cd payments
make build          # or: go build -o ./bin/ ./...
./bin/payment -conf ./configs
```

Database and port settings are in each service's `configs/config.yaml`.
