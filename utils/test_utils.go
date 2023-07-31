package testutils

import (
	"gifthub/products"

	"github.com/go-faker/faker/v4"
)

// FakeProduct génère un produit factice.
func FakeProduct() products.Product {
	return products.Product{
		ID:          faker.RandomInt(1, 1000), // Génère un ID aléatoire entre 1 et 1000.
		Title:       faker.Sentence(), // Génère une phrase aléatoire pour le titre.
		Description: faker.Paragraph(), // Génère un paragraphe aléatoire pour la description.
		Price:       faker.Price(), // Génère un prix aléatoire.
		Slug:        faker.Slug(), // Génère un slug aléatoire.
		Image:       faker.ImageURL(), // Génère une URL d'image aléatoire.
		MerchantID:  faker.RandomString(10), // Génère une chaîne aléatoire de 10 caractères pour le MerchantID.
		Links:       []string{faker.URL(), faker.URL()}, // Génère deux URLs aléatoires pour les liens.
		Meta:        map[string]string{"key": faker.Word(), "value": faker.Sentence()}, // Génère un dictionnaire avec une clé et une valeur aléatoires.
	}
}
