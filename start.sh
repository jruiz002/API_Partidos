#!/bin/sh

# Esperar hasta que la base de datos esté lista
until pg_isready -h db -U postgres; do
  echo "Esperando por la base de datos..."
  sleep 2
done

# Ejecutar la aplicación
echo "Base de datos lista, arrancando la API..."
go run main.go
