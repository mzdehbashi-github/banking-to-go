#!/bin/sh

# Function to wait for the database to be available
waitForDB() {
    retries=5
    timeout=1
    while ! nc -z -w 1 $DB_HOST 5432; do
        if [ $retries -eq 0 ]; then
            echo "Error: Database not available after multiple attempts, exiting."
            exit 1
        fi
        sleep $timeout
        timeout=$((timeout * 2))
        retries=$((retries - 1))
    done
}

# Command to execute migrations
migrateup() {
    waitForDB
    migrate -path db/migrations/ -database "$DB_SOURCE" -verbose up
}

# Command to start the application
startapp() {
    waitForDB
    /app/main
}

# Check for the provided command
case "$1" in
    migrateup)
        migrateup
        ;;
    startapp)
        startapp
        ;;
    *)
        echo "Usage: $0 {migrateup|startapp}"
        exit 1
        ;;
esac
