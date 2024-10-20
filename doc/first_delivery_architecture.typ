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


#show: codly-init.with()


= V-System

== //Group Information

#align(center)[
#linebreak()
#linebreak()
*Cybersecurity - 2024/2025*\

Cosimo Giraldi \
Giacomo Grassi\
Michele Ivan Bruna
]


== Player Information
#align(center)[
#raw(lang: "go", block: true,
`type User struct {
    UserID          []byte
    Salt            []byte
    Username        string
    Email           string
    PasswordHash    []byte
    Gachas          []Gacha
    Transactions    []Transaction
}`.text)]


== Gacha Collection

#align(center)[
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
#figure(
  image("../assets/drawio/monolithic_architecture.svg"),
)


== Database contents


== User’s interactions
