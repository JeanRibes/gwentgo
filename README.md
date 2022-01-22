# plan

- [x] logique de jeu (API gwent)
- [x] matchmaking
- [ ] jeu à deux
- [x] deckbuilding
    - [ ] carte(s) offertes par victoires
    - [ ] starter deck (possiblement autre que northern realms)
    - [x] interface pour choisir les cartes
- [ ] stockage des données (DB, json...)
- [ ] sécurité signature HMAC des decks des gens ?
- [ ] IA
    - [ ] IA conventionnelle, programmée
    - [ ] deep learning ? entraintée à jouer contre l'ordi conventionnel, puis contre elle-même

## Stratégie de l'IA

* jouer tous les espions au début puis passer si l'enemi est trop fort
* utiliser les decoys pour rejouer des espions
* utiliser les Medics uniquement quand elle peut récup des cartes
* utiliser Scorch quand ça fait plus de dégâts à l'ennemi, et vers la fin
* reset la météo à la fin du tour
* utiliser les cornes de guerre à la fin du tour
* abattre toutes ses cartes si gagner ce tour lui permet de remporter la partie

# Data model

All the possible cards are stored in a CSV file. Some rows are duplicating, indicating the number of times a card can
appear. Cards have an ID which is unique (position in CSV file), in a deck.

Upon creating a game, the cards inside a reindexed, because there can be the same cards on both sides, so we need to
differenciate them.

# Issues

* fix Villentremerth: it can be scorched by itself
* tests for scorch