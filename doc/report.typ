#import "@preview/wrap-it:0.1.0": wrap-content
#import "lib/frontpage.typ": report

#show: doc => report(
  title: "2nd Delivery Project",
  subtitle: "Advanced Software Engineering a.y. 24/25",
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

To make the web-application more scalable, the microservices have been designed to be independent and stateless. Microservices that need to store data use their own dedicated database, which they access directly.
However, if a service needs to access data managed by another service, must use the internal API which is only accessible within the internal network.


= Microservices
The list of the microservices implemented in the system are:

== admin-service
This service provides the administrator with the necessary tools to manage the system in a controlled manner, allowing operations on users, gacha, and transactions and operation carried out in the market.

So this service allows to retrieve a specific admin-user with *_FindByID_*.
It's also possible to fetch the list of users with their associated information calling *_GetAllUsers_* API and perform detailed searching using *_FindUserByID_*. Can also manage users profile calling *_UpdateUserProfile_*, and can viewing all transactions carried out by a user using *_GetUserTransactionHistory_*.

This service also enables adding new gacha items to the system *_AddGacha_*, retrieving a complete list of all gacha items *_GetAllGachas_*, searching for a specific gacha by ID *_FindGachaByID_*, deleting a gacha from the system *_DeleteGacha_*, and upgrade an existing one *_UpdateGacha_*.

The service provides a transaction history within the marketplace *_GetMarketHistory_*, a list of all auctions *_GetAllAuctions_*, and the permission to search for a specific auction by ID *_FIndAuctionByID_*, as well as allowing modifications to an active auction *_UpdateAuction_*.


== auth-service 
This service allows authentication within the web-app, guaranteeing access to the system.

So This service takes care for account registration *_Register_*, manages user *_Login_* and *_Logout_*, and also enables the *_Verify_* and the *_Revoke_* of the access token.

== oauth2-service 

== user-service 
This service allows a user to manage all the information that concerns them.

So this service allows users to view their personal information with *_GetUserAccountDetails_* and can update/modifie their personal profile with *_UpdateUserAccountDetails _*, with options to change their username, email address, and password. The user can also delete it's account with *_DeleteUserAccount_*.To make these changes, the current password must be entered for verification.
The users can also track their gacha collection with *_GetUserGachaList _* along with related details *_GetUserGachaDetails_*.


== gacha-service
This service allows to have control over the gachas in the user's collection. 

So this service offers the possibility to add a gacha to the user's collection with *_AddGachaToUser_*; Additionally, this service permits removes the gacha from the user's collection if it has been sold at auction with *_RemoveGachaFromUser_*.


== maket-service 

This service allows you to perform all actions that involve the acquisition of BugsCoins and their movements for operation in the system.

So this service allows the user to perform a *_RollGacha_* to obtain gacha items and add them to their collection, if they do not already own them. Can *_BuyGacha_* from the market at full price, is possible also buy Bugscoin with *_AddBugsCoins_* and added to the personal wallet. Additionally, the user can *_CreateAuction_* and can *_DeleteAuction_* and can keep track of all active auctions, and can view the history of completed auctions with *_AuctionList_*. The user can *_MakeBid_* on auctions and view each *_AuctionDetails_*, such as: Auction ID, Owner ID, Gacha ID, Start Time, End Time, Winner ID, as well as the complete list of bids placed on the gacha, with details such as: Bid ID, User ID, Bugscoin Spent, and Time of Bid.

= Architecture with Microfreshner
This architecture was enhanced and optimized using MicroFreshner. 

Through the analysis and detection of architectural smells—potential issues that could affect performance, scalability, or maintainability—we identified areas needing improvement and successfully refined the system's design.