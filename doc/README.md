# QuRe Reservation System Documentation

## Index

[Return to Root](../README.md)
- [Index](#index)
- [Deploying With Dockrer](#deploying-with-docker)
- [Deploying Without Docker](#deploying-without-docker)
- [Configuring](#configuring)
- [Using](#using)
- [Licenses](#licenses)
- [Developer Documentation](#developer-documentation)

## Deploying With Docker

```
qure
 ├─ docker-compose.yml
 ├─ config.json
 ├─ db
 │  └─ db.gob
 ├─ logo.png
 └─ images
    └─ ...
```

1. Copy [docker-compose.yml](../docker-compose.yml).
1. Copy [config.json.example](../server/config.json.example) and rename it to `config.json`.
1. [Configure](#configuring) as necessary
1. Check [docker-compose.yml](../docker-compose.yml) as necessary:
    - image version
    - mount settings
    - port
1. Run `docker compose up`
1. Read the log output for any issues
1. Copy the initial admin password from the logs
1. Log in with the admin credentials and change the password.
1. Press `D` to detach from server console

You can check the logs at any time with
```sh
docker logs qure-app-1
```

The server state is saved to `db.gob` on shutdown. The file is created automatically if it doesn't exist

Adding `logo.png` is optional. Adding one will replace the default logo.

`./images/` can be mounted ...

## Deploying Without Docker

##### This section is incomplete and will be completed once the system is production ready

```
qure/
 ├─ client
 |  ├─ dist
 |  |  ├─ assets
 |  |  └─ ...
 |  └─ ...
 ├─ server
 |  ├─ db
 |  |  └─ db.gob
 |  ├─ internal
 |  |  └─ ...
 |  ├─ config.json
 |  └─ server (executable)
```

1. ***[...insert initial steps...]***
1. In `qure/client` run `npm run build` to build the frontend.
1. In `qure/server` run `go build .` to build the backend.
1. Copy [config.json.example](../server/config.json.example) and rename it to `config.json`.
1. Run the backend executable `server`
1. Read the log output for any issues
1. Copy the initial admin password from the logs
1. Log in with the admin credentials and change the password.

## Configuring
##### This section is incomplete and will be completed once the system is production ready

Configuration is done with `config.json` file. A sample file `config.json.example` is provided in this repository:
> https://github.com/JValtteri/qure/blob/main/server/config.json.example

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

#### TLS

To comply with **GDPR** and other security legislation, it is practically **mandatory to use TLS** on any open web facing service that handles personally identifiable data (PID).

This server supports TLS directly, but you must supply valid `cert.pem` and `privkey.pem` for it to work. This can be done manually, and there are tools for that.

The recommended approach, however, is to use a reverse proxy, like [caddy](https://caddyserver.com/docs/quick-starts/reverse-proxy) or [nginx](https://docs.nginx.com/nginx/admin-guide/web-server/reverse-proxy/). They can handle certificates automatically for traffic going through them. If you are using their automatic certificate renewal, set `ENABLE_TLS` to `false`.

## Using
##### This section is incomplete and will be completed once the system is production ready

### Admin account

If an `admin` account doesn't exist, a new `admin` account is created automatically on server start.

**Check server log output for the `admin` credentials.**

- Username is `admin`
- Password is a random string of characters.

**The password should be changed on first login.**

## Maintaining
##### This section is incomplete and will be completed once the system is production ready

## Security

See [Security Document](./security.md)

## Licenses

See [doc/licenses](doc/licenses/README.md)

## Developer Documentation

See: [Developer Documentation](./dev.md)
