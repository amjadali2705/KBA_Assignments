package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Patient struct {
	AssetType      string `json:"assetType"`
	PatientID      string `json:"patientId"`
	Name           string `json:"name"`
	Age            string `json:"age"`
	MedicalHistory string `json:"medicalHistory"`
}

type PatientData struct {
	AssetType      string `json:"assetType"`
	PatientID      string `json:"patientId"`
	Name           string `json:"name"`
	Age            string `json:"age"`
	MedicalHistory string `json:"medicalHistory"`
}

type Insurance struct {
	AssetType   string `json:"assetType"`
	InsuranceID string `json:"insuranceId"`
	Policy      string `json:"policy"`
	InsuredAmt  string `json:"insuredAmt"`
	Status      string `json:"status"`
}

type InsuranceData struct {
	AssetType   string `json:"assetType"`
	InsuranceID string `json:"insuranceId"`
	Policy      string `json:"policy"`
	InsuredAmt  string `json:"insuredAmt"`
	Status      string `json:"status"`
}

type HistoryQueryResult struct {
	Record    *Patient `json:"record"`
	TxId      string   `json:"txId"`
	Timestamp string   `json:"timestamp"`
	IsDelete  bool     `json:"isDelete"`
}

func main() {
	router := gin.Default()

	router.Static("/public", "./public")
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Hospital Dashboard",
		})

	})

	router.POST("/api/patient", func(ctx *gin.Context) {

		var req Patient
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
			return
		}
		fmt.Println("request", req)
		result := submitTxnFn(
			"hospital",
			"carechannel",
			"Healthcare",
			"PatientContract",
			"invoke",
			make(map[string][]byte),
			"AddPatient",
			req.PatientID,
			req.Name,
			req.Age,
			req.MedicalHistory,
		)
		ctx.JSON(http.StatusOK, gin.H{"message": "Added new patient", "result": result})

	})

	router.GET("/api/patient/:id", func(ctx *gin.Context) {
		patientId := ctx.Param("id")

		result := submitTxnFn("hospital", "carechannel", "Healthcare", "PatientContract", "query", make(map[string][]byte), "GetPatient", patientId)

		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	router.GET("/api/insurance/all", func(ctx *gin.Context) {

		result := submitTxnFn("insurance", "carechannel", "Healthcare", "InsuranceContract", "query", make(map[string][]byte), "GetAllInsurance")

		var insurances []InsuranceData

		if len(result) > 0 {
			// Unmarshal the JSON array string into the orders slice
			if err := json.Unmarshal([]byte(result), &insurances); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	router.POST("/api/insurance", func(ctx *gin.Context) {
		var req Insurance
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
			return
		}

		fmt.Printf("insurance  %s", req)

		privateData := map[string][]byte{
			"insuranceId": []byte(req.InsuranceID),
			"policy":      []byte(req.Policy),
			"insuredAmt":  []byte(req.InsuredAmt),
			"status":      []byte(req.Status),
		}

		submitTxnFn("insurance", "carechannel", "Healthcare", "InsuranceContract", "private", privateData, "CreateInsurance", req.InsuranceID)

		ctx.JSON(http.StatusOK, req)
	})

	router.GET("/api/insurance/:id", func(ctx *gin.Context) {
		insuranceId := ctx.Param("id")

		result := submitTxnFn("insurance", "carechannel", "Healthcare", "InsuranceContract", "query", make(map[string][]byte), "ReadInsurance", insuranceId)

		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	router.GET("/patient", func(ctx *gin.Context) {
		result := submitTxnFn("patient", "carechannel", "Healthcare", "PatientContract", "query", make(map[string][]byte), "GetAllPatients")

		var patients []PatientData

		if len(result) > 0 {
			// Unmarshal the JSON array string into the cars slice
			if err := json.Unmarshal([]byte(result), &patients); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	router.GET("/api/patient/history", func(ctx *gin.Context) {
		patientId := ctx.Query("patientId")
		result := submitTxnFn("patient", "carechannel", "Healthcare", "PatientContract", "query", make(map[string][]byte), "GetPatientHistory", patientId)

		// fmt.Printf("result %s", result)

		var patients []HistoryQueryResult

		if len(result) > 0 {
			// Unmarshal the JSON array string into the orders slice
			if err := json.Unmarshal([]byte(result), &patients); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	router.Run(":3000")

}
