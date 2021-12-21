set -a
source .env.$ENVIRONMENT_NAME
set +a
sleep 10

sql-migrate status -env postgres

# dropping existing tables
sql-migrate down -env postgres -limit=0

# running migrations
sql-migrate up -env postgres
sql-migrate status -env postgres

./main