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

// EstimatedRestartTime invokes the elasticsearch.EstimatedRestartTime API synchronously
func (client *Client) EstimatedRestartTime(request *EstimatedRestartTimeRequest) (response *EstimatedRestartTimeResponse, err error) {
	response = CreateEstimatedRestartTimeResponse()
	err = client.DoAction(request, response)
	return
}

// EstimatedRestartTimeWithChan invokes the elasticsearch.EstimatedRestartTime API asynchronously
func (client *Client) EstimatedRestartTimeWithChan(request *EstimatedRestartTimeRequest) (<-chan *EstimatedRestartTimeResponse, <-chan error) {
	responseChan := make(chan *EstimatedRestartTimeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.EstimatedRestartTime(request)
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

// EstimatedRestartTimeWithCallback invokes the elasticsearch.EstimatedRestartTime API asynchronously
func (client *Client) EstimatedRestartTimeWithCallback(request *EstimatedRestartTimeRequest, callback func(response *EstimatedRestartTimeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *EstimatedRestartTimeResponse
		var err error
		defer close(result)
		response, err = client.EstimatedRestartTime(request)
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

// EstimatedRestartTimeRequest is the request struct for api EstimatedRestartTime
type EstimatedRestartTimeRequest struct {
	*requests.RoaRequest
	InstanceId string           `position:"Path" name:"InstanceId"`
	Force      requests.Boolean `position:"Query" name:"force"`
}

// EstimatedRestartTimeResponse is the response struct for api EstimatedRestartTime
type EstimatedRestartTimeResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Result    Result `json:"Result" xml:"Result"`
}

// CreateEstimatedRestartTimeRequest creates a request to invoke EstimatedRestartTime API
func CreateEstimatedRestartTimeRequest() (request *EstimatedRestartTimeRequest) {
	request = &EstimatedRestartTimeRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("elasticsearch", "2017-06-13", "EstimatedRestartTime", "/openapi/instances/[InstanceId]/estimated-time/restart-time", "elasticsearch", "openAPI")
	request.Method = requests.POST
	return
}

// CreateEstimatedRestartTimeResponse creates a response to parse from EstimatedRestartTime response
func CreateEstimatedRestartTimeResponse() (response *EstimatedRestartTimeResponse) {
	response = &EstimatedRestartTimeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
