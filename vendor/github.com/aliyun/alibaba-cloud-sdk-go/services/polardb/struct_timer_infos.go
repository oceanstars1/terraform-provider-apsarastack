package polardb

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

// TimerInfos is a nested struct in polardb response
type TimerInfos struct {
	Region               string `json:"Region" xml:"Region"`
	DbClusterStatus      string `json:"DbClusterStatus" xml:"DbClusterStatus"`
	TaskId               string `json:"TaskId" xml:"TaskId"`
	Action               string `json:"Action" xml:"Action"`
	PlannedStartTime     string `json:"PlannedStartTime" xml:"PlannedStartTime"`
	Status               string `json:"Status" xml:"Status"`
	PlannedEndTime       string `json:"PlannedEndTime" xml:"PlannedEndTime"`
	DbClusterDescription string `json:"DbClusterDescription" xml:"DbClusterDescription"`
	PlannedTime          string `json:"PlannedTime" xml:"PlannedTime"`
	DBClusterId          string `json:"DBClusterId" xml:"DBClusterId"`
	OrderId              string `json:"OrderId" xml:"OrderId"`
	TaskCancel           bool   `json:"TaskCancel" xml:"TaskCancel"`
}
