// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMISInstanceGroupManagers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISInstanceGroupManagersRead,

		Schema: map[string]*schema.Schema{

			"instance_group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance group ID",
			},

			"instance_group_managers": {
				Type:        schema.TypeList,
				Description: "List of instance group managers",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the instance group manager.",
						},

						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the instance group manager.",
						},

						"manager_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of instance group manager.",
						},

						"aggregation_window": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The time window in seconds to aggregate metrics prior to evaluation",
						},

						"cooldown": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The duration of time in seconds to pause further scale actions after scaling has taken place",
						},

						"max_membership_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum number of members in a managed instance group",
						},

						"min_membership_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum number of members in a managed instance group",
						},

						"manager_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of instance group manager.",
						},

						"policies": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "list of Policies associated with instancegroup manager",
						},

						"actions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_group_manager_action": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_group_manager_action_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resource_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISInstanceGroupManagersRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_instance_group_managers", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	instanceGroupID := d.Get("instance_group").(string)

	// Support for pagination
	start := ""
	allrecs := []vpcv1.InstanceGroupManagerIntf{}

	for {
		listInstanceGroupManagerOptions := vpcv1.ListInstanceGroupManagersOptions{
			InstanceGroupID: &instanceGroupID,
		}
		if start != "" {
			listInstanceGroupManagerOptions.Start = &start
		}
		instanceGroupManagerCollections, _, err := sess.ListInstanceGroupManagersWithContext(context, &listInstanceGroupManagerOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListInstanceGroupManagersWithContext failed %s", err), "(Data) ibm_is_instance_group_managers", "read")
			log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

		start = flex.GetNext(instanceGroupManagerCollections.Next)
		allrecs = append(allrecs, instanceGroupManagerCollections.Managers...)

		if start == "" {
			break
		}

	}

	instanceGroupMnagers := make([]map[string]interface{}, 0)
	for _, instanceGroupManagerIntf := range allrecs {
		instanceGroupManager := instanceGroupManagerIntf.(*vpcv1.InstanceGroupManager)

		if *instanceGroupManager.ManagerType == "scheduled" {
			manager := map[string]interface{}{
				"id":           fmt.Sprintf("%s/%s", instanceGroupID, *instanceGroupManager.ID),
				"manager_id":   *instanceGroupManager.ID,
				"name":         *instanceGroupManager.Name,
				"manager_type": *instanceGroupManager.ManagerType,
			}

			actions := make([]map[string]interface{}, 0)
			if instanceGroupManager.Actions != nil {
				for _, action := range instanceGroupManager.Actions {
					actn := map[string]interface{}{
						"instance_group_manager_action":      action.ID,
						"instance_group_manager_action_name": action.Name,
						"resource_type":                      action.ResourceType,
					}
					actions = append(actions, actn)
				}
				manager["actions"] = actions
			}
			instanceGroupMnagers = append(instanceGroupMnagers, manager)

		} else {
			manager := map[string]interface{}{
				"id":                   fmt.Sprintf("%s/%s", instanceGroupID, *instanceGroupManager.ID),
				"manager_id":           *instanceGroupManager.ID,
				"name":                 *instanceGroupManager.Name,
				"aggregation_window":   *instanceGroupManager.AggregationWindow,
				"cooldown":             *instanceGroupManager.Cooldown,
				"max_membership_count": *instanceGroupManager.MaxMembershipCount,
				"min_membership_count": *instanceGroupManager.MinMembershipCount,
				"manager_type":         *instanceGroupManager.ManagerType,
			}

			policies := make([]string, 0)
			if instanceGroupManager.Policies != nil {
				for i := 0; i < len(instanceGroupManager.Policies); i++ {
					policies = append(policies, string(*(instanceGroupManager.Policies[i].ID)))
				}
			}
			manager["policies"] = policies
			instanceGroupMnagers = append(instanceGroupMnagers, manager)
		}

	}
	if err = d.Set("instance_group_managers", instanceGroupMnagers); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting instance_group_managers %s", err), "(Data) ibm_is_instance_group_managers", "read", "instance_group_managers-set").GetDiag()
	}
	d.SetId(dataSourceIBMISInstanceGroupManagersID(d))
	return nil
}

// dataSourceIBMISInstanceGroupManagersID returns a reasonable ID for a instance group manager list.
func dataSourceIBMISInstanceGroupManagersID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
