generate-proto:
	protoc --go_out=pkg/api --go_opt=paths=source_relative \
	--go-grpc_out=pkg/api --go-grpc_opt=paths=source_relative \
	proto/mentor.proto

migrate-up:
	migrate -path ./migrations -database "postgres://admin:123@localhost:5432/mentors?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://admin:123@localhost:5432/mentors?sslmode=disable" down