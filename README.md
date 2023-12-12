# Prérequis

## Go

Installer la version Go 1.19 (suivre la documentation officielle).

## Redis

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

### Tests unitaires

Pour lancer les tests, utiliser la commande suivante:

```
go test ./...
```

Il est possible d'utiliser le mode verbose:

```
go test ./... -v
```

Il est possible de lancer sans le cache:

```
go test ./... -count=1
```

Lors de l'écriture de tests, les commandes `redis` doivent être évitées au maximum. Il faut privilégier les données ajoutées au script de peuplement.

La convention de nommage des tests est la suivante:

TestXXXReturnsYYYWhenZZZ

XXX étant le nom de la function testée.
YYY étant la donnée renvoyé par la function lors du test.
ZZZ étant la cause du renvoi de cette donnée par le test.

### Tests functionnels

```
hurl --variables-file web/hurl/admin/.env --test web/hurl/admin/*.hurl
```

## Terminal

### Importation de CSV

Pour importer un csv, lancer la commande:

```
go run console/console.go import
```

Par défault, le path est `./web/data/data.csv`. Il est possible de préciser un fichier en utilisant le flag `--file` suivi du chemin du fichier.

### Peuplement

Pour peupler les données, lancer la commande suivante:

```
go run console/console.go populate
```

Les données disponibles sont:

- a sample product with `test` as id, `skutest` as sku et `100.5` as price
- a sample user with `test` as sid and `1` as id
- a sample order with `test` as id
- a sample cart with `test` as id
- a sample expired user with `expired` as sid and `2` as id
- a sample blog article with `1` as id
- a sample blog article with `2` as id
- multiple tags:
  - mens => tshirts, books, clothes
  - womens => tshirts, clothes
  - books => arabic
  - en => womens, men
  - games => kids

Afin de pouvoir utiliser la recherche, il faut lancer le script de migration après chaque peuplement:

```
go run console/console.go migration
```

## Profiter

Siroter un bon café.

## Wiki

### Utiliser les logs

Les logs doivent être renseignés avec le package `slog`. Un contexte doit être passé pour connaître l'identifiant de la requête. Example:

```go
func Add(c context.Context, cid, id string, quantity int64) error {
    // ...
    l.LogAttrs(c, slog.LevelInfo, "adding a product to the cart")
    // ...
}
```

Il est possible de créer un log qui contiendra une valeur utilisée pour chaque log. Example:

```
l := slog.With(slog.String("cid", cid))
```

Des logs doivent être insérés en début et fin de fonction. Pour chaque erreur, il faut logger le message d'erreur. Si l'erreur est d'un type `error`, le niveau de log est `ERROR`, sinon le niveau `INFO` est utilisé. Tous les logs d'erreur doivent commencer par `cannot`.

Les logs doivent être affichés immédiatement dans le code afin d'avoir un contexte précis de l'erreur.

Pour éviter les doubles logs, il ne faut pas faire un log d'une erreur déjà traitée par une de nos fonctions.

### Contexte

Le contexte doit être utilisé dans la majorité des cas (sauf les très petites fonctions), afin d'afficher l'identifiant de la requête et potentiellement d'autres éléments. Les données disponibles dans le contexte sont :

- la langue
- l'utilisateur
- l'identifiant de la requête

### Guideline

#### CSS

Le CSS sémantique est privilié, cela signifie que le style sera définit selon les éléments html, par exemple `button`, `article`..., ou par des attributs, par exemple `aria-busy`...

Les classes doivent définir des fonctionnalités et non des utilitaires. Les différentes parties du nom sont séparées par un tiret. Par exemple `product-card`, `header-menu`, et non pas `md-6 ft-3`.

### Messages

Les messages doivent être prévus pour la traduction. Les erreurs contiennent les clés de traduction qui sont ensuite traduites dans les controlleurs. Les clés sont constituées de trois parties:

- Le type de message
- La fonctionnalité
- La description de l'erreur

Example: `error_http_csrf`, `input_firstname_required`, `url_product_details`

Il est important de garder cette convention de nommage et de garder les clés triées par order alphabétique. Les clés commençant par `input_` sont utilisées pour afficher des erreurs `inline`.

### CSP

Une politique CSP stricte est mise en place, ce qui signifie que les styles et javascripts `inline` ne peuvent être executés sans avoir ajouté une empreinte dans l'entête de la réponse http.  
C'est pour cela que lors du démarrage du serveur, une empreinte et calculée et stockée en mémoire pour chaque script se trouvants dans `web/views/admin/js`. L'empreinte peut être accédée en utilisant `security.CSP["XXX"]`, example:

```
policy := fmt.Sprintf("default-src 'self'; script-src-elem 'self' %s;", security.CSP["otp"])
w.Header().Set("Content-Security-Policy", policy)
```
