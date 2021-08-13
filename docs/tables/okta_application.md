# Table: okta_application

Applications are web-based services that provide any number of specific tasks that require user authentication.

Note: This table supports an optional `filter` column to query results based on Okta supported [filters](https://developer.okta.com/docs/reference/api/apps/#filters).

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

### List SAML 2.0 apps

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
  sign_on_mode = 'SAML_2_0';
```

### List apps assigned to a specific user using a filter

```sql
select
  id,
  label,
  name,
  sign_on_mode,
  status
from
  okta_application as app
where
  filter = 'user.id eq "00u1e5eizrjQKTWMA5d7"';
```

### List apps assigned to a specific group using a filter

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
