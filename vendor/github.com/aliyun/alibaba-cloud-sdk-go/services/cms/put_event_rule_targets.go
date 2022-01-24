package cms

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

// PutEventRuleTargets invokes the cms.PutEventRuleTargets API synchronously
func (client *Client) PutEventRuleTargets(request *PutEventRuleTargetsRequest) (response *PutEventRuleTargetsResponse, err error) {
	response = CreatePutEventRuleTargetsResponse()
	err = client.DoAction(request, response)
	return
}

// PutEventRuleTargetsWithChan invokes the cms.PutEventRuleTargets API asynchronously
func (client *Client) PutEventRuleTargetsWithChan(request *PutEventRuleTargetsRequest) (<-chan *PutEventRuleTargetsResponse, <-chan error) {
	responseChan := make(chan *PutEventRuleTargetsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.PutEventRuleTargets(request)
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

// PutEventRuleTargetsWithCallback invokes the cms.PutEventRuleTargets API asynchronously
func (client *Client) PutEventRuleTargetsWithCallback(request *PutEventRuleTargetsRequest, callback func(response *PutEventRuleTargetsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *PutEventRuleTargetsResponse
		var err error
		defer close(result)
		response, err = client.PutEventRuleTargets(request)
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

// PutEventRuleTargetsRequest is the request struct for api PutEventRuleTargets
type PutEventRuleTargetsRequest struct {
	*requests.RpcRequest
	WebhookParameters *[]PutEventRuleTargetsWebhookParameters `position:"Query" name:"WebhookParameters"  type:"Repeated"`
	ContactParameters *[]PutEventRuleTargetsContactParameters `position:"Query" name:"ContactParameters"  type:"Repeated"`
	OpenApiParameters *[]PutEventRuleTargetsOpenApiParameters `position:"Query" name:"OpenApiParameters"  type:"Repeated"`
	SlsParameters     *[]PutEventRuleTargetsSlsParameters     `position:"Query" name:"SlsParameters"  type:"Repeated"`
	RuleName          string                                  `position:"Query" name:"RuleName"`
	MnsParameters     *[]PutEventRuleTargetsMnsParameters     `position:"Query" name:"MnsParameters"  type:"Repeated"`
	FcParameters      *[]PutEventRuleTargetsFcParameters      `position:"Query" name:"FcParameters"  type:"Repeated"`
}

// PutEventRuleTargetsWebhookParameters is a repeated param struct in PutEventRuleTargetsRequest
type PutEventRuleTargetsWebhookParameters struct {
	Protocol string `name:"Protocol"`
	Method   string `name:"Method"`
	Id       string `name:"Id"`
	Url      string `name:"Url"`
}

// PutEventRuleTargetsContactParameters is a repeated param struct in PutEventRuleTargetsRequest
type PutEventRuleTargetsContactParameters struct {
	Level            string `name:"Level"`
	Id               string `name:"Id"`
	ContactGroupName string `name:"ContactGroupName"`
}

// PutEventRuleTargetsOpenApiParameters is a repeated param struct in PutEventRuleTargetsRequest
type PutEventRuleTargetsOpenApiParameters struct {
	Product string `name:"Product"`
	Role    string `name:"Role"`
	Action  string `name:"Action"`
	Id      string `name:"Id"`
	Arn     string `name:"Arn"`
	Region  string `name:"Region"`
	Version string `name:"Version"`
}

// PutEventRuleTargetsSlsParameters is a repeated param struct in PutEventRuleTargetsRequest
type PutEventRuleTargetsSlsParameters struct {
	Project  string `name:"Project"`
	Id       string `name:"Id"`
	Region   string `name:"Region"`
	LogStore string `name:"LogStore"`
}

// PutEventRuleTargetsMnsParameters is a repeated param struct in PutEventRuleTargetsRequest
type PutEventRuleTargetsMnsParameters struct {
	Id     string `name:"Id"`
	Region string `name:"Region"`
	Queue  string `name:"Queue"`
}

// PutEventRuleTargetsFcParameters is a repeated param struct in PutEventRuleTargetsRequest
type PutEventRuleTargetsFcParameters struct {
	FunctionName string `name:"FunctionName"`
	ServiceName  string `name:"ServiceName"`
	Id           string `name:"Id"`
	Region       string `name:"Region"`
}

// PutEventRuleTargetsResponse is the response struct for api PutEventRuleTargets
type PutEventRuleTargetsResponse struct {
	*responses.BaseResponse
	Code                    string                  `json:"Code" xml:"Code"`
	Message                 string                  `json:"Message" xml:"Message"`
	RequestId               string                  `json:"RequestId" xml:"RequestId"`
	Success                 bool                    `json:"Success" xml:"Success"`
	FailedParameterCount    string                  `json:"FailedParameterCount" xml:"FailedParameterCount"`
	FailedContactParameters FailedContactParameters `json:"FailedContactParameters" xml:"FailedContactParameters"`
	FailedMnsParameters     FailedMnsParameters     `json:"FailedMnsParameters" xml:"FailedMnsParameters"`
	FailedFcParameters      FailedFcParameters      `json:"FailedFcParameters" xml:"FailedFcParameters"`
}

// CreatePutEventRuleTargetsRequest creates a request to invoke PutEventRuleTargets API
func CreatePutEventRuleTargetsRequest() (request *PutEventRuleTargetsRequest) {
	request = &PutEventRuleTargetsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cms", "2019-01-01", "PutEventRuleTargets", "cms", "openAPI")
	request.Method = requests.POST
	return
}

// CreatePutEventRuleTargetsResponse creates a response to parse from PutEventRuleTargets response
func CreatePutEventRuleTargetsResponse() (response *PutEventRuleTargetsResponse) {
	response = &PutEventRuleTargetsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
