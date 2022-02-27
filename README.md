# Open Banking Demo API

Gather all your greek banking accounts into one app

## Setup

Run the following command to run the server

```bash
go run server.go
```

Config schema

```json
{
  "app_url": "string",
  "db_uri": "mongodb://username:password@localhost:27017/",
  "db_name": "string",
  "redis": {
    "address": "localhost:6379",
    "password": "super_secret_password"
  },
  "jwt_secret": "very strong password",
  "providers": {
    "alpha": {
      "base_url": "https://gw.api.alphabank.eu/sandbox",
      "base_api_url": "https://gw.api.alphabank.eu/api/sandbox",
      "client_id": "client_id",
      "client_secret": "client_secret",
      "subscription_key": "GR - Accounts subscription key"
    }
  }
}
```
