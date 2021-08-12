---
organization: Turbot
category: ["saas"]
icon_url: "/images/plugins/turbot/okta.svg"
brand_color: "#00297A"
display_name: "Okta"
name: "okta"
description: "Steampipe plugin for querying resource users, groups, applications and more from Okta."
og_description: "Query Okta with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/okta-social-graphic.png"
---

# Okta + Steampipe

[Okta](https://www.okta.com/) is the leading independent identity provider. The Okta Identity enables organizations to securely connect the right people to the right technologies at the right time.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

For example:

```sql
select
  login,
  id,
  email,
  created
from
  okta_user;
```

```
+---------------------+----------------------+---------------------+---------------------+
| login               | id                   | email               | created             |
+---------------------+----------------------+---------------------+---------------------+
| subhajit@turbot.com | 00u1e63jiqAHskqSd5d7 | subhajit@turbot.com | 2021-08-02 13:35:54 |
| lalit@turbot.com    | 00u1e5eizrjQKTWMA5d7 | lalit@turbot.com    | 2021-08-02 10:57:05 |
+---------------------+----------------------+---------------------+---------------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/okta/tables)**

## Get started

### Install

Download and install the latest Okta plugin:

```bash
steampipe plugin install okta
```

### Credentials

| Item        | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| ----------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Credentials | Okta requires an [API token](https://developer.okta.com/docs/guides/create-an-api-token/create-the-token/) and domain for all requests.                                                                                                                                                                                                                                                                                                                                                                                      |
| Permissions | API tokens have the same permissions as the user who creates them, and if the user permissions change, the API token permissions also change.                                                                                                                                                                                                                                                                                                                                                                                |
| Radius      | Each connection represents a single Okta Organization.                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| Resolution  | 1. With configuration provided in connection in steampipe _**.spc**_ config file.<br /> a) [API Token](https://developer.okta.com/docs/guides/create-an-api-token/create-the-token/)<br /> b) [Private Key](https://github.com/okta/okta-sdk-golang#oauth-20)<br />2. With okta environment variables.<br />3. An okta.yaml file in a .okta folder in the current user's home directory _**(~/.okta/okta.yaml or %userprofile\.okta\okta.yaml)**_.<br />4. A .okta.yaml file in the application or project's root directory. |

### Configuration

Installing the latest okta plugin will create a config file (~/.steampipe/config/okta.spc) with a single connection named okta:

```hcl
connection "okta" {
  plugin = "okta"

  # 1. With API TOKEN(https://developer.okta.com/docs/guides/create-an-api-token/create-the-token/)
  # domain = "https://<your_okta_domain>.okta.com"
  # token  = "this_not_real_token"


  # 2. With Private Key(https://github.com/okta/okta-sdk-golang#oauth-20)
  # domain      = "https://<your_okta_domain>.okta.com"
  # client_id   = "your_client_id"
  # private_key = "----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAmyX8wdrHK1ycOMeXNg3NOMQvebnfQp+3L5OaaiX16/+tLbwb\nJTZDYh0EXLySMVsduRxC/1PQdPuI6x50TdkoB3C4JMuU968uJqkFp7fXXy5SMAej\nHAyF67cY51dx15ztvakRNJPhhI5WaC20RfR/eow0IH5lGI3czcvTCChGau5qLue3\nHqNDYFY+U3xhOlavSDdtmuxpIFsDycn/OjYjsV4lzyRrOArqtVV/kXHKx04T6A1x\nSc99999999999999999999999999999999999999999999999999EGekHlUAIUpw\n-----END RSA PRIVATE KEY-----"
}
```

By default, all options are commented out in the default connection, thus Steampipe will resolve your credentials using the same order as mentioned in [Credentials](#credentials). This provides a quick way to get started with Steampipe, but you will probably want to customize your experience using configuration options for querying multiple organizations, configuring credentials from your okta configuration files, [environment variables](#credentials-from-environment-variables), etc.

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-okta
- Community: [Slack Channel](https://join.slack.com/t/steampipe/shared_invite/zt-oij778tv-lYyRTWOTMQYBVAbtPSWs3g)

## Configuring Okta Credentials

### Credentials from Environment Variables

The Okta plugin will use the standard Okta environment variables to obtain credentials **only if other arguments (`domain`, `token`, `client_id`, `private_key`) are not specified** in the connection:

#### API Token

```sh
export OKTA_CLIENT_ORGURL=https://<your_okta_domain>.okta.com
export OKTA_CLIENT_TOKEN=this_not_real_token
```

#### Private Key

```sh
export OKTA_CLIENT_ORGURL=https://<your_okta_domain>.okta.com
export OKTA_CLIENT_CLIENTID=your_client_id
export OKTA_CLIENT_PRIVATEKEY="----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAmyX8wdrHK1ycOMeXNg3NOMQvebnfQp+3L5OaaiX16/+tLbwb\nJTZDYh0EXLySMVsduRxC/1PQdPuI6x50TdkoB3C4JMuU968uJqkFp7fXXy5SMAej\nHAyF67cY51dx15ztvakRNJPhhI5WaC20RfR/eow0IH5lGI3czcvTCChGau5qLue3\nHqNDYFY+U3xhOlavSDdtmuxpIFsDycn/OjYjsV4lzyRrOArqtVV/kXHKx04T6A1x\nSc99999999999999999999999999999999999999999999999999EGekHlUAIUpw\n-----END RSA PRIVATE KEY-----"
```
