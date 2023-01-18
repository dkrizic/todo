module github.com/dkrizic/todo/echo

go 1.19

require (
	github.com/cloudevents/sdk-go/v2 v2.13.0
	github.com/dkrizic/todo/api/todo v0.0.0-20230110224611-543385c5427d
	github.com/gorilla/mux v1.8.0
	github.com/pytimer/mux-logrus v0.0.0-20200505085744-ce5a5e748151
	github.com/sirupsen/logrus v1.9.0
	go.opentelemetry.io/otel/exporters/jaeger v1.11.2
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.14.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	go.opentelemetry.io/otel v1.11.2 // indirect
	go.opentelemetry.io/otel/sdk v1.11.2 // indirect
	go.opentelemetry.io/otel/trace v1.11.2 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/net v0.2.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	google.golang.org/genproto v0.0.0-20221118155620-16455021b5e6 // indirect
	google.golang.org/grpc v1.51.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace github.com/dkrizic/todo/api/todo => ../api/todo
