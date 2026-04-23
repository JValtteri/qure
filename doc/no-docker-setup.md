## Deploying Without Docker

**This is not recommended. Using a Docker image enhances the security by isolating the server in a bare virtual enviroment. Instructions for setting up with docker are [here](./README.md#deploying-with-docker)**

**For developers, the recommended approach is to use the included Vite project. The complete setup is explained [here](./dev.md#setting-up-the-dev-enviroment)**


### I know what I'm doing

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
