#!/bin/bash

# Variables
DOCKER_IMAGE_NAME="torlinks_builder"
PROJECT_DIR="/home/raoul/go/torlinks"
BUILD_DIR="/home/raoul/go/torlinks/build"

# Construire l'image Docker localement
docker build -t $DOCKER_IMAGE_NAME $PROJECT_DIR

# Exécuter le conteneur pour compiler le programme
docker run --rm -v $BUILD_DIR:/app/builds $DOCKER_IMAGE_NAME /bin/bash -c "ls -l /app && go build -buildvcs=false -o /app/build/torlinks /app/cmd/torlinks && ls -l /app/build"

echo "Compilation terminée. Le fichier de sortie est placé dans $BUILD_DIR"
