package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/scottyw/lyra-bridge/cmd/provider"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scottyw/lyra-bridge/pkg/generated"
	"github.com/terraform-providers/terraform-provider-aws/aws"
)

var config *terraform.ResourceConfig

func ptr(s string) *string {
	return &s
}

func init() {
	// Fields derived from Schema
	config = &terraform.ResourceConfig{
		Config: map[string]interface{}{
			"region": "eu-west-1",
		},
	}
}

func main() {

	// Configure the provider
	p := aws.Provider().(*schema.Provider)
	err := p.Configure(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create VPC
	vpcHandler := provider.Aws_vpcHandler{Provider: p}
	vpc := &generated.Aws_vpc{
		Cidr_block:       "192.168.0.0/16",
		Instance_tenancy: ptr("default"),
		Tags: &map[string]interface{}{
			"Name": "lyra-test",
		},
	}
	vpc, vid, err := vpcHandler.Create(vpc)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("CREATED VPC:", spew.Sdump(vpc))
	fmt.Println("CREATED VPC ID:", vid)

	// Read VPC
	vpc, err = vpcHandler.Read(vid)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("READ VPC:", spew.Sdump(vpc))

	// Delete VPC
	vpcHandler.Delete(vid)
}
