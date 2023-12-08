---
title: "Steampipe Table: okta_trusted_origin - Query Okta Trusted Origins using SQL"
description: "Allows users to query Trusted Origins in Okta, providing insights into the trusted origins for the Okta domains."
---

# Table: okta_trusted_origin - Query Okta Trusted Origins using SQL

Okta Trusted Origins is a service within Okta that allows you to manage the origins that are trusted to start browser-based authentication flows or permitted to make CORS requests. It provides a centralized way to manage these trusted origins for various Okta domains. Okta Trusted Origins helps you ensure the security of your Okta domains by controlling the origins that are allowed to interact with them.

## Table Usage Guide

The `okta_trusted_origin` table provides insights into trusted origins within Okta. As a security engineer, explore trusted origin-specific details through this table, including origin names, origin types, and associated metadata. Utilize it to uncover information about trusted origins, such as those with CORS or redirect permissions, and the verification of these permissions.

**Important Notes**
- This table supports an optional `filter` column to query results based on Okta supported [filters](https://developer.okta.com/docs/reference/api/apps/#filters).

## Examples

### Basic info
Explore which trusted origins in your Okta environment have been recently updated or created. This helps keep track of changes and maintain the security of your applications and APIs.

```sql+postgres
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

```sql+sqlite
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

### List trusted origins last updated 30 days ago
Determine the trusted origins that have not been updated in the past 30 days. This is useful for maintaining security by ensuring all trusted origins are up-to-date.

```sql+postgres
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
  last_updated < current_timestamp - interval '30 days';
```

```sql+sqlite
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
  last_updated < datetime('now', '-30 day');
```

### List CORS scoped trusted origins
Explore which trusted origins have been scoped for Cross-Origin Resource Sharing (CORS) to understand the security measures in place for data requests from different origins. This can help in assessing potential vulnerabilities and ensuring appropriate CORS policies are implemented.

```sql+postgres
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
  scopes @> '[{"type":"CORS"}]'::jsonb;
```

```sql+sqlite
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
  json_extract(scopes, '$[0].type') = 'CORS';
```