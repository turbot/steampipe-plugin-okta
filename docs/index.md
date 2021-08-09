---
organization: Turbot
category: ["saas"]
icon_url: "/images/plugins/turbot/okta.svg"
brand_color: "#0089D6"
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

| Item        | Description                                                                                                                                   |
| ----------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| Credentials | Okta requires an [API token](https://developer.okta.com/docs/guides/create-an-api-token/create-the-token/) and domain for all requests.       |
| Permissions | API tokens have the same permissions as the user who creates them, and if the user permissions change, the API token permissions also change. |
| Radius      | Each connection represents a single Okta domain.                                                                                              |

### Configuration

Installing the latest okta plugin will create a config file (~/.steampipe/config/okta.spc) with a single connection named okta:

```hcl
connection "okta" {
  plugin  = "okta"
  # domain = "https://<your_okta_domain>.okta.com"
  # token  = "kvcThtthis_is_not_real_tokenzTbZO"
}

```

By default, all options are commented out in the default connection, thus Steampipe will resolve your credentials using the same order as mentioned in [Credentials](#credentials). This provides a quick way to get started with Steampipe, but you will probably want to customize your experience using configuration options for querying multiple tenants, configuring credentials from your Azure CLI, Client Certificate, etc.

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-okta
- Community: [Slack Channel](https://join.slack.com/t/steampipe/shared_invite/zt-oij778tv-lYyRTWOTMQYBVAbtPSWs3g)
