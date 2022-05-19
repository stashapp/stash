#  Integration tests

These tests run a gqlgen server against the apollo client to test real world connectivity.

First start the go server
```bash
go run integration/server/server.go
```

And in another terminal:
```bash
cd integration
npm ci
npm run test
```
