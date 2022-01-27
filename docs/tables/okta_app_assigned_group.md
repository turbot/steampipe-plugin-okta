# Table: okta_app_assigned_group

Application integrations can be assigned to groups. When app integrations share the same group they are "linked". This can be helpful if you need to add provisioning functionality in an SSO-enabled app integration.

## Examples

### Basic info

```sql
select
  id as group_id,
  app_id,
  last_updated,
  priority
from
  okta_app_assigned_group;
```

### List groups that are not assigned to any application

```sql
select
  grp.name as name,
  grp.description as description,
  jsonb_pretty(grp.group_members) as group_members
from
  okta_group grp
full outer join okta_app_assigned_group ag on grp.id = ag.id
where
  grp.id is null or ag.id is null;
```

### List applications with assigned group details

```sql
select
  app.name as app_name,
  app.id as app_id,
  app.created as app_created,
  app.status as app_status,
  ag.id as group_id,
  grp.name as group_name,
  grp.description as group_description,
  jsonb_pretty(grp.group_members) as group_members
from 
  okta_application app 
inner join okta_app_assigned_group ag on app.id = ag.app_id
inner join okta_group grp on ag.id = grp.id;
```
