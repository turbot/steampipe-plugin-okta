---
title: "Steampipe Table: okta_factor - Query Okta Factors using SQL"
description: "Allows users to query Okta Factors, specifically the authentication methods used by Okta users, providing insights into security protocols and potential vulnerabilities."
---

# Table: okta_factor - Query Okta Factors using SQL

Okta Factors are the different methods of authentication used by Okta users. These can range from password-based authentication to more advanced methods like biometric authentication. Understanding these factors is crucial for maintaining the security and integrity of an Okta environment.

## Table Usage Guide

The `okta_factor` table provides insights into the authentication methods used within Okta. As a security engineer, explore factor-specific details through this table, including the type of factor, status, and associated metadata. Utilize it to uncover information about factors, such as those that are less secure, the distribution of factor types among users, and potential vulnerabilities in authentication methods.

## Examples

### Basic info
Explore which security factors have been created for each user and their current status. This can be useful to understand the security measures in place for each user and if they are active or not.

```sql+postgres
select
  id,
  user_id,
  factor_type,
  created,
  status
from
  okta_factor;
```

```sql+sqlite
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
Explore which security factors are awaiting activation, enabling you to take necessary actions to ensure user accounts are secure and fully functional. This helps in maintaining the security standards and smooth operation of your system.

```sql+postgres
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

```sql+sqlite
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
Explore which factors are provided by Okta to gain insights into the status and creation dates of these factors. This can be useful in managing user security settings and assessing the elements within your Okta environment.

```sql+postgres
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

```sql+sqlite
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

### Get factor by factor ID and user ID
Determine the security factor settings associated with a specific user, which can be useful in understanding the user's security setup and status. This can be particularly helpful in troubleshooting or auditing security compliance.

```sql+postgres
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

```sql+sqlite
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