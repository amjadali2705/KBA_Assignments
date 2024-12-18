package contracts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type PatientContract struct {
	contractapi.Contract
}

type PaginatedQueryResult struct {
	Records             []*Patient `json:"records"`
	FetchedRecordsCount int32  `json:"fetchedRecordsCount"`
	Bookmark            string `json:"bookmark"`
}

type HistoryQueryResult struct {
	Record    *Patient `json:"record"`
	TxId      string   `json:"txId"`
	Timestamp string   `json:"timestamp"`
	IsDelete  bool     `json:"isDelete"`
}

type Patient struct {
	AssetType      string `json:"assetType"`
	PatientID      string `json:"patientId"`
	Name           string `json:"name"`
	Age            string `json:"age"`
	MedicalHistory string `json:"medicalHistory"`
}

type Claim struct {
	ClaimID       string  `json:"claimId"`
	PatientID     string  `json:"patientId"`
	Treatment     string  `json:"treatment"`
	ClaimAmount   string `json:"claimAmount"`
	ClaimStatus   string  `json:"claimStatus"`
}

type EventData struct{
	Type string
	Name string
}

func (c *PatientContract) PatientExists(ctx contractapi.TransactionContextInterface, patientId string) (bool, error) {
	data, err := ctx.GetStub().GetState(patientId)

	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return data != nil, nil
}

func (c *PatientContract) AddPatient(ctx contractapi.TransactionContextInterface, patientId, name, age, medicalHistory string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()

	if err != nil {
		return "", fmt.Errorf("could not fetch client identity. %s", err)
	}

	if clientOrgID == "HospitalMSP" {

		exists, err := c.PatientExists(ctx, patientId)
		if err != nil {
			return "", fmt.Errorf("could not fetch the details from world state.%s", err)
		} else if exists {
			return "", fmt.Errorf("the car, %s already exists", patientId)
		}

		patient := Patient{
			AssetType:      "patient",
			PatientID:      patientId,
			Name:           name,
			Age:            age,
			MedicalHistory: medicalHistory,
		}

		bytes, _ := json.Marshal(patient)

		err = ctx.GetStub().PutState(patientId, bytes)
		if err != nil {
			return "", fmt.Errorf("could not add patient. %s", err)
		} else {

			addPatientEventData := EventData{
				Type:  "Patient Addition",
				Name: name,
			}

			eventDataByte, _ := json.Marshal(addPatientEventData)

			ctx.GetStub().SetEvent("Add Patient", eventDataByte)

			return fmt.Sprintf("successfully added patient %v", patientId), nil
		}

	} else {
		return "", fmt.Errorf("user under following MSPID: %v can't perform this action", clientOrgID)
	}
}

func (c *PatientContract) SubmitClaim(ctx contractapi.TransactionContextInterface, claimId, patientId, treatment, claimAmount string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("could not fetch client identity. %s", err)
	}

	if clientOrgID == "HospitalMSP" {

		exists, err := c.PatientExists(ctx, patientId)
		if err != nil {
			return "", fmt.Errorf("could not fetch patient details from world state. %s", err)
		}
		if !exists {
			return "", fmt.Errorf("the patient with ID %s does not exist", patientId)
		}

		claim := Claim{
			ClaimID:     claimId,
			PatientID:   patientId,
			Treatment:   treatment,
			ClaimAmount: claimAmount,
			ClaimStatus: "Pending", // Initial status is "Pending"
		}

		claimBytes, _ := json.Marshal(claim)
		err = ctx.GetStub().PutState(claimId, claimBytes)
		if err != nil {
			return "", fmt.Errorf("could not submit claim. %s", err)
		}

		return fmt.Sprintf("successfully submitted claim with ID %v for patient %v", claimId, patientId), nil
	} else {
		return "", fmt.Errorf("user under following MSPID: %v can't perform this action", clientOrgID)
	}
}

// GetClaim retrieves a claim based on the claimID
func (c *PatientContract) GetClaim(ctx contractapi.TransactionContextInterface, claimId string) (*Claim, error) {
	bytes, err := ctx.GetStub().GetState(claimId)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if bytes == nil {
		return nil, fmt.Errorf("the claim with ID %s does not exist", claimId)
	}

	var claim Claim
	err = json.Unmarshal(bytes, &claim)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal world state data to type Claim")
	}

	return &claim, nil
}

func (c *PatientContract) GetPatient(ctx contractapi.TransactionContextInterface, patientId string) (*Patient, error) {

	bytes, err := ctx.GetStub().GetState(patientId)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if bytes == nil {
		return nil, fmt.Errorf("the patient %s does not exist", patientId)
	}

	var patient Patient

	err = json.Unmarshal(bytes, &patient)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal world state data to type Patient")
	}

	return &patient, nil
}

func (c *PatientContract) RemovePatient(ctx contractapi.TransactionContextInterface, patientId string) (string, error) {

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("could not fetch client identity. %s", err)
	}

	if clientOrgID == "HospitalMSP" {

		exists, err := c.PatientExists(ctx, patientId)
		if err != nil {
			return "", fmt.Errorf("%s", err)
		} else if !exists {
			return "", fmt.Errorf("the patient, %s does not exist", patientId)
		}

		err = ctx.GetStub().DelState(patientId)
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("patient with id %v is deleted from the world state.", patientId), nil
		}

	} else {
		return "", fmt.Errorf("user under following MSPID: %v can't perform this action", clientOrgID)
	}
}

func (c *PatientContract) GetPatientByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*Patient, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the  data by range. %s", err)
	}
	defer resultsIterator.Close()

	return patientResultIteratorFunction(resultsIterator)
}

func (c *PatientContract) GetAllPatients(ctx contractapi.TransactionContextInterface) ([]*Patient, error) {

	queryString := `{"selector":{"assetType":"patient"}}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the query result. %s", err)
	}
	defer resultsIterator.Close()
	return patientResultIteratorFunction(resultsIterator)
}

func patientResultIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*Patient, error) {
	var patients []*Patient
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not fetch the details of the result iterator. %s", err)
		}
		var patient Patient
		err = json.Unmarshal(queryResult.Value, &patient)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal the data. %s", err)
		}
		patients = append(patients, &patient)
	}

	return patients, nil
}

func (c *PatientContract) GetPatientHistory(ctx contractapi.TransactionContextInterface, patientId string) ([]*HistoryQueryResult, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(patientId)
	if err != nil {
		return nil, fmt.Errorf("could not get the data. %s", err)
	}
	defer resultsIterator.Close()

	var records []*HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not get the value of resultsIterator. %s", err)
		}

		var patient Patient
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &patient)
			if err != nil {
				return nil, err
			}
		} else {
			patient = Patient{
				PatientID: patientId,
			}
		}

		timestamp := response.Timestamp.AsTime()

		formattedTime := timestamp.Format(time.RFC1123)

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: formattedTime,
			Record:    &patient,
			IsDelete:  response.IsDelete,
		}
		records = append(records, &record)
	}

	return records, nil
}

func (c *PatientContract) GetPatientsWithPagination(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {

	queryString := `{"selector":{"assetType":"patient"}}`

	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, fmt.Errorf("could not get the patient records. %s", err)
	}
	defer resultsIterator.Close()

	patients, err := patientResultIteratorFunction(resultsIterator)
	if err != nil {
		return nil, fmt.Errorf("could not return the patient records %s", err)
	}

	return &PaginatedQueryResult{
		Records:             patients,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}