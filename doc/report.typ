#import "@preview/wrap-it:0.1.0": wrap-content
#import "lib/frontpage.typ": report

#show: doc => report(
  title: "BeetleQuest",
  subtitle: "Advanced Software Engineering - Project Delivery",
  authors: ("Cosimo Giraldi", "Giacomo Grassi", "Michele Ivan Bruna"),
  date: "2024/2025",
  doc,
  imagepath: "marchio_unipi_black.svg"
)

// Code block style
#show raw.where(block: true): block.with(
  fill: luma(240),
  inset: 10pt,
  radius: 10pt,
)

#let makesubparagraph(title, level: 4) = heading(numbering: none, outlined: false, level: level)[#title]



= Introduction
The goal of this project is to develop a web app and define its architecture for creating a web-based gacha game. So the users will be able to engage in all the standard activities found in a gacha game like: _`roll`, `buy coin`, `create auctions`, `bid`_.

All these actions will be implemented with *_Go_* language and through a _microservices_ architecture.

= Gacha Collection

The gachas are fictional creatures inspired by #text(red)[_todo:..._]. Below are a few examples of these imaginative beings.


#[
  #set align(center)
  #set table(
    stroke: (x, y) => if y == 0 {
      (bottom: 0.0pt + black)
    },
  )

  #table(
    columns: (auto, auto),
    column-gutter: -8em,
    align: horizon,
    //fill: (x, _) =>
    //  if calc.odd(x) { luma(240) }
    //  else { white },
    table.cell(rowspan: 2)[#image("beetle-quest-microfreshner-architecture-v2.png", width: 50%)],
    image("beetle-quest-microfreshner-architecture-v2.png", width: 50%),
    image("beetle-quest-microfreshner-architecture-v2.png", width: 50%),
  )
]

#figure(
  image("../assets/images/currency_cut.png", width: 45%),
  caption: [The currency used within the game is called `BugsCoins`],
)

/* TO DO
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
        image("../assets/images/png/assassin_mosquito_rare.png",      width: 60%),
        image("../assets/images/png/assassin_peacock_epic.png",      width: 60%),
        image("../assets/images/png/demoniac_mosquito_legendary.png",      width: 60%),
    )],
    [#set align(horizon)
    #figure(
        // image("../assets/images/png/coin.png", width: 80%),
        image("../assets/images/currency_cut.png", width: 80%),
        caption: [Currency],
        //caption-pos: top, // V 0.12
        supplement: [],
        numbering: none,
    )
    ]
    )
]
*/


= Architecture

The microservices architecture defined for this project is the result of a process of analysis and detection of the smells present in the original monolithic prototype, carried out using MicroFreshner.

#figure(
  image("beetle-quest-microfreshner-architecture-v2.png", width: 110%),
  caption: [
    BeetleQuest architecture
  ],
)

== Design Choices

The architectural analysis of our initial system, carried out using MicroFreshner, revealed smell between the microservices. To isolate potential failures and improve the system's resilience, we introduced Circuit Breakers (CBs).

The introduced Circuit Breakers effectively address the issues caused by continuous failures of a microservice, preventing the cascading propagation of errors that could slow down or completely halt the entire system.

To achieve more effective control over the system we have introduced *_Timeouts_* on database connections. This solution significantly improves resilience and reliability. If a connection or query exceeds the maximum time defined by the timeout, the system considers the operation as failed and immediately activates error-handling mechanisms, ensuring a quick response and preventing bottlenecks or slowdowns.

We have also used a reverse proxy called *_Traefik_*, which acts as an intermediary between external users and the system's internal services. In this architecture, Traefik functions as an access gateway, managing and routing requests to the appropriate microservices, ensuring efficient and centralized traffic handling.

== Microservices

The main idea was to divide a monolithic system into a series of microservices, each of which handles a specific functionality.
This fragmentation allows for greater modularity and control of the system. To make the web-application more scalable, the microservices have been designed to be independent and stateless.
Microservices that need to store data use their own dedicated database, which they access directly.
However, if a service needs to access data managed by another service, it must use the internal API which is only accessible within the internal network.

#linebreak()

In the following pharagrap we will examinate the implemented services, and their functionalities.
Eachone of the services, except for the _Static service_, has his own #link("https://www.postgresql.org/")[PostgreSQL] DB. Furthermore user sessions and market-timed-events, that will be discussed later, are stored in #link("https://redis.io/")[Redis] DBs.

#makesubparagraph([_Auth_], level: 5)
User registration, login and logout are all managed by the Auth service, which also checks the validity of access tokens, allowing authentication and authorization within the application.

#makesubparagraph([_User_], level: 5)
This service is responsable for managing user's account informations. A user, once logged in, can access it's account details, modify them or delete the account itself.

#makesubparagraph([_Gacha_], level: 5)
The Gacha service manages collections, providing users with a list of available gachas and details about each one, as well as access to inspect the personal inventories of various players.

#makesubparagraph([_Market_], level: 5)
The Market service allows users to perform actions involving the acquisition of `BugsCoins` and gachas. It manages auctions lifetime and transactions in the system. Through this service users can obtain gachas by either buying or rolling for a random gacha based on rarity.

#makesubparagraph([_Static_], level: 5)
This service is responsible for serving the static content of the web-app, like the images, the _css_ and the _html_ files.

#makesubparagraph([_Admin_], level: 5)
This service provides the administrator with the necessary tools to manage the system in a controlled manner: allowing operations on users, gacha, and transactions and market events.

#makesubparagraph([_API gateway_], level: 5)
There are two reverse proxies that implent the circuit breakers, the load balancers and the the API gateway. One is exclusive for the admin's operation the other for the clients' ones. Reverse proxies are implemented with #link("https://traefik.io/")[Traefik].

#linebreak()

=== Microservices connections

#makesubparagraph([Admin-Service ↔ Gacha-Service:], level: 5)
The _admin service_ connects with the _gacha service_ to manage gacha, such as adding/delete/modify gachas.

#makesubparagraph([Admin-Service ↔ Market-Service:], level: 5)
_Admin service_ interacts with _market service_ to regulate or manage the marketplace, including listing auctions, listing transactions, or update/modify auction.

#makesubparagraph([Admin-Service ↔ User-Service:], level: 5)
The _admin service_ connects with _user service_ to manage user accounts, such as listing users, modifying user profiles, or checking user transaction history and the user auction list.

#makesubparagraph([Auth-Service ↔ User-Service:], level: 5)
The _auth service_ relies on _user service_ for user data, such as validating credentials.

#makesubparagraph([Gacha-Service ↔ User-Service:], level: 5)
The _gacha service_ connects with the _user service_ to manage the user's gacha collection, such as listing the user's gacha collection or checking the gacha details.

#makesubparagraph([Market-Service ↔ User-Service:], level: 5)
The _market service_ connects with the _user service_ to manage the user's currency and transactions, such as checking the user's currency and adding currency.

#makesubparagraph([Market-Service ↔ Gacha-Service:], level: 5) The _market service_ connects with the _gacha service_ to manage the gacha collection, such as listing the gacha collection or checking the gacha details or when a gacha is sold in the market.

#linebreak()


= User Stories: Player

Evrey request has to pass through the _gateway_ and the _auth-service_ and _session-db_, to check if it is a valid request. So those services are omitted in the list of the microservice(s) involved for the following requests.

== Account

- I want to be able to register to the system, so that I can access the game.
  - `/auth/register` (_user-service,user-db_)

- I want to be able to delete my account, so that I can remove my information to the game.
  - `/user/account/{{userId}}` (_user-service,user-db/gacha-service,gacha-db/market-service,market-db_)

- I want to be able to modify my account information, so that I can update my profile.
  - `/user/account/{{userId}}` (_user-service,user-db_)

- I want to be able to login and logout, so that I can access and leave the game.
  #linebreak()
  I want be safe from unauthorized access, so that my account access information is protected.
  - `/auth/logout`, `/auth/login` (_auth-service,auth-db,session-db_)

== Collection

-  I want to see my gacha collection, so that I can see what I have.
  - `/gacha/user/{{userId}}/list` (_gacha-service,gacha-db_)

- I want to see the info of a gacha in my collection, so that I can see the details of a gacha.
  - `/gacha/{{gachaId}}/user/{{userId}}` (_gacha-service,gacha-db_)

- I want to see the system gacha collection, so that I can see what I can get.
  - `/gacha/list` (_gacha-service,gacha-db_)

- I want to see the info of a gacha in the system collection, so that I can see the details of a gacha.
  - `/gacha/{{gachaId}}` (_gacha-service,gacha-db_)

== Currency

- I want to use in-game currency for roll a gacha, so that I can get a random gacha.
  - `market/gacha/roll` (_user-service,user-db/market-service,market-db_)

- I want to buy in-game currency, so that I can get more gachas.
  - `/market/bugscoin/buy` (_market-service,market-db/user-service,user-db_)

- I want to be safe about the in-game currency transactions, so that my money is protected.
  - `/auth/logout`, `/auth/login` (_auth-service,session-db_)


== Market

- I want to see the auction market, so that i can evaluate if buy/sell a gacha.
  - `/market/auction/list` (_market-service,market-db_)

- I want to set an auction for one of my gacha, so that I can sell it.
  - `/market/auction/` (_gacha-service,gacha-db/market-service,market-db,market-timed-events_)

- I want to bid for a gacha from the market, so that I can buy it.
  #linebreak()
  I want to receive a gacha when i win an auction, so that I receive a gacha.
  #linebreak()
  I want to receive in-game currency when someone win my auction, so that I sell work as I expect.
  #linebreak()
  I want to receive my in-game currency back when i lost an auction, so that my in-game currency.
  #linebreak()
  I want to that the auctions cannot be temperes, so that my in-game currency and collection are safe.
  - `/market/auction/{{auctionId}}/bid` (_user-service,user-db/market-service,market-db/gacha-service,gacha-db,market-timed-events_)

- I want to view my transaction history, so that I can track my market movements.
  - `/internal/market/get_transaction_history` (_market-service,market-db_)


= User Stories: Admin

All the following endpoints requests involve the _admin-service_ and _admin-db_.

== Account

- I want to login and logout as admin from the system, so that I can access and leave the game.
  - `/auth/admin/login`,`/auth/logout` (_auth-service,auth-db,session-db_)

- I want to check all users account/profile, so that I can monitor all the users accounts/profiles.
  - `/admin/user/get_all` (_user-service,user-db_)

- I want to check a specific user account/profile, so that I can monitor user account/profile.
  #linebreak()
  I want to modify a specific user account/profile, so that I can update a specific user account/profile.
  - `/admin/user/{{userId}}` (_user-service,user-db_)

- I want to check a specific player currency transaction history, so that I can monitor the transactions of a player.
  - `/admin/user/{{userId}}/transaction_history` (_user-service,user-db/market-service,market-db_)

- I want to check a specific player market history, so that I can monitor the market of a player.
  - `/admin/user/{{userId}}/auction/get_all` (_user-service,user-db/market-service,market-db_)


== Gacha

- I want to check all the gacha collection, so that I can check all the collection.
  - `/admin/gacha/get_all` (_gacha-service,gacha-db_)

- I want to modify the gacha collection, so that I can add gachas.
  - `/admin/gacha/add` (_gacha-service,gacha-db_)

- I want to modify the gacha collection, so that I can delete gachas.
  #linebreak()
  I want to check a specific gacha, so that I can check the status of a gacha.
  #linebreak()
  I want to modify a specific gacha information, so that I can modify the status of a gacha.
  - `/admin/gacha/{{gachaId}}` (_gacha-service,gacha-db_)


== Market

- I want to see the auction market, so that I can monitor the auction market.
  - `/admin/market/auction/get_all` (_market-service,market-db_)

- I want to see a specific auction, so that I can monitor a specific auction of the market.
  #linebreak()
  I want to modify a specific auction, so that I can update the status of a specific auction.
  - `/admin/market/auction/{{auction_id}}` (_market-service,market-db/gacha-service,gacha-db_)

- I want to see the market history, so that I can check the market old auctions.
  - `/admin/market/transaction_history` (_market-service,market-db_)

#linebreak()

= Market rules

The market service has been implemented whit the following rules in mind:

- The user has the permission to create and delete it's own auctions but can not bid to them, he/she can bid to other's auctions

- When a user places a higher bid than the previous one, the currency of the previous highest bid is returned to the user after the finish of the auction.

- If someone places a bid at the very last second of the auction, they will win the gacha as the last valid bidder.

- It's also possible to bid on an auction where you are already the highest bidder. However, the user cannot place a bid if they do not have the required amount of coins to bid.

- Additionally, the owner of an auction cannot bid on their own auction. As the owner, you can delete the auction at any time before it expires, but you need to confirm the action by entering your password. The maximum duration of an auction is 24 hours.

- All bids that are expired will be refunded at the end of the auction.

- All auctions remain visible to users, along with all the auction details. Additionally, all bids made are displayed showing the bidder details.

#linebreak()

= Testing

The tests were conducted using mocks that allowed for the isolated testing of individual services. These mocks simulated the behavior of external components, enabling the verification of each service's functionality without relying on real external resources. Both unit and integration tests where carried out with #link("https://www.postman.com/")[Postman].

To conduct the test a performance testing tool, #link("https://locust.io/")[Locust], was used that allowed for load simulation and analysis of the service responses in various scenarios, ensuring an accurate assessment of the performance and robustness of each component

#linebreak()

= Security

== Data

#text(red)[TODO:] Select one input that you had to sanitize, describe what it represent, which microservice(s) use it and how you sanitize it.

#linebreak()

In the application, the databases are implemented using PostgreSQL or Redis. For PostgreSQL, Transparent Data Encryption (TDE) is used. TDE is a technology that protect sensitive data by encrypting the database files at rest. It ensures that data stored on disk is encrypted, making it inaccessible to unauthorized users or applications, while it automatically encrypts the data before it is written to disk and decrypts it when it is read.

On the other hand Redis data is not encrypted. This decision is mainly driven by its architecture as an in-memory database, which means that the data is not stored persistently on disk.

#linebreak()

All connections between databases/services and services use mutual TLS (mTLS), ensuring secure communication and authentication between the involved parties .


== Authentication and Authorization

#text(red)[TODO]

/*TO DO Describe the scenario you selected (centralized vs distributed) by indicating
the basic steps to validate a token and how the keys to sign the token are
used and stored.
Try to describe it as schematic as possible (support it with lists, tables or
figures)
Put the payload format of your Access Token (bullet list, table or image)*/


== Analyses

For the static code analysis, #link("https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck")[`govulncheck`] was used, which identifies vulnerabilities in Go dependencies by checking against the Go vulnerability database.

#figure(
  image("beetle-quest-microfreshner-architecture-v2.png", width: 55%),
  caption: text(red)[TODO: govulncheck report],
)

#linebreak()

Meanwhile, for the analysis of Docker images, #link("https://trivy.dev/")[`trivy`] was employed. In addition to analyzing images for CVEs using commands like `docker scan`, it also allows for the examination of Go binaries for vulnerable dependencies, misconfigurations, and potential leaks of secrets.

#linebreak()

The following resuls can be obtained executing `./scan-images.sh`, which is placed in `beetle-quest/tests/`, the output will be found inside `trivy_scan_results/` in the same folder.

#figure(
  image("beetle-quest-microfreshner-architecture-v2.png", width: 55%),
  caption: text(red)[TODO: trivy report],
)



= Additional features

#text(red)[TODO: describe]

- mTLS
- Shared CA
- Web Gui
- OAuth2.0
- Buy gacha
- ...

/* TO DO
Describe here any additional feature you implemented.
• What is this feature?
• Why is it useful?
• How is it implemented?
*/

/*
  // Not needed anymore

  = Interesting Flows

  Now we proceed analysing a few use case scenarios, to show the flow on the backend.


  == Registration and login:

  When a player wants to register he sends a `POST` request to the API Gateway at\ `/auth/register` containig the user's _username_, _email_ and  _password_.

    - The Gateway forwards the request to the `auth` service.
    - The `auth` service checks for the validity of the provided data;
    - If no error occurs it creates the new user;
    - It sends a request to the `user` service to create the new user data.
    - The user service store the new user data in the `user DB`.
    - The `user` service forward the action's status to the `auth` service.
    - If `auth` service gets no error it returns to the API Gateway, a success message.

  Now the user can login trough a `POST` request to the API Gateway at `/auth/login` containig the the _username_ and the _password_.

    - The Gateway forwards the request to the `auth` service.
    - The `auth` service checks for the validity of the provided data comunicating with the `user` service.
    - The `user` service checks if the user exist in the `user DB`, and return it's data to the `auth` service.
    - If the user exist and the provided data is correct the `auth` service returns, to the API Gateway, a response containig a `token` that  authenticates the user.


  #linebreak()

  From now on we assume that all the requests contain the authentification `token` and that every microservice will obtain the `user_id` from it.


  == Roll gacha:

  To roll for a gacha the user must send a `GET` request to the API Gateway at\ `/market/gacha/roll`

  - The Gateway send the request to the `auth` service.
  - The `auth` service checks for the validity of the `token`.
  - If the `token` id is valid then the request gets fowareded to the `market` service.
  - The `market` service will ask the `user` service if the user exists and it's data, then it checks if it has at least 1000 `BugsCoins`,
  - If so it removes that amount of money form the user, saving the transaction in the `market DB`.
  - The  `market` service will request the `gacha` service to get the list of all the gachas.
  - At this point the `market` service will extract randomly a gacha and, in the case that the user does not own that gacha, forward to the `gacha`   service a save operation of the gacha to the user in question.
  - If no error appears it returns a success message to the API Gateway.

  == Create auction:

  A user can create an acution sending a `POST` request to the API Gateway at\ `/market/auction` containig the _gacha id_ and the _expiration time_ of  the action.

    - The Gateway send the request to the `auth` service.
    - The `auth` service checks for the validity of the `token`.
    - If the `token` id is valid then the request gets fowareded to the `market` service.
    - The `market` service will check if the user has the specified gacha in his inventory, comunicating with the `gacha` service, then it will check   if the _expiration time_ is valid.
    - Then it will save the acution in the `market DB` and set a timed event in the `timed event DB` to close the auction.
    - If no error appears it returns a success message to the API Gateway.
    - If no other user bid to this auction before the _expiration time_, the `market` service will automatically remove the auction from the DB.

  == Bid an auction:

  To bid an auction a user has to send a `POST` request to the API Gateway at\ `/market/auction/<auctionId>/bid`, where `<auctionId>` is the id of the  auction. The request has to include the amount the user wants to bid.

    - The Gateway send the request to the `auth` service.
    - The `auth` service checks for the validity of the `token`.
    - If the `token` id is valid then the request gets fowarded to the `market` service.
    - Now the `market` service will check with the `user` service if the user has the amount of `BugsCoins` he wants to bid.
    - If the check passes the `market` service will comunicate the `user` service to remove the amount from the bidder, and store the bid in the  `market DB`.
    - If no error appears it returns a success message to the API Gateway.
  rite here any particular fact about the testing.
  For example:
  • I tested in isolation the DBManager_x together with the DB_x
  • I used a third-party service that interacts with Service_y and I put them
  together in the isolation test because mock it is hard for reason z.
*/
