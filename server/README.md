# Backend Server
[![Backend Tests](https://github.com/JValtteri/qure/actions/workflows/backend-tests.yml/badge.svg)](https://github.com/JValtteri/qure/actions/workflows/backend-tests.yml)
![Backend Coverage](https://github.com/JValtteri/qure/blob/badges/.badges/main/backend-coverage-badge.svg)

## Index

- [Requirements](#requirements)
- [Setup](#setup)
    - [Configure](#configure)
        - [Example config.json](#example-configjson)
- [Run tests](#run-tests)
- [Running dev server](#running-a-dev-server)
- [Build the server](#build-the-server)
- [Run the server](#run-the-server)
- [Admin credentials](#admin-credentials)
- [Data persistence](#data-persistence)

#### See also:

- [QuRe](../README.md)
    - [Frontend documentation](../client/README.md)
    - **Backend documentation**


## Requirements

- [**Go**](https://go.dev/) 1.24 or newer

## Setup

### Configure

Configuration is done with `config.json` file. A sample file `config.json.example` is provided.
A valid `config.json` is required to start the server, but no one key is mandatory. Listed defaults are used for any missing keys. Technically the minimum valid config is `{}`.

#### Example config.json

```json
{
    "ORIGIN_URL": "localhost",
    "SERVER_PORT": "8000",
    "ENABLE_TLS": false,
    "CERT_FILE": "cert.pem",
    "PRIVATE_KEY_FILE": "privkey.pem"
}
```
If a value is not present in config.json the server uses default value. Default values should be sufficient for most cases

| Key | default | description |
| :--: | :--: | -- |
| `ORIGIN_URL` | `"localhost"` |  |
| `SERVER_PORT` | `"8080"` | port the server is listening on |
| `ENABLE_TLS` | `"false"` | Should the server use HTTPS |
| `CERT_FILE` | `"cert.pem"` | Certificate file for use with HTTPS (optional) |
| `PRIVATE_KEY_FILE` | `"privkey.pem"` | Private key file for use with HTTPS (optional) |
| `SOURCE_DIR` | `"../client/dist"` | Directory where the frontend files are stored |
| `DB_FILE_NAME` | `"./db/db.gob"` | Filename for the saving server data |
| `MIN_USERNAME_LENGTH` | `4` | Minimum allowed username length |
| `MIN_PASSWORD_LENGTH` | `8` | Minimum allowed password length |
| `MAX_SESSION_AGE` | `2592000` (30 d) | Maximum session cookie lifetime |
| `MAX_PENDIG_RESERVATION_TIME` | `600` (10 min) | Seconds to wait for reservation confirmation |
| `RESERVATION_OVERTIME` | `3600` (1h) | How much past reservation start time reservation is kept in system |
| `EXTRA_STRICT_SESSIONS` | `false` | Detect session key forgery. High resource useage |
| `MAX_THREADS` | `0` | Maximum allowed concurrent threads. Zero = automatic |
| `RATE_LIMIT_PER_MINUTE` | `60` | Maximum allowed requests/minute per client  |
| `RATE_LIMIT_PER_MINUTE_EVENT` | `120` | Maximum allowed requests/minute per client (event info requests) |
| `RATE_LIMIT_BURST` | `5` |  Maximum allowed burst (on top of base limit) |
| `RATE_LIMIT_RESET_MINUTES` | `60` |  Interval to clear reset limiters (to purge old clients and reset counters) |
| `RATE_LIMIT_ALERT` | `300`  |  Exceeding this number of blocked requests triggers an alert in log with offending IP address and blocked request count at last limit reset |


###### Times are expressed in seconds unless stated otherwise

## Run tests

This command will
- run tests
- output a coverage report
- open a detailed coverage report as a web page
```
go test -coverprofile cover.out && go tool cover -html=cover.out
```

## Running a dev server

You can run a dev server without building by using this command.
```
go run .
```
The [frontend](../client) dev server is configured to use this backend server for API calls. Follow the instructions in the frontend [README.md](../client)

## Build the server

Run this command in the `server` folder
```
go build
```

## Run the server

Run the binary created in the previous step
```
./server
```

## Admin credentials

If an `admin` account doesn't exist, a new `admin` account is created automatically on server start.

**Check server log output for the `admin` credentials.**

- Username is `admin`
- Password is a random string of characters.

**The password should be changed on first login.**

## Data persistence

State of the server is saved automatically into `server/db/db.gob`
