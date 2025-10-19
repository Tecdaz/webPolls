#!/bin/bash

echo "Eliminando contenedores anteriores..."
docker rm -f $(docker ps -a -q --filter "name=webpolls-*")

echo "Construyendo imagen..."
docker compose build --no-cache

echo " Levantando contenedores..."
docker compose up -d

echo "Esperando a que la aplicación se inicie..."
while ! curl -s --fail http://localhost:8080/ > /dev/null; do
    echo -n "."
    sleep 1
done
echo "La aplicación se ha iniciado correctamente."

echo "
Ejecutando pruebas..."
hurl --test tests/*
