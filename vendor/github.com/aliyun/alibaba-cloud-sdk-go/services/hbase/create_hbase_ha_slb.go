package hbase

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

// CreateHbaseHaSlb invokes the hbase.CreateHbaseHaSlb API synchronously
func (client *Client) CreateHbaseHaSlb(request *CreateHbaseHaSlbRequest) (response *CreateHbaseHaSlbResponse, err error) {
	response = CreateCreateHbaseHaSlbResponse()
	err = client.DoAction(request, response)
	return
}

// CreateHbaseHaSlbWithChan invokes the hbase.CreateHbaseHaSlb API asynchronously
func (client *Client) CreateHbaseHaSlbWithChan(request *CreateHbaseHaSlbRequest) (<-chan *CreateHbaseHaSlbResponse, <-chan error) {
	responseChan := make(chan *CreateHbaseHaSlbResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateHbaseHaSlb(request)
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

// CreateHbaseHaSlbWithCallback invokes the hbase.CreateHbaseHaSlb API asynchronously
func (client *Client) CreateHbaseHaSlbWithCallback(request *CreateHbaseHaSlbRequest, callback func(response *CreateHbaseHaSlbResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateHbaseHaSlbResponse
		var err error
		defer close(result)
		response, err = client.CreateHbaseHaSlb(request)
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

// CreateHbaseHaSlbRequest is the request struct for api CreateHbaseHaSlb
type CreateHbaseHaSlbRequest struct {
	*requests.RpcRequest
	ClientToken string `position:"Query" name:"ClientToken"`
	HaTypes     string `position:"Query" name:"HaTypes"`
	HbaseType   string `position:"Query" name:"HbaseType"`
	BdsId       string `position:"Query" name:"BdsId"`
	HaId        string `position:"Query" name:"HaId"`
}

// CreateHbaseHaSlbResponse is the response struct for api CreateHbaseHaSlb
type CreateHbaseHaSlbResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateCreateHbaseHaSlbRequest creates a request to invoke CreateHbaseHaSlb API
func CreateCreateHbaseHaSlbRequest() (request *CreateHbaseHaSlbRequest) {
	request = &CreateHbaseHaSlbRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("HBase", "2019-01-01", "CreateHbaseHaSlb", "hbase", "openAPI")
	request.Method = requests.POST
	return
}

// CreateCreateHbaseHaSlbResponse creates a response to parse from CreateHbaseHaSlb response
func CreateCreateHbaseHaSlbResponse() (response *CreateHbaseHaSlbResponse) {
	response = &CreateHbaseHaSlbResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
