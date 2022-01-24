package dds

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// CreateDBInstance invokes the dds.CreateDBInstance API synchronously
func (client *Client) CreateDBInstance(request *CreateDBInstanceRequest) (response *CreateDBInstanceResponse, err error) {
	response = CreateCreateDBInstanceResponse()
	err = client.DoAction(request, response)
	return
}

// CreateDBInstanceWithChan invokes the dds.CreateDBInstance API asynchronously
func (client *Client) CreateDBInstanceWithChan(request *CreateDBInstanceRequest) (<-chan *CreateDBInstanceResponse, <-chan error) {
	responseChan := make(chan *CreateDBInstanceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateDBInstance(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// CreateDBInstanceWithCallback invokes the dds.CreateDBInstance API asynchronously
func (client *Client) CreateDBInstanceWithCallback(request *CreateDBInstanceRequest, callback func(response *CreateDBInstanceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateDBInstanceResponse
		var err error
		defer close(result)
		response, err = client.CreateDBInstance(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// CreateDBInstanceRequest is the request struct for api CreateDBInstance
type CreateDBInstanceRequest struct {
	*requests.RpcRequest
	ResourceOwnerId       requests.Integer `position:"Query" name:"ResourceOwnerId"`
	DBInstanceStorage     requests.Integer `position:"Query" name:"DBInstanceStorage"`
	CouponNo              string           `position:"Query" name:"CouponNo"`
	EngineVersion         string           `position:"Query" name:"EngineVersion"`
	NetworkType           string           `position:"Query" name:"NetworkType"`
	ResourceGroupId       string           `position:"Query" name:"ResourceGroupId"`
	SecurityToken         string           `position:"Query" name:"SecurityToken"`
	DBInstanceDescription string           `position:"Query" name:"DBInstanceDescription"`
	BusinessInfo          string           `position:"Query" name:"BusinessInfo"`
	Period                requests.Integer `position:"Query" name:"Period"`
	BackupId              string           `position:"Query" name:"BackupId"`
	OwnerId               requests.Integer `position:"Query" name:"OwnerId"`
	DBInstanceClass       string           `position:"Query" name:"DBInstanceClass"`
	SecurityIPList        string           `position:"Query" name:"SecurityIPList"`
	VSwitchId             string           `position:"Query" name:"VSwitchId"`
	AutoRenew             string           `position:"Query" name:"AutoRenew"`
	ZoneId                string           `position:"Query" name:"ZoneId"`
	ClientToken           string           `position:"Query" name:"ClientToken"`
	ReadonlyReplicas      string           `position:"Query" name:"ReadonlyReplicas"`
	ReplicationFactor     string           `position:"Query" name:"ReplicationFactor"`
	StorageEngine         string           `position:"Query" name:"StorageEngine"`
	DatabaseNames         string           `position:"Query" name:"DatabaseNames"`
	Engine                string           `position:"Query" name:"Engine"`
	RestoreTime           string           `position:"Query" name:"RestoreTime"`
	ResourceOwnerAccount  string           `position:"Query" name:"ResourceOwnerAccount"`
	SrcDBInstanceId       string           `position:"Query" name:"SrcDBInstanceId"`
	OwnerAccount          string           `position:"Query" name:"OwnerAccount"`
	ClusterId             string           `position:"Query" name:"ClusterId"`
	AccountPassword       string           `position:"Query" name:"AccountPassword"`
	VpcId                 string           `position:"Query" name:"VpcId"`
	ChargeType            string           `position:"Query" name:"ChargeType"`
}

// CreateDBInstanceResponse is the response struct for api CreateDBInstance
type CreateDBInstanceResponse struct {
	*responses.BaseResponse
	RequestId    string `json:"RequestId" xml:"RequestId"`
	DBInstanceId string `json:"DBInstanceId" xml:"DBInstanceId"`
	OrderId      string `json:"OrderId" xml:"OrderId"`
}

// CreateCreateDBInstanceRequest creates a request to invoke CreateDBInstance API
func CreateCreateDBInstanceRequest() (request *CreateDBInstanceRequest) {
	request = &CreateDBInstanceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Dds", "2015-12-01", "CreateDBInstance", "dds", "openAPI")
	request.Method = requests.POST
	return
}

// CreateCreateDBInstanceResponse creates a response to parse from CreateDBInstance response
func CreateCreateDBInstanceResponse() (response *CreateDBInstanceResponse) {
	response = &CreateDBInstanceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
