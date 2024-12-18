package main

import "fmt"

func main() {
	result := submitTxnFn(
		"hospital",
		"carechannel",
		"Healthcare",
		"PatientContract",
		"invoke",
		make(map[string][]byte),
		"AddPatient",
		"Patient-25",
		"Nikita",
		"22",
		"Dengue",
	)

	// result := submitTxnFn("hospital", "carechannel", "Healthcare", "PatientContract", "query", make(map[string][]byte), "GetPatient", "Patient-20")

	// privateData := map[string][]byte{
	// 	"policy":       []byte("Accident"),
	// 	"insuredAmt":      []byte("200000"),
	// 	"status":      []byte("pending"),
	// }

	// result := submitTxnFn("insurance", "carechannel", "Healthcare", "InsuranceContract", "private", privateData, "CreateInsurance", "Ins-22")

	//result := submitTxnFn("insurance", "carechannel", "Healthcare", "InsuranceContract", "query", make(map[string][]byte), "ReadInsurance", "Ins-03")

	// result := submitTxnFn("hospital", "carechannel", "Healthcare", "PatientContract", "query", make(map[string][]byte), "GetAllPatients")

	// result := submitTxnFn("manufacturer", "autochannel", "KBA-Automobile", "OrderContract", "query", make(map[string][]byte), "GetAllOrders")

	// result := submitTxnFn("manufacturer", "autochannel", "KBA-Automobile", "CarContract", "query", make(map[string][]byte), "GetMatchingOrders", "Car-06")

	// result := submitTxnFn("manufacturer", "autochannel", "KBA-Automobile", "CarContract", "invoke", make(map[string][]byte), "MatchOrder", "Car-06", "ORD-03")

	// result := submitTxnFn("mvd", "autochannel", "KBA-Automobile", "CarContract", "invoke", make(map[string][]byte), "RegisterCar", "Car-06", "Dani", "KL-01-CD-01")

	// result := submitTxnFn("manufacturer", "autochannel", "KBA-Automobile", "CarContract", "query", make(map[string][]byte), "ReadCar", "Car-06")

	fmt.Println(result)
}
