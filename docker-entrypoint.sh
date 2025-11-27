#!/bin/sh
set -e

# Esperar a que la base de datos esté lista (aunque depends_on ayuda, esto es más robusto)
echo "Esperando a PostgreSQL en $DB_HOST..."
until PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c '\q'; do
  echo "Unavailable - sleeping"
  sleep 1
done
echo "PostgreSQL está listo."

# Aplicar Schema
echo "Aplicando esquema..."
PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f /app/db/schema/schema.sql

# Poblar datos (Seed)
echo "Poblando base de datos..."
PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" <<EOF
INSERT INTO users (username, email, password) VALUES 
    ('agus', 'agus2@gmail.com', '123456')
ON CONFLICT (username) DO NOTHING;
EOF

echo "Iniciando aplicación..."
exec "$@"
