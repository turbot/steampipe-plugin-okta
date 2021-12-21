# Table: okta_application_user

Application integrations can be assigned to individual users. This enables assigned users to access them.

## Examples

### Basic info

```sql
select
  id as user_id,
  app_id,
  user_name,
  created,
  status
from
  okta_application_user;
```

### List users that are not assigned to any application

```sql
select
  u.id as id,
  u.login as login,
  u.created as created,
  u.status as status
from
  okta_user u
full outer join okta_application_user au on u.id = au.id
where
  u.id is null or au.id is null;
```

### List applications with assigned user details

```sql
select
  a.name as app_name,
  a.id as app_id,
  a.label as app_label,
  a.created as app_created,
  a.status as app_status,
  au.id as user_id,
  u.login as user_login,
  u.created as user_created,
  u.status as user_status
from
  okta_application a
left join okta_application_user au on a.id = au.app_id
left join okta_user u on au.id = u.id;
```
