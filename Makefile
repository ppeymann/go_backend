swagger:
	swag init --parseDependency --parseInternal -g /server/server.go

build:
	docker build --tag example:latest .

compose:
	docker-compose up -d

createdb:
	docker exec -it eg_postgres createdb --username=postgres --owner=postgres eg_db

dropdb:
	docker exec -it eg_postgres dropdb eg_db