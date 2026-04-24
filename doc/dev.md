# QuRe Reservation System Documentation

[Return Documentation Index](./README.md)

---

## Index

- [Index](#index)
- [Dependencies](#dependencies)
- [Setting Up a Dev Enviroment](#setting-up-dev-enviroment)
- [Security](#security)
- [Data Persistance](#data-persistance)
- [Architecture](#architecture)
- [Open Issues](#open-issues)

## Dependencies

|       | [Frontend](./client/README.md) | [Backend](./server/README.md) |
| ----- | :-------------------: | :----: |
| Lang  | TypeScript <br> React |   Go   |
| vers. | >=5.8.3 <br> >=19.1.1 | >=1.24 |
| Tools | npm, Vite             |        |

## Setting up dev enviroment

READMEs of [frontend](./client) and [backend](./server) have specific instructions for how to setup each dev envirnoment. The frontend dev server and backend server are preconfigured to communicate with eachother, to allow testing of the API calls between them.

1. In client/ folder run:
    1. `npm install`
    1. `npm run dev`
1. In server/ folder run:
    1. `go run .`

Access the app at `http://localhost:5173/`

**See also detailed documentation:**
- [Frontend Documentation](../client/README.md)
- [Backend Documentation](../server/README.md)

## Security

See [Security Document](./security.md)

## Data Persistance

As of writing, presistance is hanled by `internal/saveload`. Saveload uses flat GOB-files. Implementation is used by used by `internal/state/persistance_api.go`. It should be easy to swap the implementation to a more serious databace, if need be.

Currently the entire state is held in memory and is only saved in file for precistance before shutdowns.

## Architecture

![](../images/QuReDiagram.drawio.svg)
###### *Status on Feb 04th 2026*

## Open Issues

See [issues.md](./issues.md)
