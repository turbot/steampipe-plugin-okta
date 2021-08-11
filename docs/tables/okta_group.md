# Table: okta_group

A Group is made up of users. Groups are useful for representing roles, relationships, and can even be used for subscription tiers.

Note: This table supports optional `filter` column to query results based on okta supported [filters](https://developer.okta.com/docs/reference/api/groups/#filters).

## Examples

### Basic info

```sql
select
  name,
  id,
  type,
  description,
  jsonb_pretty(profile) as profile
from
  okta_group;
```

### list groups for which membership has not been changed for more than 30 days

```sql
select
  name,
  id,
  type,
  age(current_timestamp, last_membership_updated) as last_membership_updated
from
  okta_group
where
  last_membership_updated < current_timestamp - interval '30 days';
```

### Use filters to get Okta Groups with profile or memberships updated after 05/05/2021

```sql
select
  name,
  id,
  type,
  last_updated,
  last_membership_updated
from
  okta_group
where
  filter = 'type eq "OKTA_GROUP" and (lastUpdated gt "2021-05-05T00:00:00.000Z" or lastMembershipUpdated gt "2021-05-05T00:00:00.000Z")';
```

### Get group member details for groups

```sql
select
  name,
  id,
  jsonb_pretty(group_members) as group_members
from
  okta_group
```
