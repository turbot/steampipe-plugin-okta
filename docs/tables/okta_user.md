---
title: "Steampipe Table: okta_user - Query Okta Users using SQL"
description: "Allows users to query Okta Users in, specifically user profiles and statuses, providing insights into user management and access control."
---

# Table: okta_user - Query Okta Users using SQL

Okta User is a resource within that represents an authenticated entity in the Okta service. A user can be an end user (person) or a service user (software). Each user has a profile that stores the user’s data.

## Table Usage Guide

The `okta_user` table provides insights into user profiles within Okta. As a security analyst, explore user-specific details through this table, including user status, last login, and assigned roles. Utilize it to uncover information about users, such as those with high-risk access levels, inactive users, and the verification of user profiles.

**Important Notes**
- This table supports an optional `filter` column to query results based on Okta supported [filters](https://developer.okta.com/docs/reference/api/apps/#filters).

## Examples

### Basic info
Explore the basic user information in your Okta system to understand the status and type of each user. This can help in managing user accounts and ensuring the correct access levels are granted.

```sql+postgres
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

```sql+sqlite
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
Explore the various roles, profiles, and group affiliations of each user to understand their access levels and responsibilities within the system. This can assist in managing user permissions and ensuring appropriate access control.

```sql+postgres
select
  id,
  email,
  jsonb_pretty(profile) as profile,
  jsonb_pretty(user_groups) as user_groups,
  jsonb_pretty(assigned_roles) as assigned_roles
from
  okta_user;
```

```sql+sqlite
select
  id,
  email,
  profile,
  user_groups,
  assigned_roles
from
  okta_user;
```

### List users with SUPER_ADMIN role access
Explore which users have been granted the highest level of access, the SUPER_ADMIN role, in order to maintain a secure and controlled environment. This is especially useful in managing system security and monitoring potential risks.

```sql+postgres
select
  id,
  login,
  jsonb_pretty(assigned_roles) as assigned_roles
from
  okta_user
where
  assigned_roles @> '[{"type":"SUPER_ADMIN"} ]'::jsonb;
```

```sql+sqlite
select
  id,
  login,
  assigned_roles
from
  okta_user
where
  json_extract(assigned_roles, '$[*].type') = 'SUPER_ADMIN';
```

### List users who have not logged in for more than 30 days
Identify users who may not be actively using the service by pinpointing those who haven't logged in for over a month. This can be useful in engagement analysis or for conducting user clean-ups.

```sql+postgres
select
  id,
  email,
  last_login
from
  okta_user
where
  last_login < current_timestamp - interval '30 days';
```

```sql+sqlite
select
  id,
  email,
  last_login
from
  okta_user
where
  last_login < datetime('now','-30 days');
```

### List active users that have been last updated before a specific date using a filter
Analyze the active users who have last updated their details before a certain date. This can be useful to pinpoint users who may need to update their information, improving account security and accuracy.

```sql+postgres
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

```sql+sqlite
select
  id,
  email,
  created,
  status
from
  okta_user
where
  filter = 'lastUpdated lt "2021-08-05T00:00:00.000Z" and status = "ACTIVE"';
```