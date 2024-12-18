package main

import (
	"healthcare/contracts"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	patientContract := new(contracts.PatientContract)
	insuranceContract := new(contracts.InsuranceContract)

	chaincode, err := contractapi.NewChaincode(patientContract, insuranceContract)

	if err != nil {
		log.Panicf("Could not create chaincode : %v", err)
	}

	err = chaincode.Start()

	if err != nil {
		log.Panicf("Failed to start chaincode : %v", err)
	}
}
