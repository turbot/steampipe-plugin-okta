---
title: "Steampipe Table: okta_device - Query Okta Devices using SQL"
description: "Allows users to query Okta Devices, specifically enhances security and user experience by leveraging device identity and context as part of the organization's overall identity and access management strategy"
---

# Table: okta_device - Query Okta Devices using SQL

Okta Devices enhances security and user experience by leveraging device identity and context as part of the organization's overall identity and access management strategy. It enables businesses to enforce security policies while providing flexible and convenient access to applications and services across any device.

## Table Usage Guide

The `okta_device` table delivers a comprehensive view of the devices that interact with the Okta ecosystem, presenting crucial information for IT and security professionals. As a part of your security and device management strategy, leverage this table to dive into device-specific data, encompassing device type, status, registration details, and security posture. This table is instrumental for assessing the landscape of devices accessing corporate resources, enabling the identification of unmanaged or insecure devices, evaluating the adherence of devices to corporate security policies, and pinpointing potential security risks tied to device access.

## Examples

### Basic info
Explore which device have been created and their current status. This can be useful to understand the security measures in place for each device and if they are active or not.

```sql+postgres
select
  display_name,
  id,
  created,
  last_updated,
  resource_type,
  status
from
  okta_device;
```

```sql+sqlite
select
  display_name,
  id,
  created,
  last_updated,
  resource_type,
  status
from
  okta_device;
```

### Get device by device ID
Determine the security factor settings associated with a specific user, which can be useful in understanding the user's security setup and status. This can be particularly helpful in troubleshooting or auditing security compliance.

```sql+postgres
select
  id,
  display_name,
  resource_type,
  created,
  status
from
  okta_device
where
  id = 'ost1l5cklwIRvLzUY5d7';
```

```sql+sqlite
select
  id,
  display_name,
  resource_type,
  created,
  status
from
  okta_device
where
  id = 'ost1l5cklwIRvLzUY5d7';
```

### List embedded user details of the devices
Explore which factors are provided by Okta to gain insights into the status and creation dates of these factors. This can be useful in managing user security settings and assessing the elements within your Okta environment.

```sql+postgres
select
  d.id,
  d.display_name,
  d.status,
  u ->> 'Created' as user_created,
  u ->> 'ManagementStatus' as user_management_status,
  u ->> 'ScreenLockType' as screen_lock_type,
  u -> 'User' as user_info
from
  okta_device as d,
  jsonb_array_elements(d.embedded -> 'Users') as u;
```

```sql+sqlite
select
  d.id,
  d.display_name,
  d.status,
  json_extract(u.value, '$.Created') as user_created,
  json_extract(u.value, '$.ManagementStatus') as user_management_status,
  json_extract(u.value, '$.ScreenLockType') as screen_lock_type,
  json_extract(u.value, '$.User') as user_info
from
  okta_device d,
  json_each(json_extract(d.embedded, '$.Users')) as u;
```

### Get device profile details
Explore the valuable insights that can enhance device security management, compliance monitoring, and operational decision-making within an organization.

```sql+postgres
select
  display_name,
  id,
  profile ->> 'DiskEncryptionType' as disk_encryption_type,
  profile ->> 'DisplayName' as display_name,
  profile ->> 'Imei' as imei,
  profile ->> 'IntegrityJailbreak' as integrity_jailbreak,
  profile ->> 'Manufacturer' as manufacturer,
  profile ->> 'Meid' as meid,
  profile ->> 'Model' as model,
  profile ->> 'OsVersion' as os_version,
  profile ->> 'Platform' as platform,
  profile ->> 'Registered' as registered,
  profile ->> 'SecureHardwarePresent' as secure_hardware_present,
  profile ->> 'SerialNumber' as serial_number,
  profile ->> 'Sid' as sid,
  profile ->> 'TpmPublicKeyHash' as tpm_public_key_hash
from
  okta_device;
```

```sql+sqlite
select
  display_name,
  id,
  json_extract(profile, '$.DiskEncryptionType') as disk_encryption_type,
  json_extract(profile, '$.DisplayName') as display_name,
  json_extract(profile, '$.Imei') as imei,
  json_extract(profile, '$.IntegrityJailbreak') as integrity_jailbreak,
  json_extract(profile, '$.Manufacturer') as manufacturer,
  json_extract(profile, '$.Meid') as meid,
  json_extract(profile, '$.Model') as model,
  json_extract(profile, '$.OsVersion') as os_version,
  json_extract(profile, '$.Platform') as platform,
  json_extract(profile, '$.Registered') as registered,
  json_extract(profile, '$.SecureHardwarePresent') as secure_hardware_present,
  json_extract(profile, '$.SerialNumber') as serial_number,
  json_extract(profile, '$.Sid') as sid,
  json_extract(profile, '$.TpmPublicKeyHash') as tpm_public_key_hash
from
  okta_device;
```