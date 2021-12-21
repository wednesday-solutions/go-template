local:
	docker-compose --env-file ./.env.local \
	-f docker-compose.yml \
	-f docker-compose.yml down

	docker-compose --env-file ./.env.local \
	-f docker-compose.yml \
	-f docker-compose.yml build

	docker-compose --env-file ./.env.local \
	-f docker-compose.yml \
	-f docker-compose.yml up