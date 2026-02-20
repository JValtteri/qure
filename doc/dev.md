# QuRe Reservation System Documentation

## Index

- [Index](#index)
- [Data Structure Graph](#data-structure-graph)
- [Security]
- [Data Persistance](#data-persistance)

## Data Structure Graph

![](../images/QuReDiagram.drawio.svg)

###### *Status on Feb 04th 2026*

## Security

See [Security Document](./security.md)

## Data Persistance

As of writing, presistance is hanled by `internal/saveload`. Saveload uses flat GOB-files. Implementation is used by used by `internal/state/persistance_api.go`. It should be easy to swap the implementation to a more serious databace, if need be.

Currently the entire state is held in memory and is only saved in file for precistance before shutdowns.
