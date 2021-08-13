# Table: okta_user

A user can be granted access to applications, devices, and groups.

Note: This table supports an optional `filter` column to query results based on Okta supported [filters](https://developer.okta.com/docs/reference/api/users/#list-users-with-a-filter).

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

### Get profile, group, and assigned role details for each user

```sql
select
  id,
  email,
  jsonb_pretty(profile) as profile,
  jsonb_pretty(user_groups) as user_groups,
  jsonb_pretty(assigned_roles) as assigned_roles
from
  okta_user;
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

### List users who have not logged in for more than 30 days

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

### List active users that have been last updated before a specific date using a filter

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
