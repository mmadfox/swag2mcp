# validate

## Objectif

Vérifier le fichier de configuration et tous les fichiers de spécification référencés pour détecter les erreurs. C'est une commande de diagnostic **en lecture seule** — elle ne modifie jamais rien.

## Quand l'utiliser

- Après avoir modifié `swag2mcp.yaml` manuellement
- Avant d'exécuter `mcp` ou `update` pour détecter les problèmes tôt
- Lors du dépannage pour savoir pourquoi une spec ne se charge pas
- Dans les pipelines CI/CD pour valider les modifications de configuration

## Syntaxe

```bash
swag2mcp validate [chemin] [drapeaux]
```

## Arguments

| Argument | Position | Requis | Description |
|----------|----------|--------|-------------|
| `chemin` | 1 | Non | Répertoire de l'espace de travail. S'il est omis, résolution via les règles de résolution de chemin. |

## Drapeaux

| Drapeau | Raccourci | Type | Défaut | Description |
|---------|-----------|------|--------|-------------|
| `--tags` | `-t` | `string` | `""` | Valider uniquement les specs avec des étiquettes correspondantes (séparées par des virgules) |

## Comment cela fonctionne

```bash
swag2mcp validate
swag2mcp validate ./mon-espace-travail
swag2mcp validate --tags=public
```

## Ce qui est vérifié

| Vérification | Description |
|--------------|-------------|
| Syntaxe YAML | Le fichier de configuration doit être un YAML valide |
| Structure de la configuration | Tous les champs requis présents, les types sont corrects |
| Unicité des domaines | Pas de domaines en double |
| Format du domaine | Uniquement minuscules, chiffres, tirets |
| Existence du fichier de spécification | Le fichier ou l'URL `location` doit être accessible |
| Format de la spécification | Le fichier doit être un OpenAPI 3.x, Swagger 2.0 ou Postman valide |
| Paramètres d'authentification | Le type et la configuration d'auth sont valides pour la méthode sélectionnée |
| Client HTTP | Les paramètres du client HTTP sont valides |

## Ce qui n'est PAS vérifié

| Non vérifié | Raison |
|-------------|--------|
| Points d'accès d'authentification | `validate` vérifie la syntaxe de la configuration d'auth mais ne teste pas la connexion/l'échange de jetons |
| Disponibilité des points d'accès API | Seule l'URL du fichier de spécification est vérifiée, pas la `base_url` |
| Exactitude de `base_url` | Le format est validé, mais aucune requête de test n'est effectuée |
| Configuration du serveur mock | `base_mock_url` n'est pas vérifié pour la connectivité |

## Exemple de sortie

```
✅ La configuration est valide.
✓ Spec animalerie : OK
✓ Spec meteo : OK
✗ Spec ancienne-api : fichier introuvable
```

## Vérification post-commande

Si la validation réussit, la configuration est prête pour `mcp`, `update` ou `run`.

## Nuances

- **Pas d'auto-initialisation :** Contrairement à `add`, `ls` ou `run`, `validate` ne s'auto-initialise **pas** si la configuration est manquante. Il retourne une erreur : « configuration introuvable à &lt;chemin&gt; ».
- **Accès réseau :** Les URL de spécification distantes sont récupérées pendant la validation. La commande peut prendre plus de temps si les specs sont hébergées sur des serveurs lents.
- **Filtrage par étiquettes :** Lorsque `--tags` est défini, seules les specs correspondant aux étiquettes spécifiées sont validées. Les autres specs sont ignorées.
