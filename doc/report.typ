#import "@preview/wrap-it:0.1.0": wrap-content
#import "lib/frontpage.typ": report

#show: doc => report(
  title: "2nd Delivery Project",
  subtitle: "Advanced Software Engineering",
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
The goal of this project is to develop a web app and define its architecture for creating a web-based gacha game. So the users will be able to engage in all the standard activities found in a gacha game like: *_Roll, Buy Coin, Auction, Bid_*.   

All these actions will be implemented through a _microservices_ architecture..


= Microservices
The main idea was to divide a monolithic system into a series of microservices, each of which handles a specific functionality.
This fragmentation allows for greater modularity and control of the system. To make the web-application more scalable, the microservices have been designed to be independent and stateless. 
Microservices that need to store data use their own dedicated database, which they access directly.
However, if a service needs to access data managed by another service, it must use the internal API which is only accessible within the internal network.


The list of the microservices implemented in the system are:

== _Admin_
This service provides the administrator with the necessary tools to manage the system in a controlled manner, allowing operations on users, gacha, and transactions and operation carried out in the market.

It can fetch the list of users with their associated information, performs detailed searching, modifie users profile, can view all the transactions carried out by a user and can ispect user's auction list. 

It can perform global actions on the gachas, like: add new one, modify/delete an existing one and get information on the system gachas. The service provides similar actions also on transactions and auctions. 



== _Auth_ 
Allows authentication and authorizations within the web-app.

This service takes care for account *Registration*, user *Login* and *Logout*. it also has helper endpoints to check the validity of the submitted access token.


== _User_ 
This service manages user's information. The actions which can be performed are:
getting user's account details, modify user's account details, delete user's account.


== _Gacha_
This service allows control over the gacha collection, it enables a user to get a view on the system's gachas, as well as getting information on a specific gacha. Other than the previous actions, the user can inspect the personal inventory of different players, and also get details about a specific gacha in the retrieved inventory.


== _Market_ 
The Market service allows users to perform  actions involving the acquisition of BugsCoins and gachas. It manages auctions lifetime and transactions in the system.

Through this service users can obtain gachas by performing two actions: buy and roll. To roll the user has to pay 1000 BugsCoins, he/she will obtain a random gacha from the system with: the probability depends on the rarity of the gacha.

The user has the permission to create and delete it's own auctions but can not bid to them, he/she can bid to other's auctions.


== _Static_
This service is responsible for serving the static content of the web-app, like the images and the css files.


= Architecture with Microfreshner
The microservices architecture defined for this project is the result of a process of analysis and detection of the smells present in the original monolithic prototype, carried out using MicroFreshner.

#figure(
  image("beetle-quest-microfreshner-architecture-v2.png", width: 125%),
  caption: [
    beetle-quest architecture
  ],
)
= Architectural Design Choices

The architectural analysis of our initial system, carried out using MicroFreshner, revealed smell between the microservices. To isolate potential failures and improve the system's resilience, we introduced Circuit Breakers (CBs).
 
The introduced Circuit Breakers effectively address the issues caused by continuous failures of a microservice, preventing the cascading propagation of errors that could slow down or completely halt the entire system.

Moreover, to achieve more effective control over the system, we have introduced *_Timeouts_* on database connections. This solution significantly improves resilience and reliability. If a connection or query exceeds the maximum time defined by the timeout, the system considers the operation as failed and immediately activates error-handling mechanisms, ensuring a quick response and preventing bottlenecks or slowdowns.

We have also used a reverse proxy called *_Traefik_*, which acts as an intermediary between external users and the system's internal services. In this architecture, Traefik functions as an access gateway, managing and routing requests to the appropriate microservices, ensuring efficient and centralized traffic handling.


//TODO: static service in microfreshner to do
= Interesting Flows

= Conclusion
