include .env
export

migrate-up:
	@migrate -path migrations -database $(CONNECTION_STRING) up

migrate-down:
	@migrate -path migrations -database $(CONNECTION_STRING) down

service-deploy:
	docker-compose up -d

service-undeploy:
	docker-compose down