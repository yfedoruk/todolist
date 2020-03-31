dev:
	@docker-compose down && \
		docker-compose \
			-f docker-compose.yml \
			-f docker-compose.dev.yml \
			up -d --remove-orphans
build:
	docker exec todolist_server_1 go install && \
	docker-compose.exe restart server
heroku:
	heroku container:login && \
	heroku container:push --app nameless-brook-78889 web && \
	heroku container:release --app nameless-brook-78889 web && \
	heroku logs --tail --app nameless-brook-78889