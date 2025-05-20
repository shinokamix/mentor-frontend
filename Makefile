.PHONY: migrate-up migrate-down migrate-auth-up migrate-auth-down migrate-review-up migrate-review-down migrate-mentor-up migrate-mentor-down

## AUTH SERVICE
migrate-auth-up:
	migrate -path ./authorization/migrations \
		-database "postgres://admin:123@localhost:5432/mentors?sslmode=disable" up

migrate-auth-down:
	migrate -path ./authorization/migrations \
		-database "postgres://admin:123@localhost:5432/mentors?sslmode=disable" down

## REVIEW SERVICE
migrate-review-up:
	migrate -path ./review/migrations \
		-database "postgres://admin1:1234@localhost:5433/reviews?sslmode=disable" up

migrate-review-down:
	migrate -path ./review/migrations \
		-database "postgres://admin1:1234@localhost:5433/reviews?sslmode=disable" down

## MENTOR SERVICE
migrate-mentor-up:
	migrate -path ./mentor/migrations \
		-database "postgres://admin3:12345@localhost:5435/mentorss?sslmode=disable" up

migrate-mentor-down:
	migrate -path ./mentor/migrations \
		-database "postgres://admin3:12345@localhost:5435/mentorss?sslmode=disable" down

## ALL
migrate-up-all: migrate-auth-up migrate-review-up migrate-mentor-up

migrate-down-all: migrate-auth-down migrate-review-down migrate-mentor-down
