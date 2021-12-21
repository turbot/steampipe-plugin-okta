# Table: okta_application_group

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
  okta_application_group;
```

### List groups that are not assigned to any application

```sql
select
  g.name as name,
  g.description as description,
  jsonb_pretty(g.group_members) as group_members
from
  okta_group g
full outer join okta_application_group ag on g.id = ag.id
where
  g.id is null or ag.id is null;
```

### List applications with assigned group details

```sql
select
  a.name as app_name,
  a.id as app_id,
  a.created as app_created,
  a.status as app_status,
  ag.id as group_id,
  g.name as group_name,
  g.description as group_description,
  jsonb_pretty(g.group_members) as group_members
from
  okta_application a
left join okta_application_group ag on a.id = ag.app_id
left join okta_group g on ag.id = g.id;
```
