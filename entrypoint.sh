#!/bin/bash

# Clean the coverage-reports folder
rm -rf ./coverage-reports/*

# Run your application or migration script
exec bash ./migrate-and-run.sh
