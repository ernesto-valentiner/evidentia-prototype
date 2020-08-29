/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-chaincode-go/shim"

	sc "github.com/hyperledger/fabric-protos-go/peer"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
)

type EvidentiaContract struct {
}

const (
	Searching = "SEARCHING"
	Pending   = "PENDING"
	Completed = "COMPLETED"
)

type ServiceProvider struct {
	ServiceProviderDocType string   `json:"docType"`
	ID                     string   `json:"id"`
	IP                     string   `json:"ip"`
	Port                   string   `json:"port"`
	Services               []string `json:"services"`
}

const ServiceProviderDocType = "ServiceProvider"

type ServiceExecution struct {
	ServiceExecutionDocType string `json:"docType"`
	SourceID                string `json:"sourceID"`
	ServiceName             string `json:"serviceName"`
	ServiceParameters       string `json:"serviceParameters"`
	TargetID                string `json:"targetID"`
	Evidence                string `json:"evidence"`
	Response                string `json:"response"`
	Status                  string `json:"status"`
}

const ServiceExecutionDocType = "ServiceExecution"

func (s *EvidentiaContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *EvidentiaContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := stub.GetFunctionAndParameters()
	fmt.Printf("Invoke: %s\n", function)

	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "initLedger" {
		return s.Init(stub)
	} else if function == "addServiceProvider" {
		return s.AddServiceProvider(stub, args)
	} else if function == "requestServiceProvider" {
		return s.RequestServiceProvider(stub, args)
	} else if function == "updateServiceExecutionResponse" {
		return s.UpdateServiceExecutionResponse(stub, args)
	} else if function == "updateServiceExecutionTarget" {
		return s.UpdateServiceExecutionTarget(stub, args)
	} else if function == "getServiceExecution" {
		return s.GetServiceExecution(stub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *EvidentiaContract) AddServiceProvider(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4 - [username, ip, port, services] ")
	}

	mspid, _ := cid.GetMSPID(stub)
	fmt.Println("MSPID: " + mspid)
	if mspid != "Org1MSP" {
		return shim.Error("Error - user not allowed to add service provider")
	}

	serviceProviderID := args[0]
	ip := args[1]
	port := args[2]
	servicesString := args[3]
	servicesString = strings.Replace(servicesString, "[", "", -1)
	servicesString = strings.Replace(servicesString, "]", "", -1)
	services := strings.Split(servicesString, ",")

	serviceProvider := ServiceProvider{
		ServiceProviderDocType: ServiceProviderDocType,
		ID:                     serviceProviderID,
		IP:                     ip,
		Port:                   port,
		Services:               services}

	serviceProviderAsBytes, _ := json.Marshal(serviceProvider)
	key := ip + "_" + port
	stub.PutState(key, serviceProviderAsBytes)

	fmt.Println("Service Provider ADDED = " + string(serviceProviderAsBytes))
	return shim.Success(nil)
}

func (s *EvidentiaContract) RequestServiceProvider(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	// var jsonResp string

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2 - [service_name, params] ")
	}

	serviceName := args[0]
	serviceParams := args[1]

	// List of available services
	availableServiceProviders := queryProvidersByServiceName(stub, serviceName)
	// if len(availableServiceProviders) == 0 {
	// 	jsonResp = "{\"Error\":\"No service Provider available for the specified service\"}"
	// 	return shim.Error(jsonResp)
	// }
	// buffer is an array containing IPs and PORTs of available services
	var buffer bytes.Buffer
	bArrayMemberAlreadyWritten := false
	for _, serviceProvider := range availableServiceProviders {
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString(serviceProvider.IP)
		buffer.WriteString(",")
		buffer.WriteString(serviceProvider.Port)
		buffer.WriteString(",")
		buffer.WriteString(serviceProvider.ID)
		bArrayMemberAlreadyWritten = true
		//}
	}
	if len(availableServiceProviders) != 0 {
		createServiceExecution(stub, serviceName, serviceParams)
	}

	return shim.Success(buffer.Bytes())
}

func createServiceExecution(stub shim.ChaincodeStubInterface, serviceName string, serviceParams string) sc.Response {

	// Get the id of the caller
	x509, _ := cid.GetX509Certificate(stub)
	sourceID := x509.Subject.CommonName
	fmt.Println("sourceID =" + sourceID)

	serviceExecution := ServiceExecution{
		ServiceExecutionDocType: ServiceExecutionDocType,
		SourceID:                sourceID,
		ServiceName:             serviceName,
		ServiceParameters:       serviceParams,
		TargetID:                "",
		Evidence:                "",
		Response:                "",
		Status:                  Searching}

	keyString := sourceID + serviceName + serviceParams
	fmt.Println("KEY = " + keyString)
	sha256Key := sha256.Sum256([]byte(keyString))
	key := hex.EncodeToString(sha256Key[:])
	fmt.Println("hash(key)", key)

	serviceExecutionAsBytes, _ := json.Marshal(serviceExecution)
	stub.PutState(key, serviceExecutionAsBytes)

	fmt.Println("Service Execution CREATED = " + string(serviceExecutionAsBytes))
	return shim.Success(nil)
}

func (s *EvidentiaContract) UpdateServiceExecutionResponse(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4 - service name, query params, response")
	}

	x509, _ := cid.GetX509Certificate(stub)
	sourceID := x509.Subject.CommonName
	fmt.Println("sourceID =" + sourceID)
	serviceName := args[0]
	serviceParams := args[1]

	keyString := sourceID + serviceName + serviceParams
	fmt.Println("KEY = " + keyString)
	sha256Key := sha256.Sum256([]byte(keyString))
	key := hex.EncodeToString(sha256Key[:])
	fmt.Println("hash(key) = ", key)

	serviceExecutionAsBytes, er := stub.GetState(key)
	if er != nil {
		return shim.Error(er.Error())
	}

	serviceExecution := ServiceExecution{}
	json.Unmarshal(serviceExecutionAsBytes, &serviceExecution)
	fmt.Println(serviceExecution)

	targetID := serviceExecution.TargetID
	fmt.Println("targetID =" + targetID)

	if sourceID == serviceExecution.SourceID && serviceExecution.Status == Pending {

		serviceExecution.Evidence = args[2]
		serviceExecution.Response = args[3]
		serviceExecution.Status = Completed
		serviceExecutionAsBytes, _ = json.Marshal(serviceExecution)
		stub.PutState(key, serviceExecutionAsBytes)

		fmt.Println("Service Execution UPDATED with response = " + string(serviceExecutionAsBytes))
		return shim.Success(nil)
	} else {
		fmt.Println("Not allowed to update response")
		return shim.Error("Not allowed to update response")
	}

}

func (s *EvidentiaContract) UpdateServiceExecutionTarget(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	var jsonResp string

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3 - service name, query params")
	}

	x509, _ := cid.GetX509Certificate(stub)
	sourceID := x509.Subject.CommonName
	fmt.Println("Source ID: " + sourceID)
	serviceName := args[0]
	serviceParams := args[1]
	targetID := args[2]
	fmt.Println("targetID =" + targetID)

	//verify serviceProvider

	// List of available services
	availableServiceProviders := queryProvidersByServiceName(stub, serviceName)

	if len(availableServiceProviders) == 0 {
		jsonResp = "{\"Error\":\"No service Provider available for the specified service\"}"
		return shim.Error(jsonResp)
	}

	//Verify that the targetID has the right to execute the service
	verify := false
	for _, serviceProvider := range availableServiceProviders {
		if serviceProvider.ID == targetID {
			verify = true
		}
	}

	keyString := sourceID + serviceName + serviceParams
	fmt.Println("KEY = " + keyString)
	sha256Key := sha256.Sum256([]byte(keyString))
	key := hex.EncodeToString(sha256Key[:])
	fmt.Println("hash(key) = ", key)

	serviceExecutionAsBytes, er := stub.GetState(key)
	if er != nil {
		return shim.Error(er.Error())
	}

	serviceExecution := ServiceExecution{}
	json.Unmarshal(serviceExecutionAsBytes, &serviceExecution)
	fmt.Println(string(serviceExecutionAsBytes))

	if verify && serviceExecution.Status == Searching {

		serviceExecution.TargetID = targetID
		serviceExecution.Status = Pending
		serviceExecutionAsBytes, _ = json.Marshal(serviceExecution)
		stub.PutState(key, serviceExecutionAsBytes)

		fmt.Println("Service Execution UPDATED with target = " + targetID)
		return shim.Success(nil)
	} else {
		fmt.Println("Not allowed to update response")
		return shim.Error("Not allowed to update response")
	}

}

func (s *EvidentiaContract) GetServiceExecution(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3 - [sourceID, , serviceName, serviceParams] ")
	}

	sourceID := args[0]
	serviceName := args[1]
	serviceParams := args[2]

	keyString := sourceID + serviceName + serviceParams
	fmt.Println("KEY = " + keyString)
	sha256Key := sha256.Sum256([]byte(keyString))
	key := hex.EncodeToString(sha256Key[:])
	fmt.Println("hash(key) = ", key)

	serviceExecutionAsBytes, er := stub.GetState(key)
	if er != nil {
		return shim.Error(er.Error())
	}
	return shim.Success(serviceExecutionAsBytes)
}

func queryProvidersByServiceName(stub shim.ChaincodeStubInterface, serviceName string) []ServiceProvider {
	//Query all serviceProviders offering the serviceName
	queryString := "{\"selector\":{\"docType\":\"ServiceProvider\",\"services\":{\"$elemMatch\":{\"$eq\":\"" + serviceName + "\"}}}}"
	fmt.Println(queryString)
	resultsIterator, _ := stub.GetQueryResult(queryString)
	defer resultsIterator.Close()

	// Marshall them to ServiceProvider Objects
	var availableServiceProviders []ServiceProvider
	for resultsIterator.HasNext() {
		queryResponse, _ := resultsIterator.Next()
		//Verify if service provider contains the desired service
		var serviceProviderJSON ServiceProvider
		json.Unmarshal(queryResponse.Value, &serviceProviderJSON)
		availableServiceProviders = append(availableServiceProviders, serviceProviderJSON)
	}
	return availableServiceProviders
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(EvidentiaContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
