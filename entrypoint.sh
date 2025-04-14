#!/bin/bash

#run migrations
echo "Running migrations..."
migrate -path /app/migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" up

if [ $? -ne 0 ]; then
    echo "Migrations failed!"
    exit 1
fi
echo "Migrations applied successfully!"

echo "Starting application..."
exec "$@"