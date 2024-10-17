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

== Group Information

#align(center)[
#linebreak()
#linebreak()
*Cybersecurity - 2024/2025*\

Giraldi Cosimo\
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
  columns: (1fr, 1fr, 1fr, 1fr, 1fr),
  column-gutter: 0pt,
  row-gutter: 5pt,
  [*Common (C)*], [*Uncommon (U)*], [*Rare (R)*], [*Epic (E)*], [*Legendary (L)*],
  image("../assets/images/png/butterfly_legendary.png",      width: 60%),
  image("../assets/images/png/hercule_beetle_legendary.png", width: 60%),
  image("../assets/images/png/hercule_beetle_legendary.png", width: 60%),
  image("../assets/images/png/hercule_beetle_legendary.png", width: 60%),
  image("../assets/images/png/hercule_beetle_legendary.png", width: 60%),
  //
  image("../assets/images/png/butterfly_legendary.png",      width: 60%),
  image("../assets/images/png/hercule_beetle_legendary.png", width: 60%),
  image("../assets/images/png/hercule_beetle_legendary.png", width: 60%),
  image("../assets/images/png/hercule_beetle_legendary.png", width: 60%),
  image("../assets/images/png/hercule_beetle_legendary.png", width: 60%),
  // 
  image("../assets/images/png/butterfly_legendary.png",      width: 60%),
  image("../assets/images/png/hercule_beetle_legendary.png", width: 60%),
  image("../assets/images/png/hercule_beetle_legendary.png", width: 60%),
  image("../assets/images/png/hercule_beetle_legendary.png", width: 60%),
  image("../assets/images/png/hercule_beetle_legendary.png", width: 60%),
  //
  image("../assets/images/png/butterfly_legendary.png",      width: 60%),
  image("../assets/images/png/hercule_beetle_legendary.png", width: 60%),
  image("../assets/images/png/hercule_beetle_legendary.png", width: 60%),
  image("../assets/images/png/hercule_beetle_legendary.png", width: 60%),
  image("../assets/images/png/hercule_beetle_legendary.png", width: 60%),
)]

= Architecture

== 

#align(center)[
]

