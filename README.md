# alpaca-lambda

# Instructions

Prereqs:

- Make sure to have go installed
- Sign up for Alpaca Account and get your API keys

1. Pull down the repo
2. Update ENV variables
3. Replace the name of example.env to .env
4. Download all dependencies using:

```
go get .
```

5. Execute code using:

```
go run local_testing.go
docker build -t alpaca-testings .
docker run --env-file ./.env -p 9000:8080 alpaca-testings
```
