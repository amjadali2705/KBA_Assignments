package main

import (
	"insuranceclaim/contracts"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	accidentReportContract := new(contracts.AccidentReportContract)
	insurancePolicyContract := new(contracts.InsurancePolicyContract)

	chaincode, err := contractapi.NewChaincode(accidentReportContract, insurancePolicyContract)

	if err != nil {
		log.Panicf("Could not create chaincode : %v", err)
	}

	err = chaincode.Start()

	if err != nil {
		log.Panicf("Failed to start chaincode : %v", err)
	}
}
