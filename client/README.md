# Frontend
[![Frontend Tests](https://github.com/JValtteri/qure/actions/workflows/frontend-tests.yml/badge.svg)](https://github.com/JValtteri/qure/actions/workflows/frontend-tests.yml)
![Frontend Coverage](https://github.com/JValtteri/qure/blob/badges/.badges/main/frontend-coverage-badge.svg)

## Index

- [Description](#description)
- [Installing dependencies](#installing-dependencies)
- [Testing](#testing)
    - [Run unit tests](#run-unit-tests)
    - [Run coverage tests](#run-coverage-tests)
    - [Run linter](#run-linter)
- [Running a dev server](#running-a-dev-server)
    - [Without networking](#without-networking)
    - [With networking](#with-networking)
- [Building the frontend](#building-the-frontend)

#### See also:

- [QuRe](../README.md)
    - **Frontend documentation**
    - [Backend documentation](../server/README.md)

## Description

The frontend uses
- `TypeScript`
- `React`

As a prerequasite you need to install
- `node` for npm
- `npm` for package management

## Installing dependencies
Once `npm` is installed, you can install everything else with
```
npm install
```

---

## Testing

#### Run unit tests
```
npm run test
```

#### Run coverage tests
```
npm run cover
```

#### Run linter
```
npm run lint
```

## Running a dev server

#### Without networking
```
npm run dev
```
#### With networking
```
npm run host
```

Access server at `http://localhost:5173/`

The dev server has been configured in `vite.config.ts` to proxy any API calls to the real backend server. To run the backend server, follow the instructions in [server/README.md](../server)

## Building the frontend

```
npm run build
```
