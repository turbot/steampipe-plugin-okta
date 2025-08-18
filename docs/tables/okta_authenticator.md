---
title: "Steampipe Table: okta_authenticator - Query Okta Authenticators using SQL"
description: "Allows users to query Okta Authenticators configured in an organization, including type, key, status, timestamps, and related configuration."
---

# Table: okta_authenticator - Query Okta Authenticators using SQL

An Okta Authenticator represents a method users can enroll in and use to verify their identity (for example: email, password, phone, security key, or security question). Authenticators are configured at the org level and can be enabled, disabled, and customized.

## Table Usage Guide

The `okta_authenticator` table provides information about each authenticator configured in your Okta organization. As a security or identity engineer, use this table to review authenticator inventory, enablement status, and configuration details. Common use cases include listing all authenticators, identifying disabled authenticators, and reviewing when authenticators were created or last updated.

## Examples

### Basic info
List all authenticators with their type, key, status, and creation time.

```sql+postgres
select
  id,
  name,
  type,
  key,
  status,
  created
from
  okta_authenticator;
```

```sql+sqlite
select
  id,
  name,
  type,
  key,
  status,
  created
from
  okta_authenticator;
```

### List inactive authenticators
Identify authenticators that are configured but not currently active.

```sql+postgres
select
  id,
  name,
  type,
  key,
  status,
  last_updated
from
  okta_authenticator
where
  status = 'INACTIVE';
```

```sql+sqlite
select
  id,
  name,
  type,
  key,
  status,
  last_updated
from
  okta_authenticator
where
  status = 'INACTIVE';
```

### Show recently updated authenticators
Review authenticators that have been modified in the last 30 days.

```sql+postgres
select
  id,
  name,
  type,
  key,
  status,
  last_updated
from
  okta_authenticator
where
  last_updated > now() - interval '30 days'
order by
  last_updated desc;
```

```sql+sqlite
select
  id,
  name,
  type,
  key,
  status,
  last_updated
from
  okta_authenticator
where
  last_updated > datetime('now', '-30 days')
order by
  last_updated desc;
```

### Get authenticator by ID
Retrieve details for a specific authenticator by its unique ID.

```sql+postgres
select
  id,
  name,
  type,
  key,
  status,
  created,
  last_updated,
  settings,
  provider
from
  okta_authenticator
where
  id = 'aut_12345';
```

```sql+sqlite
select
  id,
  name,
  type,
  key,
  status,
  created,
  last_updated,
  settings,
  provider
from
  okta_authenticator
where
  id = 'aut_12345';
```

### Inspect settings for app-type authenticators
Explore provider and settings JSON for deeper configuration insight.

```sql+postgres
select
  name,
  type,
  key,
  provider,
  settings
from
  okta_authenticator
order by
  name;
```

```sql+sqlite
select
  name,
  type,
  key,
  provider,
  settings
from
  okta_authenticator
order by
  name;
```

