compose:
	rm -rf ./postgres
	docker compose down
	docker compose up --build
sh-postgres:
	docker exec -it postgres sh
sh-go:
	docker exec -it go sh