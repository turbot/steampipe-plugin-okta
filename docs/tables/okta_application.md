# Table: okta_application

Application are web-based services that provide any number of specific tasks that require user authentication.

Note: This table supports optional `filter` column to query results based on okta supported [filters](<(https://developer.okta.com/docs/reference/api/apps/#filters)>).

## Examples

### Basic info

```sql
select
  name,
  id,
  label,
  created,
  status,
  sign_on_mode
from
  okta_application;
```

### Get apps with `SAML_2_0` sign_on_mode

```sql
select
  name,
  id,
  label,
  created,
  status,
  sign_on_mode
from
  okta_application
where
  sign_on_mode = 'SAML_2_0'
```

### List aps assigned to a specific user using [filter](https://developer.okta.com/docs/reference/api/apps/#list-applications-assigned-to-a-user)

```sql
select
  id,
  label,
  name,
  sign_on_mode,
  status
from
  okta_application as app,
where
  filter = 'user.id eq "00u1e5eizrjQKTWMA5d7"';
```

### List apps assigned to a specific group using [filter](https://developer.okta.com/docs/reference/api/apps/#list-applications-assigned-to-a-group)

```sql
select
  id,
  label,
  name,
  sign_on_mode,
  status
from
  okta_application
where
  filter = 'group.id eq "00u1e5eizrjQKTWMA5d7"';
```
