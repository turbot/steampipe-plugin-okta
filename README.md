![image](https://hub.steampipe.io/images/plugins/turbot/okta-social-graphic.png)

# Okta Plugin for Steampipe

Use SQL to query infrastructure including users, groups, applications and more from Okta.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/okta)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/okta/tables)

- Community: [Slack Channel](https://steampipe.io/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-okta/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install okta
```

Run a query:

```sql
select login, id, email, created from okta_user;
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-okta.git
cd steampipe-plugin-okta
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/okta.spc
```

Try it!

```
steampipe query
> .inspect okta
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-azuread/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Okta Plugin](https://github.com/turbot/steampipe-plugin-okta/labels/help%20wanted)
