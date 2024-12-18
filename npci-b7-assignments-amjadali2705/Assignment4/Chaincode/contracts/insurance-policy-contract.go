package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type InsurancePolicyContract struct {
	contractapi.Contract
}

type InsurancePolicy struct {
	AssetType     string `json:"assetType"`
	PolicyID      string `json:"policyID"`
	PolicyHolder  string `json:"policyHolder"`
	InsuredAmount string `json:"insuredAmount"`
	VehicleID     string `json:"vehicleID"`
	Status        string `json:"status"`
}

const collectionName string = "InsurancePolicyCollection"

func (i *InsurancePolicyContract) InsurancePolicyExists(ctx contractapi.TransactionContextInterface, policyID string) (bool, error) {

	data, err := ctx.GetStub().GetPrivateDataHash(collectionName, policyID)

	if err != nil {
		return false, fmt.Errorf("could not fetch the private data hash. %s", err)
	}

	return data != nil, nil
}

func (i *InsurancePolicyContract) CreateInsurancePolicy(ctx contractapi.TransactionContextInterface, policyID string) (string, error) {

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("could not fetch client identity. %s", err)
	}

	// if clientOrgID == "insurancecompany-insuranceclaim-com" {
	if clientOrgID == "Org2MSP" {
		exists, err := i.InsurancePolicyExists(ctx, policyID)
		if err != nil {
			return "", fmt.Errorf("could not read from world state. %s", err)
		} else if exists {
			return "", fmt.Errorf("the asset %s already exists", policyID)
		}

		var policy InsurancePolicy

		transientData, err := ctx.GetStub().GetTransient()
		if err != nil {
			return "", fmt.Errorf("could not fetch transient data. %s", err)
		}

		if len(transientData) == 0 {
			return "", fmt.Errorf("please provide the private data of policyHolder, insuredAmount, vehicleID, status")
		}

		policyHolder, exists := transientData["policyHolder"]
		if !exists {
			return "", fmt.Errorf("the policyHolder was not specified in transient data. Please try again")
		}
		policy.PolicyHolder = string(policyHolder)

		insuredAmount, exists := transientData["insuredAmount"]
		if !exists {
			return "", fmt.Errorf("the insuredAmount was not specified in transient data. Please try again")
		}
		policy.InsuredAmount = string(insuredAmount)

		vehicleID, exists := transientData["vehicleID"]
		if !exists {
			return "", fmt.Errorf("the vehicleID was not specified in transient data. Please try again")
		}
		policy.VehicleID = string(vehicleID)

		status, exists := transientData["status"]
		if !exists {
			return "", fmt.Errorf("the status was not specified in transient data. Please try again")
		}
		policy.Status = string(status)

		policy.AssetType = "InsurancePolicy"
		policy.PolicyID = policyID

		bytes, _ := json.Marshal(policy)
		err = ctx.GetStub().PutPrivateData(collectionName, policyID, bytes)
		if err != nil {
			return "", fmt.Errorf("could not able to write the data")
		}
		return fmt.Sprintf("Insurance Policy with ID %v added successfully", policyID), nil
	} else {
		return fmt.Sprintf("Insurance policy cannot be created by organisation with MSPID %v ", clientOrgID), nil
	}
}

func (i *InsurancePolicyContract) ReadInsurancePolicy(ctx contractapi.TransactionContextInterface, policyID string) (*InsurancePolicy, error) {
	exists, err := i.InsurancePolicyExists(ctx, policyID)
	if err != nil {
		return nil, fmt.Errorf("could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("the asset %s does not exist", policyID)
	}

	bytes, err := ctx.GetStub().GetPrivateData(collectionName, policyID)
	if err != nil {
		return nil, fmt.Errorf("could not get the private data. %s", err)
	}
	var policy InsurancePolicy

	err = json.Unmarshal(bytes, &policy)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal private data collection data to type InsurancePolicy")
	}

	return &policy, nil
}

func (c *InsurancePolicyContract) ApproveInsurancePolicy(ctx contractapi.TransactionContextInterface, policyID string) (string, error) {
	// Get the client organization's MSPID to ensure only the insurance company can approve the policy
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("could not fetch client identity. %s", err)
	}

	// Check if the client MSPID is "insurancecompany-insuranceclaim-com"
	// if clientOrgID == "insurancecompany-insuranceclaim-com" {
	if clientOrgID == "Org2MSP" {
		// Fetch the transient data (the status in this case)
		transientData, err := ctx.GetStub().GetTransient()
		if err != nil {
			return "", fmt.Errorf("could not retrieve transient data. %s", err)
		}

		// Check if 'status' is in the transient data
		statusBytes, exists := transientData["status"]
		if !exists {
			return "", fmt.Errorf("status not found in transient data")
		}

		// Unmarshal the status value from transient data
		var status string
		err = json.Unmarshal(statusBytes, &status)
		if err != nil {
			return "", fmt.Errorf("could not unmarshal status from transient data. %s", err)
		}

		// Fetch the private data (policy) from the private collection (e.g., InsurancePolicyCollection)
		policyBytes, err := ctx.GetStub().GetPrivateData("InsurancePolicyCollection", policyID)
		if err != nil {
			return "", fmt.Errorf("could not fetch private data for insurance policy %s. %s", policyID, err)
		}
		if policyBytes == nil {
			return "", fmt.Errorf("insurance policy %s does not exist", policyID)
		}

		// Unmarshal the policy data
		var policy InsurancePolicy
		err = json.Unmarshal(policyBytes, &policy)
		if err != nil {
			return "", fmt.Errorf("could not unmarshal policy data for policy %s. %s", policyID, err)
		}

		// Check if the policy is already approved or rejected
		if policy.Status == "Approved" {
			return "", fmt.Errorf("policy %s is already approved", policyID)
		} else if policy.Status == "Rejected" {
			return "", fmt.Errorf("policy %s is already rejected", policyID)
		}

		// Update the policy status with the value from transient data
		policy.Status = status

		// Marshal the updated policy back into JSON
		updatedPolicyBytes, err := json.Marshal(policy)
		if err != nil {
			return "", fmt.Errorf("could not marshal the updated policy data for %s. %s", policyID, err)
		}

		// Store the updated policy in the private collection
		err = ctx.GetStub().PutPrivateData("InsurancePolicyCollection", policyID, updatedPolicyBytes)
		if err != nil {
			return "", fmt.Errorf("could not approve the insurance policy %s. %s", policyID, err)
		}

		// Return success message
		return fmt.Sprintf("insurance policy %v has been successfully %s", policyID, status), nil
	} else {
		return "", fmt.Errorf("user under MSPID: %v cannot approve an insurance policy", clientOrgID)
	}
}


func (i *InsurancePolicyContract) DeleteInsurancePolicy(ctx contractapi.TransactionContextInterface, policyID string) error {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("could not read the client identity. %s", err)
	}

	// if clientOrgID == "insurancecompany-insuranceclaim-com" {
	if clientOrgID == "Org2MSP" {

		exists, err := i.InsurancePolicyExists(ctx, policyID)

		if err != nil {
			return fmt.Errorf("could not read from world state. %s", err)
		} else if !exists {
			return fmt.Errorf("the asset %s does not exist", policyID)
		}

		return ctx.GetStub().DelPrivateData(collectionName, policyID)
	} else {
		return fmt.Errorf("organisation with %v cannot delete the insurance policy", clientOrgID)
	}
}

func (i *InsurancePolicyContract) GetAllInsurancePolicies(ctx contractapi.TransactionContextInterface) ([]*InsurancePolicy, error) {
	queryString := `{"selector":{"assetType":"InsurancePolicy"}}`
	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult(collectionName, queryString)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the query result. %s", err)
	}
	defer resultsIterator.Close()
	return InsurancePolicyResultIteratorFunction(resultsIterator)
}

func (i *InsurancePolicyContract) GetInsurancePoliciesByRange(ctx contractapi.TransactionContextInterface, startKey string, endKey string) ([]*InsurancePolicy, error) {
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(collectionName, startKey, endKey)

	if err != nil {
		return nil, fmt.Errorf("could not fetch the private data by range. %s", err)
	}
	defer resultsIterator.Close()

	return InsurancePolicyResultIteratorFunction(resultsIterator)
}

func InsurancePolicyResultIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*InsurancePolicy, error) {
	var policies []*InsurancePolicy
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not fetch the details of result iterator. %s", err)
		}
		var policy InsurancePolicy
		err = json.Unmarshal(queryResult.Value, &policy)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal the data. %s", err)
		}
		policies = append(policies, &policy)
	}

	return policies, nil
}
