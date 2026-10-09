module payments

go 1.26.5

require (
	github.com/bwmarrin/snowflake v0.3.0
	github.com/go-kratos/kratos/contrib/otel/v3 v3.0.0-20260626125723-668db92c2c00
	github.com/go-kratos/kratos/v3 v3.0.0
	github.com/google/wire v0.7.0
	github.com/nats-io/nats.go v1.53.1
	go.einride.tech/aip v0.86.3
	go.uber.org/automaxprocs v1.6.0
	google.golang.org/genproto/googleapis/api v0.0.0-20260810153831-ec0a7760b754
	google.golang.org/grpc v1.83.0
	google.golang.org/protobuf v1.36.12
	gorm.io/driver/postgres v1.6.2
	gorm.io/gorm v1.31.2
	transaction v0.0.0-00010101000000-000000000000
)

require (
	github.com/klauspost/compress v1.18.5 // indirect
	github.com/nats-io/nkeys v0.4.15 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.51.0 // indirect
)

require (
	account v0.0.0
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/fsnotify/fsnotify v1.10.1 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/form/v4 v4.3.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.10.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.44.0 // indirect
	go.opentelemetry.io/otel/metric v1.44.0 // indirect
	go.opentelemetry.io/otel/trace v1.44.0 // indirect
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260807164820-c8921c73eeea // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace account => ../account

replace transaction => ../transaction
