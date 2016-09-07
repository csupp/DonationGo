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
    "log"
    "strconv"
    "os/exec"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Donation struct {
    id string `json:"id"`
    who string `json:"who"`
    time string `json:"time"`
    rid string  `json:"rid"`
    money int   `json:"money"`
}

type Request struct {
    id string `json:"id"`
    name string  `json:"name"`
    description string `json:"description"`
    expectedMoney int `json:"expectedMoney"`
    currentMoney int  `json:"currentMoney"`
    donationList []string `json:"donationList"`
}


type Person struct {
    id string `json:"id"`
    name string `json:"name"`
    myRequests []string `json:"myRequests"`
    myDonations []string `json:"myDonations"`
}



func main() {
    err := shim.Start(new(SimpleChaincode))
    
    if err != nil {
        fmt.Printf("Error starting Simple chaincode: %s", err)
    }
}

func(t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting 1")
    }

    err := stub.PutState("hello_world", []byte(args[0]))
    if err != nil {
        return nil, err
    }

    //init some requests to test createDonation function
    //requestid, err := exec.Command("uuidgen").Output()
    // if err != nil {
    //     return nil, err
    // }
    var request Request
    var donationLts []string
    request = Request{id: "rid", name: "Donation Go", description: "Wanna to go to University", expectedMoney: 10000, currentMoney: 0, donationList: donationLts}
    rjson, err := json.Marshal(&request)
    if err != nil {
        return nil, err
    }
    stub.PutState("requestid", rjson)
    log.Println("init function has done!")
    return nil, nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
   
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

func (t *SimpleChaincode) createDonation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
     //args: ["jack", "toRequestId", money] 
     var from, toRid string
     var money int
     var djson, pjson []byte
     var myDonations []string
    // var personByte []byte
     var err error

     if len(args) != 3 {
         return nil, errors.New("Incorrect number of arguments. Expecting 3")
     }
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

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("query is running " + function)
    // Handle different functions
    if function == "read" {                            //read a variable
        return t.read(stub, args)
    }
    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query")
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

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    log.Println("Get into read function")
    var key, jsonResp string
    var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    key = args[0]
    valAsbytes, err := stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }
    if valAsbytes == nil {
        return []byte("cannot find the key's value of the chaincode"), nil
    }
    // var re Request
    // err = json.Unmarshal(valAsbytes, &re)

    return valAsbytes, nil
}
