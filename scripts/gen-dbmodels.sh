#!/bin/bash
set -a
source .env.local
set +a

# generate your database models
sqlboiler psql --no-hooks