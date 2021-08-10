# Table: okta_password_policy

The [Password Policy](https://developer.okta.com/docs/reference/api/policy/#password-policy) determines the requirements for a user's password length and complexity, as well as the frequency with which a password must be changed. This Policy also governs the recovery operations that may be performed by the User, including change password, reset (forgot) password, and self-service password unlock.

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
  okta_password_policy
order by
  priority;
```

### List system password policies

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

### List inactive policies

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

### Get policy password settings details guest users

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

### Get policy password rules details

```sql
select
  name,
  id,
  r -> 'name' as rule_name,
  r -> 'system' as rule_system,
  r -> 'status' as rule_status,
  r -> 'priority' as rule_priority,
  jsonb_pretty(r -> 'actions') as rule_actions,
  jsonb_pretty(r -> 'priority') as rule_conditions
from
  okta_password_policy,
  jsonb_array_elements(rules) as r;
```
