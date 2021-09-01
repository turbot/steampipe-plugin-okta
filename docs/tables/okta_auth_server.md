# Table: okta_auth_server

An authorization server defines your security boundary, and is used to mint access and identity tokens for use with OIDC clients and OAuth 2.0 service accounts when accessing your resources via API. Within each authorization server you can define your own OAuth scopes, claims, and access policies.

## Examples

### Basic info

```sql
select
  name,
  id,
  audiences,
  created,
  last_updated,
  status
from
  okta_auth_server;
```

### List authorization servers where manual rotation signing keys are not rotated in more than 90 days

```sql
select
  name,
  id,
  audiences,
  created,
  last_updated, 
  credentials -> 'signing' ->> 'lastRotated' as last_rotated,
  status
from
  okta_auth_server
where
  credentials -> 'signing' ->> 'rotationMode' = 'MANUAL' 
  and CAST(credentials -> 'signing' ->> 'lastRotated' as date) < current_timestamp - interval '90 days';
```

### List inactive authorization servers

```sql
select
  name,
  id,
  audiences,
  created,
  last_updated,
  status
from
  okta_auth_server
where
  status = 'INACTIVE';
```
