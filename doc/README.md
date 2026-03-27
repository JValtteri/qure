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
 ├─ db         (auto-generated)
 │  └─ db.gob  (auto-generated)
 ├─ logo.png   (optional)
 └─ images     (optional)
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

### Customization

#### Logo

Adding `logo.png` is optional. Adding one will replace the default logo accross the app.

#### Adding local assets

If you want to add images or other assets, you can place them in `./images/`. They are mounted to a folder of the same name inside the docker container. You can embed to them in your markdown as follows.
```md
![](./images/sample-image.png)
```

## Deploying Without Docker

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
1. In `qure/client` run `npx vite build` to build the frontend.
1. In `qure/server` run `go build .` to build the backend.
1. Copy [config.json.example](../server/config.json.example) and rename it to `config.json`.
1. Run the backend executable `server`
1. Read the log output for any issues
1. Copy the initial admin password from the logs
1. Log in with the admin credentials and change the password.

## Configuring

Configuration is done with `config.json` file. A sample file `config.json.example` is provided in this repository:

> https://github.com/JValtteri/qure/blob/main/server/config.json.example

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

##### Advanced Server Settings

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

###### Times are expressed in seconds unless stated otherwise

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

## Using

### Admin account

If an `admin` account doesn't exist, a new `admin` account is created automatically on server start.

**Check server log output for the `admin` credentials.**

- Username is `admin`
- Password is a random string of characters.

**The password should be changed on first login.**

> It is good policy to not share the `admin` account, but create individual accounts and promote them with the necessary roles. This way any policy violations can be tracked and offending accounts can be demoted or removed without affecting other administrators.

### Admin tools

You can access account settings, reservation and admin tools by clicking your user name at the top right corner of your screen. The last tab is **Admin Tools**.

#### Reservations
The first tool is Reservations. You can use it to inspect the reservations of a given event. This is useful for verifying customers' reservations.

Click **Reservations** then click on the event to inspect. The list will populate with reservations sorted by their time.

You can search for a specific ID by typeing it in the search field.

**This action is GDPR safe. No PII is shown.**

#### All Users
The second tab is **All Users**. It is only available to full **admin users**. **Opening the tab counts as PII access**. To comply with GDPR, there must be a good reason to access the list. **This action is logged**.

From the list, you can select a user. A detail card of the user is shown. From the card you can **delete the user** or **change its role**.

#### Deleting a User
To delete a user, you must be a full admin user. Open **Settings -> Admin tools -> All Users** and select the user you want to delete. Click **Delete**. You are asked to enter your admin password to confirm the deletion. **This action is logged.** When a user is deleted by the user or an admin, all data related to the user is removed from the system.

#### Changeing a User's Role
To change the role of a user, you must be a full admin user. Open **Settings -> Admin tools -> All Users** and select the user you want to modify. Click the ✏️ icon next to the role. Select the desired role from the drop down menu. You are asked to enter your admin password to confirm. **This action is logged.**

### Creating an Event
To create an event, you need to be logged in as an administrator. Click the card with the ➕ symbol on it to open a new event editor. All fields (except *short description*) must be filled to save the event. Event can be either **Published** or **Saved as a Draft**. A draft is visible only to administrators, while Published events are visible to anyone.

To add groups/timeslots to the event, set the group size. You can add more groups by clicking the plus (+) symbol next to the group. You can remove a group by pressing the minus (-) symbol next to it.

### Editing and Event
To create an event, you need to be logged in as an administrator. Select the event you want to edit and click Edit Event. Edit any fields you need to and either **Publish** or **Save as a Draft**. You can hide a published event by saving it as a draft, and you can publish a draft event by publishing it.

## Maintanance

#### Using docker version, you can start the server with
```sh
docker compose up
```

#### and stop the server with
```sh
docker compose down
```
As the server is stopped, the container is removed and **logs are wiped.**

#### You can view server logs (including access logs) with
```sh
docker logs [container_name]
```

#### To update the server
You should run the following commands:
```sh
docker logs [container_name]    # to view the logs before they are wiped
docker compose down             # stop the server
docker compose pull             # download the update
docker compose up               # start the server
```
It is recommended to make a script of these commands.

## Security

See [Security Document](./security.md)

## Licenses

See [doc/licenses](doc/licenses/README.md)

## Developer Documentation

See: [Developer Documentation](./dev.md)
