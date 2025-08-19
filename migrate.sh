#!/bin/bash

set -e

if [ -f .env ]; then
    echo "📄 Loading environment variables..."
    set -a
    source .env
    set +a
fi

echo "📦 Creating database if not exists..."
psql $DATABASE_URL -tc "SELECT 1 FROM pg_database WHERE datname = 'ingresso_go'" | grep -q 1 || psql -c "CREATE DATABASE ingresso_go"

echo "🔨 Running migrations..."
psql $DATABASE_URL -f db/init.sql

echo "✅ Database OK"
