# Table: okta_sign_on_policy

[Okta Sign On Policy](https://developer.okta.com/docs/reference/api/policy/#okta-sign-on-policy) controls the manner in which a user is allowed to sign on to Okta, including whether they are challenged for multifactor authentication (MFA) and how long they are allowed to remain signed in before re-authenticating.

## Examples

### Basic info

```sql
select
  name,
  id,
  created,
  status,
  priority,
  system,
  jsonb_pretty(rules) as rules
from
  okta_sign_on_policy
order by
  priority;
```

### List custom sign on policies

```sql
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_sign_on_policy
where
  not system;
```

### List inactive sign on policies

```sql
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_sign_on_policy
where
  status = 'INACTIVE';
```

### Get sign on policy rules details

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
  okta_sign_on_policy,
  jsonb_array_elements(rules) as r;
```
