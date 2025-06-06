// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.94.1-71478489-20240820-161623
 */

package codeengine

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIbmCodeEngineFunction() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmCodeEngineFunctionRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of your function.",
			},
			"code_binary": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Specifies whether the code is binary or not. Defaults to false when `code_reference` is set to a data URL. When `code_reference` is set to a code bundle URL, this field is always true.",
			},
			"code_main": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the name of the function that should be invoked.",
			},
			"code_reference": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies either a reference to a code bundle or the source code itself. To specify the source code, use the data URL scheme and include the source code as base64 encoded. The data URL scheme is defined in [RFC 2397](https://tools.ietf.org/html/rfc2397).",
			},
			"code_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the secret that is used to access the specified `code_reference`. The secret is used to authenticate with a non-public endpoint that is specified as`code_reference`.",
			},
			"computed_env_variables": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "References to config maps, secrets or literal values, which are defined and set by Code Engine and are exposed as environment variables in the function.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The key to reference as environment variable.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the environment variable.",
						},
						"prefix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A prefix that can be added to all keys of a full secret or config map reference.",
						},
						"reference": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the secret or config map.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specify the type of the environment variable.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The literal value of the environment variable.",
						},
					},
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the resource was created.",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL to invoke the function.",
			},
			"endpoint_internal": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL to function that is only visible within the project.",
			},
			"entity_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the function instance, which is used to achieve optimistic locking.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When you provision a new function, a relative URL path is created identifying the location of the instance.",
			},
			"function_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier of the resource.",
			},
			"managed_domain_mappings": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional value controlling which of the system managed domain mappings will be setup for the function. Valid values are 'local_public', 'local_private' and 'local'. Visibility can only be 'local_private' if the project supports function private visibility.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region of the project the resource is located in. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de', 'eu-gb', 'jp-osa', 'jp-tok', 'us-east', 'us-south'.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the function.",
			},
			"run_env_variables": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "References to config maps, secrets or literal values, which are defined by the function owner and are exposed as environment variables in the function.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The key to reference as environment variable.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the environment variable.",
						},
						"prefix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A prefix that can be added to all keys of a full secret or config map reference.",
						},
						"reference": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the secret or config map.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specify the type of the environment variable.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The literal value of the environment variable.",
						},
					},
				},
			},
			"runtime": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The managed runtime used to execute the injected code.",
			},
			"scale_concurrency": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of parallel requests handled by a single instance, supported only by Node.js, default is `1`.",
			},
			"scale_cpu_limit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional amount of CPU set for the instance of the function. For valid values see [Supported memory and CPU combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo).",
			},
			"scale_down_delay": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Optional amount of time in seconds that delays the scale down behavior for a function.",
			},
			"scale_max_execution_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timeout in secs after which the function is terminated.",
			},
			"scale_memory_limit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional amount of memory set for the instance of the function. For valid values see [Supported memory and CPU combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo). The units for specifying memory are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the function.",
			},
			"status_details": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The detailed status of the function.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Provides additional information about the status of the function.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmCodeEngineFunctionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_function", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getFunctionOptions := &codeenginev2.GetFunctionOptions{}

	getFunctionOptions.SetProjectID(d.Get("project_id").(string))
	getFunctionOptions.SetName(d.Get("name").(string))

	function, _, err := codeEngineClient.GetFunctionWithContext(context, getFunctionOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetFunctionWithContext failed: %s", err.Error()), "(Data) ibm_code_engine_function", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *getFunctionOptions.ProjectID, *getFunctionOptions.Name))

	if err = d.Set("code_binary", function.CodeBinary); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting code_binary: %s", err), "(Data) ibm_code_engine_function", "read", "set-code_binary").GetDiag()
	}

	if !core.IsNil(function.CodeMain) {
		if err = d.Set("code_main", function.CodeMain); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting code_main: %s", err), "(Data) ibm_code_engine_function", "read", "set-code_main").GetDiag()
		}
	}

	if err = d.Set("code_reference", function.CodeReference); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting code_reference: %s", err), "(Data) ibm_code_engine_function", "read", "set-code_reference").GetDiag()
	}

	if !core.IsNil(function.CodeSecret) {
		if err = d.Set("code_secret", function.CodeSecret); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting code_secret: %s", err), "(Data) ibm_code_engine_function", "read", "set-code_secret").GetDiag()
		}
	}

	if !core.IsNil(function.ComputedEnvVariables) {
		computedEnvVariables := []map[string]interface{}{}
		for _, computedEnvVariablesItem := range function.ComputedEnvVariables {
			computedEnvVariablesItemMap, err := DataSourceIbmCodeEngineFunctionEnvVarToMap(&computedEnvVariablesItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_function", "read", "computed_env_variables-to-map").GetDiag()
			}
			computedEnvVariables = append(computedEnvVariables, computedEnvVariablesItemMap)
		}
		if err = d.Set("computed_env_variables", computedEnvVariables); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting computed_env_variables: %s", err), "(Data) ibm_code_engine_function", "read", "set-computed_env_variables").GetDiag()
		}
	}

	if !core.IsNil(function.CreatedAt) {
		if err = d.Set("created_at", function.CreatedAt); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_code_engine_function", "read", "set-created_at").GetDiag()
		}
	}

	if !core.IsNil(function.Endpoint) {
		if err = d.Set("endpoint", function.Endpoint); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting endpoint: %s", err), "(Data) ibm_code_engine_function", "read", "set-endpoint").GetDiag()
		}
	}

	if !core.IsNil(function.EndpointInternal) {
		if err = d.Set("endpoint_internal", function.EndpointInternal); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting endpoint_internal: %s", err), "(Data) ibm_code_engine_function", "read", "set-endpoint_internal").GetDiag()
		}
	}

	if err = d.Set("entity_tag", function.EntityTag); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting entity_tag: %s", err), "(Data) ibm_code_engine_function", "read", "set-entity_tag").GetDiag()
	}

	if !core.IsNil(function.Href) {
		if err = d.Set("href", function.Href); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_code_engine_function", "read", "set-href").GetDiag()
		}
	}

	if !core.IsNil(function.ID) {
		if err = d.Set("function_id", function.ID); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting function_id: %s", err), "(Data) ibm_code_engine_function", "read", "set-function_id").GetDiag()
		}
	}

	if err = d.Set("managed_domain_mappings", function.ManagedDomainMappings); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting managed_domain_mappings: %s", err), "(Data) ibm_code_engine_function", "read", "set-managed_domain_mappings").GetDiag()
	}

	if !core.IsNil(function.Region) {
		if err = d.Set("region", function.Region); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting region: %s", err), "(Data) ibm_code_engine_function", "read", "set-region").GetDiag()
		}
	}

	if !core.IsNil(function.ResourceType) {
		if err = d.Set("resource_type", function.ResourceType); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_code_engine_function", "read", "set-resource_type").GetDiag()
		}
	}

	runEnvVariables := []map[string]interface{}{}
	for _, runEnvVariablesItem := range function.RunEnvVariables {
		runEnvVariablesItemMap, err := DataSourceIbmCodeEngineFunctionEnvVarToMap(&runEnvVariablesItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_function", "read", "run_env_variables-to-map").GetDiag()
		}
		runEnvVariables = append(runEnvVariables, runEnvVariablesItemMap)
	}
	if err = d.Set("run_env_variables", runEnvVariables); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting run_env_variables: %s", err), "(Data) ibm_code_engine_function", "read", "set-run_env_variables").GetDiag()
	}

	if err = d.Set("runtime", function.Runtime); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting runtime: %s", err), "(Data) ibm_code_engine_function", "read", "set-runtime").GetDiag()
	}

	if err = d.Set("scale_concurrency", flex.IntValue(function.ScaleConcurrency)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting scale_concurrency: %s", err), "(Data) ibm_code_engine_function", "read", "set-scale_concurrency").GetDiag()
	}

	if err = d.Set("scale_cpu_limit", function.ScaleCpuLimit); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting scale_cpu_limit: %s", err), "(Data) ibm_code_engine_function", "read", "set-scale_cpu_limit").GetDiag()
	}

	if err = d.Set("scale_down_delay", flex.IntValue(function.ScaleDownDelay)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting scale_down_delay: %s", err), "(Data) ibm_code_engine_function", "read", "set-scale_down_delay").GetDiag()
	}

	if err = d.Set("scale_max_execution_time", flex.IntValue(function.ScaleMaxExecutionTime)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting scale_max_execution_time: %s", err), "(Data) ibm_code_engine_function", "read", "set-scale_max_execution_time").GetDiag()
	}

	if err = d.Set("scale_memory_limit", function.ScaleMemoryLimit); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting scale_memory_limit: %s", err), "(Data) ibm_code_engine_function", "read", "set-scale_memory_limit").GetDiag()
	}

	if !core.IsNil(function.Status) {
		if err = d.Set("status", function.Status); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status: %s", err), "(Data) ibm_code_engine_function", "read", "set-status").GetDiag()
		}
	}

	statusDetails := []map[string]interface{}{}
	statusDetailsMap, err := DataSourceIbmCodeEngineFunctionFunctionStatusToMap(function.StatusDetails)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_function", "read", "status_details-to-map").GetDiag()
	}
	statusDetails = append(statusDetails, statusDetailsMap)
	if err = d.Set("status_details", statusDetails); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status_details: %s", err), "(Data) ibm_code_engine_function", "read", "set-status_details").GetDiag()
	}

	return nil
}

func DataSourceIbmCodeEngineFunctionEnvVarToMap(model *codeenginev2.EnvVar) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Prefix != nil {
		modelMap["prefix"] = *model.Prefix
	}
	if model.Reference != nil {
		modelMap["reference"] = *model.Reference
	}
	modelMap["type"] = *model.Type
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func DataSourceIbmCodeEngineFunctionFunctionStatusToMap(model *codeenginev2.FunctionStatus) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Reason != nil {
		modelMap["reason"] = *model.Reason
	}
	return modelMap, nil
}
