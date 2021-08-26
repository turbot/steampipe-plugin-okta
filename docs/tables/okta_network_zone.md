# Table: okta_network_zone

The Okta Zones provides operations to manage Zones in your organization. There are two usage Zone types: Policy Network Zones and Block List Network Zones. Policy Network Zones are used to guide policy decisions. Block List Network Zones are used to deny access from certain IP addresses, locations, proxy types, or Autonomous System Numbers (ASNs) before policy evaluation.

## Examples

### Basic info

```sql
select
  name,
  id,
  created,
  status,
  system,
  locations
from
  okta_network_zone
order by
  priority;
```

### List TorAnonymizer proxy type network zone

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