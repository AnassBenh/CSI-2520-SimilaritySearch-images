//Anass Benharbit 300324339

import java.io.File;
import java.util.Map;
import java.util.HashMap;
import java.util.PriorityQueue;

public class SimilaritySearch {
    public static void main(String[] args) {
        // Vérifier que deux arguments sont fournis en ligne de commande
        if (args.length != 2) {
            System.out.println("Usage: java SimilaritySearch <image_requête> <répertoire_base_de_données>");
            return;
        }

        String imageRequeteFichier = args[0];
        String repertoireBaseDeDonnees = args[1];



        ColorImage imageRequete = new ColorImage(imageRequeteFichier);
        ColorHistogram histogrammeRequete = new ColorHistogram(3);
        histogrammeRequete.setImage(imageRequete);

        Map<String, Double> scoresSimilarite = new HashMap<>();

        // Parcourir les fichiers du répertoire de la base de données
        File repertoire = new File(repertoireBaseDeDonnees);
        for (File fichier : repertoire.listFiles()) {
            if (fichier.isFile() && (fichier.getName().endsWith(".jpg"))) {
                if (!fichier.getName().equals(imageRequeteFichier)) {
                    ColorImage imageBaseDonnees = new ColorImage(fichier.getPath());
                    ColorHistogram histogrammeBaseDonnees = new ColorHistogram(3);
                    histogrammeBaseDonnees.setImage(imageBaseDonnees);
                    double similarite = histogrammeRequete.compare(histogrammeBaseDonnees);
                    // Stocker le score de similarité avec le nom du fichier de la base de données
                    scoresSimilarite.put(fichier.getName(), similarite);
                }
            }
        }

        // Trier les scores de similarité par ordre décroissant
        PriorityQueue<Map.Entry<String, Double>> fileAttentePriorite = new PriorityQueue<>(
            (a, b) -> Double.compare(b.getValue(), a.getValue())
        );
        fileAttentePriorite.addAll(scoresSimilarite.entrySet());

        // Vérifier si aucune image similaire n'a été trouvée
        if (scoresSimilarite.isEmpty()) {
            System.out.println("Aucune image similaire trouvée dans le répertoire spécifié.");
            return;
        }

        // Afficher les noms des 5 images les plus similaires
        System.out.println("Voici les 5 images les plus similaires: ");
        int compteur = 0;
        while (!fileAttentePriorite.isEmpty() && compteur < 5) {
            Map.Entry<String, Double> entrée = fileAttentePriorite.poll();
            System.out.println(entrée.getKey() + ": " + entrée.getValue());
            compteur++;
        }
    }
}
