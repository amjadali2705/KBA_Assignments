package main

type Config struct {
	CertPath     string `json:"certPath"`
	KeyDirectory string `json:"keyPath"`
	TLSCertPath  string `json:"tlsCertPath"`
	PeerEndpoint string `json:"peerEndpoint"`
	GatewayPeer  string `json:"gatewayPeer"`
	MSPID        string `json:"mspID"`
}

var profile = map[string]Config{

	"hospital": {
		CertPath:     "../Healthcare-network/organizations/peerOrganizations/hospital.care.com/users/User1@hospital.care.com/msp/signcerts/cert.pem",
		KeyDirectory: "../Healthcare-network/organizations/peerOrganizations/hospital.care.com/users/User1@hospital.care.com/msp/keystore/",
		TLSCertPath:  "../Healthcare-network/organizations/peerOrganizations/hospital.care.com/peers/peer0.hospital.care.com/tls/ca.crt",
		PeerEndpoint: "localhost:7051",
		GatewayPeer:  "peer0.hospital.care.com",
		MSPID:        "HospitalMSP",
	},

	"insurance": {
		CertPath:     "../Healthcare-network/organizations/peerOrganizations/insurance.care.com/users/User1@insurance.care.com/msp/signcerts/cert.pem",
		KeyDirectory: "../Healthcare-network/organizations/peerOrganizations/insurance.care.com/users/User1@insurance.care.com/msp/keystore/",
		TLSCertPath:  "../Healthcare-network/organizations/peerOrganizations/insurance.care.com/peers/peer0.insurance.care.com/tls/ca.crt",
		PeerEndpoint: "localhost:9051",
		GatewayPeer:  "peer0.insurance.care.com",
		MSPID:        "InsuranceMSP",
	},

	"patient": {
		CertPath:     "../Healthcare-network/organizations/peerOrganizations/patient.care.com/users/User1@patient.care.com/msp/signcerts/cert.pem",
		KeyDirectory: "../Healthcare-network/organizations/peerOrganizations/patient.care.com/users/User1@patient.care.com/msp/keystore/",
		TLSCertPath:  "../Healthcare-network/organizations/peerOrganizations/patient.care.com/peers/peer0.patient.care.com/tls/ca.crt",
		PeerEndpoint: "localhost:11051",
		GatewayPeer:  "peer0.patient.care.com",
		MSPID:        "PatientMSP",
	},
}
