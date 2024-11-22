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
However, if a service needs to access data managed by another service, must use the internal API which is only accessible within the internal network.


The list of the microservices implemented in the system are:

== _admin-service_
This service provides the administrator with the necessary tools to manage the system in a controlled manner, allowing operations on users, gacha, and transactions and operation carried out in the market.

So this service allows to retrieve a specific admin-user with *FindByID*.
It's also possible to fetch the list of users with their associated information calling *GetAllUsers* API and perform detailed searching using *FindUserByID*. Can also manage users profile calling *UpdateUserProfile*, and can viewing all transactions carried out by a user using *GetUserTransactionHistory*.

This service also enables adding new gacha items to the system *AddGacha*, retrieving a complete list of all gacha items *GetAllGachas*, searching for a specific gacha by ID *FindGachaByID*, deleting a gacha from the system *DeleteGacha*, and upgrade an existing one *UpdateGacha*.

The service provides a transaction history within the marketplace *GetMarketHistory*, a list of all auctions *GetAllAuctions*, and the permission to search for a specific auction by ID *FIndAuctionByID*, as well as allowing modifications to an active auction *UpdateAuction*.


== _auth-service_ 
This service allows authentication within the web-app, guaranteeing access to the system.

So This service takes care for account registration *Register*, manages user *Login* and *Logout*, and also enables the *Verify* and the *Revoke* of the access token.

== oauth2-service

== _user-service_ 
This service allows a user to manage all the information that concerns them.

So this service allows users to view their personal information with *GetUserAccountDetails* and can update/modifie their personal profile with *UpdateUserAccountDetails*, with options to change their username, email address, and password. The user can also delete it's account with *DeleteUserAccount*.To make these changes, the current password must be entered for verification.
The users can also track their gacha collection with *GetUserGachaList* along with related details *GetUserGachaDetails*.


== _gacha-service_
This service allows to have control over the gachas in the user's collection. 

So this service offers the possibility to add a gacha to the user's collection with *AddGachaToUser*; Additionally, this service permits removes the gacha from the user's collection if it has been sold at auction with *RemoveGachaFromUser*.


== _maket-service_ 
This service allows you to perform all actions that involve the acquisition of BugsCoins and their movements for operation in the system.

So this service allows the user to perform a *RollGacha* to obtain gacha items and add them to their collection, if they do not already own them. Can *BuyGacha* from the market at full price, is possible also buy Bugscoin with *AddBugsCoins* and added to the personal wallet. Additionally, the user can *CreateAuction* and can *DeleteAuction* and can keep track of all active auctions, and can view the history of completed auctions with *AuctionList*. The user can *MakeBid* on auctions and view each *AuctionDetails*, such as: Auction ID, Owner ID, Gacha ID, Start Time, End Time, Winner ID, as well as the complete list of bids placed on the gacha, with details such as: Bid ID, User ID, Bugscoin Spent, and Time of Bid.

= Architecture with Microfreshner
The microservices architecture defined for this project is the result of a process of analysis and detection of the smells present in the original architectural prototype, carried out using MicroFreshner.

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


Sono stati introdotto anche dei timeOUt che permettono di 
Chiedere jack traefink