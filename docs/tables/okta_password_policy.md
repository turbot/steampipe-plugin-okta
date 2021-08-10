# Table: okta_password_policy

The [Password Policy](https://developer.okta.com/docs/reference/api/policy/#password-policy) determines the requirements for a user's password length and complexity, as well as the frequency with which a password must be changed. This Policy also governs the recovery operations that may be performed by the User, including change password, reset (forgot) password, and self-service password unlock.

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
  okta_password_policy;
```

### List guest users

```sql
select
  display_name,
  user_principal_name,
  id,
  mail
from
  okta_password_policy
where
  user_type = 'Guest';
```
