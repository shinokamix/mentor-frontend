generate-proto:
	protoc --go_out=pkg/api --go_opt=paths=source_relative \
	--go-grpc_out=pkg/api --go-grpc_opt=paths=source_relative \
	proto/review.proto

migrate-up:
	migrate -path ./migrations -database "postgres://admin1:1234@localhost:5433/reviews?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://admin1:1234@localhost:5433/reviews?sslmode=disable" down