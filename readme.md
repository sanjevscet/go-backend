go mod init github.com/sanjevscet/go-backend.git

go install github.com/air-verse/air@latest

air init

go get github.com/lib/pq

brew install golang-migrate

// to remove docker with volume
docker compose down -v

migrate create -seq -ext sql -dir ./cmd/migrate/migrations create_users
migrate create -seq -ext sql -dir ./cmd/migrate/migrations create_posts

migrate -path ./cmd/migrate/migrations -database "postgres://sanjeev:sanjeev@localhost:11432/social?sslmode=disable" up  
migrate -path ./cmd/migrate/migrations -database "postgres://sanjeev:sanjeev@localhost:11432/social?sslmode=disable" down
