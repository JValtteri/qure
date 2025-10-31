# Frontend
[![Frontend Tests](https://github.com/JValtteri/qure/actions/workflows/frontend-tests.yml/badge.svg)](https://github.com/JValtteri/qure/actions/workflows/frontend-tests.yml)
![Frontend Coverage](https://github.com/JValtteri/qure/blob/badges/.badges/main/frontend-coverage-badge.svg)


## Description

The frontend uses `TypeScript` and `React`. Package management is handled with `npm`.

## Installing dependencies

```
npm install
```


## Building the frontend

```
npm run build
```

---

## Dev commands

### Testing

#### Run unit tests
```
npm run test
```

#### Run coverage tests
```
npm run cover
```

### Running a dev server

#### Without networking
```
npm run dev
```
#### With networking
```
npm run host
```

Access server at `http://localhost:5173/`

The dev server has been configured in `vite.config.ts` to proxy any API calls to the real backend server. To run the backend server, follow the instructions in `qure/server/README.md`
