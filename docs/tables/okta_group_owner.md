---
title: "Steampipe Table: okta_group_owner - Query Okta Group Owners using SQL"
description: "Allows users to query Okta Group Owners, providing insights into group ownership details within Okta."
---

# Table: okta_group_owner - Query Okta Group Owners using SQL

Okta Group Owners are individuals responsible for managing and overseeing specific groups within the Okta identity and access management platform. They have the authority to add or remove group members, set group policies, and ensure that the group's access permissions are correctly configured to maintain security and compliance.

## Table Usage Guide

The `okta_group_owner` table provides detailed information about group owners within Okta. As an IT administrator, you can explore ownership details through this table, including the display name, origin type, and last updated timestamp. Utilize it to manage group ownership, identify responsible individuals, and ensure proper group administration.

## Examples

### Basic Info
Retrieve basic information about group owners in Okta to understand their roles and responsibilities. This is useful for managing group ownership and ensuring proper oversight.

```sql+postgres
select
  group_id,
  id,
  display_name,
  type,
  origin_type,
  last_updated
from
  okta_group_owner;
```

```sql+sqlite
select
  group_id,
  id,
  display_name,
  type,
  origin_type,
  last_updated
from
  okta_group_owner;
```

### List group owners by domain
Identify group owners based on their Okta domain. This helps in organizing and managing owners within specific domains.

```sql+postgres
select
  domain,
  group_id,
  id,
  display_name
from
  okta_group_owner
where
  domain = 'example.com';
```

```sql+sqlite
select
  domain,
  group_id,
  id,
  display_name
from
  okta_group_owner
where
  domain = 'example.com';
```

### Group owners with application origin
Find group owners whose ownership is managed by applications. This can help in understanding the source of group management and reconciling ownership details.

```sql+postgres
select
  group_id,
  id,
  display_name,
  origin_id,
  origin_type,
  resolved
from
  okta_group_owner
where
  origin_type = 'APPLICATION';
```

```sql+sqlite
select
  group_id,
  id,
  display_name,
  origin_id,
  origin_type,
  resolved
from
  okta_group_owner
where
  origin_type = 'APPLICATION';
```

### List recently updated group owners
Get a list of group owners who were recently updated. This helps in tracking recent changes and updates to group ownership.

```sql+postgres
select
  group_id,
  id,
  display_name,
  last_updated
from
  okta_group_owner
where
  last_updated > current_timestamp - interval '30 days';
```

```sql+sqlite
select
  group_id,
  id,
  display_name,
  last_updated
from
  okta_group_owner
where
  last_updated > datetime('now', '-30 days');
```