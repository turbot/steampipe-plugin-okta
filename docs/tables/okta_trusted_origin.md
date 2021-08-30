# Table: okta_trusted_origin

The Okta Trusted Origins API provides operations to manage Trusted Origins and sources.

When external URLs are requested during sign-in, sign-out, or recovery operations, Okta checks those URLs against the allowed list of Trusted Origins. Trusted Origins also enable browser-based applications to access Okta APIs from JavaScript (CORS). If the origins aren't specified, the related operation (redirect or Okta API access) isn't permitted.

Note: This table does not support the optional `filter` column to query results based on Okta supported [filters](https://developer.okta.com/docs/reference/api/trusted-origins/#list-trusted-origins-with-a-filter).

## Examples

### Basic info

```sql
select
  name,
  id,
  created,
  last_updated,
  origin,
  scopes,
  status
from
  okta_trusted_origin;
```

### List trusted origins links

```sql
select
  name,
  id,
  status,
  jsonb_pretty(links -> 'deactivate') as link_deactivate,
  jsonb_pretty(links -> 'self') as link_self
from
  okta_trusted_origin;
```

### Get authorization server by ID

```sql
select
  name,
  id,
  created,
  last_updated,
  origin,
  scopes,
  status
from
  okta_trusted_origin
where
  id = 'tos1l3v1djOJMSQkh5d7';
```
