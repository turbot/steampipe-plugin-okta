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
