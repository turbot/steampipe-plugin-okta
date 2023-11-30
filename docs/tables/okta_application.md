---
title: "Steampipe Table: okta_application - Query Okta Applications using SQL"
description: "Allows users to query Okta Applications, specifically retrieving information about the applications configured within an Okta organization."
---

# Table: okta_application - Query Okta Applications using SQL

Okta Applications are integral components of the Okta identity management service. These applications represent the software applications that are linked to Okta for single sign-on, provisioning, or API access management. They enable seamless and secure access to all the applications your users need, from Microsoft Office 365 to custom applications built in-house.

## Table Usage Guide

The `okta_application` table provides insights into applications configured within an Okta organization. As a Security Analyst, explore application-specific details through this table, including application type, status, and associated metadata. Utilize it to uncover information about applications, such as those with specific accessibility, the users assigned to each application, and the verification of application settings.

## Examples

### Basic info
Explore which applications are currently active in your system by identifying their status and creation dates. This can help in managing the applications effectively and ensuring timely updates.

```sql
select
  name,
  id,
  label,
  created,
  status,
  sign_on_mode
from
  okta_application;
```

### List SAML 2.0 apps
Identify applications that use SAML 2.0 as their sign-on mode. This can be useful in understanding the security protocols of your applications.

```sql
select
  name,
  id,
  label,
  created,
  status,
  sign_on_mode
from
  okta_application
where
  sign_on_mode = 'SAML_2_0';
```

### List apps assigned to a specific user using a filter
Explore which applications are assigned to a specific user by filtering based on user ID. This is useful for understanding the scope of access and permissions granted to individual users within your system.

```sql
select
  id,
  label,
  name,
  sign_on_mode,
  status
from
  okta_application as app
where
  filter = 'user.id eq "00u1e5eizrjQKTWMA5d7"';
```

### List apps assigned to a specific group using a filter
Explore which applications are assigned to a specific group, useful for understanding application accessibility and managing group permissions. This can aid in maintaining security protocols and ensuring appropriate access rights.

```sql
select
  id,
  label,
  name,
  sign_on_mode,
  status
from
  okta_application
where
  filter = 'group.id eq "00u1e5eizrjQKTWMA5d7"';
```