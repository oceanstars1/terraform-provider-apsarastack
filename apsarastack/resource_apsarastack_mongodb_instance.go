package apsarastack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApsaraStackMongoDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackMongoDBInstanceCreate,
		Read:   resourceApsaraStackMongoDBInstanceRead,
		Update: resourceApsaraStackMongoDBInstanceUpdate,
		Delete: resourceApsaraStackMongoDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"engine_version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			}, // EngineVersion
			"audit_policy": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_audit_policy": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"storage_period": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  30,
						},
					},
				},
			}, //AuditPolicy
			"db_instance_class": {
				Type:     schema.TypeString,
				Required: true,
			}, // DBInstanceClass
			"new_connection_string": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"db_instance_storage": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(10, 2000),
				Required:     true,
			}, // DBInstanceStorage
			"replication_factor": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{3, 5, 7}),
				Optional:     true,
				Computed:     true,
			},
			"storage_engine": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"WiredTiger", "RocksDB"}, false),
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
			}, // Engine
			"instance_charge_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{string(PrePaid), string(PostPaid)}, false),
				Optional:     true,
				Default:      PostPaid,
			}, //ChargeType
			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			}, //Period
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			}, // ZoneId
			"vswitch_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			}, // VSwitchId
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"security_ip_list": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			}, //SecurityIPList
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"account_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			}, //AccountPassword
			"kms_encrypted_password": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kmsDiffSuppressFunc,
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
			"backup_period": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"backup_time": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
				Optional:     true,
				Computed:     true,
			},
			"ssl_action": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Open", "Close", "Update"}, false),
				Optional:     true,
				Computed:     true,
			},
			//Computed
			"retention_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"replica_set_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tde_status": {
				Type: schema.TypeString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old != "" || d.Get("engine_version").(string) < "4.0"
				},
				ValidateFunc: validation.StringInSlice([]string{"enabled"}, false),
				Optional:     true,
				ForceNew:     true,
			},
			"maintain_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"maintain_end_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ssl_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func buildMongoDBCreateRequest(d *schema.ResourceData, meta interface{}) (*dds.CreateDBInstanceRequest, error) {
	client := meta.(*connectivity.ApsaraStackClient)

	request := dds.CreateCreateDBInstanceRequest()
	request.RegionId = string(client.Region)
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "dds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.EngineVersion = Trim(d.Get("engine_version").(string))
	request.Engine = "MongoDB"
	request.DBInstanceStorage = requests.NewInteger(d.Get("db_instance_storage").(int))
	request.DBInstanceClass = Trim(d.Get("db_instance_class").(string))
	request.DBInstanceDescription = d.Get("name").(string)

	request.AccountPassword = d.Get("account_password").(string)
	if request.AccountPassword == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return request, WrapError(err)
			}
			request.AccountPassword = decryptResp.Plaintext
		}
	}

	request.ZoneId = d.Get("zone_id").(string)
	request.StorageEngine = d.Get("storage_engine").(string)

	if replication_factor, ok := d.GetOk("replication_factor"); ok {
		request.ReplicationFactor = strconv.Itoa(replication_factor.(int))
	}

	request.NetworkType = string(Classic)
	vswitchId := Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		// check vswitchId in zone
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitch(vswitchId)
		if err != nil {
			return nil, WrapError(err)
		}

		if request.ZoneId == "" {
			request.ZoneId = vsw.ZoneId
		} else if strings.Contains(request.ZoneId, MULTI_IZ_SYMBOL) {
			zonestr := strings.Split(strings.SplitAfter(request.ZoneId, "(")[1], ")")[0]
			if !strings.Contains(zonestr, string([]byte(vsw.ZoneId)[len(vsw.ZoneId)-1])) {
				return nil, WrapError(fmt.Errorf("The specified vswitch %s isn't in the multi zone %s.", vsw.VSwitchId, request.ZoneId))
			}
		} else if request.ZoneId != vsw.ZoneId {
			return nil, WrapError(fmt.Errorf("The specified vswitch %s isn't in the zone %s.", vsw.VSwitchId, request.ZoneId))
		}
		request.VSwitchId = vswitchId
		request.NetworkType = strings.ToUpper(string(Vpc))
		request.VpcId = vsw.VpcId
	}

	request.ChargeType = d.Get("instance_charge_type").(string)
	period, ok := d.GetOk("period")
	if PayType(request.ChargeType) == PrePaid && ok {
		request.Period = requests.NewInteger(period.(int))
	}

	request.SecurityIPList = LOCAL_HOST_IP
	if len(d.Get("security_ip_list").(*schema.Set).List()) > 0 {
		request.SecurityIPList = strings.Join(expandStringList(d.Get("security_ip_list").(*schema.Set).List())[:], COMMA_SEPARATED)
	}

	request.ClientToken = buildClientToken(request.GetActionName())
	return request, nil
}

func resourceApsaraStackMongoDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ddsService := MongoDBService{client}

	request, err := buildMongoDBCreateRequest(d, meta)
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	if err != nil {
		return WrapError(err)
	}

	raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.CreateDBInstance(request)
	})

	if err != nil {
		return WrapError(err)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, ok := raw.(*dds.CreateDBInstanceResponse)
	if !ok {
		return WrapErrorf(err, "Error in Parsing CreateDBInstanceResponse")
	}
	d.SetId(response.DBInstanceId)
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 2*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), "Instance", []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	auditPolicy, ok := d.Get("audit_policy").(map[string]interface{})
	if ok {
		auditPolicyreq := dds.CreateModifyAuditPolicyRequest()
		if auditPolicy["enable_audit_policy"].(string) == "true" {
			auditPolicyreq.AuditStatus = "Enable"
		}
		storagePeriod, _ := strconv.Atoi(auditPolicy["storage_period"].(string))
		auditPolicyreq.StoragePeriod = requests.NewInteger(storagePeriod)
		auditPolicyreq.DBInstanceId = d.Id()
		auditPolicyreq.RegionId = string(client.Region)
		auditPolicyreq.Headers = map[string]string{"RegionId": client.RegionId}
		auditPolicyreq.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "dds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		audit, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
			return client.ModifyAuditPolicy(auditPolicyreq)
		})

		if err != nil {
			return WrapError(err)
		}

		addDebug(auditPolicyreq.GetActionName(), audit, auditPolicyreq)
	}
	if okay := func() bool {
		if _, ok := d.GetOk("backup_period"); ok {
			return ok
		}
		if _, ok := d.GetOk("backup_time"); ok {
			return ok
		}
		return false
	}; okay() {
		err := ddsService.MotifyMongoDBBackupPolicy(d, "Instance")
		if err != nil {
			return WrapError(err)
		}
	}

	if _, sslok := d.GetOk("ssl_action"); sslok {
		sslrequest := dds.CreateModifyDBInstanceSSLRequest()
		sslrequest.DBInstanceId = d.Id()
		sslrequest.RegionId = client.RegionId
		sslrequest.Headers = map[string]string{"RegionId": client.RegionId}
		sslrequest.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "dds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

		sslrequest.SSLAction = d.Get("ssl_action").(string)

		sslraw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.ModifyDBInstanceSSL(sslrequest)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), sslraw, request.RpcRequest, request)
		d.SetPartial("ssl_action")
	}
	return resourceApsaraStackMongoDBInstanceUpdate(d, meta)
}

func resourceApsaraStackMongoDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ddsService := MongoDBService{client}
	instance, err := ddsService.DescribeMongoDBInstance(d.Id(), "Instance")
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	backupPolicy, err := ddsService.DescribeMongoDBBackupPolicy(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("backup_time", backupPolicy.PreferredBackupTime)
	d.Set("backup_period", strings.Split(backupPolicy.PreferredBackupPeriod, ","))
	d.Set("retention_period", backupPolicy.BackupRetentionPeriod)
	ips, err := ddsService.DescribeMongoDBSecurityIps(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("security_ip_list", ips)

	//groupIp, err := ddsService.DescribeMongoDBSecurityGroupId(d.Id())
	//if err != nil {
	//	return WrapError(err)
	//}
	//if len(groupIp.Items.RdsEcsSecurityGroupRel) > 0 {
	//	d.Set("security_group_id", groupIp.Items.RdsEcsSecurityGroupRel[0].SecurityGroupId)
	//}

	d.Set("name", instance.DBInstanceDescription)
	d.Set("engine_version", instance.EngineVersion)
	d.Set("db_instance_class", instance.DBInstanceClass)
	d.Set("db_instance_storage", instance.DBInstanceStorage)
	d.Set("zone_id", instance.ZoneId)
	d.Set("instance_charge_type", instance.ChargeType)
	//if instance.ChargeType == "PrePaid" {
	//	period, err := computePeriodByUnit(instance.CreationTime, instance.Ti, d.Get("period").(int), "Month")
	//	if err != nil {
	//		return WrapError(err)
	//	}
	//	d.Set("period", period)
	//}
	d.Set("vswitch_id", instance.VSwitchId)
	d.Set("storage_engine", instance.StorageEngine)
	d.Set("maintain_start_time", instance.MaintainStartTime)
	d.Set("maintain_end_time", instance.MaintainEndTime)
	d.Set("replica_set_name", instance.ReplicaSetName)

	sslAction, err := ddsService.DescribeDBInstanceSSL(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("ssl_status", sslAction.SSLStatus)

	if replication_factor, err := strconv.Atoi(instance.ReplicationFactor); err == nil {
		d.Set("replication_factor", replication_factor)
	}
	tdeInfo, err := ddsService.DescribeMongoDBTDEInfo(d.Id(), "Instance")
	if err != nil {
		return WrapError(err)
	}
	d.Set("tde_Status", tdeInfo.TDEStatus)

	d.Set("tags", instance.Tags.Tag)
	return nil
}

func resourceApsaraStackMongoDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ddsService := MongoDBService{client}

	d.Partial(true)

	if !d.IsNewResource() && (d.HasChange("instance_charge_type") && d.Get("instance_charge_type").(string) == "PrePaid") {
		prePaidRequest := dds.CreateTransformToPrePaidRequest()
		prePaidRequest.InstanceId = d.Id()
		prePaidRequest.AutoPay = requests.NewBoolean(true)
		prePaidRequest.Period = requests.NewInteger(d.Get("period").(int))
		raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
			return client.TransformToPrePaid(prePaidRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), prePaidRequest.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(prePaidRequest.GetActionName(), raw, prePaidRequest.RpcRequest, prePaidRequest)
		// wait instance status is running after modifying
		stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), "Instance", []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
		d.SetPartial("instance_charge_type")
		d.SetPartial("period")
	}

	if d.HasChange("backup_time") || d.HasChange("backup_period") {
		if err := ddsService.MotifyMongoDBBackupPolicy(d, "Instance"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("backup_time")
		d.SetPartial("backup_period")
	}

	if d.HasChange("tde_status") {
		request := dds.CreateModifyDBInstanceTDERequest()
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "dds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.DBInstanceId = d.Id()
		request.TDEStatus = d.Get("tde_status").(string)
		raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
			return client.ModifyDBInstanceTDE(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("tde_status")
	}

	if d.HasChange("maintain_start_time") || d.HasChange("maintain_end_time") {
		request := dds.CreateModifyDBInstanceMaintainTimeRequest()
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "dds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

		request.DBInstanceId = d.Id()
		request.MaintainStartTime = d.Get("maintain_start_time").(string)
		request.MaintainEndTime = d.Get("maintain_end_time").(string)

		raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
			return client.ModifyDBInstanceMaintainTime(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("maintain_start_time")
		d.SetPartial("maintain_end_time")
	}

	if d.HasChange("security_group_id") {
		request := dds.CreateModifySecurityGroupConfigurationRequest()
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "dds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

		request.DBInstanceId = d.Id()
		request.SecurityGroupId = d.Get("security_group_id").(string)

		raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
			return client.ModifySecurityGroupConfiguration(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("security_group_id")
	}

	if err := ddsService.setInstanceTags(d); err != nil {
		return WrapError(err)
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceApsaraStackMongoDBInstanceRead(d, meta)
	}

	if d.HasChange("name") {
		request := dds.CreateModifyDBInstanceDescriptionRequest()
		request.DBInstanceId = d.Id()
		request.DBInstanceDescription = d.Get("name").(string)

		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.ModifyDBInstanceDescription(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("name")
	}

	if d.HasChange("security_ip_list") {
		ipList := expandStringList(d.Get("security_ip_list").(*schema.Set).List())
		ipstr := strings.Join(ipList[:], COMMA_SEPARATED)
		// default disable connect from outside
		if ipstr == "" {
			ipstr = LOCAL_HOST_IP
		}

		if err := ddsService.ModifyMongoDBSecurityIps(d.Id(), "Instance", ipstr); err != nil {
			return WrapError(err)
		}
		d.SetPartial("security_ip_list")
	}

	if d.HasChange("account_password") || d.HasChange("kms_encrypted_password") {
		var accountPassword string
		if accountPassword = d.Get("account_password").(string); accountPassword != "" {
			d.SetPartial("account_password")
		} else if kmsPassword := d.Get("kms_encrypted_password").(string); kmsPassword != "" {
			kmsService := KmsService{meta.(*connectivity.ApsaraStackClient)}
			decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			accountPassword = decryptResp.Plaintext
			d.SetPartial("kms_encrypted_password")
			d.SetPartial("kms_encryption_context")
		}

		err := ddsService.ResetAccountPassword(d, accountPassword)
		if err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("ssl_action") {
		request := dds.CreateModifyDBInstanceSSLRequest()
		request.DBInstanceId = d.Id()
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "dds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

		request.SSLAction = d.Get("ssl_action").(string)

		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.ModifyDBInstanceSSL(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("ssl_action")
	}

	if d.HasChange("db_instance_storage") ||
		d.HasChange("db_instance_class") ||
		d.HasChange("replication_factor") {

		request := dds.CreateModifyDBInstanceSpecRequest()
		request.DBInstanceId = d.Id()

		request.DBInstanceClass = d.Get("db_instance_class").(string)
		request.DBInstanceStorage = strconv.Itoa(d.Get("db_instance_storage").(int))
		request.ReplicationFactor = strconv.Itoa(d.Get("replication_factor").(int))

		// wait instance status is running before modifying
		stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), "Instance", []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.ModifyDBInstanceSpec(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}

		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("db_instance_class")
		d.SetPartial("db_instance_storage")
		d.SetPartial("replication_factor")

		// wait instance status is running after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
	}
	newconnectionString, ok := d.Get("new_connection_string").(string)
	currconnectionString, ok1 := d.Get("connection_string").(string)

	if ok && ok1 {
		var connFound bool
		if !d.HasChange("new_connection_string") {
			goto contd
		}
		DBInstance, err := ddsService.DescribeMongoDBInstanceAttribute(d.Id())
		if err != nil {
			return WrapError(err)
		}
	db:
		for _, x := range DBInstance.ReplicaSets.ReplicaSet {
			if currconnectionString == x.ConnectionDomain {
				connFound = true
				break db
			}
		}
		if !connFound {
			errs := Error("CurrentConnectionString Not Found")
			return WrapError(errs)
		}
		conn := DBInstanceConnectionString(d, meta)
		conn.NewConnectionString = newconnectionString
		conn.CurrentConnectionString = currconnectionString
		audit, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
			return client.ModifyDBInstanceConnectionString(conn)
		})
		if err != nil {
			return WrapError(err)
		}

		addDebug(conn.GetActionName(), audit, conn)
	}
contd:
	if _, okay := d.GetOk("audit_policy"); okay {
		if d.HasChange("audit_policy") {
			auditPolicy, ok := d.Get("audit_policy").(map[string]interface{})
			if ok {
				auditPolicyreq := dds.CreateModifyAuditPolicyRequest()
				if auditPolicy["enable_audit_policy"].(string) == "true" {
					auditPolicyreq.AuditStatus = "Enable"
				} else if auditPolicy["enable_audit_policy"].(string) == "false" {
					auditPolicyreq.AuditStatus = "Disabled"
				}

				storagePeriod, _ := strconv.Atoi(auditPolicy["storage_period"].(string))
				auditPolicyreq.StoragePeriod = requests.NewInteger(storagePeriod)
				auditPolicyreq.DBInstanceId = d.Id()
				auditPolicyreq.RegionId = string(client.Region)
				auditPolicyreq.Headers = map[string]string{"RegionId": client.RegionId}
				auditPolicyreq.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "dds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
				audit, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
					return client.ModifyAuditPolicy(auditPolicyreq)
				})

				if err != nil {
					return WrapError(err)
				}

				addDebug(auditPolicyreq.GetActionName(), audit, auditPolicyreq)
			} else {
				auditPolicyreq := dds.CreateModifyAuditPolicyRequest()
				auditPolicyreq.AuditStatus = "Disabled"
				auditPolicyreq.DBInstanceId = d.Id()
				auditPolicyreq.RegionId = string(client.Region)
				auditPolicyreq.Headers = map[string]string{"RegionId": client.RegionId}
				auditPolicyreq.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "dds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
				audit, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
					return client.ModifyAuditPolicy(auditPolicyreq)
				})

				if err != nil {
					return WrapError(err)
				}

				addDebug(auditPolicyreq.GetActionName(), audit, auditPolicyreq)

			}
		}
	} else if !okay {
		auditPolicyreq := dds.CreateModifyAuditPolicyRequest()
		auditPolicyreq.AuditStatus = "Disabled"
		auditPolicyreq.DBInstanceId = d.Id()
		auditPolicyreq.RegionId = string(client.Region)
		auditPolicyreq.Headers = map[string]string{"RegionId": client.RegionId}
		auditPolicyreq.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "dds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		audit, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
			return client.ModifyAuditPolicy(auditPolicyreq)
		})

		if err != nil {
			return WrapError(err)
		}

		addDebug(auditPolicyreq.GetActionName(), audit, auditPolicyreq)
	}

	d.Partial(false)
	return resourceApsaraStackMongoDBInstanceRead(d, meta)
}

func resourceApsaraStackMongoDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ddsService := MongoDBService{client}
	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{"AccessKeyId": client.AccessKey, "AccessKeySecret": client.SecretKey, "Product": "Dds", "RegionId": client.RegionId, "Action": "DeleteDBInstance", "Version": "2015-12-01", "Department": client.Department, "ResourceGroup": client.ResourceGroup, "Forwardedregionid": client.RegionId}
	request.Method = "POST"
	request.Product = "Dds"
	request.Version = "2015-12-01"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ServiceCode = "Dds"
	request.ApiName = "DeleteDBInstance"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.RegionId = client.RegionId
	request.Domain = client.Domain
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"Product":           "Dds",
		"Version":           "2015-12-01",
		"RegionId":          client.RegionId,
		"AccessKeyId":       client.AccessKey,
		"AccessKeySecret":   client.SecretKey,
		"Department":        client.Department,
		"ResourceGroup":     client.ResourceGroup,
		"DBInstanceId":      d.Id(),
		"Action":            "DeleteDBInstance",
		"Format":            "JSON",
		"Forwardedregionid": client.RegionId,
	}
	raw, err := client.WithEcsClient(func(client *ecs.Client) (interface{}, error) {
		return client.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	response := raw.(*responses.CommonResponse)
	if !response.IsSuccess() {
		return Error("DeleteDBInstance Failed", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	_, err = ddsService.DescribeMongoDBInstance(d.Id(), "Instance")
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil
		}
		return WrapError(err)
	}

	return WrapError(err)
}
func DBInstanceConnectionString(d *schema.ResourceData, meta interface{}) *dds.ModifyDBInstanceConnectionStringRequest {
	client := meta.(*connectivity.ApsaraStackClient)
	conn := dds.CreateModifyDBInstanceConnectionStringRequest()
	conn.DBInstanceId = d.Id()
	conn.RegionId = string(client.Region)
	conn.Headers = map[string]string{"RegionId": client.RegionId}
	conn.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "dds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	return conn
}
