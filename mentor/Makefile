generate-proto:
	protoc --go_out=pkg/api --go_opt=paths=source_relative \
	--go-grpc_out=pkg/api --go-grpc_opt=paths=source_relative \
	proto/mentor.proto

migrate-up:
	migrate -path ./migrations -database "postgres://admin3:12345@localhost:5435/mentorss?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://admin3:12345@localhost:5435/mentorss?sslmode=disable" down