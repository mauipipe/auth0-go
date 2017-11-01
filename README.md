# auth0-go

[![GoDoc](https://godoc.org/github.com/credcap/auth0-go?status.svg)](https://godoc.org/github.com/credcap/auth0-go)

Helper library for [Auth0 Management API](https://auth0.com/docs/api/management/v2)

# Basics

All requests to the Management API need to use a token, you can get one using the `GetToken` method.

```
token, err := auth0.GetToken("client_id", "client_secret", "https://mock.auth0.com/api/v2/")

fmt.Println(token.AccessToken)
// outputs eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IlFqUTVORG[...]
```

**OR**

You can create a `Client` which will keep the token fresh and automatically correctly add in the authorization header to requests.

For example, to call the [Create Users](https://auth0.com/docs/api/management/v2#!/Users/post_users) endpoint...

```
c, _ := auth0.NewClient("client_id", "client_secret", "https://mock.auth0.com/api/v2/")

var res interface{}

c.POST("users", struct {
    Connection    string
    Email         string
    PhoneNumber   string `json:"phone_number"`
    EmailVerified bool   `json:"email_verified"`
    VerifyEmail   bool   `json:"verify_email"`
    PhoneVerified bool   `json:"phone_verified"`
}{
    Connection:    "email",
    Email:         "nick.bradshaw@us.navy.mil",
    PhoneNumber:   "+121255512121",
    EmailVerified: true,
    VerifyEmail:   false,
    PhoneVerified: true,
}, &res)

fmt.Println(res["user_id"])
// outputs email|59f8db607cd312629715dbb2
```

The `POST` and `GET` methods are such that you can still use this lib for methods we haven't wired up.

# Getting dependencies (if you need to)

`auth0-go` uses the excellent [dep](https://github.com/golang/dep) tool and the `vendor` directory is not included in this repository.

Install the dependencies using `dep`, e.g. `dep ensure`.

