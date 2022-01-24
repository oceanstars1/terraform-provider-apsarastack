package ecs

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

// VirtualBorderRouterForPhysicalConnectionType is a nested struct in ecs response
type VirtualBorderRouterForPhysicalConnectionType struct {
	CreationTime    string `json:"CreationTime" xml:"CreationTime"`
	CircuitCode     string `json:"CircuitCode" xml:"CircuitCode"`
	RecoveryTime    string `json:"RecoveryTime" xml:"RecoveryTime"`
	TerminationTime string `json:"TerminationTime" xml:"TerminationTime"`
	ActivationTime  string `json:"ActivationTime" xml:"ActivationTime"`
	VbrOwnerUid     int64  `json:"VbrOwnerUid" xml:"VbrOwnerUid"`
	VbrId           string `json:"VbrId" xml:"VbrId"`
	VlanId          int    `json:"VlanId" xml:"VlanId"`
}
