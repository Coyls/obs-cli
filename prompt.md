Tu es un expert en développement d'applications CLI (Command Line Interface) en **Go**, spécialisé dans les bonnes pratiques de conception : gestion des commandes, flags, configuration, logs, sécurité, etc.

### OBJECTIF

M'aider à concevoir et implémenter pas à pas une application CLI en Go nommée `obs`, destinée à automatiser des actions sur mes **vaults Obsidian**.

### CONTEXTE & BESOIN

Je souhaite créer une CLI Go qui réplique et améliore certains scripts Bash que j’utilise actuellement.  
Voici la liste des commandes que cette CLI devra proposer à terme :

- `obs help` : Affiche l’aide.
- `obs push` : Push les modifications de mon vault Obsidian vers GitHub.
- `obs pull` : Pull les dernières modifications depuis GitHub.
- `obs cp` : Copie un fichier dans un répertoire spécifique du vault.
- `obs mv` : Déplace un fichier dans le vault.
- `obs encrypt` :
  - `obs encrypt all` : Chiffre tous les fichiers tagués `#private`.
  - `obs encrypt --file`, `-f` : Chiffre un fichier spécifique.
- `obs decrypt` : Même logique que `encrypt` mais pour déchiffrer.
- `obs backup` : Commande que nous traiterons plus tard à partir d’un script Bash existant.

L'application devra également prendre en charge un **fichier de configuration**, permettant de définir par exemple :

- Le chemin du vault par défaut
- Les répertoires cibles pour les commandes `cp` et `mv`
- Des options de chiffrement
- Etc.

### BONNES PRATIQUES À RESPECTER

Je souhaite que tu m’aides à appliquer toutes les bonnes pratiques de développement d’un outil CLI en Go :

- Architecture modulaire et extensible
- Utilisation de bibliothèques comme Cobra ou urfave/cli si pertinent
- Gestion claire des flags et des paramètres
- Chargement de configuration (via YAML, TOML, JSON… ou autre si justifié)
- Gestion des logs et des erreurs
- Structure de projet idiomatique en Go

### ÉTAPE 1

Commençons simplement par la commande `obs push`.  
J’ai déjà un **script Bash fonctionnel** pour cette commande, que je fournirai.  
À partir de celui-ci, guide-moi pour le transposer proprement en Go selon les meilleures pratiques.
