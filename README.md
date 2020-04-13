# Modélisation d'une épidémie

---
## Contexte 

Adapté d'un article du [Washington Post](https://www.washingtonpost.com/graphics/2020/world/corona-simulator/?fbclid=IwAR1QLrFRcQZ7BNc87RjHbX3V6e9J1dAcKDvPQGA8LfQpfcYMXpGLAgWAa08)

---
## Méthode

### Paramètres

TODO

### Modélisation des collisions avec les murs

TODO

### Modélisation des collisions entre les agents

TODO

---
## Installation



### Pré-requis

* [Python](https://www.anaconda.com/distribution/) >= 3.7
* [Go](https://golang.org/doc/install) >=1.13

Packages Python :
* numpy
* pandas
* matplotlib
* json

### Installation

```
git clone https://github.com/AureClai/modelisation-epidemie
```

---
## Utilisation

### Paramétrer une simulation

Modifier le fichier `settings.json` à la racine du dossier

### Lancer le coeur de calcul

```
go run ./Go-Core/.
```

### Observer un graphique seul (étape de validation)

```
python -m instant_grah.py /chemin_vers_le_dossier_de_résultats
```

### Produire une vidéo

```
python -m make_video.py /chemin_vers_le_dossier_de_résultats
```

--- 
## Licence

[GNU General Public License v3.0](https://github.com/AureClai/modelisation-epidemie/blob/master/LICENSE)

---
## Pistes d'améliorations

* Ajouter d'autres paramètres
* Modifier les paramètres vidéos sur demande
* etc...
