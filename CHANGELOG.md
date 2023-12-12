## v0.11.0 [2023-12-12]

_What's new?_

- The plugin can now be downloaded and used with the [Steampipe CLI](https://steampipe.io/docs), as a [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/overview), as a [SQLite extension](https://steampipe.io/docs//steampipe_sqlite/overview) and as a standalone [exporter](https://steampipe.io/docs/steampipe_export/overview). ([#105](https://github.com/turbot/steampipe-plugin-okta/pull/105))
- The table docs have been updated to provide corresponding example queries for Postgres FDW and SQLite extension. ([#105](https://github.com/turbot/steampipe-plugin-okta/pull/105))
- Docs license updated to match Steampipe [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-okta/blob/main/docs/LICENSE). ([#105](https://github.com/turbot/steampipe-plugin-okta/pull/105))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.8.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v580-2023-12-11) that includes plugin server encapsulation for in-process and GRPC usage, adding Steampipe Plugin SDK version to `_ctx` column, and fixing connection and potential divide-by-zero bugs. ([#104](https://github.com/turbot/steampipe-plugin-okta/pull/104))

## v0.10.1 [2023-10-04]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v562-2023-10-03) which prevents nil pointer reference errors for implicit hydrate configs. ([#94](https://github.com/turbot/steampipe-plugin-okta/pull/94))

## v0.10.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#91](https://github.com/turbot/steampipe-plugin-okta/pull/91))
- Recompiled plugin with Go version `1.21`. ([#91](https://github.com/turbot/steampipe-plugin-okta/pull/91))

## v0.9.0 [2023-04-15]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v530-2023-03-16) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#82](https://github.com/turbot/steampipe-plugin-okta/pull/82))

## v0.8.0 [2022-09-27]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.7](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v417-2022-09-08) which includes several caching and memory management improvements. ([#76](https://github.com/turbot/steampipe-plugin-okta/pull/76))
- Recompiled plugin with Go version `1.19`. ([#76](https://github.com/turbot/steampipe-plugin-okta/pull/76))

## v0.7.0 [2022-06-01]

_Enhancements_

- Added additional optional key quals and filter support to the following tables: ([#66](https://github.com/turbot/steampipe-plugin-okta/pull/66))
  - okta_app_assigned_group
  - okta_app_assigned_user
  - okta_application
  - okta_auth_server
  - okta_factor
  - okta_group
  - okta_idp_discovery_policy

## v0.6.1 [2022-05-23]

_Bug fixes_

- Fixed the Slack community links in README and docs/index.md files. ([#71](https://github.com/turbot/steampipe-plugin-okta/pull/71))

## v0.6.0 [2022-04-27]

_Enhancements_

- Added support for native Linux ARM and Mac M1 builds. ([#69](https://github.com/turbot/steampipe-plugin-okta/pull/69))
- Recompiled plugin with [steampipe-plugin-sdk v3.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v310--2022-03-30) and Go version `1.18`. ([#68](https://github.com/turbot/steampipe-plugin-okta/pull/68))

## v0.5.0 [2022-01-27]

_Enhancements_

- Added limit and context cancellation handling to the following tables ([#61](https://github.com/turbot/steampipe-plugin-okta/pull/61)) ([#63](https://github.com/turbot/steampipe-plugin-okta/pull/63))
  - okta_app_assigned_group
  - okta_app_assigned_user
  - okta_application
  - okta_auth_server
  - okta_factor
  - okta_group
  - okta_idp_discovery_policy
  - okta_network_zone
  - okta_password_policy
  - okta_signon_policy
  - okta_trusted_origin
  - okta_user
  - okta_user_type

_Bug fixes_

- Updated `okta_app_assigned_group` and `okta_app_assigned_user` tables to hydrate more efficiently and reduce the number of API calls the tables make for accounts with a large number of applications ([#61](https://github.com/turbot/steampipe-plugin-okta/pull/61))
- Fixed the `okta_app_assigned_group` and `okta_app_assigned_user` tables to correctly return assigned groups and users instead of empty results ([#61](https://github.com/turbot/steampipe-plugin-okta/pull/61))

## v0.4.0 [2022-01-20]

_Enhancements_

- Added column `profile` to the `okta_app_assigned_group` table ([#58](https://github.com/turbot/steampipe-plugin-okta/pull/58))

## v0.3.0 [2022-01-12]

_What's new?_

- New tables added
  - [okta_app_assigned_group](https://hub.steampipe.io/plugins/turbot/okta/tables/okta_app_assigned_group) ([#50](https://github.com/turbot/steampipe-plugin-okta/pull/50))
  - [okta_app_assigned_user](https://hub.steampipe.io/plugins/turbot/okta/tables/okta_app_assigned_user) ([#51](https://github.com/turbot/steampipe-plugin-okta/pull/51))

_Enhancements_

-  Recompiled plugin with [steampipe-plugin-sdk-v1.8.3](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v183--2021-12-23) ([#55](https://github.com/turbot/steampipe-plugin-okta/pull/55))

_Bug fixes_

- Removed columns `assigned_users` and `assigned_groups` from `okta_application` table ([#53](https://github.com/turbot/steampipe-plugin-okta/pull/53))
- Fixed the `okta_application` table to correctly return okta application details instead of throwing an error ([#46](https://github.com/turbot/steampipe-plugin-okta/pull/46))

## v0.2.0 [2021-12-15]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk-v1.8.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v182--2021-11-22) ([#44](https://github.com/turbot/steampipe-plugin-okta/pull/44))
- Recompiled plugin with Go version 1.17 ([#44](https://github.com/turbot/steampipe-plugin-okta/pull/44))
- Added `assigned_users` and `assigned_groups` columns to `okta_application` table ([#41](https://github.com/turbot/steampipe-plugin-okta/pull/41))
- Removed the `profile` column from `okta_application` table ([#42](https://github.com/turbot/steampipe-plugin-okta/pull/42))

## v0.1.0 [2021-09-02]

_What's new?_

- New tables added
  - [okta_auth_server](https://hub.steampipe.io/plugins/turbot/okta/tables/okta_auth_server) ([#28](https://github.com/turbot/steampipe-plugin-okta/pull/28))
  - [okta_factor](https://hub.steampipe.io/plugins/turbot/okta/tables/okta_factor) ([#30](https://github.com/turbot/steampipe-plugin-okta/pull/30))
  - [okta_idp_discovery_policy](https://hub.steampipe.io/plugins/turbot/okta/tables/okta_idp_discovery_policy) ([#25](https://github.com/turbot/steampipe-plugin-okta/pull/25))
  - [okta_mfa_policy](https://hub.steampipe.io/plugins/turbot/okta/tables/okta_mfa_policy) ([#26](https://github.com/turbot/steampipe-plugin-okta/pull/26))
  - [okta_network_zone](https://hub.steampipe.io/plugins/turbot/okta/tables/okta_network_zone) ([#29](https://github.com/turbot/steampipe-plugin-okta/pull/29))
  - [okta_signon_policy](https://hub.steampipe.io/plugins/turbot/okta/tables/okta_signon_policy) ([#27](https://github.com/turbot/steampipe-plugin-okta/pull/27))
  - [okta_trusted_origin](https://hub.steampipe.io/plugins/turbot/okta/tables/okta_trusted_origin) ([#31](https://github.com/turbot/steampipe-plugin-okta/pull/31))

_Enhancements_

- The documentation now includes additional information on `okta_user_type` and `okta_network_zone` tables

_Bug fixes_

- The `title` column for `okta_user` table will no longer have `nil` values ([#35](https://github.com/turbot/steampipe-plugin-okta/pull/35))

## v0.0.2 [2021-08-25]

_What's new?_

- Add support for service application and private key authentication ([#22](https://github.com/turbot/steampipe-plugin-okta/pull/22))
- Update the config and docs/index.md with more information about service app credentials

## v0.0.1 [2021-08-13]

_What's new?_

- New tables added

  - [okta_application](https://hub.steampipe.io/plugins/turbot/okta/tables/okta_application)
  - [okta_group](https://hub.steampipe.io/plugins/turbot/okta/tables/okta_group)
  - [okta_password_policy](https://hub.steampipe.io/plugins/turbot/okta/tables/okta_password_policy)
  - [okta_user](https://hub.steampipe.io/plugins/turbot/okta/tables/okta_user)
  - [okta_user_type](https://hub.steampipe.io/plugins/turbot/okta/tables/okta_user_type)
