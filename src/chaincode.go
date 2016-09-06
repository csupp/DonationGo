/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
    "errors"
    "encoding/json"
    "fmt"
    "strconv"
    "os/exec"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Donation struct {
    id string
    who string
    time string
    rid string
    money int
}

type Request struct {
	id string
	name string
	projectName string
	description string
	expectedMoney int
	currentMoney int
	donationList []string
}


type Person struct {
    id string
    name string
    myRequests []string
    myDonations []string
}



func main() {
    err := shim.Start(new(SimpleChaincode))
    
    if err != nil {
        fmt.Printf("Error starting Simple chaincode: %s", err)
    }
}

func(t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting 1")
    }

    err := stub.PutState("hello_world", []byte(args[0]))
    if err != nil {
        return nil, err
    }

    return nil, nil
}

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
   
     //1. create a donation
     if function == "createDonation" {
     	return t.createDonation(stub, args)
     }
	 //2. create a donation request

	 // if function == "createDonationRequest" {
	 // 	return t.createDonationRequest(stub, args)
	 // }

     return nil, errors.New("Received unknown function invocation")
}

func (t *SimpleChaincode) createDonation(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
     //args: ["jack", "toRequestId", money] 
     var from, toRid string
     var money int
     var djson, pjson []byte
     var myDonations []string
    // var personByte []byte
     var err error

     from = args[0]
     toRid = args[1]
     money, err = strconv.Atoi(args[2])
     if err != nil {
        return nil, err
     }
     //generate donation bean
     donationId, err:= exec.Command("uuidgen").Output()
     if err != nil {
        return nil, err
     }
     donation := new(Donation)
     donation.rid = toRid
     donation.id = string(donationId)
     donation.who = from
     donation.money = money
     djson, err = json.Marshal(donation)
     if err != nil {
        return nil, err
     }
     stub.PutState(donation.id, djson)
  
     var person Person
     // update person data
     personByte, err := stub.GetState(from)
     if err != nil {
        return nil, err
     }
     if personByte == nil {
        person := new(Person)
        pid, err:= exec.Command("uuidgen").Output()
        if err != nil {
           return nil, err
        }
        person.id = string(pid)
        person.name = from
        pjson, err = json.Marshal(person)
        stub.PutState(from, pjson)
     } else {
        err := json.Unmarshal(personByte, &person)
        if err != nil {
           return nil, err
        }
     }
     
    myDonations = person.myDonations
    if myDonations == nil {
        myDonations = make([]string, 0)
    }
    myDonations = append(myDonations, donation.id)
    person.myDonations = myDonations
    

    requestByte, err := stub.GetState(toRid)
    if err != nil {
           return nil, err
    }
    if requestByte == nil {
        return nil, errors.New("request did not exist")
    }
    var request Request
    err = json.Unmarshal(requestByte, &request)
    if err != nil {
           return nil, err
    }
    donationList := request.donationList
    if donationList == nil {
        donationList = make([]string, 0)
    }
    donationList = append(donationList, donation.id)
    request.donationList = donationList
    return nil, nil     
}

// func (t *SimpleChaincode) createDonationRequest(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
//      //args: [jack, projectName, description, expectedMoney]
//      var name, projectName, description string
//      var expectedMoney int
//      name = args[0]
//      projectName = args[1]
//      description = args[2]
//      expectedMoney =args[3]
     
//      //generate a unique id for request id
//  	 requestId, err:= exec.Command("uuidgen").Output()
//      // generate a request bean
//      request := new(Request)
//      request.id = requestId
//      request.name =  name
//      request.projectName = projectName
//      requst.description = description
//      request.expectedMoney = expectedMoney
//      request.currentMoney = 0
   
//      /**
//      update person data
//      **/
//      if person, err := stub.GetState(name); person ==nil {
//      	person := new(Person)
//      	pid, err:= exec.Command("uuidgen").Output()
//      	person.id = pid
//      	person.name = name;
//      	stub.PutState(name, person)
//      }

//      if requestList, err := person.myRequests; requestList == nil {
//  	 	requestList :=make([]int, 0)
//  	 }
//  	 requestList = append(requestList, request.id)
//  	 person.myRequests = requestList

//      // update allRequest
//      if requests, err := stub.GetState("allRequest"); err !=nil {
//      	requests := []*Request{}
//      }
//      requests = append(requests, request)

//      stub.PutState(requestId, request)
//      stub.PutState("allRequest", requests)
// }

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
     
	 // if function != "queryAll" {
	 // 	requests, err := stub.GetState("allRequest")
  //       if err != nil {
  //       	return nil, errors.New("error happened")
  //       }

  //       if requests == nil {
  //       	return nil, errors.New("The all requests are empty!")
  //       }

         //return requests, nil
	 // }
     return nil, nil
}