#lang racket

; Obtenir la liste de tous les fichiers texte dans un répertoire
(define (list-text-files-in-directory directory-path)
  (filter (lambda (file)
            (string-suffix? file ".txt"))
          (map (lambda (file) (path->string (build-path directory-path file)))
               (directory-list directory-path))))

; Lire un fichier texte contenant un histogramme et renvoyer les valeurs dans une liste
(define (read-hist-file filename) 
  (call-with-input-file filename
    (lambda (p)
      (letrec ((f (lambda (x)
                    (if (eof-object? x) '() (cons x (f (read p)))))))
        (cdr (f (read p)))))))

; Calculer la similarité entre deux histogrammes
(define (calculate-similarity hist1 hist2)
  (if (or (null? hist1) (null? hist2))
      0
      (+ (abs (- (car hist1) (car hist2))) (calculate-similarity (cdr hist1) (cdr hist2)))))

; Trouver les cinq images les plus similaires ainsi que leurs noms de fichiers et scores de similarité
(define (find-most-similar-images query-histogram all-histograms)
  (define sorted-similarities
    (sort all-histograms
          (lambda (a b)
            (< (calculate-similarity (cdr query-histogram) (cdr a))
               (calculate-similarity (cdr query-histogram) (cdr b))))))

  (take sorted-similarities 5))

; Lire tous les histogrammes d'un répertoire
(define (read-all-histograms directory)
  (map (lambda (filename)
         (cons filename (read-hist-file filename)))
       (list-text-files-in-directory directory)))

; La fonction principale
(define (similaritySearch queryHistogramFilename imageDatasetDirectory)
  (let* ((query-histogram (cons queryHistogramFilename (read-hist-file queryHistogramFilename)))
         (all-histograms (read-all-histograms imageDatasetDirectory))
         (similar-images (find-most-similar-images query-histogram all-histograms)))
    (map car similar-images)))
