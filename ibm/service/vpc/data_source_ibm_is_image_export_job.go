// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsImageExport() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIBMIsImageExportRead,

		Schema: map[string]*schema.Schema{
			"image": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The image identifier.",
			},
			"image_export_job": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The image export job identifier.",
			},
			"completed_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the image export job was completed.If absent, the export job has not yet completed.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the image export job was created.",
			},
			"encrypted_data_key": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A base64-encoded, encrypted representation of the key that was used to encrypt the data for the exported image. This key can be unwrapped with the image's `encryption_key` root key using either Key Protect or Hyper Protect Crypto Service.If absent, the export job is for an unencrypted image.",
			},
			"format": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The format of the exported image.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this image export job.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user-defined name for this image export job.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of resource referenced.",
			},
			"started_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the image export job started running.If absent, the export job has not yet started.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of this image export job:- `deleting`: Export job is being deleted- `failed`: Export job could not be completed successfully- `queued`: Export job is queued- `running`: Export job is in progress- `succeeded`: Export job was completed successfullyThe exported image object is automatically deleted for `failed` jobs.",
			},
			"status_reasons": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current status (if any).The enumerated reason code values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected reason code was encountered.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the status reason.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the status reason.",
						},
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about this status reason.",
						},
					},
				},
			},
			"storage_bucket": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Cloud Object Storage bucket of the exported image object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN of this Cloud Object Storage bucket.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name of this Cloud Object Storage bucket.",
						},
					},
				},
			},
			"storage_href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Cloud Object Storage location of the exported image object. The object at this location may not exist until the job is started, and will be incomplete while the job is running.After the job completes, the exported image object is not managed by the IBM VPC service, and may be removed or replaced with a different object by any user or service with IAM authorization to the bucket.",
			},
			"storage_object": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Cloud Object Storage object for the exported image. This object may not exist untilthe job is started, and will not be complete until the job completes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of this Cloud Object Storage object. Names are unique within a Cloud Object Storage bucket.",
						},
					},
				},
			},
		},
	}
}

func DataSourceIBMIsImageExportRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_export_job", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getImageExportJobOptions := &vpcv1.GetImageExportJobOptions{}

	getImageExportJobOptions.SetImageID(d.Get("image").(string))
	getImageExportJobOptions.SetID(d.Get("image_export_job").(string))

	imageExportJob, _, err := vpcClient.GetImageExportJobWithContext(context, getImageExportJobOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetImageExportJobWithContext failed: %s", err.Error()), "(Data) ibm_is_image_export_job", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *getImageExportJobOptions.ImageID, *getImageExportJobOptions.ID))

	if !core.IsNil(imageExportJob.CompletedAt) {
		if err = d.Set("completed_at", flex.DateTimeToString(imageExportJob.CompletedAt)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting completed_at: %s", err), "(Data) ibm_is_image_export_job", "read", "set-completed_at").GetDiag()
		}
	}

	if err = d.Set("created_at", flex.DateTimeToString(imageExportJob.CreatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_is_image_export_job", "read", "set-created_at").GetDiag()
	}

	if !core.IsNil(imageExportJob.EncryptedDataKey) {
		if err = d.Set("encrypted_data_key", base64.StdEncoding.EncodeToString(*imageExportJob.EncryptedDataKey)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting encrypted_data_key: %s", err), "(Data) ibm_is_image_export_job", "read", "set-encrypted_data_key").GetDiag()
		}
	}

	if err = d.Set("format", imageExportJob.Format); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting format: %s", err), "(Data) ibm_is_image_export_job", "read", "set-format").GetDiag()
	}

	if err = d.Set("href", imageExportJob.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_is_image_export_job", "read", "set-href").GetDiag()
	}

	if err = d.Set("name", imageExportJob.Name); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_is_image_export_job", "read", "set-name").GetDiag()
	}

	if err = d.Set("resource_type", imageExportJob.ResourceType); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_is_image_export_job", "read", "set-resource_type").GetDiag()
	}

	if !core.IsNil(imageExportJob.StartedAt) {
		if err = d.Set("started_at", flex.DateTimeToString(imageExportJob.StartedAt)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting started_at: %s", err), "(Data) ibm_is_image_export_job", "read", "set-started_at").GetDiag()
		}
	}

	if err = d.Set("status", imageExportJob.Status); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status: %s", err), "(Data) ibm_is_image_export_job", "read", "set-status").GetDiag()
	}
	statusReasons := []map[string]interface{}{}
	if imageExportJob.StatusReasons != nil {
		for _, modelItem := range imageExportJob.StatusReasons {
			modelMap, err := DataSourceIBMIsImageExportImageExportJobStatusReasonToMap(&modelItem)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_export_job", "read", "status_reasons-to-map").GetDiag()
			}
			statusReasons = append(statusReasons, modelMap)
		}
	}
	if err = d.Set("status_reasons", statusReasons); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status_reasons: %s", err), "(Data) ibm_is_image_export_job", "read", "set-status_reasons").GetDiag()
	}

	storageBucket := []map[string]interface{}{}
	if imageExportJob.StorageBucket != nil {
		modelMap, err := DataSourceIBMIsImageExportCloudObjectStorageBucketReferenceToMap(imageExportJob.StorageBucket)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_export_job", "read", "storage_bucket-to-map").GetDiag()
		}
		storageBucket = append(storageBucket, modelMap)
	}
	if err = d.Set("storage_bucket", storageBucket); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting storage_bucket: %s", err), "(Data) ibm_is_image_export_job", "read", "set-storage_bucket").GetDiag()
	}

	if err = d.Set("storage_href", imageExportJob.StorageHref); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting storage_href: %s", err), "(Data) ibm_is_image_export_job", "read", "set-storage_href").GetDiag()
	}

	storageObject := []map[string]interface{}{}
	if imageExportJob.StorageObject != nil {
		modelMap, err := DataSourceIBMIsImageExportCloudObjectStorageObjectReferenceToMap(imageExportJob.StorageObject)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_export_job", "read", "storage_object-to-map").GetDiag()
		}
		storageObject = append(storageObject, modelMap)
	}
	if err = d.Set("storage_object", storageObject); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting storage_object: %s", err), "(Data) ibm_is_image_export_job", "read", "set-storage_object").GetDiag()
	}

	return nil
}

func DataSourceIBMIsImageExportImageExportJobStatusReasonToMap(model *vpcv1.ImageExportJobStatusReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Code != nil {
		modelMap["code"] = *model.Code
	}
	if model.Message != nil {
		modelMap["message"] = *model.Message
	}
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func DataSourceIBMIsImageExportCloudObjectStorageBucketReferenceToMap(model *vpcv1.CloudObjectStorageBucketReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}

func DataSourceIBMIsImageExportCloudObjectStorageObjectReferenceToMap(model *vpcv1.CloudObjectStorageObjectReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}
