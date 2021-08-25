# Table: okta_mfa_policy

The [MFA Policy](https://developer.okta.com/docs/reference/api/policy/#multifactor-mfa-enrollment-policy) controls which MFA methods are 
available for a User, as well as when a User may enroll in a particular Factor.

## Examples

### Basic info

```sql
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

```sql
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

### List inactive mfa policies

```sql
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

### Get rules details for each mfa policy

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
  okta_mfa_policy,
  jsonb_array_elements(rules) as r;
```
