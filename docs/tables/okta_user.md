# Table: okta_user

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

### Users who have not logged in for more than 30 days

```sql
select
  id,
  email,
  last_login
from
  okta_user
where
  last_login < current_timestamp - interval '1 days';
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
