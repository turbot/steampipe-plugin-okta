# Table: okta_app_user

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
  okta_app_user;
```

### List users that are not assigned to any application

```sql
select
  usr.id as id,
  usr.login as login,
  usr.created as created,
  usr.status as status
from
  okta_user usr
full outer join okta_app_user au on usr.id = au.id
where
  usr.id is null or au.id is null;
```

### List applications with assigned user details

```sql
select
  app.name as app_name,
  app.id as app_id,
  app.label as app_label,
  app.created as app_created,
  app.status as app_status,
  au.id as user_id,
  usr.login as user_login,
  usr.created as user_created,
  usr.status as user_status
from
  okta_application app
left join okta_app_user au on app.id = au.app_id
left join okta_user usr on au.id = usr.id;
```
