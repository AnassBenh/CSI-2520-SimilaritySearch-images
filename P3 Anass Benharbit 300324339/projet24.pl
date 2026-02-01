% dataset(NomDuRépertoire)
% c'est l'endroit où se trouve l'ensemble de données d'images
dataset('C:\\Users\\benha\\OneDrive\\Bureau\\Uni\\Hiver2024\\CSI2520\\P3\\imageDataset2_15_20\\').

% directory_textfiles(NomDuRépertoire, ListeDeFichiersTexte)
% produit la liste des fichiers texte dans un répertoire
directory_textfiles(D,Textfiles):- directory_files(D,Files), include(isTextFile, Files, Textfiles).
isTextFile(Filename):-string_concat(_,'.txt',Filename).

% read_hist_file(NomDuFichier,ListOfNumbers)
% lit un fichier d'histogramme et produit une liste de nombres (valeurs de bin)
read_hist_file(Filename,Numbers):- open(Filename,read,Stream),read_line_to_string(Stream,_),
                                   read_line_to_string(Stream,String), close(Stream),
								   atomic_list_concat(List, ' ', String),atoms_numbers(List,Numbers).
								   
% similarity_search(FichierRequête,ListeImagesSimilaires)
% renvoie la liste des images similaires à l'image de requête
% les images similaires sont spécifiées comme (NomImage, ScoreSimilarité)
% le prédicat dataset/1 fournit l'emplacement de l'ensemble d'images
similarity_search(QueryFile,SimilarList) :- dataset(D), directory_textfiles(D,TxtFiles),
                                            similarity_search(QueryFile,D,TxtFiles,SimilarList).
											
% similarity_search(FichierRequête, RépertoireJeuDeDonnées, ListeFichiersHisto, ListeImagesSimilaires)
similarity_search(QueryFile,DatasetDirectory, DatasetFiles,Best):- read_hist_file(QueryFile,QueryHisto), 
                                            compare_histograms(QueryHisto, DatasetDirectory, DatasetFiles, Scores), 
                                            sort(2,@>,Scores,Sorted),take(Sorted,5,Best).

% compare_histograms(HistoRequête,RépertoireJeuDeDonnées,FichiersJeuDeDonnées,Scores)
% compare un histogramme de requête avec une liste de fichiers d'histogramme
compare_histograms(_, _, [], []).
compare_histograms(QueryHisto, DatasetDirectory, [H|T], [(H, Score)|Scores]) :-
    string_concat(DatasetDirectory, H, Path),
    read_hist_file(Path, Histo),
    histogram_intersection(QueryHisto, Histo, Score),
    compare_histograms(QueryHisto, DatasetDirectory, T, Scores).

% histogram_intersection(Histogramme1, Histogramme2, Score)
% calcule le score de similarité d'intersection entre deux histogrammes
% Le score est compris entre 0.0 et 1.0 (1.0 pour des histogrammes identiques)
histogram_intersection(H1, H2, Score) :-
    sum_list(H1, Total),
    intersection_score(H1, H2, Intersection),
    Score is Intersection / Total.

% intersection_score(Histogramme1, Histogramme2, Score)
% prédicat auxiliaire pour calculer le score d'intersection
intersection_score([], [], 0).
intersection_score([H1|T1], [H2|T2], Score) :-
    intersection_score(T1, T2, RestScore),
    Min is min(H1, H2),
    Score is Min + RestScore.


% take(Liste,K,KListe)
% extrait les K premiers éléments d'une liste
take(Src,N,L) :- findall(E, (nth1(I,Src,E), I =< N), L).

% atoms_numbers(ListeD'Atomes,ListeDeNombres)
% convertit une liste d'atomes en une liste de nombres
atoms_numbers([],[]).
atoms_numbers([X|L],[Y|T]):- atom_number(X,Y), atoms_numbers(L,T).
atoms_numbers([X|L],T):- \+atom_number(X,_), atoms_numbers(L,T).
