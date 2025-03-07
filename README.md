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
│  ├── jwt.env
│  ├── oauth2.env
│  ├── postgress.env
│  └── redis.env
...
└──compose.yml
```

To generate the needed SSL/TLS certificates execute:

```bash
cd beetle-quest/deploy/cacerts
./generate-db-certs.sh

cd ../traefik/certs/
./create_internal_traefik_cert.sh
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

> [!NOTE]
> The service gui is available at `https://localhost/static`.

> [!WARNING]
> The admin interface is not available, but it can be accessed with tools like `curl`. The admin's
> login requires `otp`, the seed can be retrieved from the qrcode inside `assets/admin_otp_qrcode.png`.
> Admins CANNOT be created, the credentials for the development one are:
>
> - AdminID: `09087f45-5209-4efa-85bd-761562a6df53`
> - Password: `admin`
> - OTP: retrieved from the qrcode

An example request to login as an admin:

```bash
curl -v -k -X POST \
  -H "Content-Type: application/json" \
  -d '{"admin_id": "09087f45-5209-4efa-85bd-761562a6df53", "password": "admin", "otp_code": "<OTP_CODE>"}' \
  https://localhost/api/v1/auth/admin/login
```

### Stopping the application

To stop the application, run the following command:

```bash
cd beetle-quest/deploy
docker compose down -v
```

## Tests

> [!IMPORTANT]
> The tests needs to be executed in a clean environment. If the system has been used before we recommend using a
> a `docker compose down -v` to remove all volumes and start from scratch.

### Postman

You fill find the Postman collection file `collection.json` inside `beetle-quest/tests/postman/`, you can execute them with Postman Newman:

```sh
cd ./beetle-quest/tests/postman/
docker run --rm --net beetle-quest_internal --net beetle-quest_admin -v ./beetle-quest.json:/collection.json postman/newman run /collection.json --bail --insecure --ignore-redirects --color on
```

> [!NOTE]
> If you are using the desktop postaman app remember to:
>
> - allow the programmatic modification of cookies on the `localhost` domain;
> - use `https://localhost/api/v1/` as the `baseUrl` and `adminUrl`;
> - use `https://localhost/` as the `hostUrl` and `adminHostUrl`.

### Locust

Locust's tests can be found inside `beetle-quest/tests/locust/` folder, to start Locust, and run the tests, these commands have to be executed:

```sh
cd beetle-quest/tests/locust/
docker build -t beetle-quest-locust:latest .
docker run --rm --network=beetle-quest_external -p 127.0.0.1:8089:8089 beetle-quest-locust:latest
```

### How to run unit tests on a single service

> [!IMPORTANT]
> Being in test mode the microservices which are started in this mode have the following characteristics:
>
> - mTLS is disabled.
> - The service is available at `http://localhost:8080`.
> - Authorization is provided by an API key instead of Oauth2 (to use mock tokens in testing environment).
> - Data are stored in memory.
> - The services are not behind a reverse proxy.
> - The services will not be able to communicate with the other services, they simulate the other services with a memory storage thanks to the mock repositories.

> [!NOTE]
> The unit tests are developed to test one service at the time. The right procedure would be to run a specific microservice and
> its test, found in `./beetle-quest/tests/unit/postman/beetle_quest-<service-name>_service-unit_tests.json`.
> Then stop the service, run the next service with its tests.

The first thing to do is to build a test image of a service (from the deploy folder), for example the `auth` service:

```sh
cd beetle-quest/deploy/
docker build -t beetle-quest-auth:test -f ./auth/Dockerfile ..
```

Now we can run the image using the correct parameters:

```sh
docker run --rm -it -p 8080:8080 -e JWT_SECRET_KEY="e6df59f91871f2229a0296c6b5ffaf44cef6af30cd05057857b9f0a74b0d28c1" beetle-quest-auth:test
```

#### Utility script

To start services in test mode there is `deploy/beetle_quest_unit_tests_utility.sh` script, it can be used to start a service in test mode.

Help message:

```sh
Usage: ./beetle_quest_unit_tests_utility.sh <service-name> [options]
Options:
  -p, --port PORT      Container port mapping (default: 8080)
  -e, --env KEY=VALUE  Environment variables
  -h, --help           Show this help message
```

```sh
./beetle_quest_unit_tests_utility.sh -p 8080 --env JWT_SECRET_KEY="e6df59f91871f2229a0296c6b5ffaf44cef6af30cd05057857b9f0a74b0d28c1" <service-name>
```

> [!NOTE]
> 'service-name' value can be:
>
> - `auth`
> - `gacha`
> - `user`
> - `market`
> - `admin`

## Project structure

The project is structured as follows:

```bash
.
├── assets/
│   ├── admin_otp_qrcode.png
│   ├── drawio
│   ├── gimp
│   └── images
├── beetle-quest/
│   ├── api/
│   ├── cmd/
│   ├── deploy/
│   ├── go.mod
│   ├── go.sum
│   ├── internal/
│   ├── pkg/
│   ├── templates/
│   └── tests/
├── doc/
└── README.md
```
