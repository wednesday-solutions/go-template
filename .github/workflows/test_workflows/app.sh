# Build the project locally
sudo docker-compose --env-file ./.env.keploy -f docker-compose.yml build

echo "Project built successfully"

docker volume ls

# Start the app in the background
sudo docker-compose --env-file ./.env.keploy -f docker-compose.yml up &

# Wait for the service to be ready
echo "Waiting for the service to be ready..."
until curl -s -o /dev/null -w "%{http_code}" http://0.0.0.0:9000/graphql/ | grep -q "200"; do
  echo "Service is not ready yet. Waiting..."
  sleep 5
done

echo "Service is ready!"

# Perform curl requests
curl --request POST \
  --url http://localhost:9000/graphql \
  --header 'Host: localhost:9000' \
  --header 'Accept-Encoding: gzip, deflate, br' \
  --header 'Connection: keep-alive' \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: PostmanRuntime/7.39.0' \
  --header 'Accept: */*' \
  --header 'Postman-Token: 146a5779-ec95-40d3-aeb3-eec5e046c7d1' \
  --data '{"query":"mutation Login($loginUsername: String!, $loginPassword: String!) {\n  login(username: $loginUsername, password: $loginPassword) {\n    token\n    refreshToken\n  }\n}\n","variables":{
  "loginUsername": "khareyash05",
  "loginPassword": "password123"
}
}'

# Bring down the app
sudo docker-compose --env-file ./.env.keploy -f docker-compose.yml down