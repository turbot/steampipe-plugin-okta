---
title: "Steampipe Table: okta_auth_server - Query Okta Authorization Servers using SQL"
description: "Allows users to query Okta Authorization Servers, providing insights into the authorization server configurations and policies."
---

# Table: okta_auth_server - Query Okta Authorization Servers using SQL

An Okta Authorization Server is a component within Okta that provides developers with different types of security tokens such as JSON Web Tokens (JWT) and access tokens. It is responsible for maintaining a set of resources, defining access policies, and performing token generation and distribution. The server is critical for managing and securing access to web applications and APIs.

## Table Usage Guide

The `okta_auth_server` table provides insights into the configuration and policies of Okta Authorization Servers. As a security engineer or developer, explore server-specific details through this table, including the server's name, audience, issuer mode, and associated metadata. Utilize it to uncover information about the authorization servers, such as the server's status, the creation and last modified dates, and the verification of access policies.

## Examples

### Basic info
Explore the status and update history of your authentication servers. This can be useful to track changes over time and ensure all servers are functioning as expected.

```sql
select
  name,
  id,
  audiences,
  created,
  last_updated,
  status
from
  okta_auth_server;
```

### List authorization servers where manual rotation signing keys are not rotated in more than 90 days
Determine the areas in which authorization servers have not had their manual rotation signing keys rotated in more than 90 days. This is useful for maintaining security standards and ensuring regular key rotation.

```sql
select
  name,
  id,
  audiences,
  created,
  last_updated, 
  credentials -> 'signing' ->> 'lastRotated' as last_rotated,
  status
from
  okta_auth_server
where
  credentials -> 'signing' ->> 'rotationMode' = 'MANUAL' 
  and CAST(credentials -> 'signing' ->> 'lastRotated' as date) < current_timestamp - interval '90 days';
```

### List inactive authorization servers
Analyze the settings to understand which authorization servers are currently inactive. This is useful for maintaining server efficiency and ensuring all resources are optimally utilized.

```sql
select
  name,
  id,
  audiences,
  created,
  last_updated,
  status
from
  okta_auth_server
where
  status = 'INACTIVE';
```