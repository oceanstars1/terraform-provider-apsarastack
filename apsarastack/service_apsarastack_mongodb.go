package apsarastack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
)

type MongoDBService struct {
	client *connectivity.ApsaraStackClient
}

func (s *MongoDBService) DescribeMongoDBInstance(id, dbType string) (instance dds.DBInstance, err error) {
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "Dds"
	request.Version = "2015-12-01"
	//request.Scheme = "http"
	request.ServiceCode = "Dds"
	request.ApiName = "DescribeDBInstanceAttribute"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeyId": s.client.AccessKey, "AccessKeySecret": s.client.SecretKey, "Product": "Dds", "RegionId": s.client.RegionId, "Action": "DescribeDBInstanceAttribute", "Version": "2015-12-01", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.RegionId = s.client.RegionId
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	if dbType == "Sharding" {
		request.QueryParams = map[string]string{
			"Product":         "Dds",
			"Action":          "DescribeDBInstanceAttribute",
			"Version":         "2015-12-01",
			"RegionId":        s.client.RegionId,
			"AccessKeyId":     s.client.AccessKey,
			"AccessKeySecret": s.client.SecretKey,
			"Department":      s.client.Department,
			"ResourceGroup":   s.client.ResourceGroup,
			"DBInstanceId":    id,
			"DBInstanceType":  "sharding",
		}
	} else {
		request.QueryParams = map[string]string{
			"Product":         "Dds",
			"Action":          "DescribeDBInstanceAttribute",
			"Version":         "2015-12-01",
			"RegionId":        s.client.RegionId,
			"AccessKeyId":     s.client.AccessKey,
			"AccessKeySecret": s.client.SecretKey,
			"Department":      s.client.Department,
			"ResourceGroup":   s.client.ResourceGroup,
			"DBInstanceId":    id,
		}
	}
	var raw interface{}

	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithEcsClient(func(client *ecs.Client) (interface{}, error) {
			return client.ProcessCommonRequest(request)
		})

		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
				return resource.RetryableError(err)
			}
			return resource.RetryableError(err)
		}

		return nil
	}); err != nil {
		return instance, WrapErrorf(err, DefaultErrorMsg, "apsarastack_mongodb_instance", "DescribeDBInstanceAttribute", ApsaraStackSdkGoERROR)
	}
	var Dbresponse dds.DescribeDBInstanceAttributeResponse
	response, ok := raw.(*responses.CommonResponse)
	if !ok {
		return instance, WrapErrorf(err, "Error in parsing DescribeDBInstance Response")
	}

	addDebug(request.GetActionName(), raw, request)
	err = json.Unmarshal(response.GetHttpContentBytes(), &Dbresponse)
	if err != nil {
		panic(err)
	}
	if len(Dbresponse.DBInstances.DBInstance) == 0 {
		return instance, WrapErrorf(Error(GetNotFoundMessage("MongoDB Instance", id)), NotFoundMsg, ApsaraStackSdkGoERROR)
	}
	return Dbresponse.DBInstances.DBInstance[0], nil
}
func (s *MongoDBService) DescribeMongoDBInstanceAttribute(id string) (instance dds.DBInstance, err error) {
	request := dds.CreateDescribeDBInstanceAttributeRequest()
	request.Method = "POST"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeyId": s.client.AccessKey, "AccessKeySecret": s.client.SecretKey, "RegionId": s.client.RegionId, "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.RegionId = s.client.RegionId
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.DBInstanceId = id
	var raw interface{}

	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
			return client.DescribeDBInstanceAttribute(request)
		})

		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
				return resource.RetryableError(err)
			}
			return resource.RetryableError(err)
		}

		return nil
	}); err != nil {
		return instance, WrapErrorf(err, DefaultErrorMsg, "apsarastack_mongodb_instance", "DescribeDBInstanceAttribute", ApsaraStackSdkGoERROR)
	}
	response, ok := raw.(*dds.DescribeDBInstanceAttributeResponse)

	if !ok {
		return instance, WrapErrorf(err, "Error in parsing DescribeDBInstance Response")
	}

	addDebug(request.GetActionName(), raw, request)

	return response.DBInstances.DBInstance[0], nil
}

// WaitForInstance waits for instance to given statusid
func (s *MongoDBService) WaitForMongoDBInstance(instanceId, dbType string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		instance, err := s.DescribeMongoDBInstance(instanceId, dbType)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if instance.DBInstanceStatus == string(status) {
			return nil
		}

		if status == Updating {
			if instance.DBInstanceStatus == "NodeCreating" ||
				instance.DBInstanceStatus == "NodeDeleting" ||
				instance.DBInstanceStatus == "DBInstanceClassChanging" {
				return nil
			}
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, instanceId, GetFunc(1), timeout, instance.DBInstanceStatus, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *MongoDBService) RdsMongodbDBInstanceStateRefreshFunc(id, dbType string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMongoDBInstance(id, dbType)
		if err != nil {
			//if NotFoundError(err) {
			//	// Set this to nil as if we didn't find anything.
			//	return nil, "", WrapError(err)
			//}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if object.DBInstanceStatus == failState {
				return object, object.DBInstanceStatus, WrapError(Error(FailedToReachTargetStatus, object.DBInstanceStatus))
			}
		}
		return object, object.DBInstanceStatus, nil
	}
}

func (s *MongoDBService) DescribeMongoDBSecurityIps(instanceId string) (ips []string, err error) {
	request := dds.CreateDescribeSecurityIpsRequest()
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "dds", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}

	request.DBInstanceId = instanceId

	raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.DescribeSecurityIps(request)
	})
	if err != nil {
		return ips, WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	response, _ := raw.(*dds.DescribeSecurityIpsResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	var ipstr, separator string
	ipsMap := make(map[string]string)
	for _, ip := range response.SecurityIpGroups.SecurityIpGroup {
		if ip.SecurityIpGroupAttribute == "hidden" {
			continue
		}
		ipstr += separator + ip.SecurityIpList
		separator = COMMA_SEPARATED
	}

	for _, ip := range strings.Split(ipstr, COMMA_SEPARATED) {
		ipsMap[ip] = ip
	}

	var finalIps []string
	if len(ipsMap) > 0 {
		for key := range ipsMap {
			finalIps = append(finalIps, key)
		}
	}

	return finalIps, nil
}

func (s *MongoDBService) ModifyMongoDBSecurityIps(instanceId, ips, dbType string) error {
	request := dds.CreateModifySecurityIpsRequest()
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "dds", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.ModifyMode = "Append"
	request.DBInstanceId = instanceId
	request.SecurityIps = ips

	raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.ModifySecurityIps(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err := s.WaitForMongoDBInstance(instanceId, dbType, Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	return nil
}

func (s *MongoDBService) DescribeMongoDBSecurityGroupId(id, dbType string) (*dds.DescribeSecurityGroupConfigurationResponse, error) {
	response := &dds.DescribeSecurityGroupConfigurationResponse{}
	request := dds.CreateDescribeSecurityGroupConfigurationRequest()
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "dds", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}

	request.DBInstanceId = id
	if err := s.WaitForMongoDBInstance(id, dbType, Running, DefaultTimeoutMedium); err != nil {
		return response, WrapError(err)
	}
	raw, err := s.client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
		return ddsClient.DescribeSecurityGroupConfiguration(request)
	})
	if err != nil {
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*dds.DescribeSecurityGroupConfigurationResponse)

	return response, nil
}

func (server *MongoDBService) ModifyMongodbShardingInstanceNode(
	instanceID string, nodeType MongoDBShardingNodeType, stateList, diffList []interface{}) error {
	client := server.client

	err := server.WaitForMongoDBInstance(instanceID, "Sharding", Running, DefaultLongTimeout)
	if err != nil {
		return WrapError(err)
	}

	//create node
	if len(stateList) < len(diffList) {
		createList := diffList[len(stateList):]
		diffList = diffList[:len(stateList)]

		for _, item := range createList {
			node := item.(map[string]interface{})

			request := dds.CreateCreateNodeRequest()
			request.RegionId = server.client.RegionId
			request.Headers = map[string]string{"RegionId": server.client.RegionId}
			request.QueryParams = map[string]string{"AccessKeySecret": server.client.SecretKey, "Product": "dds", "Department": server.client.Department, "ResourceGroup": server.client.ResourceGroup}

			request.DBInstanceId = instanceID
			request.NodeClass = node["node_class"].(string)
			request.NodeType = string(nodeType)
			request.ClientToken = buildClientToken(request.GetActionName())

			if nodeType == MongoDBShardingNodeShard {
				request.NodeStorage = requests.NewInteger(node["node_storage"].(int))
			}

			raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
				return ddsClient.CreateNode(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, instanceID, request.GetActionName(), ApsaraStackSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)

			err = server.WaitForMongoDBInstance(instanceID, "Sharding", Updating, DefaultLongTimeout)
			if err != nil {
				return WrapError(err)
			}

			err = server.WaitForMongoDBInstance(instanceID, "Sharding", Running, DefaultLongTimeout)
			if err != nil {
				return WrapError(err)
			}
		}
	} else if len(stateList) > len(diffList) {
		deleteList := stateList[len(diffList):]
		stateList = stateList[:len(diffList)]

		for _, item := range deleteList {
			node := item.(map[string]interface{})

			request := dds.CreateDeleteNodeRequest()
			request.RegionId = server.client.RegionId
			request.Headers = map[string]string{"RegionId": server.client.RegionId}
			request.QueryParams = map[string]string{"AccessKeySecret": server.client.SecretKey, "Product": "dds", "Department": server.client.Department, "ResourceGroup": server.client.ResourceGroup}

			request.DBInstanceId = instanceID
			request.NodeId = node["node_id"].(string)
			request.ClientToken = buildClientToken(request.GetActionName())

			raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
				return ddsClient.DeleteNode(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, instanceID, request.GetActionName(), ApsaraStackSdkGoERROR)
			}

			addDebug(request.GetActionName(), raw, request.RpcRequest, request)

			err = server.WaitForMongoDBInstance(instanceID, "Sharding", Running, DefaultLongTimeout)
			if err != nil {
				return WrapError(err)
			}
		}
	}

	//motify node
	for key := 0; key < len(stateList); key++ {
		state := stateList[key].(map[string]interface{})
		diff := diffList[key].(map[string]interface{})

		if state["node_class"] != diff["node_class"] ||
			state["node_storage"] != diff["node_storage"] {
			request := dds.CreateModifyNodeSpecRequest()
			request.RegionId = server.client.RegionId
			request.Headers = map[string]string{"RegionId": server.client.RegionId}
			request.QueryParams = map[string]string{"AccessKeySecret": server.client.SecretKey, "Product": "dds", "Department": server.client.Department, "ResourceGroup": server.client.ResourceGroup}

			request.DBInstanceId = instanceID
			request.NodeClass = diff["node_class"].(string)
			request.ClientToken = buildClientToken(request.GetActionName())

			if nodeType == MongoDBShardingNodeShard {
				request.NodeStorage = requests.NewInteger(diff["node_storage"].(int))
			}
			request.NodeId = state["node_id"].(string)

			raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
				return ddsClient.ModifyNodeSpec(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, instanceID, request.GetActionName(), ApsaraStackSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			err = server.WaitForMongoDBInstance(instanceID, "Sharding", Updating, DefaultLongTimeout)
			if err != nil {
				return WrapError(err)
			}
			err = server.WaitForMongoDBInstance(instanceID, "Sharding", Running, DefaultLongTimeout)
			if err != nil {
				return WrapError(err)
			}
		}
	}
	return nil
}

func (s *MongoDBService) DescribeMongoDBBackupPolicy(id string) (*dds.DescribeBackupPolicyResponse, error) {
	response := &dds.DescribeBackupPolicyResponse{}
	request := dds.CreateDescribeBackupPolicyRequest()
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "dds", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}

	request.DBInstanceId = id
	raw, err := s.client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
		return ddsClient.DescribeBackupPolicy(request)
	})
	if err != nil {
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*dds.DescribeBackupPolicyResponse)
	return response, nil
}

func (s *MongoDBService) DescribeMongoDBTDEInfo(id, dbType string) (*dds.DescribeDBInstanceTDEInfoResponse, error) {

	response := &dds.DescribeDBInstanceTDEInfoResponse{}
	request := dds.CreateDescribeDBInstanceTDEInfoRequest()
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "dds", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}

	request.DBInstanceId = id
	statErr := s.WaitForMongoDBInstance(id, dbType, Running, DefaultLongTimeout)
	if statErr != nil {
		return response, WrapError(statErr)
	}
	raw, err := s.client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
		return ddsClient.DescribeDBInstanceTDEInfo(request)
	})
	if err != nil {
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*dds.DescribeDBInstanceTDEInfoResponse)
	return response, nil
}

func (s *MongoDBService) DescribeDBInstanceSSL(id string) (*dds.DescribeDBInstanceSSLResponse, error) {
	response := &dds.DescribeDBInstanceSSLResponse{}
	request := dds.CreateDescribeDBInstanceSSLRequest()
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "dds", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}

	request.DBInstanceId = id
	raw, err := s.client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
		return ddsClient.DescribeDBInstanceSSL(request)
	})
	if err != nil {
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*dds.DescribeDBInstanceSSLResponse)
	return response, nil
}

func (s *MongoDBService) MotifyMongoDBBackupPolicy(d *schema.ResourceData, dbType string) error {
	if err := s.WaitForMongoDBInstance(d.Id(), dbType, Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	periodList := expandStringList(d.Get("backup_period").(*schema.Set).List())
	backupPeriod := fmt.Sprintf("%s", strings.Join(periodList[:], COMMA_SEPARATED))
	backupTime := d.Get("backup_time").(string)

	request := dds.CreateModifyBackupPolicyRequest()
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "dds", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}

	request.DBInstanceId = d.Id()
	request.PreferredBackupPeriod = backupPeriod
	request.PreferredBackupTime = backupTime
	raw, err := s.client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
		return ddsClient.ModifyBackupPolicy(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if err := s.WaitForMongoDBInstance(d.Id(), dbType, Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	return nil
}

func (s *MongoDBService) ResetAccountPassword(d *schema.ResourceData, password string) error {
	request := dds.CreateResetAccountPasswordRequest()
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "dds", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}

	request.DBInstanceId = d.Id()
	request.AccountName = "root"
	request.AccountPassword = password
	raw, err := s.client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
		return ddsClient.ResetAccountPassword(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return err
}

func (s *MongoDBService) setInstanceTags(d *schema.ResourceData) error {
	oraw, nraw := d.GetChange("tags")
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})

	create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

	if len(remove) > 0 {
		var tagKey []string
		for _, v := range remove {
			tagKey = append(tagKey, v.Key)
		}
		request := dds.CreateUntagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = "INSTANCE"
		request.TagKey = &tagKey
		request.RegionId = s.client.RegionId
		request.Headers = map[string]string{"RegionId": s.client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "dds", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}

		raw, err := s.client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.UntagResources(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if len(create) > 0 {
		request := dds.CreateTagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.Tag = &create
		request.ResourceType = "INSTANCE"
		request.RegionId = s.client.RegionId
		request.Headers = map[string]string{"RegionId": s.client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "dds", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}

		raw, err := s.client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.TagResources(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	d.SetPartial("tags")
	return nil
}

func (s *MongoDBService) tagsToMap(tags []dds.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func (s *MongoDBService) ignoreTag(t dds.Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func (s *MongoDBService) tagsInAttributeToMap(tags []dds.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTagInAttribute(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func (s *MongoDBService) ignoreTagInAttribute(t dds.Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func (s *MongoDBService) diffTags(oldTags, newTags []dds.TagResourcesTag) ([]dds.TagResourcesTag, []dds.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []dds.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *MongoDBService) tagsFromMap(m map[string]interface{}) []dds.TagResourcesTag {
	result := make([]dds.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, dds.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

type DescribeDBInstance struct {
	*responses.BaseResponse
	ServerRole      string `json:"serverRole"`
	EagleEyeTraceID string `json:"eagleEyeTraceId"`
	AsapiSuccess    bool   `json:"asapiSuccess"`
	AsapiRequestID  string `json:"asapiRequestId"`
	RequestID       string `json:"RequestId"`
	Domain          string `json:"domain"`
	API             string `json:"api"`
	DBInstances     struct {
		DBInstance []struct {
			ReplicaSetName              string `json:"ReplicaSetName" xml:"ReplicaSetName"`
			DBInstanceDescription       string `json:"DBInstanceDescription" xml:"DBInstanceDescription"`
			Engine                      string `json:"Engine" xml:"Engine"`
			ChargeType                  string `json:"ChargeType" xml:"ChargeType"`
			ReadonlyReplicas            string `json:"ReadonlyReplicas" xml:"ReadonlyReplicas"`
			DBInstanceClass             string `json:"DBInstanceClass" xml:"DBInstanceClass"`
			VpcAuthMode                 string `json:"VpcAuthMode" xml:"VpcAuthMode"`
			CapacityUnit                string `json:"CapacityUnit" xml:"CapacityUnit"`
			DestroyTime                 string `json:"DestroyTime" xml:"DestroyTime"`
			LastDowngradeTime           string `json:"LastDowngradeTime" xml:"LastDowngradeTime"`
			RegionId                    string `json:"RegionId" xml:"RegionId"`
			MaxConnections              int    `json:"MaxConnections" xml:"MaxConnections"`
			ResourceGroupId             string `json:"ResourceGroupId" xml:"ResourceGroupId"`
			CloudType                   string `json:"CloudType" xml:"CloudType"`
			DBInstanceType              string `json:"DBInstanceType" xml:"DBInstanceType"`
			MaintainEndTime             string `json:"MaintainEndTime" xml:"MaintainEndTime"`
			ExpireTime                  string `json:"ExpireTime" xml:"ExpireTime"`
			DBInstanceId                string `json:"DBInstanceId" xml:"DBInstanceId"`
			NetworkType                 string `json:"NetworkType" xml:"NetworkType"`
			ReplicationFactor           string `json:"ReplicationFactor" xml:"ReplicationFactor"`
			MaxIOPS                     int    `json:"MaxIOPS" xml:"MaxIOPS"`
			DBInstanceReleaseProtection bool   `json:"DBInstanceReleaseProtection" xml:"DBInstanceReleaseProtection"`
			ReplacateId                 string `json:"ReplacateId" xml:"ReplacateId"`
			EngineVersion               string `json:"EngineVersion" xml:"EngineVersion"`
			VPCId                       string `json:"VPCId" xml:"VPCId"`
			VSwitchId                   string `json:"VSwitchId" xml:"VSwitchId"`
			VPCCloudInstanceIds         string `json:"VPCCloudInstanceIds" xml:"VPCCloudInstanceIds"`
			MaintainStartTime           string `json:"MaintainStartTime" xml:"MaintainStartTime"`
			CreationTime                string `json:"CreationTime" xml:"CreationTime"`
			DBInstanceStorage           int    `json:"DBInstanceStorage" xml:"DBInstanceStorage"`
			StorageEngine               string `json:"StorageEngine" xml:"StorageEngine"`
			DBInstanceStatus            string `json:"DBInstanceStatus" xml:"DBInstanceStatus"`
			CurrentKernelVersion        string `json:"CurrentKernelVersion" xml:"CurrentKernelVersion"`
			ZoneId                      string `json:"ZoneId" xml:"ZoneId"`
			ProtocolType                string `json:"ProtocolType" xml:"ProtocolType"`
			KindCode                    string `json:"KindCode" xml:"KindCode"`
			LockMode                    string `json:"LockMode" xml:"LockMode"`
			ReplicaSets                 struct {
				ReplicaSet []struct {
					ReplicaSetRole     string `json:"ReplicaSetRole"`
					ConnectionDomain   string `json:"ConnectionDomain"`
					VPCCloudInstanceID string `json:"VPCCloudInstanceId"`
					ConnectionPort     string `json:"ConnectionPort"`
					VPCID              string `json:"VPCId"`
					NetworkType        string `json:"NetworkType"`
					VSwitchID          string `json:"VSwitchId"`
				} `json:"ReplicaSet"`
			} `json:"ReplicaSets"`
			MongosList struct {
				MongosAttribute []interface{} `json:"MongosAttribute"`
			} `json:"MongosList"`
			Tags struct {
				Tag []interface{} `json:"Tag"`
			} `json:"Tags"`
			ConfigserverList struct {
				ConfigserverAttribute []interface{} `json:"ConfigserverAttribute"`
			} `json:"ConfigserverList"`
			ShardList struct {
				ShardAttribute []interface{} `json:"ShardAttribute"`
			} `json:"ShardList"`
		} `json:"DBInstance"`
	} `json:"DBInstances"`
}

type DBInstance struct {
	ReplicaSetName              string `json:"ReplicaSetName" xml:"ReplicaSetName"`
	DBInstanceDescription       string `json:"DBInstanceDescription" xml:"DBInstanceDescription"`
	Engine                      string `json:"Engine" xml:"Engine"`
	ChargeType                  string `json:"ChargeType" xml:"ChargeType"`
	ReadonlyReplicas            string `json:"ReadonlyReplicas" xml:"ReadonlyReplicas"`
	DBInstanceClass             string `json:"DBInstanceClass" xml:"DBInstanceClass"`
	VpcAuthMode                 string `json:"VpcAuthMode" xml:"VpcAuthMode"`
	CapacityUnit                string `json:"CapacityUnit" xml:"CapacityUnit"`
	DestroyTime                 string `json:"DestroyTime" xml:"DestroyTime"`
	LastDowngradeTime           string `json:"LastDowngradeTime" xml:"LastDowngradeTime"`
	RegionId                    string `json:"RegionId" xml:"RegionId"`
	MaxConnections              int    `json:"MaxConnections" xml:"MaxConnections"`
	ResourceGroupId             string `json:"ResourceGroupId" xml:"ResourceGroupId"`
	CloudType                   string `json:"CloudType" xml:"CloudType"`
	DBInstanceType              string `json:"DBInstanceType" xml:"DBInstanceType"`
	MaintainEndTime             string `json:"MaintainEndTime" xml:"MaintainEndTime"`
	ExpireTime                  string `json:"ExpireTime" xml:"ExpireTime"`
	DBInstanceId                string `json:"DBInstanceId" xml:"DBInstanceId"`
	NetworkType                 string `json:"NetworkType" xml:"NetworkType"`
	ReplicationFactor           string `json:"ReplicationFactor" xml:"ReplicationFactor"`
	MaxIOPS                     int    `json:"MaxIOPS" xml:"MaxIOPS"`
	DBInstanceReleaseProtection bool   `json:"DBInstanceReleaseProtection" xml:"DBInstanceReleaseProtection"`
	ReplacateId                 string `json:"ReplacateId" xml:"ReplacateId"`
	EngineVersion               string `json:"EngineVersion" xml:"EngineVersion"`
	VPCId                       string `json:"VPCId" xml:"VPCId"`
	VSwitchId                   string `json:"VSwitchId" xml:"VSwitchId"`
	VPCCloudInstanceIds         string `json:"VPCCloudInstanceIds" xml:"VPCCloudInstanceIds"`
	MaintainStartTime           string `json:"MaintainStartTime" xml:"MaintainStartTime"`
	CreationTime                string `json:"CreationTime" xml:"CreationTime"`
	DBInstanceStorage           int    `json:"DBInstanceStorage" xml:"DBInstanceStorage"`
	StorageEngine               string `json:"StorageEngine" xml:"StorageEngine"`
	DBInstanceStatus            string `json:"DBInstanceStatus" xml:"DBInstanceStatus"`
	CurrentKernelVersion        string `json:"CurrentKernelVersion" xml:"CurrentKernelVersion"`
	ZoneId                      string `json:"ZoneId" xml:"ZoneId"`
	ProtocolType                string `json:"ProtocolType" xml:"ProtocolType"`
	KindCode                    string `json:"KindCode" xml:"KindCode"`
	LockMode                    string `json:"LockMode" xml:"LockMode"`
	ReplicaSets                 struct {
		ReplicaSet []struct {
			ReplicaSetRole     string `json:"ReplicaSetRole"`
			ConnectionDomain   string `json:"ConnectionDomain"`
			VPCCloudInstanceID string `json:"VPCCloudInstanceId"`
			ConnectionPort     string `json:"ConnectionPort"`
			VPCID              string `json:"VPCId"`
			NetworkType        string `json:"NetworkType"`
			VSwitchID          string `json:"VSwitchId"`
		} `json:"ReplicaSet"`
	} `json:"ReplicaSets"`
	MongosList struct {
		MongosAttribute []interface{} `json:"MongosAttribute"`
	} `json:"MongosList"`
	Tags struct {
		Tag []interface{} `json:"Tag"`
	} `json:"Tags"`
	ConfigserverList struct {
		ConfigserverAttribute []interface{} `json:"ConfigserverAttribute"`
	} `json:"ConfigserverList"`
	ShardList struct {
		ShardAttribute []interface{} `json:"ShardAttribute"`
	} `json:"ShardList"`
}
