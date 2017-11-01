# auth0-go
Helper library for Auth0 API

# basics

All requests to the Management API need to use a token, you can get one using the `GetToken` method.

```
    token, err := auth0.GetToken("client_id", "client_secret", "https://mock.auth0.com/api/v2/")

    fmt.Println(token.AccessToken)
    // outputs eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IlFqUTVORG[...]
```

