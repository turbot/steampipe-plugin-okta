# Table: okta_sign_on_policy

[Okta Sign On Policy](https://developer.okta.com/docs/reference/api/policy/#okta-sign-on-policy) controls the manner in which a user is allowed to sign on to Okta, including whether they are challenged for multifactor authentication (MFA) and how long they are allowed to remain signed in before re-authenticating.

## Examples

### Basic info

```sql
select
  display_name,
  user_principal_name,
  id,
  given_name,
  mail
from
  okta_sign_on_policy;
```

### List guest users

```sql
select
  display_name,
  user_principal_name,
  id,
  mail
from
  okta_sign_on_policy
where
  user_type = 'Guest';
```
