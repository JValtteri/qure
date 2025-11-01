# QuRe Registration system with Queuing

[![Frontend Tests](https://github.com/JValtteri/qure/actions/workflows/frontend-tests.yml/badge.svg)](https://github.com/JValtteri/qure/actions/workflows/frontend-tests.yml)
[![Backend Tests](https://github.com/JValtteri/qure/actions/workflows/backend-tests.yml/badge.svg)](https://github.com/JValtteri/qure/actions/workflows/backend-tests.yml)
![Frontend Coverage](https://github.com/JValtteri/qure/blob/badges/.badges/main/frontend-coverage-badge.svg)
![Backend Coverage](https://github.com/JValtteri/qure/blob/badges/.badges/main/backend-coverage-badge.svg)

## Description

Reservation System Template. A free, [open source](LICENSE) implementation of a reservation system. Easily adaptable to the needs of the user or organization.

A particular focus is in ensuring the system is provably compliant with GDPR and relevant law.

![screenshot](images/Screenshot.png)

## Components

|       | [Frontend](client/README.md) | [Backend](server/README.md) |
| ----- | :-------------------: | :----: |
| Lang  | TypeScript <br> React |   Go   |
| vers. | >=5.8.3 <br> >=19.1.1 | >=1.24 |

## Status

###### This is an ongoing project. The goal is to have an MVP implementation by 2026

### Completion estimate

|                      | Backend | Frontend | Total |
| -------------------- | :-----: | :------: | :---: |
| Reservation          |   82%   |   30%    |  60%  |
| Resuming & Modifying |   33%   |   0%     |  19%  |
| Event Creation       |   57%   |   0%     |  29%  |
| Security             |   50%   |   0%     |  50%  |
| User management      |    0%   |   0%     |   0%  |
|                      |         |          |**36%**|

##### Projected to reach 100% on Jan. 2026

## Setting up dev enviroment

READMEs of [frontend](./client) and [backend](./server) have specific instructions for how to setup each dev envirnoment. The frontend dev server and backend server are pre-configured to communicate with eachother, to allow testing of the API calls between them.

Once setup, all you need to do is run

### In `client/` folder
```
npm run dev
```

### In `server/` folder
```
go run .
```

## Deploying

```diff
- This project is not ready to deploy at this time.
```

Once production ready, the idea is to automatically package a release, likely a docker image for the entire thing.

[Documentation](./doc), containing instructions for deploying, configuring and using the system will be in [`doc/`](./doc) folder

## Copyright Notice

See [doc/licenses](doc/licenses/README.md)
