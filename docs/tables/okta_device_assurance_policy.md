---
title: "Steampipe Table: okta_device_assurance_policy - Query Okta Device Assurance Policies using SQL"
description: "Allows users to query Okta device assurance policies, including the platform, minimum OS version, disk encryption, and screen lock requirements."
---

# Table: okta_device_assurance_policy - Query Okta Device Assurance Policies using SQL

A device assurance policy in Okta defines the posture a device must meet before it can access resources protected by an authentication policy. Each policy targets one platform (Android, ChromeOS, iOS, macOS, or Windows) and can require a minimum OS version, disk encryption, a screen lock, secure hardware, and more.

## Table Usage Guide

The `okta_device_assurance_policy` table provides the requirements configured in each device assurance policy. As a security or identity engineer, use this table to audit device posture requirements, for example to confirm that every platform requires disk encryption and a screen lock, or to review the OS version floor. Requirements that do not apply to a platform, or that are not configured, are returned as null.

**Important Notes**
- Reading device assurance policies requires an API token or OAuth scope that permits reading policies (`okta.policies.read`).

## Examples

### Basic info
List all device assurance policies with the platform they apply to.

```sql+postgres
select
  name,
  id,
  platform,
  created_date,
  last_update
from
  okta_device_assurance_policy;
```

```sql+sqlite
select
  name,
  id,
  platform,
  created_date,
  last_update
from
  okta_device_assurance_policy;
```

### List policies that do not require disk encryption
Find policies on platforms that support disk encryption checks but do not require an encryption type.

```sql+postgres
select
  name,
  platform
from
  okta_device_assurance_policy
where
  platform in ('ANDROID', 'MACOS', 'WINDOWS')
  and disk_encryption_type is null;
```

```sql+sqlite
select
  name,
  platform
from
  okta_device_assurance_policy
where
  platform in ('ANDROID', 'MACOS', 'WINDOWS')
  and disk_encryption_type is null;
```

### List policies that do not require a screen lock
Identify policies that do not enforce a screen lock on the device.

```sql+postgres
select
  name,
  platform
from
  okta_device_assurance_policy
where
  screen_lock_type is null;
```

```sql+sqlite
select
  name,
  platform
from
  okta_device_assurance_policy
where
  screen_lock_type is null;
```

### Show the minimum OS version required per platform
Review the operating system floor each policy enforces. A null value means no minimum is set.

```sql+postgres
select
  name,
  platform,
  os_version ->> 'minimum' as minimum_os_version
from
  okta_device_assurance_policy
order by
  platform;
```

```sql+sqlite
select
  name,
  platform,
  json_extract(os_version, '$.minimum') as minimum_os_version
from
  okta_device_assurance_policy
order by
  platform;
```

### List policies last modified more than a year ago
Find policies that may be out of date.

```sql+postgres
select
  name,
  platform,
  last_update
from
  okta_device_assurance_policy
where
  last_update < now() - interval '1 year';
```

```sql+sqlite
select
  name,
  platform,
  last_update
from
  okta_device_assurance_policy
where
  last_update < datetime('now', '-1 year');
```
