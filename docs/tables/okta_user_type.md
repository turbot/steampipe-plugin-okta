---
title: "Steampipe Table: okta_user_type - Query Okta User Types using SQL"
description: "Allows users to query Okta User Types, specifically providing access to user type properties and metadata for better user management and classification."
---

# Table: okta_user_type - Query Okta User Types using SQL

Okta User Types represent the different categories of users within an Okta organization. They allow for variation in the properties defined for users of different types, thereby enabling more granular control over user access and permissions. User Types are integral to Okta's identity and access management service, providing a flexible and customizable framework for user management.

## Table Usage Guide

The `okta_user_type` table provides insights into User Types within Okta's identity and access management service. As an IT administrator or security analyst, explore user type-specific details through this table, including properties, created and last updated timestamps, and associated metadata. Utilize it to uncover information about user types, such as their specific properties, the time of their creation and last update, and any other relevant details.

## Examples

### Basic info
Explore the different types of users within your Okta environment. This can help you understand how your system is structured and who has access to what, providing valuable insights for security and access management.
**Note:** _default_ is a reserved word and has to be double-quoted when used as identifier.


```sql+postgres
select
  name,
  id,
  "default",
  description,
  created,
  created_by
from
  okta_user_type;
```

```sql+sqlite
select
  name,
  id,
  "default",
  description,
  created,
  created_by
from
  okta_user_type;
```