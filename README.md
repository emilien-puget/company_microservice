# fizzbuzz

Expose two ports

- 8080: for business endpoints
- 2112: for internal endpoints

## env vars

- PORT: business endpoints port, default 8080
- INTERNAL_PORT: internal endpoints port, default 2112
- BASE_URL: base url for the business endpoints
- MONGODB_CONNECTION_STRING: mongodb connection string

## Business endpoints

login is hardcoded to ``jon`` and password to ``shhh!``

an open api 3 specification is available in ``cmd/api/openapi.yaml``

## Internal endpoints

- ping: for health check probe
- metrics: a prometheus endpoint

