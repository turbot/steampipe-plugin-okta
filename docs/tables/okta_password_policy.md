---
title: "Steampipe Table: okta_password_policy - Query Okta Password Policies using SQL"
description: "Allows users to query Okta Password Policies, providing insights into the password policy configurations and rules within an Okta organization."
---

# Table: okta_password_policy - Query Okta Password Policies using SQL

Okta Password Policy is a set of rules and settings within Okta that governs the complexity requirements for user passwords and the actions to take when users violate these rules. It helps organizations to enforce strong password practices, enhancing security by reducing the risk of password-based attacks. Okta Password Policy provides a way to customize these rules and settings to meet the specific security needs of an organization.

## Table Usage Guide

The `okta_password_policy` table provides insights into the password policies within Okta. As a security analyst, explore policy-specific details through this table, including the complexity requirements, lockout settings, and associated metadata. Utilize it to uncover information about the password policies, such as those with weak complexity requirements, the number of failed attempts before a lockout, and the duration of the lockout period.

## Examples

### Basic info
Explore which password policies have been implemented, understanding their creation dates, status, and priority. This can be useful for assessing the security measures in place and their relative importance.

```sql
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_password_policy
order by
  priority;
```

### List system password policies
Analyze the settings to understand the system's password policies, enabling you to assess their creation, status, and priority. This is beneficial for maintaining security standards and prioritizing system updates.

```sql
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_password_policy
where
  system;
```

### List inactive password policies
Explore which password policies have been marked as inactive, allowing you to identify and review any outdated or unused policies that could potentially impact system security. This is particularly useful in maintaining security standards and ensuring all policies are up-to-date.

```sql
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_password_policy
where
  status = 'INACTIVE';
```

### Get policy details for each password policy
Explore the specifics of each password policy, including its age, complexity, lockout details, recovery factors, and delegation options. This can help you understand and manage the security standards across different policies.

```sql
select
  name,
  id,
  status,
  jsonb_pretty(settings -> 'password' -> 'age') as password_age,
  jsonb_pretty(settings -> 'password' -> 'complexity') as password_complexity,
  jsonb_pretty(settings -> 'password' -> 'lockout') as password_lockout,
  jsonb_pretty(settings -> 'recovery' -> 'factors') as recovery_factors,
  jsonb_pretty(settings -> 'delegation' -> 'options') as delegation_options
from
  okta_password_policy;
```

### Get rules details for each password policy
Explore the specific rules associated with each password policy to gain insights into their statuses, priorities, and conditions. This can help in understanding and managing security measures more effectively.

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
  okta_password_policy,
  jsonb_array_elements(rules) as r;
```