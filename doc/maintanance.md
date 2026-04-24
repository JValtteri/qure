## Server Maintanance Manual

[Return Documentation Index](./README.md)

---

| Task                  | Command               |
| :--                   | :--                   |
| start the server      | `docker compose up`   |
| stop the server       | `docker compose down` |
| view server logs<br>(including access logs) | `docker logs qure` |
| update and restart    | `./update.sh`         |
| pull latest container | `docker compose pull` |

**Note:** As the server is stopped, the container is removed and **logs are wiped.**
