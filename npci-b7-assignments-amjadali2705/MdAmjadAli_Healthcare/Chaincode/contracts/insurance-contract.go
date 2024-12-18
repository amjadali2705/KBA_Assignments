package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type InsuranceContract struct {
	contractapi.Contract
}

type Insurance struct {
	AssetType   string `json:"assetType"`
	InsuranceID string `json:"insuranceId"`
	Policy      string `json:"policy"`
	InsuredAmt  string `json:"insuredAmt"`
	Status      string `json:"status"`
}

type InsuranceClaim struct {
	ClaimID     string  `json:"claimId"`
	PatientID   string  `json:"patientId"`
	Treatment   string  `json:"treatment"`
	ClaimAmount string  `json:"claimAmount"`
	ClaimStatus string  `json:"claimStatus"` // "Pending", "Approved", "Rejected"
}

const collectionName string = "InsuranceCollection"

func (c *InsuranceContract) InsuranceExists(ctx contractapi.TransactionContextInterface, insuranceId string) (bool, error) {

	data, err := ctx.GetStub().GetPrivateDataHash(collectionName, insuranceId)

	if err != nil {
		return false, fmt.Errorf("could not fetch the private data hash. %s", err)
	}

	return data != nil, nil
}

func (c *InsuranceContract) CreateInsurance(ctx contractapi.TransactionContextInterface, insuranceId string) (string, error) {

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("could not fetch client identity. %s", err)
	}

	if clientOrgID == "InsuranceMSP" {
		exists, err := c.InsuranceExists(ctx, insuranceId)
		if err != nil {
			return "", fmt.Errorf("could not read from world state. %s", err)
		} else if exists {
			return "", fmt.Errorf("the asset %s already exists", insuranceId)
		}

		var insurance Insurance

		transientData, err := ctx.GetStub().GetTransient()
		if err != nil {
			return "", fmt.Errorf("could not fetch transient data. %s", err)
		}

		if len(transientData) == 0 {
			return "", fmt.Errorf("please provide the private data of policy, insuredAmt")
		}

		policy, exists := transientData["policy"]
		if !exists {
			return "", fmt.Errorf("the policy was not specified in transient data. Please try again")
		}
		insurance.Policy = string(policy)

		insuredAmt, exists := transientData["insuredAmt"]
		if !exists {
			return "", fmt.Errorf("the insuredAmt was not specified in transient data. Please try again")
		}
		insurance.InsuredAmt = string(insuredAmt)

		status, exists := transientData["status"]
		if !exists {
			return "", fmt.Errorf("the status was not specified in transient data. Please try again")
		}
		insurance.Status = string(status)

		insurance.AssetType = "insurance"
		insurance.InsuranceID = insuranceId

		bytes, _ := json.Marshal(insurance)
		err = ctx.GetStub().PutPrivateData(collectionName, insuranceId, bytes)
		if err != nil {
			return "", fmt.Errorf("could not able to write the data")
		}
		return fmt.Sprintf("insurance with id %v added successfully", insuranceId), nil
	} else {
		return fmt.Sprintf("insurance cannot be created by organisation with MSPID %v ", clientOrgID), nil
	}
}

// SubmitClaimApproval function to approve or reject a claim
func (c *InsuranceContract) SubmitClaimApproval(ctx contractapi.TransactionContextInterface, claimId, status string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("could not fetch client identity. %s", err)
	}

	if clientOrgID == "InsuranceMSP" {

		claimBytes, err := ctx.GetStub().GetState(claimId)
		if err != nil {
			return "", fmt.Errorf("failed to read claim from world state: %v", err)
		}
		if claimBytes == nil {
			return "", fmt.Errorf("claim with ID %s does not exist", claimId)
		}

		var insuranceClaim InsuranceClaim
		err = json.Unmarshal(claimBytes, &insuranceClaim)
		if err != nil {
			return "", fmt.Errorf("could not unmarshal claim data: %v", err)
		}

		if insuranceClaim.ClaimStatus != "Pending" {
			return "", fmt.Errorf("the claim is already %s", insuranceClaim.ClaimStatus)
		}

		if status != "Approved" && status != "Rejected" {
			return "", fmt.Errorf("invalid status. Allowed values are 'Approved' or 'Rejected'")
		}

		insuranceClaim.ClaimStatus = status

		claimBytes, err = json.Marshal(insuranceClaim)
		if err != nil {
			return "", fmt.Errorf("could not marshal updated claim data: %v", err)
		}

		err = ctx.GetStub().PutState(claimId, claimBytes)
		if err != nil {
			return "", fmt.Errorf("could not update claim status in world state: %v", err)
		}

		return fmt.Sprintf("Claim %s has been %s", claimId, status), nil
	} else {
		return "", fmt.Errorf("user under following MSPID: %v cannot perform this action", clientOrgID)
	}
}

func (c *InsuranceContract) ReadInsurance(ctx contractapi.TransactionContextInterface, insuranceId string) (*Insurance, error) {
	exists, err := c.InsuranceExists(ctx, insuranceId)
	if err != nil {
		return nil, fmt.Errorf("could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("the asset %s does not exist", insuranceId)
	}

	bytes, err := ctx.GetStub().GetPrivateData(collectionName, insuranceId)
	if err != nil {
		return nil, fmt.Errorf("could not get the private data. %s", err)
	}
	var insurance Insurance

	err = json.Unmarshal(bytes, &insurance)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal private data collection data to type Insurance")
	}

	return &insurance, nil
}

func (c *InsuranceContract) DeleteInsurance(ctx contractapi.TransactionContextInterface, insuranceId string) error {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("could not read the client identity. %s", err)
	}

	if clientOrgID == "InsuranceMSP" {

		exists, err := c.InsuranceExists(ctx, insuranceId)

		if err != nil {
			return fmt.Errorf("could not read from world state. %s", err)
		} else if !exists {
			return fmt.Errorf("the asset %s does not exist", insuranceId)
		}

		return ctx.GetStub().DelPrivateData(collectionName, insuranceId)
	} else {
		return fmt.Errorf("organisation with %v cannot delete the insurance", clientOrgID)
	}
}

func (c *InsuranceContract) GetAllInsurance(ctx contractapi.TransactionContextInterface) ([]*Insurance, error) {
	queryString := `{"selector":{"assetType":"insurance"}}`
	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult(collectionName, queryString)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the query result. %s", err)
	}
	defer resultsIterator.Close()
	return InsuranceResultIteratorFunction(resultsIterator)
}

func (c *InsuranceContract) GetInsuranceByRange(ctx contractapi.TransactionContextInterface, startKey string, endKey string) ([]*Insurance, error) {
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(collectionName, startKey, endKey)

	if err != nil {
		return nil, fmt.Errorf("could not fetch the private data by range. %s", err)
	}
	defer resultsIterator.Close()

	return InsuranceResultIteratorFunction(resultsIterator)
}

func InsuranceResultIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*Insurance, error) {
	var insurances []*Insurance
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not fetch the details of result iterator. %s", err)
		}
		var insurance Insurance
		err = json.Unmarshal(queryResult.Value, &insurance)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal the data. %s", err)
		}
		insurances = append(insurances, &insurance)
	}

	return insurances, nil
}
