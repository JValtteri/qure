# Backend Server
[![Backend Tests](https://github.com/JValtteri/qure/actions/workflows/backend-tests.yml/badge.svg)](https://github.com/JValtteri/qure/actions/workflows/backend-tests.yml)
![Backend Coverage](https://github.com/JValtteri/qure/blob/badges/.badges/main/backend-coverage-badge.svg)

## Requirements

- [**Go**](https://go.dev/) 1.19 or newer

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

## Run server

```
./server
```

## Build instructions

Run this command in the `server` folder
```
go build
```
