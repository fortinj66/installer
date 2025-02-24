// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func dataSourceIBMIsInstanceNetworkInterface() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsInstanceNetworkInterfaceRead,

		Schema: map[string]*schema.Schema{
			"instance_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The instance name.",
			},
			"network_interface_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network interface name.",
			},
			"allow_ip_spoofing": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether source IP spoofing is allowed on this interface. If false, source IP spoofing is prevented on this interface. If true, source IP spoofing is allowed on this interface.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the network interface was created.",
			},
			"floating_ips": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The floating IPs associated with this network interface.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique IP address.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this floating IP.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this floating IP.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this floating IP.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this floating IP.",
						},
					},
				},
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this network interface.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user-defined name for this network interface.",
			},
			"port_speed": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The network interface port speed in Mbps.",
			},
			"primary_ipv4_address": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The primary IPv4 address.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"security_groups": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of security groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The security group's CRN.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The security group's canonical URL.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this security group.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this security group. Names must be unique within the VPC the security group resides in.",
						},
					},
				},
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the network interface.",
			},
			"subnet": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The associated subnet.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this subnet.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this subnet.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this subnet.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this subnet.",
						},
					},
				},
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of this network interface as it relates to an instance.",
			},
		},
	}
}

func dataSourceIBMIsInstanceNetworkInterfaceRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	instance_name := d.Get("instance_name").(string)
	listInstancesOptions := &vpcv1.ListInstancesOptions{}

	start := ""
	allrecs := []vpcv1.Instance{}
	for {

		if start != "" {
			listInstancesOptions.Start = &start
		}

		instances, response, err := vpcClient.ListInstancesWithContext(context, listInstancesOptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error Fetching Instances %s\n%s", err, response))
		}
		start = GetNext(instances.Next)
		allrecs = append(allrecs, instances.Instances...)
		if start == "" {
			break
		}
	}

	ins_id := ""
	for _, instance := range allrecs {
		if *instance.Name == instance_name {
			ins_id = *instance.ID
			listInstanceNetworkInterfacesOptions := &vpcv1.ListInstanceNetworkInterfacesOptions{
				InstanceID: &ins_id,
			}
			networkInterfaceCollection, response, err := vpcClient.ListInstanceNetworkInterfacesWithContext(context, listInstanceNetworkInterfacesOptions)

			if err != nil {
				log.Printf("[DEBUG] ListSecurityGroupNetworkInterfacesWithContext failed %s\n%s", err, response)
				return diag.FromErr(fmt.Errorf("ListSecurityGroupNetworkInterfacesWithContext failed %s\n%s", err, response))
			}
			network_interface_name := d.Get("network_interface_name").(string)
			for _, networkInterface := range networkInterfaceCollection.NetworkInterfaces {
				if *networkInterface.Name == network_interface_name {
					d.SetId(fmt.Sprintf("%s/%s", ins_id, *networkInterface.ID))
					if err = d.Set("allow_ip_spoofing", networkInterface.AllowIPSpoofing); err != nil {
						return diag.FromErr(fmt.Errorf("Error setting allow_ip_spoofing: %s", err))
					}
					if err = d.Set("created_at", dateTimeToString(networkInterface.CreatedAt)); err != nil {
						return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
					}

					if networkInterface.FloatingIps != nil {
						err = d.Set("floating_ips", dataSourceNetworkInterfaceFlattenFloatingIps(networkInterface.FloatingIps))
						if err != nil {
							return diag.FromErr(fmt.Errorf("Error setting floating_ips %s", err))
						}
					}
					if err = d.Set("href", networkInterface.Href); err != nil {
						return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
					}
					if err = d.Set("name", networkInterface.Name); err != nil {
						return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
					}
					if err = d.Set("port_speed", intValue(networkInterface.PortSpeed)); err != nil {
						return diag.FromErr(fmt.Errorf("Error setting port_speed: %s", err))
					}
					if err = d.Set("primary_ipv4_address", networkInterface.PrimaryIpv4Address); err != nil {
						return diag.FromErr(fmt.Errorf("Error setting primary_ipv4_address: %s", err))
					}
					if err = d.Set("resource_type", networkInterface.ResourceType); err != nil {
						return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
					}

					if networkInterface.SecurityGroups != nil {
						err = d.Set("security_groups", dataSourceNetworkInterfaceFlattenSecurityGroups(networkInterface.SecurityGroups))
						if err != nil {
							return diag.FromErr(fmt.Errorf("Error setting security_groups %s", err))
						}
					}
					if err = d.Set("status", networkInterface.Status); err != nil {
						return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
					}

					if networkInterface.Subnet != nil {
						err = d.Set("subnet", dataSourceNetworkInterfaceFlattenSubnet(*networkInterface.Subnet))
						if err != nil {
							return diag.FromErr(fmt.Errorf("Error setting subnet %s", err))
						}
					}
					if err = d.Set("type", networkInterface.Type); err != nil {
						return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
					}
					return nil
				}
			}
			return diag.FromErr(fmt.Errorf("Network interface %s not found.", network_interface_name))
		}
	}

	return diag.FromErr(fmt.Errorf("Instance %s not found. %s", instance_name, err))
}

func dataSourceNetworkInterfaceFlattenFloatingIps(result []vpcv1.FloatingIPReference) (floatingIps []map[string]interface{}) {
	for _, floatingIpsItem := range result {
		floatingIps = append(floatingIps, dataSourceNetworkInterfaceFloatingIpsToMap(floatingIpsItem))
	}

	return floatingIps
}

func dataSourceNetworkInterfaceFloatingIpsToMap(floatingIpsItem vpcv1.FloatingIPReference) (floatingIpsMap map[string]interface{}) {
	floatingIpsMap = map[string]interface{}{}

	if floatingIpsItem.Address != nil {
		floatingIpsMap["address"] = floatingIpsItem.Address
	}
	if floatingIpsItem.CRN != nil {
		floatingIpsMap["crn"] = floatingIpsItem.CRN
	}
	if floatingIpsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceNetworkInterfaceFloatingIpsDeletedToMap(*floatingIpsItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		floatingIpsMap["deleted"] = deletedList
	}
	if floatingIpsItem.Href != nil {
		floatingIpsMap["href"] = floatingIpsItem.Href
	}
	if floatingIpsItem.ID != nil {
		floatingIpsMap["id"] = floatingIpsItem.ID
	}
	if floatingIpsItem.Name != nil {
		floatingIpsMap["name"] = floatingIpsItem.Name
	}

	return floatingIpsMap
}

func dataSourceNetworkInterfaceFloatingIpsDeletedToMap(deletedItem vpcv1.FloatingIPReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceNetworkInterfaceFlattenSecurityGroups(result []vpcv1.SecurityGroupReference) (securityGroups []map[string]interface{}) {
	for _, securityGroupsItem := range result {
		securityGroups = append(securityGroups, dataSourceNetworkInterfaceSecurityGroupsToMap(securityGroupsItem))
	}

	return securityGroups
}

func dataSourceNetworkInterfaceSecurityGroupsToMap(securityGroupsItem vpcv1.SecurityGroupReference) (securityGroupsMap map[string]interface{}) {
	securityGroupsMap = map[string]interface{}{}

	if securityGroupsItem.CRN != nil {
		securityGroupsMap["crn"] = securityGroupsItem.CRN
	}
	if securityGroupsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceNetworkInterfaceSecurityGroupsDeletedToMap(*securityGroupsItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		securityGroupsMap["deleted"] = deletedList
	}
	if securityGroupsItem.Href != nil {
		securityGroupsMap["href"] = securityGroupsItem.Href
	}
	if securityGroupsItem.ID != nil {
		securityGroupsMap["id"] = securityGroupsItem.ID
	}
	if securityGroupsItem.Name != nil {
		securityGroupsMap["name"] = securityGroupsItem.Name
	}

	return securityGroupsMap
}

func dataSourceNetworkInterfaceSecurityGroupsDeletedToMap(deletedItem vpcv1.SecurityGroupReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceNetworkInterfaceFlattenSubnet(result vpcv1.SubnetReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceNetworkInterfaceSubnetToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceNetworkInterfaceSubnetToMap(subnetItem vpcv1.SubnetReference) (subnetMap map[string]interface{}) {
	subnetMap = map[string]interface{}{}

	if subnetItem.CRN != nil {
		subnetMap["crn"] = subnetItem.CRN
	}
	if subnetItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceNetworkInterfaceSubnetDeletedToMap(*subnetItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		subnetMap["deleted"] = deletedList
	}
	if subnetItem.Href != nil {
		subnetMap["href"] = subnetItem.Href
	}
	if subnetItem.ID != nil {
		subnetMap["id"] = subnetItem.ID
	}
	if subnetItem.Name != nil {
		subnetMap["name"] = subnetItem.Name
	}

	return subnetMap
}

func dataSourceNetworkInterfaceSubnetDeletedToMap(deletedItem vpcv1.SubnetReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
