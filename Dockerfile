# Utiliser une image de base contenant Go
FROM golang:1.19

# Définir le répertoire de travail à l'intérieur du conteneur
WORKDIR /app

# Copier le contenu du projet dans le conteneur
COPY . /app

# Lister les fichiers dans le répertoire source pour vérifier la copie
RUN echo "Contenu du répertoire /app avant compilation:" && ls -l /app

# Construire le programme
RUN go build -buildvcs=false -o /app/build/torlinks /app/cmd/torlinks

# Lister les fichiers dans le répertoire builds pour vérifier la compilation
RUN echo "Contenu du répertoire /app/builds après compilation:" &&  ls -l /app/build

# Commande par défaut pour le conteneur
CMD ["/bin/bash"]

