docker-down: 
	docker-compose --env-file ./.env.docker \
	-f docker-compose.yml \
	-f docker-compose.yml down

docker-build: 
	docker-compose --env-file ./.env.docker \
	-f docker-compose.yml \
	-f docker-compose.yml build

docker-up: 
	docker-compose --env-file ./.env.docker \
	-f docker-compose.yml \
	-f docker-compose.yml up

docker: docker-down docker-build docker-up

tests:
	./scripts/test.sh

setup-local: init
	./scripts/setup-local.sh

setup-precommit:
	./scripts/setup-pre-commit.sh

init: setup-precommit
	go install github.com/volatiletech/sqlboiler/v4@latest
	go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest
	go install github.com/rubenv/sql-migrate/...@latest
	go mod vendor; go mod download; go mod tidy;

setup-ecs:
	./scripts/setup-ecs.sh $(name) $(env)

update-ecs:
	./scripts/update-ecs.sh $(name) $(env)

deploy-ecs:
	./scripts/deploy-ecs.sh $(name) $(env)