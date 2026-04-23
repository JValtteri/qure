# QuRe Reservation System Documentation

## Index

- [Index](#index)
- [Data Structure Graph](#data-structure-graph)
- [Security]
- [Data Persistance](#data-persistance)

## Data Structure Graph

![](../images/QuReDiagram.drawio.svg)

###### *Status on Feb 04th 2026*

## Setting up the dev-enviroment

1. Install **npm**
1. Install **Go** >=1.24
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
