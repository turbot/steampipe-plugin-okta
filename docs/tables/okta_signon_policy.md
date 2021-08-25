# Table: okta_signon_policy

The [Sign On Policy](https://developer.okta.com/docs/reference/api/policy/#okta-sign-on-policy) controls the manner in which a user is allowed to sign on to Okta, including whether they are challenged for multifactor authentication (MFA) and how long they are allowed to remain signed in before re-authenticating.

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
  okta_signon_policy
order by
  priority;
```

### List system sing on policies

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
