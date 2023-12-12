---
title: "Steampipe Table: okta_group - Query Okta Groups using SQL"
description: "Allows users to query Okta Groups, specifically the group profile details, providing insights into group memberships and access control."
---

# Table: okta_group - Query Okta Groups using SQL

Okta Groups are a collection of users defined in the Okta service. They provide a way to manage users and their access to applications and resources. Groups are central to the role-based access control (RBAC) model in Okta, and they can be used for assigning roles and permissions.

## Table Usage Guide

The `okta_group` table provides insights into groups within Okta. As an IT administrator, explore group-specific details through this table, including group profile, type, and associated users. Utilize it to manage access control, identify groups with specific roles, and verify the consistency of group memberships.

**Important Notes**
- This table supports an optional `filter` column to query results based on Okta supported [filters](https://developer.okta.com/docs/reference/api/groups/#filters).

## Examples

### Basic info
Explore the basic information about user groups in Okta to understand their purpose and configuration. This is useful for managing access controls and implementing security policies.

```sql+postgres
select
  name,
  id,
  type,
  description,
  jsonb_pretty(profile) as profile
from
  okta_group;
```

```sql+sqlite
select
  name,
  id,
  type,
  description,
  profile
from
  okta_group;
```

### List groups without membership changes for more than 30 days
Determine the groups that have not undergone membership alterations in over a month. This could be useful for identifying inactive or stagnant groups and assessing the need for membership reviews or updates.

```sql+postgres
select
  name,
  id,
  type,
  age(current_timestamp, last_membership_updated) as last_membership_updated
from
  okta_group
where
  last_membership_updated < current_timestamp - interval '30 days';
```

```sql+sqlite
select
  name,
  id,
  type,
  julianday('now') - julianday(last_membership_updated) as last_membership_updated
from
  okta_group
where
  julianday('now') - julianday(last_membership_updated) > 30;
```

### List groups with profile or membership updates after a specific date using a filter
Explore which groups have had updates to their profiles or memberships after a specific date. This is useful for keeping track of recent changes in group data and ensuring up-to-date information.

```sql+postgres
select
  name,
  id,
  type,
  last_updated,
  last_membership_updated
from
  okta_group
where
  filter = 'type eq "OKTA_GROUP" and (lastUpdated gt "2021-05-05T00:00:00.000Z" or lastMembershipUpdated gt "2021-05-05T00:00:00.000Z")';
```

```sql+sqlite
select
  name,
  id,
  type,
  last_updated,
  last_membership_updated
from
  okta_group
where
  filter = 'type eq "OKTA_GROUP"' 
  and (datetime(lastUpdated) > datetime('2021-05-05T00:00:00') 
  or datetime(lastMembershipUpdated) > datetime('2021-05-05T00:00:00'));
```

### Get group member details for each group
Determine the members associated with each group within your organization. This can help in understanding the group structure and managing user access effectively.

```sql+postgres
select
  name,
  id,
  jsonb_pretty(group_members) as group_members
from
  okta_group;
```

```sql+sqlite
select
  name,
  id,
  group_members
from
  okta_group;
```