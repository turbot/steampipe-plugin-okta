# Table: okta_idp_discovery_policy

The [IdP Discovery Policy](https://developer.okta.com/docs/reference/api/policy/#idp-discovery-policy) determines where to route Users when they are attempting to sign in to your org. Users can be routed to a variety of Identity Providers (SAML2, IWA, AgentlessDSSO, X509, FACEBOOK, GOOGLE, LINKEDIN, MICROSOFT, OIDC) based on multiple conditions.

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
  okta_idp_discovery_policy
order by
  priority;
```

### List system idp discovery policies

```sql
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

```sql
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
  okta_idp_discovery_policy,
  jsonb_array_elements(rules) as r;
```
