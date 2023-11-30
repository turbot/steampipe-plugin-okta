---
title: "Steampipe Table: okta_signon_policy - Query Okta Sign-On Policies using SQL"
description: "Allows users to query Okta Sign-On Policies, providing details about the policies including their names, IDs, status, priority, and other related information."
---

# Table: okta_signon_policy - Query Okta Sign-On Policies using SQL

Okta Sign-On Policies are a set of rules that specify the actions to be taken during user sign-in based on a variety of conditions. These policies govern the authentication requirements users must meet before they are granted access to applications. They are an integral part of Okta's adaptive multi-factor authentication (MFA) and can be used to increase an organization's security.

## Table Usage Guide

The `okta_signon_policy` table provides insights into Okta Sign-On Policies. As a security analyst, you can leverage this table to understand the various sign-on policies within your organization, including their priority, status, and the conditions under which they are applied. This information is crucial for auditing security measures and ensuring that your organization's sign-on procedures are in line with best practices.

## Examples

### Basic info
Explore the priority-based organization of Okta sign-on policies. This query can be used to assess the order of policies based on their priority, providing insights into the system's security measures and configurations.

```sql
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_signon_policy
order by
  priority;
```

### List system sign on policies
Explore which system sign-on policies are currently in place. This can help in understanding the security measures in effect and prioritizing any necessary changes.

```sql
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_signon_policy
where
  system;
```

### List inactive sign on policies
Explore which sign-on policies are inactive. This is useful for maintaining security by identifying potential gaps in your active policies.

```sql
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_signon_policy
where
  status = 'INACTIVE';
```

### Get rules details for each sign on policy
This query is useful to gain insights into each sign-on policy's rules within your system. It provides a detailed view of the rules' names, systems, statuses, priorities, actions, and conditions, aiding in policy management and security assessment.

```sql
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
  okta_signon_policy,
  jsonb_array_elements(rules) as r;
```