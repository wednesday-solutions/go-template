set -a
source .env.local
set +a
sql-migrate status -env postgres

# dropping existing tables
sql-migrate down -env postgres -limit=0

# running migrations
sql-migrate up -env postgres
sql-migrate status -env postgres
