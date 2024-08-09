#!/bin/sh

echo $ENVIRONMENT_NAME

handle_int() {
     echo "SIGINT received, forwarding to child process..."
     kill -INT "$child" 2>/dev/null
     echo "Waiting for child process to exit..."
     wait "$child"
     echo "Child process exited. Waiting for coverage data to be fully written..."
     echo "Exiting after delay..."
     exit 0
 }

 # Trap SIGINT signal   
 trap 'handle_int' INT TERM

./migrations

if [[ $ENVIRONMENT_NAME == "docker" ]]; then
    echo "seeding"
    ./seeder
fi

./server &

child=$!
echo "Started process with PID $child"

 # Wait for child process to finish
wait "$child"