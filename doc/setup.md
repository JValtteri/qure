## Deploying With Docker

[Return Documentation Index](./README.md)

---

### Setup

- **Prerequisites:**
    - docker
    - curl

1. run `curl -fsSL "https://raw.githubusercontent.com/JValtteri/qure/refs/heads/main/setup.sh" | sh`
1. configure `config.json` as necessary
    - mount settings
    - port number
1. Add your custom `logo.png` *(optional)*
1. Run `docker compose up`
1. Read the log output for any issues
1. Copy the initial admin password from the logs
1. Log in with the admin credentials and change the password.
1. Press `D` to detach from server console


#### Folder structure:
```
qure
 ├─ docker-compose.yml
 ├─ config.json
 ├─ db
 │  └─ db.gob
 ├─ logo.png   (optional)
 └─ images     (optional)
    └─ ...
```

### Configuration

**See [Configuration Manual](./configuration.md)**


**For server commands, see [Server Maintanance Manual](./maintanance.md)**
