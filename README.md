migrate create -ext sql -dir migrations -seq init

migrate -database "postgres://127.0.0.1:5434/prod_service?sslmode=disable&user=postgres&password=postgres" -path ./migrations force 1

migrate -database "postgres://127.0.0.1:5434/prod_service?sslmode=disable&user=postgres&password=postgres" -path ./migrations up