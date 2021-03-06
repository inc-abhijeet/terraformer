// Copyright 2020 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aws

import (
	"context"
	"log"

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

var tgwAllowEmptyValues = []string{"tags."}

type TransitGatewayGenerator struct {
	AWSService
}

func (g *TransitGatewayGenerator) getTransitGateways(svc *ec2.Client) error {

	p := ec2.NewDescribeTransitGatewaysPaginator(svc.DescribeTransitGatewaysRequest(&ec2.DescribeTransitGatewaysInput{}))
	for p.Next(context.Background()) {
		for _, tgw := range p.CurrentPage().TransitGateways {
			g.Resources = append(g.Resources, terraform_utils.NewSimpleResource(
				aws.StringValue(tgw.TransitGatewayId),
				aws.StringValue(tgw.TransitGatewayId),
				"aws_ec2_transit_gateway",
				"aws",
				tgwAllowEmptyValues,
			))
		}
	}
	return p.Err()
}

func (g *TransitGatewayGenerator) getTransitGatewayRouteTables(svc *ec2.Client) error {

	p := ec2.NewDescribeTransitGatewayRouteTablesPaginator(svc.DescribeTransitGatewayRouteTablesRequest(&ec2.DescribeTransitGatewayRouteTablesInput{}))
	for p.Next(context.Background()) {
		for _, tgwrt := range p.CurrentPage().TransitGatewayRouteTables {
			// Default route table are automatically created on the tgw creation
			if *tgwrt.DefaultAssociationRouteTable {
				continue
			} else {
				g.Resources = append(g.Resources, terraform_utils.NewSimpleResource(
					aws.StringValue(tgwrt.TransitGatewayRouteTableId),
					aws.StringValue(tgwrt.TransitGatewayRouteTableId),
					"aws_ec2_transit_gateway_route_table",
					"aws",
					tgwAllowEmptyValues,
				))
			}
		}
	}
	return p.Err()
}

func (g *TransitGatewayGenerator) getTransitGatewayVpcAttachments(svc *ec2.Client) error {
	p := ec2.NewDescribeTransitGatewayVpcAttachmentsPaginator(svc.DescribeTransitGatewayVpcAttachmentsRequest(&ec2.DescribeTransitGatewayVpcAttachmentsInput{}))
	for p.Next(context.Background()) {
		for _, tgwa := range p.CurrentPage().TransitGatewayVpcAttachments {
			g.Resources = append(g.Resources, terraform_utils.NewSimpleResource(
				aws.StringValue(tgwa.TransitGatewayAttachmentId),
				aws.StringValue(tgwa.TransitGatewayAttachmentId),
				"aws_ec2_transit_gateway_vpc_attachment",
				"aws",
				tgwAllowEmptyValues,
			))
		}

	}
	return p.Err()
}

// Generate TerraformResources from AWS API,
// from each customer gateway create 1 TerraformResource.
// Need CustomerGatewayId as ID for terraform resource
func (g *TransitGatewayGenerator) InitResources() error {
	config, e := g.generateConfig()
	if e != nil {
		return e
	}
	svc := ec2.New(config)
	g.Resources = []terraform_utils.Resource{}
	err := g.getTransitGateways(svc)
	if err != nil {
		log.Println(err)
	}

	err = g.getTransitGatewayRouteTables(svc)
	if err != nil {
		log.Println(err)
	}

	err = g.getTransitGatewayVpcAttachments(svc)
	if err != nil {
		log.Println(err)
	}

	return nil
}
