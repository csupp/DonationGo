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

    "log"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Donation struct {
    Id string `json:"id"`
    Who string `json:"who"`
    Rid string  `json:"rid"`
    Money int   `json:"money"`
}

type Request struct {
    Id string `json:"id"`
    Name string  `json:"name"`
    Description string `json:"description"`
    ExpectedMoney int `json:"expectedMoney"`
    CurrentMoney int  `json:"currentMoney"`
    DonationList []string `json:"donationList"`
}


type Person struct {
    Id string `json:"id"`
    Name string `json:"name"`
    MyRequests []string `json:"myRequests"`
    MyDonations []string `json:"myDonations"`
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

    var request Request
    var donationLts []string
    request = Request{Id: "requestid", Name: "Donation Go", Description: "Wanna to go to University", ExpectedMoney: 10000, CurrentMoney: 0, DonationList: donationLts}
    rjson, err := json.Marshal(&request)
    if err != nil {
        return nil, err
    }
    stub.PutState("requestid", rjson)
    log.Println("init function has done!")
    return nil, nil
}

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

     if function == "createDonation" {
        return t.createDonation(stub, args)
     }
     return nil, errors.New("Received unknown function invocation")
}

func (t *SimpleChaincode) createDonation(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
     //args: ["jack", "requestid", money] 
     var from, toRid string
     var money int
     var err error
   
     if len(args) != 3 {
         return nil, errors.New("My createDonation. Incorrect number of arguments. Expecting 3")
     }
     from = args[0]
     toRid = args[1]
     money, err = strconv.Atoi(args[2])
     if err != nil {
        return nil, errors.New("money cannot convert to number")
     }

     var donation Donation
     donation = Donation{Id: "donationid", Rid: toRid, Who: from, Money: money}
     djson, err := json.Marshal(&donation)
     if err != nil {
        return nil, err
     }
     var a = donation.Id
     stub.PutState(a, djson)
     
     
     
     var person Person
     var myReqs, myDons []string
     // update person data
     personByte, err := stub.GetState(from)
     if err != nil {
        fmt.Println("No person value for " + from)
        person = Person{Id: from, Name: from, MyRequests: myReqs, MyDonations: myDons}
        var pid2 = person.Id
        pJson, err := json.Marshal(&person)
        if err != nil {
            return nil, errors.New("failed to JSON person instance")
        }
        stub.PutState(pid2, pJson)
     } else {
        err = json.Unmarshal(personByte, &person)
        if err != nil {
            return nil, errors.New("failed to Unmarshal person instance")
        }
     }
    
    myDonations := person.MyDonations
    if myDonations == nil {
        myDonations = make([]string, 0)
    }
    myDonations = append(myDonations, donation.Id)
    person.MyDonations = myDonations
    
    requestByte, err := stub.GetState(toRid)
    if err != nil {
           return nil, errors.New("request did not exist")
    }

    var request Request
    err = json.Unmarshal(requestByte, &request)
    if err != nil {
           return nil, errors.New("failed to Unmarshal request instance")
    }
    request.CurrentMoney += money
    donationList := request.DonationList 
    if donationList == nil {
        donationList = make([]string, 0)
    }
    donationList = append(donationList, donation.Id)
    request.DonationList = donationList
    return []byte("create donation has finished"), nil     
}

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    log.Println("query is running " + function)
    log.Println(function)
    log.Println(args[0])
    // Handle different functions
    if function == "read" {                            //read a variable
        return t.read(stub, args)
    }
    
    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query")
}

func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    log.Println("Get into read function")
 
    var key, jsonResp string
    var err error

    key = args[0]
    valAsbytes, err := stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }
    if valAsbytes == nil {
        return []byte("cannot find the key's value of the chaincode"), nil
    }

    return valAsbytes, nil
}
