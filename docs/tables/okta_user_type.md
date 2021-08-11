# Table: okta_user_type

A user_type is a record of information stored in Okta Universal Directory that contains specific user attributes such as the user's name and phone number, location, and role.

## Examples

### Basic info

**Note:** _default_ is a reserved word and has to be double-quoted when used as identifier.

```sql
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
