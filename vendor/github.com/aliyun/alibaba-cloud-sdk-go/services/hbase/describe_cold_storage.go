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

// DescribeColdStorage invokes the hbase.DescribeColdStorage API synchronously
func (client *Client) DescribeColdStorage(request *DescribeColdStorageRequest) (response *DescribeColdStorageResponse, err error) {
	response = CreateDescribeColdStorageResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeColdStorageWithChan invokes the hbase.DescribeColdStorage API asynchronously
func (client *Client) DescribeColdStorageWithChan(request *DescribeColdStorageRequest) (<-chan *DescribeColdStorageResponse, <-chan error) {
	responseChan := make(chan *DescribeColdStorageResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeColdStorage(request)
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

// DescribeColdStorageWithCallback invokes the hbase.DescribeColdStorage API asynchronously
func (client *Client) DescribeColdStorageWithCallback(request *DescribeColdStorageRequest, callback func(response *DescribeColdStorageResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeColdStorageResponse
		var err error
		defer close(result)
		response, err = client.DescribeColdStorage(request)
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

// DescribeColdStorageRequest is the request struct for api DescribeColdStorage
type DescribeColdStorageRequest struct {
	*requests.RpcRequest
	ClusterId string `position:"Query" name:"ClusterId"`
}

// DescribeColdStorageResponse is the response struct for api DescribeColdStorage
type DescribeColdStorageResponse struct {
	*responses.BaseResponse
	OpenStatus            string `json:"OpenStatus" xml:"OpenStatus"`
	RequestId             string `json:"RequestId" xml:"RequestId"`
	PayType               string `json:"PayType" xml:"PayType"`
	ColdStorageUsePercent string `json:"ColdStorageUsePercent" xml:"ColdStorageUsePercent"`
	ColdStorageUseAmount  string `json:"ColdStorageUseAmount" xml:"ColdStorageUseAmount"`
	ColdStorageSize       string `json:"ColdStorageSize" xml:"ColdStorageSize"`
	ColdStorageType       string `json:"ColdStorageType" xml:"ColdStorageType"`
	ClusterId             string `json:"ClusterId" xml:"ClusterId"`
}

// CreateDescribeColdStorageRequest creates a request to invoke DescribeColdStorage API
func CreateDescribeColdStorageRequest() (request *DescribeColdStorageRequest) {
	request = &DescribeColdStorageRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("HBase", "2019-01-01", "DescribeColdStorage", "hbase", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeColdStorageResponse creates a response to parse from DescribeColdStorage response
func CreateDescribeColdStorageResponse() (response *DescribeColdStorageResponse) {
	response = &DescribeColdStorageResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}