# Table: okta_system_log

Okta System Log captures a comprehensive set of system events related to users, security policies, app integrations, and other activities within an Okta organization. These logs provide an audit trail that can be used for monitoring, troubleshooting, and compliance purposes. Events in the system log range from user sign-in attempts and profile updates to admin activities and system-level actions.

- By default, this table will provide data for the last 7 days. You can give the `log_event_time` value in the below ways to fetch data in a range. The examples below can guide you.

  - log_event_time >= ‘2023-03-11T00:00:00Z’ and log_event_time <= ‘2023-03-15T00:00:00Z’
  - log_event_time between ‘2023-03-11T00:00:00Z’ and ‘2023-03-15T00:00:00Z’
  - log_event_time >= now() - interval '30 days' (The data will be fetched from the last 30 days)
  - log_event_time > ‘2023-03-15T00:00:00Z’ (The data will be fetched from the provided time to the current time)
  - log_event_time < ‘2023-03-15T00:00:00Z’ (The data will be fetched from one day before the provided time to the provided time)

- This table supports optional quals. Queries with optional quals are optimised to use Okta [filters](https://developer.okta.com/docs/reference/api/system-log/#bounded-requests). Optional quals are supported for the following columns:
  - filter
  - log_actor_id
  - log_ip_address
  - log_event_type
  - log_event_time

## Examples

### Basic info

```sql
select
  log_actor_id,
  log_actor_name,
  log_ip_address,
  display_message,
  log_event_type,
  title
from
  okta_system_log;
```

### List logs which occur in last 30 days

```sql
select
  log_actor_name,
  log_ip_address,
  display_message,
  title
from
  okta_system_log
where
  log_event_time >= now() - interval '3 days';
```

### Show logs of a particular actor

```sql
select
  log_actor_id,
  log_actor_name,
  log_ip_address,
  display_message,
  log_event_type,
  title
from
  okta_system_log
where
  log_actor_name = 'sourav';
```

### Show error logs

```sql
select
  log_actor_id,
  log_actor_name,
  log_ip_address,
  display_message,
  log_event_type,
  severity
from
  okta_system_log
where
  severity = 'ERROR';
```

### List target details of the logs

```sql
select
  title,
  severity,
  t ->> 'alternateId' as target_alternate_id,
  t ->> 'displayName' as target_name,
  t ->> 'id' as target_id,
  t ->> 'type' as target_type,
  t -> 'detailEntry' as detail_entry
from
  okta_system_log,
  jsonb_array_elements(target) as t;
```

### List transaction details of the logs

```sql
select
  title,
  severity,
  transaction ->> 'id' as transaction_id,
  transaction ->> 'type' as transaction_type,
  transaction -> 'detail' as transaction_detail
from
  okta_system_log;
```

### List events that match the filter pattern term **eventType**

```sql
select
  log_actor_id,
  log_actor_name,
  log_ip_address,
  display_message,
  log_event_type,
  severity
from
  okta_system_log
where
  filter = 'eventType eq "user.session.start"';
```
