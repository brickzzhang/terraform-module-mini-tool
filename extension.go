package main

var templateStr = `
# TencentCloud [MODIFY !!!] Module for Terraform

## terraform-tencentcloud-vpc

A terraform module used to create TencentCloud VPC, subnet and route entry[MODIFY !!!].

The following resources are included.

* [MODIFY !!!][VPC](https://www.terraform.io/docs/providers/tencentcloud/r/vpc.html)
* [MODIFY !!!][VPC Subnet](https://www.terraform.io/docs/providers/tencentcloud/r/subnet.html)
* [MODIFY !!!][VPC Route Entry](https://www.terraform.io/docs/providers/tencentcloud/r/route_table_entry.html)

## Usage


## Conditional Creation

This module can create VPC and VPC Subnet[MODIFY !!!].

## Inputs

%s

## Outputs

%s

## Authors

Created and maintained by [TencentCloud](https://github.com/terraform-providers/terraform-provider-tencentcloud)

## License

Mozilla Public License Version 2.0.
See LICENSE for full details.
`
