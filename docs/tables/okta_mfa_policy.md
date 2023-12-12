---
title: "Steampipe Table: okta_mfa_policy - Query Okta Multi-Factor Authentication Policies using SQL"
description: "Allows users to query Multi-Factor Authentication Policies in Okta, specifically the policy details and settings, providing insights into the security measures in place."
---

# Table: okta_mfa_policy - Query Okta Multi-Factor Authentication Policies using SQL

Okta Multi-Factor Authentication (MFA) is a security feature that provides an additional layer of protection for user accounts. It requires users to verify their identity with at least two forms of identification before gaining access to resources. MFA policies in Okta allow administrators to define and enforce security measures.

## Table Usage Guide

The `okta_mfa_policy` table provides insights into the MFA policies within Okta. As a security administrator, you can explore policy-specific details through this table, including policy settings, conditions, and associated metadata. Use it to uncover information about policies, such as those with specific conditions, the type of factors enforced, and the verification of policy settings.

## Examples

### Basic info
Explore the multi-factor authentication policies in your Okta system to understand their creation dates, priority levels, and statuses. This will help you assess the security strength and identify any potential vulnerabilities or areas for improvement.

```sql+postgres
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_mfa_policy
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
  okta_mfa_policy
order by
  priority;
```

### List system mfa policies
Explore the Multi-Factor Authentication (MFA) policies active within your system. This aids in assessing the priority and status of each policy, helping you maintain a secure and efficient environment.

```sql+postgres
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_mfa_policy
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
  okta_mfa_policy
where
  system = 1;
```

### List inactive mfa policies
Explore which multi-factor authentication (MFA) policies are currently inactive. This is useful for security audits to ensure all necessary policies are active and functioning as intended.

```sql+postgres
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_mfa_policy
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
  okta_mfa_policy
where
  status = 'INACTIVE';
```

### List highest priority mfa policy details
Explore the highest priority multi-factor authentication (MFA) policies in your system. This can be useful for prioritizing security measures and identifying potential vulnerabilities.

```sql+postgres
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_mfa_policy
where
  priority = 1;
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
  okta_mfa_policy
where
  priority = 1;
```

### Get rules details for each mfa policy
Explore the specific rules associated with each multi-factor authentication policy. This allows for a comprehensive understanding of the security measures in place and their prioritization, enabling more informed decisions about potential modifications or enhancements.

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
  okta_mfa_policy,
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
  okta_mfa_policy,
  json_each(rules) as r;
```