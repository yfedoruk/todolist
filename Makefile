dev:
	@docker-compose down && \
		docker-compose \
			-f docker-compose.yml \
			-f docker-compose.dev.yml \
			up -d --remove-orphans --build
build:
	docker exec todolist_server_1 go build -o /go/bin/todolist github.com/yfedoruck/todolist/cmd && \
	docker-compose.exe restart server
deb:
	@docker-compose down && \
			docker-compose \
				-f docker-compose.yml \
				-f docker-compose.debug.yml \
				up -d --remove-orphans --build
heroku:
	heroku container:login && \
	heroku container:push --app nameless-brook-78889 web && \
	heroku container:release --app nameless-brook-78889 web && \
	heroku logs --tail --app nameless-brook-78889