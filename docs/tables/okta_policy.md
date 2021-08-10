# Table: okta_policy

[Policies](https://developer.okta.com/docs/concepts/policies/) help you manage access to your applications and APIs. You can restrict access based on a number of conditions such as user and group membership, device, location, or time. You can also require more authentication steps for access to sensitive applications, such as confirmation of a push notification to a mobile device or re-authentication through an SMS one-time passcode.

Note: This table requires `type` (i.e [policy type](https://developer.okta.com/docs/reference/api/policy/#policy-types)) as required field.

Supports below policy types:

- [Okta Sign On Policy](https://developer.okta.com/docs/reference/api/policy/#okta-sign-on-policy) controls the manner in which a user is allowed to sign on to Okta, including whether they are challenged for multifactor authentication (MFA) and how long they are allowed to remain signed in before re-authenticating.

- [Multifactor (MFA) Enrollment Policy](https://developer.okta.com/docs/reference/api/policy/#multifactor-mfa-enrollment-policy) controls which MFA methods are available for a User, as well as when a User may enroll in a particular Factor.

- [IdP Discovery Policy](https://developer.okta.com/docs/reference/api/policy/#idp-discovery-policy) determines where to route Users when they are attempting to sign in to your org. Users can be routed to a variety of Identity Providers (SAML2, IWA, AgentlessDSSO, X509, FACEBOOK, GOOGLE, LINKEDIN, MICROSOFT, OIDC) based on multiple conditions.

## Examples

### List sign on policies

```sql
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_policy
where
  type = 'OKTA_SIGN_ON'
order by
  priority;
```

### List MFA enroll policies

```sql
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_policy
where
  type = 'MFA_ENROLL'
order by
  priority;
```

### List IDP Discovery policies

```sql
select
  name,
  id,
  created,
  status,
  priority,
  system
from
  okta_policy
where
  type = 'IDP_DISCOVERY';
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
  okta_policy
where
  type = 'OKTA_SIGN_ON' and
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
  okta_policy,
  jsonb_array_elements(rules) as r
where
  type = 'OKTA_SIGN_ON';
```
