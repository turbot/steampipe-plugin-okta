---
title: "Steampipe Table: okta_authentication_policy - Query Okta Authentication Policies using SQL"
description: "Allows users to query Okta Authentication Policies, providing details about the policies including their names, IDs, status, priority, and other related information."
---

# Table: okta_authentication_policy - Query Okta Authentication Policies using SQL

Okta Authentication Policies are a set of rules that specify the actions to be taken during user authentication based on a variety of conditions. These policies govern the authentication requirements users must meet before they are granted access to applications. They are an integral part of Okta's adaptive multi-factor authentication (MFA) and can be used to increase an organization's security.

## Table Usage Guide

The `okta_authentication_policy` table provides insights into Okta Authentication Policies. As a security analyst, you can leverage this table to understand the various authentication policies within your organization, including their priority, status, and the conditions under which they are applied. This information is crucial for auditing security measures and ensuring that your organization's authentication procedures are in line with best practices.

## Examples

### Basic info
Explore the priority-based organization of Okta authentication policies. This query can be used to assess the order of policies based on their priority, providing insights into the system's security measures and configurations.

```sql+postgres
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_authentication_policy
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
  okta_authentication_policy
order by
  priority;
```

### List inactive sign on policies
Explore which authentication policies are inactive. This is useful for maintaining security by identifying potential gaps in your active policies.

```sql+postgres
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_authentication_policy
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
  okta_authentication_policy
where
  status = 'INACTIVE';
```

### Get rules details for each sign on policy
This query is useful to gain insights into each authentication policy's rules within your system. It provides a detailed view of the rules' names, systems, statuses, priorities, actions, and conditions, aiding in policy management and security assessment.

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
  okta_authentication_policy,
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
  okta_authentication_policy,
  json_each(rules) as r;
```