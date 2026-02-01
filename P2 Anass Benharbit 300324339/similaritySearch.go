// Anass Benharbit 300324339

package main

import (
	"fmt"
	"image"
	_ "image/jpeg" // Importation du package image/jpeg pour la prise en charge des images JPEG
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// Structure Histo représentant un histogramme
type Histo struct {
	Nom string    // Nom de l'image
	H   []float64 // Valeurs de l'histogramme
}

// Fonction pour calculer l'histogramme d'une image
func computeHistogramme(cheminImage string, profondeur int) (Histo, error) {
	// Ouverture du fichier image
	fichier, err := os.Open(cheminImage)
	if err != nil {
		return Histo{}, err
	}
	defer fichier.Close()

	// Décodage de l'image
	img, _, err := image.Decode(fichier)
	if err != nil {
		return Histo{}, err
	}

	// Récupération des limites de l'image
	limites := img.Bounds()
	histogramme := make([]float64, 512) // Création d'un histogramme de taille 512

	// Parcours de chaque pixel de l'image pour calculer l'histogramme
	for y := limites.Min.Y; y < limites.Max.Y; y++ {
		for x := limites.Min.X; x < limites.Max.X; x++ {
			// Récupération des composantes RGB du pixel
			r, g, b, _ := img.At(x, y).RGBA()
			// Décalage des composantes pour la réduction de la profondeur
			rDécalé := uint32(int(r) >> uint(8-profondeur))
			gDécalé := uint32(int(g) >> uint(8-profondeur))
			bDécalé := uint32(int(b) >> uint(8-profondeur))
			idx := (rDécalé << uint(2*profondeur)) + (gDécalé << uint(profondeur)) + bDécalé

			if int(idx) >= len(histogramme) {
				continue
			}

			// Incrémentation du compteur d'occurrence de la valeur de l'histogramme
			histogramme[idx]++
		}
	}

	// Normalisation de l'histogramme en divisant chaque valeur par le nombre total de pixels
	totalPixels := float64((limites.Max.X - limites.Min.X) * (limites.Max.Y - limites.Min.Y))
	for i := range histogramme {
		histogramme[i] = float64(histogramme[i]) / totalPixels
	}

	// Retourne l'histogramme calculé et le nom de l'image
	return Histo{Nom: cheminImage, H: histogramme}, nil
}

// Fonction pour calculer les histogrammes de plusieurs images en parallèle
func computeHistogrammes(cheminsImages []string, profondeur int, canalHistogrammes chan<- Histo, wg *sync.WaitGroup) {
	defer wg.Done()
	// Parcours de chaque chemin d'image dans la tranche
	for _, chemin := range cheminsImages {
		// Calcul de l'histogramme de l'image
		histogramme, err := computeHistogramme(chemin, profondeur)
		if err != nil {
			// En cas d'erreur, affichage du message d'erreur et passage à l'image suivante
			log.Printf("Erreur lors du calcul de l'histogramme pour %s: %v", chemin, err)
			continue
		}
		// Envoi de l'histogramme calculé sur le canal des histogrammes
		canalHistogrammes <- histogramme
	}
}

func main() {
	// Vérification des arguments en ligne de commande
	if len(os.Args) < 3 {
		fmt.Println("Utilisation: go run similaritySearch.go nomFichierImageRequête répertoireDatasetImages")
		return
	}

	// Récupération du nom de l'image de requête et du répertoire contenant les images du dataset
	nomFichierImageRequête := os.Args[1]
	répertoireDatasetImages := os.Args[2]

	// Démarrage du chronomètre
	début := time.Now()

	// Initialisation du WaitGroup pour la synchronisation
	var wg sync.WaitGroup
	// Création du canal pour stocker les histogrammes calculés
	histogrammes := make(chan Histo)

	// Lecture des fichiers dans le répertoire du dataset
	fichiersImages, err := ioutil.ReadDir(répertoireDatasetImages)
	if err != nil {
		log.Fatal(err)
	}

	var nomsFichiers []string
	// Parcours de chaque fichier dans le répertoire
	for _, fichier := range fichiersImages {
		// Vérification si le fichier est une image JPEG
		if strings.HasSuffix(fichier.Name(), ".jpg") {
			// Ajout du chemin complet du fichier à la liste des noms de fichiers
			nomsFichiers = append(nomsFichiers, filepath.Join(répertoireDatasetImages, fichier.Name()))
		}
	}

	// Définition du nombre de tranches pour la distribution du calcul des histogrammes
	K := 1048
	// Création des tranches pour distribuer les images sur plusieurs goroutines
	slicesImages := make([][]string, K)
	for i, fichier := range nomsFichiers {
		idx := i % K
		slicesImages[idx] = append(slicesImages[idx], fichier)
	}

	// Lancement des goroutines pour calculer les histogrammes en parallèle
	for i := 0; i < K; i++ {
		wg.Add(1)
		go computeHistogrammes(slicesImages[i], 3, histogrammes, &wg)
	}

	// Fonction anonyme pour attendre la fin du calcul des histogrammes et fermer le canal
	go func() {
		wg.Wait()
		close(histogrammes)
	}()

	var résultats []Histo
	// Collecte des histogrammes calculés depuis le canal
	for histogramme := range histogrammes {
		résultats = append(résultats, histogramme)
	}

	fmt.Println("Top 5 images similaires:")

	// Calcul de l'histogramme de l'image de requête
	histogrammeRequête, err := computeHistogramme(nomFichierImageRequête, 3)
	if err != nil {
		log.Fatal(err)
	}

	// Structure pour stocker les chemins d'images similaires avec leur distance
	type ImageSimilaire struct {
		CheminImage string
		Distance    float64
	}

	var imagesSimilaires []ImageSimilaire
	// Calcul de la distance entre l'histogramme de l'image de requête et ceux du dataset
	for _, résultat := range résultats {
		distance := distanceA(histogrammeRequête.H, résultat.H)
		imagesSimilaires = append(imagesSimilaires, ImageSimilaire{CheminImage: résultat.Nom, Distance: distance})
	}

	// Tri des images similaires par distance croissante
	sort.Slice(imagesSimilaires, func(i, j int) bool {
		return imagesSimilaires[i].Distance < imagesSimilaires[j].Distance
	})

	// Affichage des 5 images les plus similaires
	for i := 0; i < 5 && i < len(imagesSimilaires); i++ {
		fmt.Printf("Image similaire %d: %s, %f\n", i+1, imagesSimilaires[i].CheminImage, imagesSimilaires[i].Distance)
	}

	// Calcul du temps écoulé et affichage
	écoulé := time.Since(début)
	fmt.Printf("Temps d'exécution: %s\n", écoulé)
}

// Fonction pour calculer la distance entre deux histogrammes
func distanceA(hist1, hist2 []float64) float64 {
	// Vérification si les deux histogrammes ont la même longueur
	if len(hist1) != len(hist2) {
		log.Fatal("Les histogrammes doivent être de même longueur")
	}

	// Calcul de la distance selon la formule de la distance d'intersection
	somme := 0.0
	for i := range hist1 {
		somme += math.Sqrt(hist1[i] * hist2[i])
	}

	// Retourne la distance calculée
	return -math.Log(somme)
}
