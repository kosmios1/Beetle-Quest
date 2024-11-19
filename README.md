# Beetle Quest

A beautifull and fullfilling web based gacha system.

## Get Started

This guide will help you set up and run the application using Docker and Docker Compose. By following these steps, you'll have a local version of the application running on your machine.
Prerequisites:

Before you begin, ensure you have the following installed:

- Docker
- Docker Compose

### Cloning the repo

Clone the repository to your local machine using Git.

```bash
git clone https://github.com/iTzMorderin/Beetle-Quest
cd Beetle-Quest
```

### Configuring the application

The application uses environment variables to configure itself. Their values can be modified inside the docker compose file located in the deploy directory. Instead the `secrets` values are stored inside inside `secrets/` folder with the following structure:

```bash
deploy/
├──secrets/
│  ├── jwt.env
│  ├── oauth2.env
│  ├── postgress.env
│  └── redis.env
...
└──compose.yml
```

### Starting the application

To start the application, run the following commands:

```bash
cd beetle-quest/deploy
docker compose up
```

If you want to run the application in the background, you can use the `-d` flag:

```bash
docker compose up -d
```

### Stopping the application

To stop the application, run the following command:

```bash
cd beetle-quest/deploy
docker compose down
```

## Tests

You fill find the Postman collection file`beetle-quest-collection.json` inside `beetle-quest/tests/postman/`, you can execute them with Postam Newman:

```sh
docker run --rm --net beetle-quest_internal -v <path/to/this/repo>/beetle-quest/tests/postman/collection.json:/collection.json postman/newman run /collection.json --insecure --ignore-redirects --color on
```

## Project structure

The project is structured as follows:

```bash
.
├── assets/
│   ├── drawio
│   ├── gimp
│   └── images
├── beetle-quest/
│   ├── api/
│   ├── cmd/
│   ├── deploy/
│   ├── go.mod
│   ├── go.sum
│   ├── internal/
│   ├── pkg/
│   └── templates/
├── doc/
└── README.md
```

## References

- [Project Structure 1](https://betterprogramming.pub/how-are-you-structuring-your-go-microservices-a355d6293932)
- [Project Structure 2](https://gochronicles.com/project-structure/)
