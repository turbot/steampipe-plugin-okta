---
title: "Steampipe Table: okta_network_zone - Query Okta Network Zones using SQL"
description: "Allows users to query Network Zones in Okta, specifically providing detailed information about each network zone's type, status, system, and conditions."
---

# Table: okta_network_zone - Query Okta Network Zones using SQL

A Network Zone in Okta represents a set of IP address ranges. It can be used to control the security behavior of end users and applications based on their IP location. Network Zones help in defining trusted IP ranges, blocking suspicious IP ranges, and setting up behavior detection for each zone.

## Table Usage Guide

The `okta_network_zone` table provides insights into Network Zones within Okta. As a security administrator, explore zone-specific details through this table, including zone type, status, system, and conditions. Utilize it to manage access control based on IP location, identify trusted IP ranges, and set up behavior detection for each zone.

## Examples

### Basic info
Gain insights into the creation, status, and usage of various network zones in your Okta system. This information can be helpful in understanding the overall network configuration and assessing any potential security risks.

```sql
select
  name,
  id,
  created,
  status,
  system,
  locations,
  proxy_type,
  usage
from
  okta_network_zone;
```

### List TorAnonymizer proxy type network zone
Explore which network zones are set up as TorAnonymizer proxies. This can be particularly useful for identifying potential security risks or for auditing your network's configuration.

```sql
select
  name,
  id,
  created,
  status,
  system,
  proxy_type
from
  okta_network_zone
where
  proxy_type = 'TorAnonymizer';
```

### List network zones location and region details
Analyze your network zones to understand their geographical distribution. This is useful when you need to pinpoint specific locations for network management or security purposes.

```sql
select
  name,
  id,
  l -> 'country' as country,
  l -> 'region' as region
from
  okta_network_zone,
  jsonb_array_elements(locations) as l;
```

### List system network zones
Explore which network zones are system-generated in your Okta environment. This can help you understand and manage the security of your system.

```sql
select
  name,
  id,
  created,
  status,
  system
from
  okta_network_zone
where
  system;
```

### List active network zones
Explore the active network zones in your system, allowing you to understand the current operational areas for better management and security planning.

```sql
select
  name,
  id,
  created,
  status,
  system,
  proxy_type
from
  okta_network_zone
where
  status = 'ACTIVE';
```