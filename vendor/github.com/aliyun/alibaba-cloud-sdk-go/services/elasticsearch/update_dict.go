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

// UpdateDict invokes the elasticsearch.UpdateDict API synchronously
func (client *Client) UpdateDict(request *UpdateDictRequest) (response *UpdateDictResponse, err error) {
	response = CreateUpdateDictResponse()
	err = client.DoAction(request, response)
	return
}

// UpdateDictWithChan invokes the elasticsearch.UpdateDict API asynchronously
func (client *Client) UpdateDictWithChan(request *UpdateDictRequest) (<-chan *UpdateDictResponse, <-chan error) {
	responseChan := make(chan *UpdateDictResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.UpdateDict(request)
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

// UpdateDictWithCallback invokes the elasticsearch.UpdateDict API asynchronously
func (client *Client) UpdateDictWithCallback(request *UpdateDictRequest, callback func(response *UpdateDictResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *UpdateDictResponse
		var err error
		defer close(result)
		response, err = client.UpdateDict(request)
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

// UpdateDictRequest is the request struct for api UpdateDict
type UpdateDictRequest struct {
	*requests.RoaRequest
	InstanceId  string `position:"Path" name:"InstanceId"`
	ClientToken string `position:"Query" name:"clientToken"`
}

// UpdateDictResponse is the response struct for api UpdateDict
type UpdateDictResponse struct {
	*responses.BaseResponse
	RequestId string     `json:"RequestId" xml:"RequestId"`
	Result    []DictList `json:"Result" xml:"Result"`
}

// CreateUpdateDictRequest creates a request to invoke UpdateDict API
func CreateUpdateDictRequest() (request *UpdateDictRequest) {
	request = &UpdateDictRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("elasticsearch", "2017-06-13", "UpdateDict", "/openapi/instances/[InstanceId]/dict", "elasticsearch", "openAPI")
	request.Method = requests.PUT
	return
}

// CreateUpdateDictResponse creates a response to parse from UpdateDict response
func CreateUpdateDictResponse() (response *UpdateDictResponse) {
	response = &UpdateDictResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
