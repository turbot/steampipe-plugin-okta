---
title: "Steampipe Table: okta_app_assigned_user - Query Okta Assigned Users using SQL"
description: "Allows users to query Assigned Users in Okta, specifically the users assigned to applications, providing insights into user-application associations and potential security implications."
---

# Table: okta_app_assigned_user - Query Okta Assigned Users using SQL

Okta is an identity and access management service that provides secure access to tools and data. It allows organizations to manage their users, groups, and applications across different systems. An Okta Assigned User is a user that has been assigned to an application within Okta, allowing them access to that application.

## Table Usage Guide

The `okta_app_assigned_user` table provides insights into the users assigned to applications within Okta. As a security analyst or administrator, explore user-application associations through this table, including the user's ID, the application's ID, and the assignment's status. Utilize it to uncover information about user access rights, such as which users have access to specific applications, and the verification of user-application associations.

## Examples

### Basic info
Explore which users are assigned to specific applications in your Okta environment, with details including their ID, username, the time they were created, and their current status. This can help you manage user access and ensure appropriate permissions are maintained.

```sql+postgres
select
  id as user_id,
  app_id,
  user_name,
  created,
  status
from
  okta_app_assigned_user;
```

```sql+sqlite
select
  id as user_id,
  app_id,
  user_name,
  created,
  status
from
  okta_app_assigned_user;
```

### List users that are not assigned to any application
Explore which users are not linked to any application, useful for identifying potential unused or inactive accounts. This can aid in optimizing resource allocation and enhancing security measures.

```sql+postgres
select
  usr.id as id,
  usr.login as login,
  usr.created as created,
  usr.status as status
from
  okta_user usr
full outer join okta_app_assigned_user au on usr.id = au.id
where
  usr.id is null or au.id is null;
```

```sql+sqlite
select
  usr.id as id,
  usr.login as login,
  usr.created as created,
  usr.status as status
from
  okta_user usr
left join okta_app_assigned_user au on usr.id = au.id
where
  usr.id is null or au.id is null
union all
select
  usr.id as id,
  usr.login as login,
  usr.created as created,
  usr.status as status
from
  okta_user usr
right join okta_app_assigned_user au on usr.id = au.id
where
  usr.id is null or au.id is null;
```

### List applications with assigned user details
This query helps you identify all applications that have users assigned to them, along with the users' details. It's useful for monitoring application usage and managing user access, ensuring security and efficiency in your system.

```sql+postgres
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
inner join okta_app_assigned_user au on app.id = au.app_id
inner join okta_user usr on au.id = usr.id;
```

```sql+sqlite
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
join okta_app_assigned_user au on app.id = au.app_id
join okta_user usr on au.id = usr.id;
```