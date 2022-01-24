package elasticsearch

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

// InterruptElasticsearchTask invokes the elasticsearch.InterruptElasticsearchTask API synchronously
func (client *Client) InterruptElasticsearchTask(request *InterruptElasticsearchTaskRequest) (response *InterruptElasticsearchTaskResponse, err error) {
	response = CreateInterruptElasticsearchTaskResponse()
	err = client.DoAction(request, response)
	return
}

// InterruptElasticsearchTaskWithChan invokes the elasticsearch.InterruptElasticsearchTask API asynchronously
func (client *Client) InterruptElasticsearchTaskWithChan(request *InterruptElasticsearchTaskRequest) (<-chan *InterruptElasticsearchTaskResponse, <-chan error) {
	responseChan := make(chan *InterruptElasticsearchTaskResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.InterruptElasticsearchTask(request)
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

// InterruptElasticsearchTaskWithCallback invokes the elasticsearch.InterruptElasticsearchTask API asynchronously
func (client *Client) InterruptElasticsearchTaskWithCallback(request *InterruptElasticsearchTaskRequest, callback func(response *InterruptElasticsearchTaskResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *InterruptElasticsearchTaskResponse
		var err error
		defer close(result)
		response, err = client.InterruptElasticsearchTask(request)
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

// InterruptElasticsearchTaskRequest is the request struct for api InterruptElasticsearchTask
type InterruptElasticsearchTaskRequest struct {
	*requests.RoaRequest
	InstanceId  string `position:"Path" name:"InstanceId"`
	ClientToken string `position:"Query" name:"clientToken"`
}

// InterruptElasticsearchTaskResponse is the response struct for api InterruptElasticsearchTask
type InterruptElasticsearchTaskResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Code      string `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	Result    bool   `json:"Result" xml:"Result"`
}

// CreateInterruptElasticsearchTaskRequest creates a request to invoke InterruptElasticsearchTask API
func CreateInterruptElasticsearchTaskRequest() (request *InterruptElasticsearchTaskRequest) {
	request = &InterruptElasticsearchTaskRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("elasticsearch", "2017-06-13", "InterruptElasticsearchTask", "/openapi/instances/[InstanceId]/actions/interrupt", "elasticsearch", "openAPI")
	request.Method = requests.POST
	return
}

// CreateInterruptElasticsearchTaskResponse creates a response to parse from InterruptElasticsearchTask response
func CreateInterruptElasticsearchTaskResponse() (response *InterruptElasticsearchTaskResponse) {
	response = &InterruptElasticsearchTaskResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
