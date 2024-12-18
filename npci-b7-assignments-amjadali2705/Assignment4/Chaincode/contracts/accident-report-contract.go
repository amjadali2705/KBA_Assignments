package contracts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type AccidentReportContract struct {
	contractapi.Contract
}

type PaginatedQueryResult struct {
	Records             []*AccidentReport `json:"records"`
	FetchedRecordsCount int32             `json:"fetchedRecordsCount"`
	Bookmark            string            `json:"bookmark"`
}

type HistoryQueryResult struct {
	Record    *AccidentReport `json:"record"`
	TxId      string          `json:"txId"`
	Timestamp string          `json:"timestamp"`
	IsDelete  bool            `json:"isDelete"`
}

type AccidentReport struct {
	AssetType         string `json:"assetType"`
	ReportDate        string `json:"reportDate"`
	ReportDescription string `json:"reportDescription"`
	ReportId          string `json:"reportId"`
	VehicleNo         string `json:"vehicleNo"`
	VehicleType       string `json:"vehicleType"`
}

func (a *AccidentReportContract) AccidentReportExists(ctx contractapi.TransactionContextInterface, ReportID string) (bool, error) {
	data, err := ctx.GetStub().GetState(ReportID)

	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)

	}
	return data != nil, nil
}

func (a *AccidentReportContract) CreateAccidentReport(ctx contractapi.TransactionContextInterface, ReportID, ReportDate, ReportDescription, VehicleNo, VehicleType string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}

	// if clientOrgID == "police-insuranceclaim-com" {
	if clientOrgID == "Org1MSP" {

		exists, err := a.AccidentReportExists(ctx, ReportID)
		if err != nil {
			return "", fmt.Errorf("%s", err)
		} else if exists {
			return "", fmt.Errorf("the report, %s already exists", ReportID)
		}

		accidentReport := AccidentReport{
			AssetType:         "accident",
			ReportDate:        ReportDate,
			ReportDescription: ReportDescription,
			ReportId:          ReportID,
			VehicleNo:         VehicleNo,
			VehicleType:       VehicleType,
		}

		bytes, _ := json.Marshal(accidentReport)

		err = ctx.GetStub().PutState(ReportID, bytes)
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("successfully created report %v", ReportID), nil
		}

	} else {
		return "", fmt.Errorf("user under following MSPID: %v can't perform this action", clientOrgID)
	}
}

func (a *AccidentReportContract) DeleteAccidentReport(ctx contractapi.TransactionContextInterface, ReportID string) (string, error) {

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}
	// if clientOrgID == "police-insuranceclaim-com" {
	if clientOrgID == "Org1MSP" {

		exists, err := a.AccidentReportExists(ctx, ReportID)
		if err != nil {
			return "", fmt.Errorf("%s", err)
		} else if !exists {
			return "", fmt.Errorf("the report, %s does not exist", ReportID)
		}

		err = ctx.GetStub().DelState(ReportID)
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("report with id %v is deleted from the world state.", ReportID), nil
		}

	} else {
		return "", fmt.Errorf("user under following MSPID: %v can't perform this action", clientOrgID)
	}
}

func (a *AccidentReportContract) ReadAccidentReport(ctx contractapi.TransactionContextInterface, ReportID string) (*AccidentReport, error) {

	bytes, err := ctx.GetStub().GetState(ReportID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if bytes == nil {
		return nil, fmt.Errorf("the car %s does not exist", ReportID)
	}

	var accidentReport AccidentReport

	err = json.Unmarshal(bytes, &accidentReport)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal world state data to type Car")
	}

	return &accidentReport, nil
}

func accidentReportResultIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*AccidentReport, error) {
	var accidentReports []*AccidentReport
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not fetch the details of the result iterator. %s", err)
		}
		var accidentReport AccidentReport
		err = json.Unmarshal(queryResult.Value, &accidentReport)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal the data. %s", err)
		}
		accidentReports = append(accidentReports, &accidentReport)
	}

	return accidentReports, nil
}

func (a *AccidentReportContract) GetAccidentReportByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*AccidentReport, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the  data by range. %s", err)
	}
	defer resultsIterator.Close()

	return accidentReportResultIteratorFunction(resultsIterator)
}

func (a *AccidentReportContract) GetAllAccidentReports(ctx contractapi.TransactionContextInterface) ([]*AccidentReport, error) {

	queryString := `{"selector":{"assetType":"accident"}}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the query result. %s", err)
	}
	defer resultsIterator.Close()
	return accidentReportResultIteratorFunction(resultsIterator)
}

func (a *AccidentReportContract) GetAccidentReportHistory(ctx contractapi.TransactionContextInterface, ReportID string) ([]*HistoryQueryResult, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(ReportID)
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

		var accidentReport AccidentReport
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &accidentReport)
			if err != nil {
				return nil, err
			}
		} else {
			accidentReport = AccidentReport{
				ReportId: ReportID,
			}
		}

		timestamp := response.Timestamp.AsTime()

		formattedTime := timestamp.Format(time.RFC1123)

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: formattedTime,
			Record:    &accidentReport,
			IsDelete:  response.IsDelete,
		}
		records = append(records, &record)
	}

	return records, nil
}


func (a *AccidentReportContract) GetAccidentReportsWithPagination(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {

	queryString := `{"selector":{"assetType":"accident"}}`

	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, fmt.Errorf("could not get the accidentReport records. %s", err)
	}
	defer resultsIterator.Close()

	accidentReports, err := accidentReportResultIteratorFunction(resultsIterator)
	if err != nil {
		return nil, fmt.Errorf("could not return the accidentReport records %s", err)
	}

	return &PaginatedQueryResult{
		Records:             accidentReports,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}