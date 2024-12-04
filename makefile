compose:
	rm -rf ./postgres
	docker compose down
	docker compose up --build
sh-postgres:
	docker exec -it postgres sh
sh-go:
	docker exec -it go sh
test-post:
	curl -i --header "Content-Type: application/json" \
	--request POST \
	--data '{"group":"ballin","name":"ballers","releaseDate":"2023-12-12","text":"texts","link":"rer"}' \
	localhost:8000/api/v1/songs