# Configuration Manual

[Return Documentation Index](./README.md)

---

## Logo

Adding `logo.png` is optional. Adding one will replace the default logo accross the app.

## Adding local assets

If you want to add images or other assets, you can place them in `./images/`. They are mounted to a folder of the same name inside the docker container. You can embed to them in your markdown as follows.
```md
![](./images/sample-image.png)
```

## Configuring

Configuration is done with `config.json` file. A valid `config.json` is required to start the server. Config with default values is automatically created when you run the setup script.

Technically the minimum valid config is`{}`.

### Example config.json

```json
{
    "ORIGIN_URL": "localhost",
    "SERVER_PORT": "8000",
    "ENABLE_TLS": false,
    "CERT_FILE": "cert.pem",
    "PRIVATE_KEY_FILE": "privkey.pem"
}
```

If a value is not present in config.json the server uses default value. Default values should be sufficient for most cases.

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

#### Advanced Server Settings

You should not need to modify these settings.

| Key | default | description |
| :--: | :--: | -- |
| `MAX_THREADS` | `0` | Maximum allowed concurrent threads. Zero = automatic |
| `RATE_LIMIT_PER_MINUTE` | `60` | Maximum allowed requests/minute per client  |
| `RATE_LIMIT_PER_MINUTE_EVENT` | `120` | Maximum allowed requests/minute per client (event info requests) |
| `RATE_LIMIT_BURST` | `5` |  Maximum allowed burst (on top of base limit) |
| `RATE_LIMIT_RESET_MINUTES` | `60` |  Interval to clear reset limiters (to purge old clients and reset counters) |
| `RATE_LIMIT_ALERT` | `300`  |  Exceeding this number of blocked requests triggers an alert in log with offending IP address and blocked request count at last limit reset |
| `REQUEST_SIZE_LIMIT` | `30720` | Maximum allowed request size in bytes. This is to protect from excessively large requests |

There are even more advanced settings that have not been documented here. You can find them in the source [config.go](./internal/config/config.go). The reason they are not documented is you really shouldn't modify them. Changeing them may be detrimental to security or break compatability with previously saved data.

##### Times are expressed in seconds unless stated otherwise

---

## Detailed Descriptions

#### TLS

To comply with **GDPR** and other security legislation, it is practically **mandatory to use TLS** on any open web facing service that handles personally identifiable data (PID).

This server supports TLS directly, but you must supply valid `cert.pem` and `privkey.pem` for it to work. This can be done manually, and there are tools for that.

The recommended approach, however, is to use a reverse proxy, like [caddy](https://caddyserver.com/docs/quick-starts/reverse-proxy) or [nginx](https://docs.nginx.com/nginx/admin-guide/web-server/reverse-proxy/). They can handle certificates automatically for traffic going through them. If you are using their automatic certificate renewal, set `ENABLE_TLS` to `false`.

#### Account Policy

`MIN_USERNAME_LENGTH` allows you to set the minimum allowed username length

`MIN_PASSWORD_LENGTH` allows you to set the minimum allowed password length. This does not affect existing passwords.

`MAX_SESSION_AGE` how long can a user stay logged in before having to login again.

#### Rate Limiter

If the server can't keep up with the incoming requests, it's memory useage will spike rapidly, resulting in eventual crash, unless requests stop. The memory useage is very slow to recover, leaving the server extra vulnerable if a heavy load resumes. For this reason it's important for the rate limiter to be configured correctly.

To protect the server from DoS attacks - intentional or otherwise - the server backend implements an IP based rate limiter. The defaults should be fine for most hardware and use cases, but you should monitor the RAM useage under high load conditions to assess if adjustments are necessary.

`RATE_LIMIT_BURST` Sets the maximum allowed initial burst of requests. This needs to be large enough to cover the simultanous requests the frontend may send. The default value should be fine.

`RATE_LIMIT_PER_MINUTE_EVENT` Maximum allowed event requests/minute per client, afte the burst quota is consumed. Exceeding this limit will block the IP, until the request rate is below the limit.

Event requests are sent, when ever a user clicks on an event in the list. This is the most common type of event the user will create and is only limited by how fast the user can click.

`RATE_LIMIT_PER_MINUTE` Maximum allowed requests/minute per client. This applies to all other API requests. These requests are much more rare and it is difficult for a legitimate user to achieve high rates, so this limit can be set much lower than `RATE_LIMIT_PER_MINUTE_EVENT`

`RATE_LIMIT_RESET_MINUTES` Rate limiter memory is cleared after this many minutes. This purges old clients and resets counters. Any IP won't be blocked longer than this.

`RATE_LIMIT_ALERT` Exceeding this number of blocked requests triggers an alert in log with offending IP address, along with the blocked request count at last limit reset. This should be high enough that a legitimate user will never trigger this, but low enough that abuse is logged. The limit is in relation to `RATE_LIMIT_RESET_MINUTES`, so they should be changed together.
