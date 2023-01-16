echo "name of the service is: $1 $2"

copilot init -a "$1" -t "Load Balanced Web Service" -n "$1-$2-svc" -d ./Dockerfile

copilot env init --name $2 --profile default --default-config

copilot env deploy --name $2

copilot storage init -n "$1-$2-cluster" -t Aurora -w "$1-$2-svc" --engine PostgreSQL --initial-db "$1_$2_db"

copilot storage init -n "$1-$2-bucket" -t S3 

aws s3 cp .env "s3://$1-$2-bucket/develop/"

# copilot deploy --name "$1-$2-svc" -e "$2"
