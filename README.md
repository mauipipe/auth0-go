# auth0-go
Helper library for Auth0 API

# Basics

All requests to the Management API need to use a token, you can get one using the `GetToken` method.

```
    token, err := auth0.GetToken("client_id", "client_secret", "https://mock.auth0.com/api/v2/")

    fmt.Println(token.AccessToken)
    // outputs eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IlFqUTVORG[...]
```

# Getting dependencies (if you need to)

`auth0-go` uses the excellent [dep](https://github.com/golang/dep) tool and the `vendor` directory is not included in this repository.

Install the dependencies using `dep`, e.g. `dep ensure`.

