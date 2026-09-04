---
title: "Steampipe Table: tencentcloud_cvm_key_pair - Query Tencent Cloud CVM Key Pairs using SQL"
description: "Allows users to query CVM key pairs in Tencent Cloud to retrieve information about SSH login credentials bound to instances."
folder: "CVM"
---

# Table: tencentcloud_cvm_key_pair - Query Tencent Cloud CVM Key Pairs using SQL

Tencent Cloud CVM key pairs are SSH credentials used to securely log in to Linux instances. A key pair consists of a public key (stored by Tencent Cloud) and a private key (kept only by the user). Key pairs can be bound to or unbound from instances, providing a more secure alternative to password-based login.

## Table Usage Guide

The `tencentcloud_cvm_key_pair` table in Steampipe provides you with detailed information about CVM key pairs within your Tencent Cloud account. This table allows you, as a DevOps engineer or cloud administrator, to query key pair details, including the key ID, name, project, public key, associated instances, creation time, and tags. You can utilize this table to audit SSH access, find key pairs not associated with any instance, and review tag-based classifications.

Note: The private key is never returned by the Tencent Cloud API after creation, so this table does not expose a private key column.

## Examples

### Basic info

Get an overview of all key pairs in your account, including their names, IDs, and associated instances. This is useful for SSH credential inventory management.

```sql+postgres
select
  key_id,
  key_name,
  project_id,
  associated_instance_ids,
  created_time
from
  tencentcloud_cvm_key_pair;
```

```sql+sqlite
select
  key_id,
  key_name,
  project_id,
  associated_instance_ids,
  created_time
from
  tencentcloud_cvm_key_pair;
```

### Find unused key pairs

Identify key pairs that are not associated with any instance. Unused credentials should be reviewed and removed to reduce the attack surface.

```sql+postgres
select
  key_id,
  key_name,
  created_time
from
  tencentcloud_cvm_key_pair
where
  associated_instance_ids is null
  or jsonb_array_length(associated_instance_ids) = 0;
```

```sql+sqlite
select
  key_id,
  key_name,
  created_time
from
  tencentcloud_cvm_key_pair
where
  associated_instance_ids is null
  or json_array_length(associated_instance_ids) = 0;
```

### Get a specific key pair

Retrieve details of a single key pair by its ID. The `key_id` qualifier is pushed down to the DescribeKeyPairs API.

```sql+postgres
select
  key_id,
  key_name,
  public_key,
  description
from
  tencentcloud_cvm_key_pair
where
  key_id = 'skey-11112222';
```

```sql+sqlite
select
  key_id,
  key_name,
  public_key,
  description
from
  tencentcloud_cvm_key_pair
where
  key_id = 'skey-11112222';
```

### Search key pairs by name

Filter key pairs by name. The filter is pushed down to the DescribeKeyPairs API, avoiding a full client-side scan.

```sql+postgres
select
  key_id,
  key_name,
  created_time
from
  tencentcloud_cvm_key_pair
where
  key_name = 'my-deploy-key';
```

```sql+sqlite
select
  key_id,
  key_name,
  created_time
from
  tencentcloud_cvm_key_pair
where
  key_name = 'my-deploy-key';
```

### List key pairs associated with a specific instance

Find which key pairs grant SSH access to a given instance, useful for access audits before decommissioning or rotating credentials.

```sql+postgres
select
  key_id,
  key_name,
  instance_id
from
  tencentcloud_cvm_key_pair,
  jsonb_array_elements_text(associated_instance_ids) as instance_id
where
  instance_id = 'ins-xxxxxxxx';
```

```sql+sqlite
select
  key_id,
  key_name,
  instance_id.value as instance_id
from
  tencentcloud_cvm_key_pair,
  json_each(associated_instance_ids) as instance_id
where
  instance_id.value = 'ins-xxxxxxxx';
```

### Find key pairs with specific tags

Query key pairs that have been tagged with a specific key-value pair, such as an environment classification.

```sql+postgres
select
  key_id,
  key_name,
  tags
from
  tencentcloud_cvm_key_pair
where
  tags ->> 'Environment' = 'Production';
```

```sql+sqlite
select
  key_id,
  key_name,
  tags
from
  tencentcloud_cvm_key_pair
where
  json_extract(tags, '$.Environment') = 'Production';
```
