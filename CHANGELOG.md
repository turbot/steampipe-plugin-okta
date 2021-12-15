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
