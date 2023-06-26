# Spécifications techniques de Cadeau

# 1 Présentation

Le​ ​but​ ​de​ ​ce​ ​document​ ​est​ ​de proposer des spécifications techniques pour le projet “cadeau”.

# 2 Représentants

​Arnaud Deville, Salim Tison et Reda Madjoub sont les développeurs de cette application.

# 3 Versions

| Version |    Date    |         Auteur |
| ------- | :--------: | -------------: |
| 1       | 11/06/2023 | Arnaud Deville |

# 4 Stack technique

Le code est réalisé avec _GO_.

Les interactions entre le client et le serveur sont réalisées à l’aide de HTMX. Lorsqu’une requête est déclenchée, une icône de chargement s’affiche dans le bouton pour indiquer qu'une requête est en cours.

Lorsque des interactions concernent uniquement le client, du code JavaScript peut être ajouté (en respectant les contraintes de sécurité). Cependant, cela doit être utilisé en dernier recours et doit être considéré comme une amélioration de l’interface client, non pas une fonctionnalité requise.

Le code du socle serveur doit être écrit dans l'optique d'être utilisé dans le cadre de plusieurs projets, e-commerce ou marketplace, en utilisant une configuration spécifique et sans modification du code source. Cependant, les templates HTML sont spécifiques à chaque site e-commerce.

Le style est en CSS sans préprocesseur (de type SASS ou POSTCSS), à l’aide du framework PicoCSS, permettant de produire un style compatible avec du code HTML sémantique.

Redis est utilisé en tant que base de données, notamment pour sa simplicité et ses performances.

Imgproxy est utilisé pour servir les images des produits.

Chaque requête contient un numéro unique.

**Remarque**: Il n’est pas question de construire ici un SPA, afin de ne pas augmenter (grandement) la complexité de l’application. Cependant, un PWA pourra être mis en place à l’aide du fichier manifest.json.

## 4.1 Redis

**Remarque**: Il est important de consulter la documentation de Redis afin de mieux comprendre son fonctionnement.

Le stockage dans Redis se fait à travers une association clé et valeur.

Les clés peuvent être composées de séparateur (par défaut: les deux points `:` ).

Les valeurs peuvent être des nombres ou des chaînes de caractères. Les objets sont stockés à l'aide de `HASH`.

Si plusieurs commandes sont nécessaires pour réaliser une action, elles doivent être réalisées à travers une transaction.

## 4.2 Architecture

L'architecture est découpée en trois couches, dont les niveaux respectifs du plus haut au plus bas sont:

- Routes: Traitement de la requête, appels des services...etc
- Services: Logique business
- Repository: Interactions avec la base de données

Une couche ne peut appeler que la couche directement en dessous d’elle. Une couche ne peut pas appeler les couches au-dessus d’elle.

# 5 Description technique

## 4.1 Layout

La layout comprend une entête disposant des éléments suivants:

- Un lien vers la page _Se connecter_
- Un lien vers la page de contact
- Une section recherche si nécessaire. Cette dernière comprendrait uniquement un champ texte qui sera utilisé par Redis Search en regardant uniquement les titres et descriptions des produits.
- Une icône panier qui redirige vers le panier et qui affiche un compteur correspondant au nombre de produits dans le panier

Le pied de page contient les liens vers les pages statiques.

## 4.1 Page d’accueil

La page d’accueil affiche les X produits les plus récents. X étant à définir dans la configuration. Les produits sont récupérés à l'aide de Redis Search. Le clic sur un élément renvoie sur la page de détails du produit.

## 4.2 Liste des produits

La page d’accueil affiche les X produits les plus récents. X étant à définir dans la configuration. Les produits sont récupérés à l'aide de Redis Search. Le clic sur un élément renvoie sur la page de détails du produit.

L’application peut activer le filtre par tags à l’aide d’une configuration.

La pagination est gérée à l'aide d'un bouton qui, lors du clic sur dernier, lance une requête vers le serveur à l’aide de HTMX, en incrémentant la page courante.

## 4.3 Détail d’un produit

Tous les champs du produit présents dans le CSV sont affichés dans le détail.

Les options sont affichées dynamiquement avec le nom de l’option à gauche et sa valeur à droite.

Un champ numérique _quantité_ est présent, avec deux boutons plus et moins pour ajuster cette dernière.

Un bouton permet d’ajouter le produit au panier. Lors du clic sur ce dernier, une requête HTMX est envoyée vers le serveur pour ajuster le panier. Lorsque la requête est finalisée, le bouton est désactivé pour quelques secondes et le compteur du panier est mis à jour.

Un autre bouton permet d’accéder au panier.

Les produits qui sont liés ([voir](#51-importation-csv-de-produits)), sont affichés en liste avec la photo et le titre. Le clic sur un produit redirige vers le détail de ce dernier.

Si le produit n'est pas en ligne, une page de type `404` est affichée.

Si le produit n'est plus en stock, il n'est pas possible d'ajouter le produit au panier et le bouton est grisé.

## 4.4 Panier

Le panier affiche la liste des produits qui ont été ajoutés par l’utilisateur. Lorsque ce dernier souhaite ajouter un produit dans un panier, le serveur va d’abord vérifier qu’un cookie contenant l’identifiant du panier `cartID`, une sorte de session, existe. S’il n’existe pas, il est créé et stocké dans les cookies.

Le préfixe utilisé pour la clé de stockage est: `cart`.

La clé de stockage est la combinaison du préfixe et du `cartID`. _Example: cart:cartID_.

Les produits sont stockés sous la forme de hash dont la clé est la combinaise du préfixe et du `cartID`. Le hash a la valeur du `PID` ([voir](#5.1-Importation-CSV-de-produits)) et sa valeur est la quantité. _Example: cart:12331 1221FD3X3_.

Si la configuration précise une durée de vie du panier, la commande `EXPIRE` de Redis sera utilisée. Dans ce cas, l'expiration sera rafraîchie à chaque nouvelle requête.

Lors de l’affichage de la page du panier, tous les identifiants et quantités sont récupérés dans Redis, puis pour chaque produit, les détails sont récupérés. Le total du panier est aussi calculé et affiché.

Le bouton permettant de valider la commande redirige sur la page de saisie de l’adresse de livraison.

**Remarque** On considère que le paiement d'une commande ne peut contenir que les produits d’une même devise.

## 4.5 Paiements

La paiement commence par la saisie de l’adresse de facturation avec les champs suivants:

- firstname
- lastname
- address
- complementary
- zipcode
- city
- email
- phone

Tous les champs sont requis, à part l’adresse complémentaire et le numéro de téléphone.

Si l’utilisateur est connecté, ces champs sont pré-remplis.

Si l'utilisateur n'est pas connecté, l'application lui propose de créer un compte en saisissant son mot de passe. Le canal de notification par email sera utilisé par défaut.
Si l'utilisateur ne crée pas de compte, sa commande sera traitée comme une commande `invitée`.

Après validation, si l’application propose le retrait sur place, l’utilisateur peut choisir entre ce mode et la livraison à domicile. Sinon cette dernière sera automatiquement sélectionnée. Si la livraison à domicile est sélectionnée, un écran lui propose d’utiliser les mêmes coordonnées que les données de facturation. S’il refuse, il peut saisir tous les champs mentionnés précédemment pour son adresse de livraison.

Après validation du mode de livraison, l'utilisateur sélectionne son mode de paiement. La liste des moyens de paiement pourra être configurable, voici une liste non exhaustive:

- Espèce
- Carte bleue
- Virement
- Bitcoin

Après validation du paiement, deux cas sont possibles:

- **Synchrone**: Le paiement est réalisé de façon synchrone. Après avoir obtenu la confirmation du paiement, un numéro de commande est généré et le statut de la commande est `payment_validated`.
- **Asynchrone**: Le paiement est réalisé de façon asynchrone. Après avoir effectué le paiement, un numéro de commande est généré et le statut de la commande est `payment_in_progress`.

Lorsque la commande est terminée, l'écran de confirmation affiche le numéro de commande. Si l'application autorise les PUSH notifications et qu'elles n'ont pas encore été proposées à l’utilisateur, un bouton s'affiche pour qu'il puisse en bénéficier. Après validation des permissions, le jeton récupéré est ajouté aux données de la commande dans Redis.

Les commandes stockées dans Redis contiennent les mêmes éléments du panier avec le statut en plus. Le panier est ensuite supprimé de Redis. Les identifiants de commande sont stockées dans un _sorted set_ dont le score est le _timestamp_.

## 4.6 Compte utilisateur

Un utilisateur peut se connecteur à l'aide du couple `email` et `password`. S'il a oublié son mot de passe, il peut cliquer sur un lien pour récupérer son mot de passe. Il recevra alors un email et/ou une PUSH notification contenant un lien magique qui, lors du clic, le redirige sur un formulaire de mise à jour de son mot de passe. Une fois cela réalisé, il est automatiquement connecté.

Il peut aussi modifier ses coordonnées de facturation et livraison, changer activer ou désactiver les PUSH notifications et consulter l'historique des commandes. Ce dernier affiche les éléments suivants:

- Numéro de la commande
- Date de la commande
- Prix de la commande
- Statut de la commande
- Un lien vers le détail de cette commande

Lors du clic sur le détail de la commande, l'utilisateur voit la liste des produits contenu dans cette commande, ainsi que les éventuelles notes ajoutées.

Le préfixe utilisé pour stocker les utilisateur est `user`. Chaque utilisateur possède un identifiant incrémenté dont la clé est `user:next_id`.

La clé de stockage est la combinaison du préfixe et de l’identifiant de l'utilisateur. _Example: user:123455_.

Le mot de passe est stocké sous la forme d'un hash généré à partir de `bcrypt`.

Le lien entre l'identifiant utilisateur et son email est stocké à l'aide de la clé de stockage `user:email`.

Les identifiants seront stockés dans un _sorted set_ dont la clé de stockage sera `users` et le score sera le _timestamp_.

Un identifiant de session est créé, `session_id`, et la relation entre le `session_id` et l'identifiant utilisateur est stocké dans redis. La session expire si aucune requête n'a faite durant un temps _T_, _T_ étant définit dans la configuration.

## 4.7 Recherche

Si l'application autorise la recherche, elle est réalisée à l'aide de Redis Search. Un simple champs texte est disponible et recherche dans dans les champs `title` et `description` des produits.

# 5 Administration

Les paramètres seront gérés à l'aide de `flags`. L'identifiant marchant est optionnel. S'il n'est pas renseigné, la valeur par défaut dans la configuration sera utilisée.

## 5.1 Importation CSV de produits

Les paramètres sont:

- --file: Chemin vers le fichier à importer

Le séparateur est celui par défaut, la virgule `,`. L’ordre des colonnes du fichier CSV n’a pas d’importance, contrairement au nom des colonnes. Voici les champs disponibles:

- **sku**: Référence unique du produit
- **title**: Le titre du produit
- **price**: Le prix du produit
- **currency**: La devise du prix du produit
- **quantity**: La quantité du produit
- **online**: `true` si le produit est en ligne, `false` sinon
- **status**: Le statut du produit: `offline` ou `online`.
- **description**: La description du produit
- **images**: Les images produit séparées par un espace.
- **option names**: Les noms de des options séparées par un espace. Si le nom possède un espace, alors il doit être entouré de guillemets (optionnel)
- **option values**: Les valeurs des options correspondantes respectivement aux noms des options, séparées par un espace. Si le nom possède un espace, alors il doit être entouré de guillemets (optionnel)
- **links**: Les identifiants des produits séparés par un espace (optionnel)
- **tags**: Les tags (ou catégories) des produits (optionnel)
- **date**: La date au format JSON incluant le fuseau horaire (optionnel)

Le modèle présenté ci-dessus essaie d’être le plus minimaliste possible. Les options sont un bon moyen d’afficher des informations spécifiques selon les différents projets. Ils seront affichés dynamiquement dans la description du produit.

Si les valeurs contiennent des virgules, alors elles doivent être entourées de guillemets.

Les tags sont des moyens plus flexibles pour grouper des produits. Il sera possible de proposer dans la recherche, des tags prédéfinis que l’utilisateur pourra sélectionner. L’application peut aussi restreindre les tags possibles.

Les images devront être téléchargées et stockées dans le dossier servi par Imgproxy.

Si un marchand souhaite ajouter ses produits dans plusieurs langues, il doit créer une ligne pour chaque langue.

Si un marchand souhaite ajouter un produit avec différentes déclinaisons, il doit créer une ligne pour chaque produit. Il peut lier les produits entre eux, à l’aide des champs `links`. Sur la fiche produit, les produits liés affichent leur photo et il est possible de les consulter en cliquant dessus. La clé de stockage est `links`, le hash est le `PID` du produit et la valeur le `PID` du lien.

Si un produit existe, les données sont écrasées par la nouvelle importation. Si des options et des liens existaient, ils sont supprimés au profit des nouveaux liens et options.

La clé de stockage est la combinaison du préfixe `product` et d'un identifiant géré aléatoirement et unique, appelé `PID`. _Example: product:X6785FD49DN_.

Pour faciliter la récupération des produits lors de l'importation, le lien entre le `PID` et le `sku` est stocké dans redis, dont la clé est la combinaison de l'identifiant du vendeur et du `sku`. _Example: 1234:1233_.

**Remarque**: L’utilisation d’un identifiant vendeur permet de généraliser le projet à un marketplace. Cet identifiant est renseigné à travers la configuration pour des sites e-commerce, tandis que pour les marketplace, il est renseigné soit manuellement lors de l'importation des produits, soit grâce à l'identifiant du marchand qui s'est connecté à son interface.

## 5.2 Liste des produits

Les paramètres sont:

- --page: Pagination
- --merchant: Identifiant merchant (optionnel)

La pagination est un nombre qui représente un coefficient multiplicateur par le nombre d'éléments à afficher par page, disponible dans la configuration.

## 5.3 Détail d'un produit

Les paramètres sont:

- --pid: Le `PID` du produit

Le détail récupère tous les éléments stocké dans Redis.

## 5.4 Liste des utilisateurs

Les paramètres sont:

- --page: Pagination

Renvoie la liste des utilisateurs donc les premiers sont les plus récents.

## 5.5 Liste des commandes

Les paramètres sont:

- --page: Pagination

Renvoie la liste des commandes donc les premières sont les plus récents.

# 5.6 Modifier le statut d'une commande

Les paramètres sont:

- --id: L'identifiant de la commande
- --statut: Le nouveau statut de la commande

Les statuts disponibles sont:

- `payment_validated`
- `payment_progress`
- `payment_refused`

# 5.6 Ajouter une note à la commande

Les paramètres sont:

- --id: L'identifiant de la commande
- --note: La note à ajouter

# 6 Performances

Les performances sont d’une importance capitale. Les requêtes serveurs doivent répondre le plus rapidement possible. Le client doit contenir le minimum de javascript et le style CSS doit être optimisé, sans sélecteur complexe.

# 7 Sécurité

Les recommandations d'[OWASP](https://cheatsheetseries.owasp.org/index.html) sont respectées au maximum.

La protection CSRF est assurée par la vérification du header `HX-Request` et les méthode POST sont en AJAX exclusivement.

Les cookies ont le niveau de sécurité maximum.

# 8 Configuration

Les éléments de configuration de la plateforme sont disponible au format json avec les éléments suivants:

- **search**: `true` si la recherche est activée. La valeur par défaut est `false`.
- **items_per_page**: Le nombre d'éléments par page. La valeur par défaut est `12`.
- **tags**: `true` si la recherche par tags est activée. La valeur par défaut est `false`.
- **tags_list**: Limite la liste des tags utilisés dans la recherche.
- **withdraw**: `true` si la livraison par retrait est activée. La valeur par défaut est `false`.
- **payments**: La liste d'object contenant la configuration pour les moyens de paiement.
- **push_notification**: `true` si les push notifications sont activées. La valeur par défaut est `false`.
- **merchant**: L'identifiant du marchant par défaut. La valeur par défaut est `me`.
- **langs**: Liste de langues supportées par l'application. La valeur par défaut est `["fr"]`.
- **session_expiration**: La durée de la session utilisateur. La valeur par défaut est 7 \*24 \* 3600.

# 9 Internationalisation

Le socle serveur gère les traductions de chaque texte dans des fichiers dédiés à la traduction dans différentes langues. Les URLs doivent être traduites.

# 10 Style du code

GoLang impose un format unique.  
Pour le css, deux modes peuvent être utilisés:

- **classless**: Cela consiste à cibler sans utiliser de classe, ou alors une seule classe parente. Cette méthode à l'avantage de laisser le code HTML très propre, mais ne doit être utilisée que si les sélecteurs sont simples.
- **BEM**: [Référence](https://en.bem.info/methodology)

Pour le reste des fichiers (HTML, JS), prettier sera utilisé pour le formatage.

# 11 Tests

Les tests les plus importants sont les tests fonctionnels. [HURL](https://hurl.dev) est utilisé pour cela.

Cependant, il est vivement recommandé d'écrire des tests unitaires en utilisant l'approche de GoLang, au fur et à mesure, car cela permet de s'assurer de la qualité du projet.

# 12 Livrables

Un exécutable sera généré en fonction de la distribution du serveur, et des fichiers statiques (HTML, JS, CSS, JPG...) seront disponibles. Idéalement, ces fichiers pourront varier selon les implémentations des sites e-commerce, sans avoir des développements spécifiques du socle serveur.

# 13 Points d’entrée

L'application intercepte les erreurs et traite le retour selon le type de requête:

- **HTTP**: Une page d'erreur est affiché avec le message spécifique de l'erreur ou un message générique. Le code de l'erreur est renvoyé par la requête.
- **HTMX**: Une popup est affichée avec le message spécifique de l'erreur ou un message générique. Le code `200` est toujours renvoyé pour être traité par `HTMX`.

Dans les deux cas, le numéro de la requête doit être affiché.

| URL                              | Méthode | Type | Paramètres                                                        | Erreur                                     |
| -------------------------------- | ------: | ---- | ----------------------------------------------------------------- | ------------------------------------------ |
| /                                |     GET | HTTP | -                                                                 | -                                          |
| /contact.html                    |     GET | HTTP | -                                                                 | -                                          |
| /contact.html                    |    POST | HTMX | email, message                                                    | -                                          |
| /cgv.html                        |     GET | HTTP | -                                                                 | -                                          |
| /se-connecter.html               |     GET | HTTP | -                                                                 | -                                          |
| /se-connecter.html               |    POST | HTMX | email, password                                                   | 200 bad_credentials                        |
| /retrouver-mon-mot-de-passe.html |     GET | HTTP | -                                                                 | -                                          |
| /retrouver-mon-mot-de-passe.html |    POST | HTMX | email                                                             | 200 bad_email                              |
| /creer-mon-compte.html           |     GET | HTTP | -                                                                 | -                                          |
| /creer-mon-compte.html           |    POST | HTMX | email,password,confirm                                            | 200 bad_confirm, 200 bad_email             |
| /mon-compte.html                 |     GET | HTTP | -                                                                 | -                                          |
| /mon-compte.html                 |    POST | HTMX | previous_password,password,confirm                                | 200 bad_confirm, 200 bad_previous_password |
| /mon-adresse.html                |     GET | HTTP | -                                                                 | -                                          |
| /mon-adresse.html                |    POST | HTMX | firstname, lastname, address, complementary, zipcode, city, phone | 200 bad_parameters                         |
| /ma-facturation.html             |    POST | HTMX | firstname, lastname, address, complementary, zipcode, city, phone | 200 bad_parameters                         |
| /mon-historique.html             |     GET | HTTP | -                                                                 |
| /mon-historique/${id}.html       |     GET | HTTP | -                                                                 | -                                          |
| /products.html                   |     GET | HTTP | page                                                              | -                                          |
| /products/${id}-${slug}.html     |     GET | HTTP | -                                                                 | 404 product_not_found                      |
| /recherche.html                  |     GET | HTTP | page, titre, description, tags                                    | -                                          |
| /panier.html                     |    POST | HTMX | id, quantity                                                      | 404 product_not_found                      |
| /panier-connexion.html           |     GET | HTTP | -                                                                 | 200 bad_credentials                        |
| /panier-connexion.html           |    POST | HTMX | email, password, guest                                            | 200 bad_credentials                        |
| /panier-livraison.html           |     GET | HTTP | -                                                                 | 200 bad_credentials                        |
| /panier-livraison.html           |    POST | HTMX | type                                                              |
| /panier-adresse.html             |     GET | HTTP | -                                                                 | 200 bad_credentials                        |
| /panier-adresse.html             |    POST | HTMX | firstname, lastname, address, complementary, zipcode, city, phone | 200 bad_parameters                         |
| /panier-facturation.html         |     GET | HTTP | -                                                                 | 200 bad_credentials                        |
| /panier-facturation.html         |    POST | HTMX | firstname, lastname, address, complementary, zipcode, city, phone | 200 bad_parameters                         |
| /panier-paiement.html            |     GET | HTTP | -                                                                 | -                                          |
| /panier-paiement.html            |    POST | HTMX | type, ...                                                         | 200 bad_payment                            |
