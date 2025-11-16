#!/bin/bash
echo "👤 Creando usuario de prueba..."

docker exec -i webpolls-postgres psql -U postgres -d webpolls << 'EOF'
INSERT INTO users (username, email, password) VALUES 
    ('agus', 'agus2@gmail.com', '123456')
ON CONFLICT (username) DO NOTHING;
EOF

echo "✅ Usuario 'agus' creado"
