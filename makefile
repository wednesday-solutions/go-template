local:
	docker-compose --env-file ./.env.local \
	-f docker-compose.yml \
	-f docker-compose.local.yml down

	docker-compose --env-file ./.env.local \
	-f docker-compose.yml \
	-f docker-compose.local.yml build

	docker-compose --env-file ./.env.local \
	-f docker-compose.yml \
	-f docker-compose.local.yml up -d

	docker exec -it go-template_server_1 /bin/bash

prod:
	docker-compose --env-file ./.env.prod \
	-f docker-compose.yml \
	-f docker-compose.prod.yml down

	docker-compose --env-file ./.env.prod \
	-f docker-compose.yml \
	-f docker-compose.prod.yml build

	docker-compose --env-file ./.env.prod \
	-f docker-compose.yml \
	-f docker-compose.prod.yml up -d

test:
	docker-compose --env-file ./.env.test \
	-f docker-compose.test.yml down

	docker-compose --env-file ./.env.test \
	-f docker-compose.test.yml build

	docker-compose --env-file ./.env.test \
	-f docker-compose.test.yml up -d

	docker exec -it go-template_server_1 ./test.sh

logs:
	docker-compose --env-file ./.env.$(env) logs -f

tear:
	if [ $(env) = "test" ]; then \
		docker-compose --env-file ./.env.$(env) \
		-f docker-compose.$(env).yml down; \
	else \
		docker-compose --env-file ./.env.$(env) \
		-f docker-compose.yml \
		-f docker-compose.$(env).yml down; \
	fi
