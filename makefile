compose: swagger-generate
	go mod tidy
	rm -rf ./postgres
	docker compose down
	docker compose up --build

swagger-generate:
	swag init --parseDependency -d "./cmd/server,./internal/routes,./internal/db,./"

sh-postgres:
	docker exec -it postgres sh
sh-go:
	docker exec -it go sh
