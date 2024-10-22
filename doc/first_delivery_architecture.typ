#import "@preview/touying:0.5.3": *
#import "@preview/codly:1.0.0": *
#import "@preview/cetz:0.3.0"

#import themes.university: *

#show: university-theme.with(
    aspect-ratio: "16-9",
    config-common(preamble: {
        codly(languages: (
            go: (
            name: "go",
            color: rgb("#CE412B"),
            ),
        ))
    })
)


// #show: codly-init.with()


= V-System

== //Group Information

#align(center+top)[
#linebreak()
#linebreak()
*Cybersecurity - 2024/2025*\

Cosimo Giraldi \
Giacomo Grassi\
Michele Ivan Bruna
]


== Player Information
#align(center+horizon)[
#raw(lang: "go", block: true,
`type User struct {
    UserID          []byte
    ApiUUID         [16]byte
    Salt            []byte
    Username        string
    Email           string
    PasswordHash    []byte
    Gachas          []Gacha
    Transactions    []Transaction
}`.text)]


== Gacha Collection

#align(center+horizon)[
#set text(size: 14.2pt)
#grid(
  columns: (3fr, .3fr),
  [
  #grid(
    columns: (1fr, 1fr, 1fr, 1fr, 1fr),
    column-gutter: 0pt,
    row-gutter: 5pt,
    [*Common (C)*], [*Uncommon (U)*], [*Rare (R)*], [*Epic (E)*], [*Legendary (L)*],
    image("../assets/images/png/warrior_cricket_common.png",      width: 60%),
    image("../assets/images/png/warrior_centipede_uncommon.png",      width: 60%),
    image("../assets/images/png/warrior_beetle_rare.png",      width: 60%),
    image("../assets/images/png/mage_moth_epic.png",      width: 60%),
    image("../assets/images/png/warrior_hercule_beetle_legendary.png",      width: 60%),
    //
    image("../assets/images/png/warrior_locust_common.png", width: 60%),
    image("../assets/images/png/priest_cicada_uncommon.png",      width: 60%),
    image("../assets/images/png/priest_moth_rare.png",      width: 60%),
    image("../assets/images/png/mage_butterfly_epic.png",      width: 60%),
    image("../assets/images/png/mage_butterfly_legendary.png",      width: 60%),
    //
    image("../assets/images/png/tank_mole-cricket_common.png", width: 60%),
    image("../assets/images/png/mage_mosquito_uncommon.png",      width: 60%),
    image("../assets/images/png/druid_butterfly_rare.png",      width: 60%),
    image("../assets/images/png/warrior_dragonfly_epic.png",      width: 60%),
    image("../assets/images/png/druid_butterfly_legendary.png",      width: 60%),
    //
    image("../assets/images/png/munich_grasshopper_common.png",      width: 60%),
    image("../assets/images/png/druid_bee_uncommon.png",      width: 60%),
    image("../assets/images/png/assassin_mosquito_rare.png",      width: 60%), // TODO: change
    image("../assets/images/png/assassin_peacock_epic.png",      width: 60%),
    image("../assets/images/png/demoniac_mosquito_legendary.png",      width: 60%),
  )],
  [#set align(horizon)
  #figure(
    image("../assets/images/png/coin.png", width: 80%),
    caption: [Currency],
    //caption-pos: top, // V 0.12
    supplement: [],
    numbering: none,
  )
  ]
)]



= Monolithic architecture

==
#align(center+horizon)[
    #figure(
    image("../assets/drawio/monolithic_architecture.svg"),
    )
]


== Database contents

#align(center + horizon)[
#set text(size: 14.2pt)
#set align(left)

#table(
    columns: (3fr, 3fr),
    gutter: 10pt,
    [```sql
    CREATE TABLE users (
        user_id BYTEA PRIMARY KEY UNIQUE,
        api_uuid BYTEA NOT NULL UNIQUE,
        salt BYTEA NOT NULL,
        username VARCHAR(255) NOT NULL UNIQUE,
        email VARCHAR(255) NOT NULL UNIQUE,
        password_hash BYTEA NOT NULL
    );```],
    [```sql
    CREATE TABLE gachas (
        gacha_id BYTEA PRIMARY KEY UNIQUE,
        api_uuid BYTEA NOT NULL UNIQUE
        name VARCHAR(255) NOT NULL,
        rarity rarity NOT NULL,
        price BIGINT
    );```],
    [```sql
    CREATE TABLE transactions (
        transaction_id BYTEA PRIMARY KEY UNIQUE,
        api_uuid BYTEA NOT NULL UNIQUE,
        transaction_type transaction_type NOT NULL,
        user_id BYTEA REFERENCES users(user_id),
        amount BIGINT NOT NULL,
        date_time TIMESTAMP NOT NULL,
        event_type VARCHAR(255) NOT NULL,
        event_id BYTEA NOT NULL
    );```],
    [```sql
    CREATE TABLE user_gachas (
        user_id BYTEA REFERENCES users(user_id),
        gacha_id BYTEA REFERENCES gachas(gacha_id),
        PRIMARY KEY (user_id, gacha_id)
    );```]
)


#text(size: 7pt)[
```sql
CREATE TYPE transaction_type AS ENUM ('Deposit', 'Withdraw');
```
```sql
CREATE TYPE rarity AS ENUM ('Common', 'Uncommon', 'Rare', 'Epic', 'Legendary');
```
]

]
