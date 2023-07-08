# Prérequis

## Go

Installer la version Go 1.19 (suivre la documentation officielle).

## Redis

Deux choix sont possibles:

1- Installer [Redis](https://redis.io/docs/getting-started/installation/install-redis-on-linux)  
2- Utiliser Docker:

```
docker run -p 6379:6379 -v data:/data redis/redis-stack:latest
```

Installer aussi [Redis Insight](https://redis.com/redis-enterprise/redis-insight) pour explorer facilement les données.

## VSCode

Installer VSCode avec les extensions [Go](https://marketplace.visualstudio.com/items?itemName=golang.Go) et [Prettier](https://marketplace.visualstudio.com/items?itemName=esbenp.prettier-vscode).

# Lancement

## Dépendances

Récupérer les dépendances à l'aide de la commande suivante:

```sh
go get
```

## Lancer le serveur

```
go run main.go
```

Le serveur est démarré sur le port `8080` par défault.

## Tester

Pour lancer les tests, utiliser la commande suivante:

```
go test ./...
```

Il est possible d'utiliser le mode verbose:

```
go test ./... -v
```

## Terminal

### Importation de CSV

Pour importer un csv, lancer la commande:

```
go run terminal/hop.go
```

Par défault, le path est `./static/fake/data.csv`. Il est possible de préciser un fichier en utilisant le flag `--file` suivi du chemin du fichier.

## Profiter

Siroter un bon café.
