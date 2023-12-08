---
title: "Steampipe Table: okta_app_assigned_group - Query Okta App Assigned Groups using SQL"
description: "Allows users to query App Assigned Groups in Okta, specifically providing details about which groups are assigned to which applications."
---

# Table: okta_app_assigned_group - Query Okta App Assigned Groups using SQL

Okta App Assigned Groups are a part of Okta's Universal Directory, which is a flexible, cloud-based user store. It allows you to manage users and their group memberships in your applications. This includes the ability to assign users to groups that are then assigned to applications.

## Table Usage Guide

The `okta_app_assigned_group` table provides insights into the App Assigned Groups within Okta's Universal Directory. As a system administrator, you can explore group-specific details through this table, including which users are part of which groups and which applications these groups have access to. This can be particularly useful for managing and auditing access controls within your applications.

## Examples

### Basic info
Explore which applications are assigned to different groups in Okta, along with their last updated status and priority. This can be useful in managing application access and prioritizing updates.

```sql+postgres
select
  id as group_id,
  app_id,
  last_updated,
  priority
from
  okta_app_assigned_group;
```

```sql+sqlite
select
  id as group_id,
  app_id,
  last_updated,
  priority
from
  okta_app_assigned_group;
```

### List groups that are not assigned to any application
Determine the groups that are not associated with any application to assess potential inefficiencies or unnecessary resources. This can aid in resource management and ensure optimal application performance.

```sql+postgres
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

```sql+sqlite
select
  grp.name as name,
  grp.description as description,
  grp.group_members as group_members
from
  okta_group grp
left join okta_app_assigned_group ag on grp.id = ag.id
where
  grp.id is null or ag.id is null
union
select
  grp.name as name,
  grp.description as description,
  grp.group_members as group_members
from
  okta_group grp
right join okta_app_assigned_group ag on grp.id = ag.id
where
  grp.id is null or ag.id is null;
```

### List applications with assigned group details
Identify applications and their associated group details to understand their status and configuration. This is useful for managing application access and ensuring appropriate group assignments.

```sql+postgres
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

```sql+sqlite
select
  app.name as app_name,
  app.id as app_id,
  app.created as app_created,
  app.status as app_status,
  ag.id as group_id,
  grp.name as group_name,
  grp.description as group_description,
  grp.group_members as group_members
from
  okta_application app
join okta_app_assigned_group ag on app.id = ag.app_id
join okta_group grp on ag.id = grp.id;
```