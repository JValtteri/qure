# QuRe Reservation System Documentation

## Index

[Return to Root](../README.md)
- [Index](#index)
- [Deploying](#deploying)
- [Configuring](#configuring)
- [Using](#using)
- [Licenses](#licenses)
- [Developer Documentation](#developer-documentation)

## Deploying
##### This section is incomplete and will be completed once the system is production ready

`insert docker compile command`

`insert docker folder structure for config file and database`

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
