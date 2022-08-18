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
	
