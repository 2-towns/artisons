# Cadeau

# 1 Présentation

Le​ ​but​ ​de​ ​ce​ ​document​ ​est​ ​de proposer des spécifications fonctionnelles pour le projet “cadeau”.

# 2 Représentants

​Arnaud Deville, Salim Tison et Reda Madjoub sont les développeurs de cette application.

# 3 Versions

| Version |    Date    |         Auteur |
| ------- | :--------: | -------------: |
| 1       | 10/06/2023 | Arnaud Deville |

# 4 Description fonctionnelle

Remarque: Ce document présente deux solutions:

- La solution 1 avec une création de compte classique email et mot de passe
- La solution 2 sans compte. La récupération des informations d’une commande se fait alors à l’aide du couple, email et numéro de commande.

## 4.1 Page d’accueil

La page d’accueil présente l'activité principale. On peut voir une partie des produits qui sont disponibles à la vente. Depuis cette page, on peut:

- Se connecter
- Voir la liste des produits
- Accès aux pages statiques
- Accéder au panier
- Accéder à l’historique
- Un champ pour faire la recherche si les produits sont relativement nombreux (plus de 30)

## 4.2 Liste des produits

La liste des produits affiche les produits sous forme de grille à l’aide d’une pagination. Il n’est pas possible de filtrer ces produits. Le clic sur un produit renvoie sur la page de détail de cet article.

## 4.3 Détail d’un produit

Le détail d’un produit affiche toutes les caractéristiques de ce produit. Si plusieurs photos sont disponibles, elles seront affichées en miniature en dessous de la photo principale. Le clic sur l’une des miniatures, remplacera la photo principale par celle de la miniature.

Un ensemble de critères sont disponibles afin de personnaliser son cadeau:

- Prix: Estimation prix du cadeau, requis
- Type d'événement: Liste libre, requis
- Catégories: Homme/Femme/Enfant, requis
- Centre d'intérêt: Trois champs libres, requis

Un champ libre est un champ proposant une liste de suggestions, mais laisse l’utilisateur entrer ce qu’il souhaite.

Un champ quantité est également présent. Un bouton est présent pour ajouter les produits au panier, avec les options sélectionnées. Une fois le produit ajouté, le bouton est désactivé. \

## 4.4 Panier

La panier panier affiche la liste des produits sélectionnés par l’utilisateur avec ces options. Il peut modifier les quantités. Le total de la commande est affiché sur la page. Lorsque l’utilisateur clique sur le bouton de validation de la commande, deux actions sont possibles:

- L’utilisateur n’est pas connecté:

  La page affiche un formation de création de compte avec les champs suivants:

  - Nom
  - Email
  - Confirmation d’email
  - Téléphone

  Après avoir rempli le formulaire, il passera à l’action suivante décrite ci-dessous.

- L’utilisateur est connecté:

  Un formulaire de paiement s’affiche. Après validation du paiement, un écran de confirmation de paiement s’affiche, mentionnant le numéro de commande ;

## 4.5 Historique de commande

### 4.5.1 Solution 1

La liste des commandes affiche l’ensemble des commandes réalisées par l’utilisateur avec les éléments suivants:

- Numéro de commande
- Date
- Statut de la commande

Lors du clique sur un élément, l’utilisateur peut voir le détail de la commande.

### 4.5.2 Solution 2

L’écran d’historique des commandes affiche un formulaire avec les champs suivants:

- Email
- Numéro de commande

Lorsqu’il valide sa commande, le détail de sa commande est affiché.

## 4.6 Pages statiques

Les pages statiques affichent du contenu classique:

- Contact: Contient un email de contact
- Conditions d’utilisation
- Conditions générales de vente

## 4.7 Connexion

### 4.7.1 Solution 1

Le formulaire de connexion affiche les champs suivants:

- Email
- Mot de passe

Le mot de passe correspond par défaut au numéro de la dernière commande qu’il a effectuée. Lorsqu’il clique sur le lien d’oubli de son numéro de commande, un formulaire remplace le formulaire précédent, lui demandant son email. Lorsqu’il valide, un email lui est envoyé avec son dernier numéro de commande.

S’il ne valide pas le formulaire mais clique sur le lien de connexion, le premier formulaire apparaît à nouveau. .

### 4.7.2 Solution 2

Il n’y a pas de connexion.

## 4.8 Administration

L’administration de l’application se sera, dans un premier temps, à l’aide du terminal. Une interface pourra être développée plus tard.

Les actions suivantes sont possibles:

- Importer les produits à l’aide d’un CSV
- Lister les produits à l’aide d’une pagination
- Afficher les détails d’un produit
- Afficher la liste des utilisateurs à l’aide d’une pagination
- Afficher la liste des commandes à l’aide d’une pagination
- Modifier le statut d’une commande

# 7 Couverture géographique

# 8 Moyens de paiements

# 9 Relations clients

# 10 Gestion des insatisfactions clientes

# 11 Références

https://www.mindmeister.com/map/2777728011?t=PPVYtObQ2k
