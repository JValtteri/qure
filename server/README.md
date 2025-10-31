# Backend Server
[![Backend Tests](https://github.com/JValtteri/qure/actions/workflows/backend-tests.yml/badge.svg)](https://github.com/JValtteri/qure/actions/workflows/backend-tests.yml)
![Backend Coverage](https://github.com/JValtteri/qure/blob/badges/.badges/main/backend-coverage-badge.svg)

## Requirements

- [**Go**](https://go.dev/) 1.24 or newer

## Setup

### Configure

Configuration is done with `config.json` file. A sample file `config.json.example` is provided.

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

| Key | default | description |
| :--: | :--: | -- |
| `ORIGIN_URL` | `"localhost"` |  |
| `SERVER_PORT` | `"8000"` | port the server is listening on |
| `ENABLE_TLS` | `"false"` | Should the server use HTTPS |
| `CERT_FILE` | `"cert.pem"` | Certificate file for use with HTTPS (optional) |
| `PRIVATE_KEY_FILE` | `"privkey.pem"` | Private key file for use with HTTPS (optional) |

## Run tests

This command will
- run tests
- output a coverage report
- open a detailed coverage report as a web page
```
go test -coverprofile cover.out && go tool cover -html=cover.out
```

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

### Admin credentials

If an `admin` account doesn't exist, a new `admin` account is created automatically on server start.

**Check server log output for the `admin` credentials.**

- Username is `admin`
- Password is a random string of characters.

**The password should be changed on first login.**
