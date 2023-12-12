---
title: "Steampipe Table: okta_idp_discovery_policy - Query Okta Identity Provider Discovery Policies using SQL"
description: "Allows users to query Okta Identity Provider Discovery Policies, providing insights into the configuration and rules associated with the policy."
---

# Table: okta_idp_discovery_policy - Query Okta Identity Provider Discovery Policies using SQL

Okta Identity Provider Discovery is a feature that allows organizations to route users to different identity providers based on certain conditions. These conditions are defined in the Identity Provider Discovery Policy. This feature helps organizations manage multiple identity providers and control user access based on their attributes or group membership.

## Table Usage Guide

The `okta_idp_discovery_policy` table provides insights into the Identity Provider Discovery Policies within Okta. As a Security or IT administrator, explore policy-specific details through this table, including conditions, actions, and associated rules. Utilize it to uncover information about policies, such as those with specific conditions, the actions associated with each policy, and the verification of rules.

## Examples

### Basic info
Explore the priority-based arrangement of identity provider discovery policies in your system, which can help you understand their creation timelines, statuses, and associated identities for better management and security.

```sql+postgres
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_idp_discovery_policy
order by
  priority;
```

```sql+sqlite
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_idp_discovery_policy
order by
  priority;
```

### List system idp discovery policies
Explore the discovery policies in your system with this query. It helps you understand the priority and status of each policy, and when it was created, providing a comprehensive view of your system's identity provider (IdP) discovery policies.

```sql+postgres
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_idp_discovery_policy
where
  system;
```

```sql+sqlite
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_idp_discovery_policy
where
  system;
```

### List inactive idp discovery policies
Explore which IDP discovery policies are inactive. This can be useful for identifying policies that are no longer in use and may need to be updated or removed.

```sql+postgres
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_idp_discovery_policy
where
  status = 'INACTIVE';
```

```sql+sqlite
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_idp_discovery_policy
where
  status = 'INACTIVE';
```

### Get rules details for each idp discovery policy
Analyze the specifics of each identity provider discovery policy to gain insights into the rules applied, including their name, system, status, and priority. This can be useful in reviewing and managing your security configurations and policies.

```sql+postgres
select
  name,
  id,
  r -> 'name' as rule_name,
  r -> 'system' as rule_system,
  r -> 'status' as rule_status,
  r -> 'priority' as rule_priority,
  jsonb_pretty(r -> 'actions') as rule_actions,
  jsonb_pretty(r -> 'conditions') as rule_conditions
from
  okta_idp_discovery_policy,
  jsonb_array_elements(rules) as r;
```

```sql+sqlite
select
  name,
  id,
  json_extract(r.value, '$.name') as rule_name,
  json_extract(r.value, '$.system') as rule_system,
  json_extract(r.value, '$.status') as rule_status,
  json_extract(r.value, '$.priority') as rule_priority,
  r.value as rule_actions,
  r.value as rule_conditions
from
  okta_idp_discovery_policy,
  json_each(rules) as r;
```