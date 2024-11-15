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

== Service

microFreshner images

1-Admin-Service: Provides administrative capabilities, for monitoring, configuring, or management of the gatcha application. 

Images of Admin-service

2-Auth-Service: Handling the login and the registration of the users,and manage the token generation.

Images of Auth-service

== Service

3-User-Service: Responsible for managing user data and profiles. It communicates with the user-db for data storage and may support features like profile updates, user details.

Images of User-Service

4-Gacha-Service: 
The service handles main gacha functionality.
It supports the main features: roll, list, and research with gachaID. 
The service connects to the gacha-db for managing available gatcha pools and it's probabilities.
And connect to user-gacha-db to store user-specific data such as new gatcha aquired and it's information.

Images of Gacha-Service

== Service

5- Market-Service: Manages the marketplace functions, such as purchasing and the selling of the gatchas directly at full price from the game or for create an auction for sell gatchas, or bid for acquired new one. 
It connects to transaction-db and bid-db for storing market-related data, like transactions and information about the bid of auctions.

Images of Market-Service
