# Gimme Token

An example script that shows how to authenticate with the [Medisanté Eliot API](https://api-docs.medisante.net).

The first step is to obtain an OAuth 2 token. This repository shows how to do that. Bear in mind that the token is only valid for `3600` seconds (1h). After that, a new token has to be obtained.

**Important**: Because security is of utmost importance at Medisanté, we only allow logging in via [SRP](https://en.wikipedia.org/wiki/Secure_Remote_Password_protocol). Direct logins with username/password are not possible because they are much less secure.

The following values are static and you will need them in your implementation when using the AWS SDK. They are also referenced in the example implementation in this repository.

- Cognito User Pool ID: `eu-central-1_P8l0OEy9K`
- OAuth 2 Client ID: `7g6sangq93eruio247tlsp883n`

As soon as a valid token has been obtained, regular HTTP requests can be made against the API. The token needs to be included in the `Authorization` header in the following way:

```
Authorization:Bearer $YOUR_TOKEN
```

The Medisanté Eliot API uses [AWS Cognito](https://aws.amazon.com/cognito) as its authentication and authorization service. We strongly recommend to use the AWS SDK for your respective programming language or platform to interact with it and to get your token. After that, you can use any HTTP client to interact with the Medisanté Eliot API.

## Usage

You need to have [Go](https://golang.org) installed in order to build this tool.

1. Run `make` to build
1. Run `./bin/gimme-token -username your@email.com` to start
1. Enter your password when prompted
1. Use [jwt.io](https://jwt.io) to inspect your token or use the token to perform API requests. Example:

   ```
   $ curl https://api.medisante.net/v1/orgs --header "Authorization:Bearer $YOUR_TOKEN" | jq
   ```
