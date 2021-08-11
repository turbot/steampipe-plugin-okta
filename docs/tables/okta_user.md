# Table: okta_user

Note: This table supports optional `filter` column to query results based on okta supported [filters](https://developer.okta.com/docs/reference/api/users/#list-users-with-a-filter).

## Examples

### Basic info

```sql
select
  email,
  id,
  login,
  created,
  status,
  type
from
  okta_user;
```

### Get profile properties of users

```sql
select
  id,
  email,
  jsonb_pretty(profile) as profile
from
  okta_user;
```

### Get groups details for users

```sql
select
  id,
  email,
  jsonb_pretty(user_groups) as user_groups
from
  okta_user;
```

### Get asssigned role details for users

```sql
select
  id,
  login,
  jsonb_pretty(assigned_roles)
from
  okta_user
```

### List users with SUPER_ADMIN role access

```sql
select
  id,
  login,
  jsonb_pretty(assigned_roles) as assigned_roles
from
  okta_user
where
  assigned_roles @> '[{"type":"SUPER_ADMIN"} ]'::jsonb;
```

### Users who have not logged in for more than 30 days

```sql
select
  id,
  email,
  last_login
from
  okta_user
where
  last_login < current_timestamp - interval '30 days';
```

### List users using [filter](https://developer.okta.com/docs/reference/api/users/#list-users-with-a-filter)

```sql
select
  id,
  email,
  created,
  status
from
  okta_user
where
  filter = 'lastUpdated lt "2021-08-05T00:00:00.000Z" and status eq "ACTIVE"';
```
