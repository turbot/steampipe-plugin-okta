# Table: okta_factor

The Okta Factor provides operations to enroll, manage, and verify factors for multifactor authentication (MFA). It allows to manage both administration and end-user accounts, or verify an individual factor at any time.

## Examples

### Basic info

```sql
select
  id,
  user_id,
  factor_type,
  created,
  status
from
  okta_factor;
```

### List factors pending activation

```sql
select
  id,
  user_id,
  factor_type,
  created,
  status
from
  okta_factor
where
  status = 'PENDING_ACTIVATION';
```

### List factors provided by Okta

```sql
select
  id,
  user_id,
  factor_type,
  created,
  provider,
  status
from
  okta_factor
where
  provider = 'OKTA';
```

### List factor links

```sql
select
  id,
  user_id,
  status,
  jsonb_pretty(links -> 'activate') as link_activate,
  jsonb_pretty(links -> 'self') as link_self,
  jsonb_pretty(links -> 'user') as link_scopes,
  jsonb_pretty(links -> 'verify') as link_rotateKey
from
  okta_factor;
```

### Get factor by factor ID and user ID

```sql
select
  id,
  user_id,
  factor_type,
  created,
  status
from
  okta_factor
where
  id = 'ost1l5cklwIRvLzUY5d7' and user_id = '00u1kcigdvWtR96HY5d7';
```