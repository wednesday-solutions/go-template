local:
	docker-compose --env-file ./.env.local down

	docker-compose --env-file ./.env.local \
	-f docker-compose.yml \
	-f docker-compose.local.yml build

	docker-compose --env-file ./.env.local \
	-f docker-compose.yml \
	-f docker-compose.local.yml up -d

	docker exec -it go-template_server_1 /bin/bash

prod:
	docker-compose --env-file ./.env.prod down

	docker-compose --env-file ./.env.prod \
	-f docker-compose.yml \
	-f docker-compose.prod.yml build

	docker-compose --env-file ./.env.prod \
	-f docker-compose.yml \
	-f docker-compose.prod.yml up -d

test:
	docker-compose --env-file ./.env.test down

	docker-compose --env-file ./.env.test \
	-f docker-compose.yml \
	-f docker-compose.test.yml build

	docker-compose --env-file ./.env.test \
	-f docker-compose.yml \
	-f docker-compose.test.yml up -d

logs:
	docker-compose --env-file ./.env.$(env) logs -f

tear:
	docker-compose --env-file ./.env.$(env) down