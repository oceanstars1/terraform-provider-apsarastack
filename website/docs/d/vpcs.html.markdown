---
subcategory: "VPC"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_vpcs"
sidebar_current: "docs-apsarastack-datasource-vpcs"
description: |-
    Provides a list of VPCs owned by an ApsaraStack Cloud account.
---

# apsarastack\_vpcs

This data source provides VPCs available to the user.

## Example Usage

```
data "apsarastack_vpcs" "vpcs_ds" {
  cidr_block = "172.16.0.0/12"
  status     = "Available"
  name_regex = "^foo"
}

output "first_vpc_id" {
  value = "${data.apsarastack_vpcs.vpcs_ds.vpcs.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Optional) Filter results by a specific CIDR block. For example: "172.16.0.0/12".
* `status` - (Optional) Filter results by a specific status. Valid value are `Pending` and `Available`.
* `name_regex` - (Optional) A regex string to filter VPCs by name.
* `is_default` - (Optional, type: bool) Indicate whether the VPC is the default one in the specified region.
* `vswitch_id` - (Optional) Filter results by the specified VSwitch.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `ids` - (Optional) A list of VPC IDs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of VPC IDs.
* `names` - A list of VPC names.
* `vpcs` - A list of VPCs. Each element contains the following attributes:
  * `id` - ID of the VPC.
  * `region_id` - ID of the region where the VPC is located.
  * `status` - Status of the VPC.
  * `vpc_name` - Name of the VPC.
  * `vswitch_ids` - List of VSwitch IDs in the specified VPC
  * `cidr_block` - CIDR block of the VPC.
  * `vrouter_id` - ID of the VRouter.
  * `route_table_id` - Route table ID of the VRouter.
  * `description` - Description of the VPC
  * `is_default` - Whether the VPC is the default VPC in the region.
  * `creation_time` - Time of creation.
  * `secondary_cidr_blocks` - The secondary IPv4 CIDR block of the VPC.


 
