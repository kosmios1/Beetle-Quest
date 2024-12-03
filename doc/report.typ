#import "@preview/wrap-it:0.1.0": wrap-content
#import "lib/frontpage.typ": report

#show: doc => report(
  title: "BeetleQuest",
  subtitle: "Advanced Software Engineering - 2nd Project Delivery",
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

#let makesubparagraph(title) = heading(numbering: none, outlined: false, level: 4)[#title]



= Introduction
The goal of this project is to develop a web app and define its architecture for creating a web-based gacha game. So the users will be able to engage in all the standard activities found in a gacha game like: _`roll`, `buy coin`, `create auctions`, `bid`_.

All these actions will be implemented through a _microservices_ architecture..


= Microservices
The main idea was to divide a monolithic system into a series of microservices, each of which handles a specific functionality.
This fragmentation allows for greater modularity and control of the system. To make the web-application more scalable, the microservices have been designed to be independent and stateless.
Microservices that need to store data use their own dedicated database, which they access directly.
However, if a service needs to access data managed by another service, it must use the internal API which is only accessible within the internal network.

#linebreak()

In the following pharagraps we will examinate the implemented services, and expose their functionalities.

== _Auth_
User registration, login and logout are all managed in a centralised manner by the same service: the Auth service.
Which also provides helper endpoints to check the validity of access tokens, allowing authentication and authorization within the web-app.

== _User_
This service is responsable for managing user's account informations. A user, once logged in, can access it's account details, modify them or delete the account itself.

== _Gacha_
Gacha collections are managed by the Gacha service. It allows users to get the list of available gachas as well as information on each one of them. User can inspect the personal inventory of different players and their personal one.

== _Market_
The Market service allows users to perform actions involving the acquisition of `BugsCoins` and gachas. It manages auctions lifetime and transactions in the system.

#linebreak()
Through this service users can obtain gachas by performing two actions: buy and roll. To roll the user has to pay 1000 `BugsCoins`, he/she will obtain a random gacha from the system with a probability which depends on the rarity of the gacha.

#linebreak()
The user has the permission to create and delete it's own auctions but can not bid to them, he/she can bid to other's auctions.

== _Static_
This service is responsible for serving the static content of the web-app, like the images, the _css_ and the _html_ files.


== _Admin_
This service provides the administrator with the necessary tools to manage the system in a controlled manner, allowing operations on users, gacha, and transactions and operation carried out in the market.

It can fetch the list of users with their associated information, performs detailed searching, modify users profile, view all the transactions carried out by a user and ispect user's auction list.

It can perform global actions on the gachas, like: add new one, modify/delete an existing one and get information on the system gachas. The service provides similar actions also on transactions and auctions.


= Architecture with MicroFreshner
The microservices architecture defined for this project is the result of a process of analysis and detection of the smells present in the original monolithic prototype, carried out using MicroFreshner.

#figure(
  //TODO: static service in microfreshner
  image("beetle-quest-microfreshner-architecture-v2.png", width: 110%),
  caption: [
    BeetleQuest architecture
  ],
)
= Architectural Design Choices

The architectural analysis of our initial system, carried out using MicroFreshner, revealed smell between the microservices. To isolate potential failures and improve the system's resilience, we introduced Circuit Breakers (CBs).

The introduced Circuit Breakers effectively address the issues caused by continuous failures of a microservice, preventing the cascading propagation of errors that could slow down or completely halt the entire system.

Moreover, to achieve more effective control over the system, we have introduced *_Timeouts_* on database connections. This solution significantly improves resilience and reliability. If a connection or query exceeds the maximum time defined by the timeout, the system considers the operation as failed and immediately activates error-handling mechanisms, ensuring a quick response and preventing bottlenecks or slowdowns.

We have also used a reverse proxy called *_Traefik_*, which acts as an intermediary between external users and the system's internal services. In this architecture, Traefik functions as an access gateway, managing and routing requests to the appropriate microservices, ensuring efficient and centralized traffic handling.


= User Stories
/*TO DO  
For each user story of the player put the endpoint(s) to realize it and the microservice(s) involved. Use a table or a bullet list.*/


= Market rules
/*TO DO 
Give a general description of the decision you took for the market.
Try to answer questions like:
• What happens to the currency of a player when someone else bids
higher?
• What happens if I bid at the last second of the auction?
• Can I bid for an auction in which I am the highest bidder?*/


= Testing
/*TO DO 
Write here any particular fact about the testing.
For example:
• I tested in isolation the DBManager_x together with the DB_x
• I used a third-party service that interacts with Service_y and I put them
together in the isolation test because mock it is hard for reason z.*/


= Secuirty

== Secuirty - Data
/*TO DO
• Select one input that you had to sanitize, describe what it represent,
which microservice(t) use it and how you sanitize it.
• List the data you encrypted at rest, describe what they represent, which
database stores them and where you en/decrypt them.*/


== Secuirty - Authentication and Authorization
/*TO DO Describe the scenario you selected (centralized vs distributed) by indicating
the basic steps to validate a token and how the keys to sign the token are
used and stored.
Try to describe it as schematic as possible (support it with lists, tables or
figures)
Put the payload format of your Access Token (bullet list, table or image)*/


== Security - Analyses
/* TO DO 
Put the screenshot of:
• The report of the static analysis tool you used (e.g. Bandit’s final table).
• The dashboard of docker scout with your (developed) images, where the
vulnerabilities are indicated.
If your language does not have static analysis tools available, specify it here.
Otherwise, put the command(s) you used to reach the results in the
screenshot.
Put also the name of the docker hub repository with the images.
Note: there is no need for docker scout for images of third-party software.*/


= Additional features
/* TO DO 
Describe here any additional feature you implemented.
• What is this feature?
• Why is it useful?
• How is it implemented?
*/

= Interesting Flows

Now we proceed analysing a few use case scenarios, to show the flow on the backend.


== Registration and login:

When a player wants to register he sends a `POST` request to the API Gateway at\ `/auth/register` containig the user's _username_, _email_ and _password_.

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
  - If the user exist and the provided data is correct the `auth` service returns, to the API Gateway, a response containig a `token` that authenticates the user.


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
- At this point the `market` service will extract randomly a gacha and, in the case that the user does not own that gacha, forward to the `gacha` service a save operation of the gacha to the user in question.
- If no error appears it returns a success message to the API Gateway.

== Create auction:

A user can create an acution sending a `POST` request to the API Gateway at\ `/market/auction` containig the _gacha id_ and the _expiration time_ of the action.

  - The Gateway send the request to the `auth` service.
  - The `auth` service checks for the validity of the `token`.
  - If the `token` id is valid then the request gets fowareded to the `market` service.
  - The `market` service will check if the user has the specified gacha in his inventory, comunicating with the `gacha` service, then it will check if the _expiration time_ is valid.
  - Then it will save the acution in the `market DB` and set a timed event in the `timed event DB` to close the auction.
  - If no error appears it returns a success message to the API Gateway.
  - If no other user bid to this auction before the _expiration time_, the `market` service will automatically remove the auction from the DB.

== Bid an auction:

To bid an auction a user has to send a `POST` request to the API Gateway at\ `/market/auction/<auctionId>/bid`, where `<auctionId>` is the id of the auction. The request has to include the amount the user wants to bid.

  - The Gateway send the request to the `auth` service.
  - The `auth` service checks for the validity of the `token`.
  - If the `token` id is valid then the request gets fowarded to the `market` service.
  - Now the `market` service will check with the `user` service if the user has the amount of `BugsCoins` he wants to bid.
  - If the check passes the `market` service will comunicate the `user` service to remove the amount from the bidder, and store the bid in the `market DB`.
  - If no error appears it returns a success message to the API Gateway.
rite here any particular fact about the testing.
For example:
• I tested in isolation the DBManager_x together with the DB_x
• I used a third-party service that interacts with Service_y and I put them
together in the isolation test because mock it is hard for reason z.