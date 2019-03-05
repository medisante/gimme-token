# Gimme Token

An example application that shows how to authenticate with the [Medisanté Eliot API](https://api-docs.medisante.net).

The first step is to obtain an OAuth 2 token. This repository shows how to do that. Bear in mind that the token is only valid for `3600` seconds. After that, a new token has to be obtained.

As soon as a valid token has been obtained, regular HTTP requests can be made against the API. The token needs to be included in the `Authorization` header in the following way:

```
Authorization:Bearer $YOUR_TOKEN
```

The Medisanté Eliot API uses [AWS Cognito](https://aws.amazon.com/cognito) as its authentication and authorization service. We strongly recommend to use the AWS SDK for your respective programming language or platform to interact with it and to get your token. After that, you can use any HTTP client to interact with the Medisanté Eliot API.

## Usage

You need to have [Go](https://golang.org) installed in order to build this tool.

1. Run `make` to build
1. Run `./bin/gimme-token -username your@email.com` to get your token
1. Use [jwt.io](https://jwt.io) to inspect your token or use the token to perform API requests. Example:

   ```
   $ curl https://api.medisante.net/v1/orgs --header "Authorization:Bearer $YOUR_TOKEN" | jq
   ```
