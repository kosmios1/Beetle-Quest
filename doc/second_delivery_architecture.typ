#import "@preview/kunskap:0.1.0": *
#import "@preview/codelst:2.0.1": sourcecode
#import "@preview/tablex:0.0.8": tablex

#import "@preview/silver-dev-cv:1.0.0": *


#show: cv.with(
  font-type: "PT Serif",
  continue-header: "false",
  address: "University of Pisa",
  name: "2nd_Delivery",
  lastupdated: "true",
  pagecount: "true",

  contacts: (
    (text: "Cosimo Giraldi", link: ""),
    (text: "Giacomo Grassi", link: ""),
    (text: "Michele Ivan Bruna", link: ""),
  ),
)

// about the project
#section[Introduction]
#descript[This report describes the architecture designed to support a Gacha game, characterized by a modular implementation based on microservices.

This architecture was enhanced and optimized using MicroFreshner. 

Through the analysis and detection of architectural smells—potential issues that could affect performance, scalability, or maintainability—we identified areas needing improvement and successfully refined the system's design.]

// about the project
#section[Architecture]
#figure(
  image("beetle-quest-microfreshner-architecture.png", width: 125%),
  caption: [
    Architecture of the system
  ],
)


#sectionsep
//Experience
#section("microservice")
    - marameo

    - marameo

    - marameo

#show: kunskap.with(
    title: [Report],
    author:"Cosimo Giraldi
Giacomo Grassi
Michele Ivan Bruna",
    header: "ASE Report",  //MODIFIE THIS
    date: datetime.today().display("[month repr:long] [day padding:zero], [year repr:full]"),

    //Paper size, fonts, and colors can optionally be customized as well

    // Paper size
    //paper-size: "a4",

    // Fonts
    //body-font: ("Noto Serif"),
    //body-font-size: 10pt,
    //raw-font: ("Hack Nerd Font", "Hack NF", "Hack", "Source Code Pro"),
    //raw-font-size: 9pt,
    //headings-font: ("Source Sans Pro", "Source Sans 3"),

    // Colors
    //link-color: link-color,
    //muted-color: muted-color,
    //block-bg-color: block-bg-color,
)

#lorem(21)

= A bit more detail

#lorem(49)

== Even more

#lorem(49)

=== Final comments
#lorem(12)

= Introduction Cloud Operations and DevOps
== The DevOps pipelines
*What is DevOps?*\

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


== How to Run the Project

1- Clone the repository to your local environment: git clone "link to the repository"

2-Move to `deploy foulder`
"cd deploy"

3- Rune the following command to build and launch the project 
"sudo docker compose up --build"
 
The gatcha application should now be running 

4- Access to the application:
Open your browser and go to:
https://localhost/static/

You can now use the "gatcha" application







== 4 interesting player operations

1- Use roll feature: 
    After the login the page of the profile is showed properly:
    1.1-Click the market Section the user sand this GET request for access to the market: 
        `GET /api/v1/market/auction/list HTTP/2`
    1.2-The system return 
        `HTTP/2 200 OK`

