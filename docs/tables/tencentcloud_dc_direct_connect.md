---
title: "Steampipe Table: tencentcloud_dc_direct_connect - Query Tencent Cloud Direct Connects using SQL"
description: "Allows users to query Direct Connect physical connections between an on-premises data center and Tencent Cloud."
folder: "DC"
---

# Table: tencentcloud_dc_direct_connect - Query Tencent Cloud Direct Connects using SQL

Tencent Cloud Direct Connect (DC) provides dedicated physical links between an on-premises data center and Tencent Cloud. Each connection records the ISP, access point, port type, VLAN, bandwidth, and billing state, along with the Tencent-side and customer-side IP addresses used for debugging. DC is the foundation for dedicated tunnels and Direct Connect gateways.

## Table Usage Guide

The `tencentcloud_dc_direct_connect` table provides you with detailed information about physical Direct Connect links. Use this table as a network administrator or DevOps engineer to inventory connections across regions, audit ISP and access point assignments, inspect bandwidth allocations, verify VLAN configuration, and review tag-based classifications.

## Examples

### Basic info

List every Direct Connect and its core attributes.

```sql+postgres
select
  direct_connect_id,
  direct_connect_name,
  state,
  line_operator,
  access_point_id,
  bandwidth,
  vlan,
  region
from
  tencentcloud_dc_direct_connect;
```

```sql+sqlite
select
  direct_connect_id,
  direct_connect_name,
  state,
  line_operator,
  access_point_id,
  bandwidth,
  vlan,
  region
from
  tencentcloud_dc_direct_connect;
```

### Filter Direct Connects by state

Identify Direct Connects that are currently available.

```sql+postgres
select
  direct_connect_id,
  direct_connect_name,
  line_operator,
  bandwidth,
  access_point_id,
  region
from
  tencentcloud_dc_direct_connect
where
  state = 'AVAILABLE';
```

```sql+sqlite
select
  direct_connect_id,
  direct_connect_name,
  line_operator,
  bandwidth,
  access_point_id,
  region
from
  tencentcloud_dc_direct_connect
where
  state = 'AVAILABLE';
```

### Count Direct Connects by ISP

Summarize how many Direct Connects are provided by each ISP.

```sql+postgres
select
  line_operator,
  count(*) as count
from
  tencentcloud_dc_direct_connect
group by
  line_operator
order by
  count desc;
```

```sql+sqlite
select
  line_operator,
  count(*) as count
from
  tencentcloud_dc_direct_connect
group by
  line_operator
order by
  count desc;
```

### Inspect Direct Connects by region and bandwidth

Find Direct Connects with bandwidth at or above a threshold, grouped by region.

```sql+postgres
select
  direct_connect_id,
  direct_connect_name,
  bandwidth,
  port_type,
  region
from
  tencentcloud_dc_direct_connect
where
  bandwidth >= 1000
order by
  bandwidth desc;
```

```sql+sqlite
select
  direct_connect_id,
  direct_connect_name,
  bandwidth,
  port_type,
  region
from
  tencentcloud_dc_direct_connect
where
  bandwidth >= 1000
order by
  bandwidth desc;
```

### List Direct Connects by tag

Filter Direct Connects by a specific tag key-value pair.

```sql+postgres
select
  direct_connect_id,
  direct_connect_name,
  tags
from
  tencentcloud_dc_direct_connect
where
  tags ->> 'Environment' = 'Production';
```

```sql+sqlite
select
  direct_connect_id,
  direct_connect_name,
  tags
from
  tencentcloud_dc_direct_connect
where
  json_extract(tags, '$.Environment') = 'Production';
```
