set -a
source .env.$ENVIRONMENT_NAME
set +a
sleep 20
echo $ENVIRONMENT_NAME
sql-migrate status -env mysql

# dropping existing tables
# sql-migrate down -env mysql -limit=0

# running migrations
sql-migrate up -env mysql
sql-migrate status -env mysql


if [[ $ENVIRONMENT_NAME == "docker" ]]; then
    echo "seeding"
    ./scripts/seed.sh
fi

./main