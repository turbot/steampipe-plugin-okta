---
title: "Steampipe Table: okta_group_rule - Query Okta Group Rules using SQL"
description: "Allows users to query Okta Group Rules, specifically the group rule details, providing insights into group rule assignments and access control."
---

# Table: okta_group_rule - Query Okta Group Rules using SQL

Okta Group Rules are a collection of rules defined in the Okta service. They provide a way to manage group rules and their access to applications and resources. Group Rules are central to the role-based access control (RBAC) model in Okta, and they can be used for assigning roles and permissions.

## Table Usage Guide

The `okta_group_rule` table provides insights into group rules within Okta. As an IT administrator, explore group rule-specific details through this table, including group rule profile, type, and associated groups. Utilize it to manage access control, identify group rules with specific roles, and verify the consistency of group rule assignments.

## Examples

### Basic info

Explore the basic information about group rules in Okta to understand their purpose and configuration. This is useful for managing access controls and implementing security policies.

```sql+postgres
select
  name,
  id,
  status,
  actions,
  last_updated,
  jsonb_pretty(conditions) as conditions
from
  okta_group_rule;
```

```sql+sqlite
select
  name,
  id,
  status,
  actions,
  last_updated,
  conditions
from
  okta_group_rule;
```

### List group rules without assignment changes for more than 30 days

Determine the group rules that have not undergone assignment alterations in over a month. This could be useful for identifying inactive or stagnant group rules and assessing the need for assignment reviews or updates.

```sql+postgres
select
  name,
  id,
  status,
  last_updated
from
  okta_group_rule
where
  last_updated < current_timestamp - interval '30 days';
```

```sql+sqlite
select
  name,
  id,
  status,
  last_updated
from
  okta_group_rule
where
  last_updated < strftime('%s', 'now') - 30*24*60*60;
```

### List active group rules

Retrieve group rules with the status 'ACTIVE' to identify active group rules and ensure that they are correctly configured. This is useful for managing access control and verifying the consistency of group rule assignments.

```sql+postgres
select
  name,
  id,
  status,
  jsonb_pretty(conditions) as conditions
from
  okta_group_rule
where
  status = 'ACTIVE';
```

```sql+sqlite
select
  name,
  id,
  status,
  conditions
from  
  okta_group_rule
where
  status = 'ACTIVE';
```

### List group rules with specific group membership actions

This query retrieves Okta group rules, including the name, ID, status, and conditions, while filtering by specific group membership actions. This is useful for identifying group rules that assign users to specific groups.

```sql+postgres
with expanded_actions as (
  select
    name,
    id,
    status,
    jsonb_pretty(conditions) as conditions,
    jsonb_array_elements_text(actions->'assignUserToGroups'->'groupIds') as group_id
  from
    okta_group_rule
)
select
  name,
  id,
  status,
  conditions
from
  expanded_actions
where
  group_id = '00gl0xw4khfR4h5qJ5d7';
```

```sql+sqlite
with expanded_actions as (
  select
    name,
    id,
    status,
    json(conditions) as conditions,
    json_each.value as group_id
  from
    okta_group_rule,
    json_each(actions->'$.assignUserToGroups.groupIds')
)
select
  name,
  id,
  status,
  conditions
from
  expanded_actions
where
  group_id = '00gl0xw4khfR4h5qJ5d7';
```
