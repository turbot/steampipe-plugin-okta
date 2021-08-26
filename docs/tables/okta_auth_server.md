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

### List authorization server links

```sql
select
  name,
  id,
  status,
  jsonb_pretty(links -> 'activate') as link_activate,
  jsonb_pretty(links -> 'claims') as link_claims,
  jsonb_pretty(links -> 'deactivate') as link_deactivate,
  jsonb_pretty(links -> 'metadata') as link_metadata,
  jsonb_pretty(links -> 'policies') as link_policies,
  jsonb_pretty(links -> 'rotateKey') as link_rotateKey,
  jsonb_pretty(links -> 'scopes') as link_scopes,
  jsonb_pretty(links -> 'self') as link_self
from
  okta_auth_server;
```

### Get authorization server by ID

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
  id = 'aus1kchdp0mdlLV7o5d7';
```
