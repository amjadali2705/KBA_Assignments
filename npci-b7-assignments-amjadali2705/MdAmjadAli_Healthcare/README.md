# Healthcare Network

### Create the folder structure
```
mkdir Healthcare-network
```
```
cd Healthcare-network
```
### Generate the certificates using fabric-ca
```
mkdir docker
```
### Register the ca admin for each organization
Build the docker-compose-ca.yaml in the docker folder
```
docker compose -f docker/docker-compose-ca.yaml up -d
```
```
docker ps -a
```
```
sudo chmod -R 777 organizations/
```
### Register and enroll the users for each organization
Build the registerEnroll.sh script file
```
chmod +x registerEnroll.sh
```
```
./registerEnroll.sh
```
### Build the infrastructure
Build the docker-compose-3org.yaml in the docker folder
```
docker compose -f docker/docker-compose-3org.yaml up -d
```
```
docker ps -a
```
### Generate the channel artifacts
mkdir config
#
Build the configtx.yaml file in the config folder
```
export FABRIC_CFG_PATH=./config

export CHANNEL_NAME=carechannel
```
```
configtxgen -profile ThreeOrgsChannel -outputBlock ./channel-artifacts/${CHANNEL_NAME}.block -channelID $CHANNEL_NAME
```
### Add the orderer to the channel
```
export ORDERER_CA=./organizations/ordererOrganizations/care.com/orderers/orderer.care.com/msp/tlscacerts/tlsca.care.com-cert.pem

export ORDERER_ADMIN_TLS_SIGN_CERT=./organizations/ordererOrganizations/care.com/orderers/orderer.care.com/tls/server.crt

export ORDERER_ADMIN_TLS_PRIVATE_KEY=./organizations/ordererOrganizations/care.com/orderers/orderer.care.com/tls/server.key
```
```
osnadmin channel join --channelID $CHANNEL_NAME --config-block ./channel-artifacts/$CHANNEL_NAME.block -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
```
```
osnadmin channel list -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
```
### Join peers to the channel
mkdir peercfg
#
Build the core.yaml in peercfg folder
#
Open another terminal with in Healthcare-network folder, let's call this terminal as peer0_Hospital terminal.
### ############## peer0_Hospital terminal ##############
```
export FABRIC_CFG_PATH=./peercfg
export CHANNEL_NAME=carechannel
export CORE_PEER_LOCALMSPID=HospitalMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/hospital.care.com/peers/peer0.hospital.care.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/hospital.care.com/users/Admin@hospital.care.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/care.com/orderers/orderer.care.com/msp/tlscacerts/tlsca.care.com-cert.pem
export HOSPITAL_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/hospital.care.com/peers/peer0.hospital.care.com/tls/ca.crt
export PATIENT_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/patient.care.com/peers/peer0.patient.care.com/tls/ca.crt
export INSURANCE_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/insurance.care.com/peers/peer0.insurance.care.com/tls/ca.crt
```
### Join peer0_Hospital to the channel
```
peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block
```
```
peer channel list
```
#
Open another terminal with in Healthcare-network folder, let's call this terminal as peer0_Insurance terminal.
### ############## peer0_Insurance terminal ##############
```
export FABRIC_CFG_PATH=./peercfg
export CHANNEL_NAME=carechannel 
export CORE_PEER_LOCALMSPID=InsuranceMSP 
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ADDRESS=localhost:9051 
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/insurance.care.com/peers/peer0.insurance.care.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/insurance.care.com/users/Admin@insurance.care.com/msp
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/care.com/orderers/orderer.care.com/msp/tlscacerts/tlsca.care.com-cert.pem
export HOSPITAL_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/hospital.care.com/peers/peer0.hospital.care.com/tls/ca.crt
export INSURANCE_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/insurance.care.com/peers/peer0.insurance.care.com/tls/ca.crt
export PATIENT_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/patient.care.com/peers/peer0.patient.care.com/tls/ca.crt
```
### Join peer0_Insurance to the channel
```
peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block

peer channel list
```
#
Open another terminal with in Healthcare-network folder, let's call this terminal as peer0_Patient terminal.
### ############## peer0_Patient terminal ##############
```
export FABRIC_CFG_PATH=./peercfg
export CHANNEL_NAME=carechannel 
export CORE_PEER_LOCALMSPID=PatientMSP 
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ADDRESS=localhost:11051 
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/patient.care.com/peers/peer0.patient.care.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/patient.care.com/users/Admin@patient.care.com/msp
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/care.com/orderers/orderer.care.com/msp/tlscacerts/tlsca.care.com-cert.pem
export HOSPITAL_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/hospital.care.com/peers/peer0.hospital.care.com/tls/ca.crt
export INSURANCE_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/insurance.care.com/peers/peer0.insurance.care.com/tls/ca.crt
export PATIENT_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/patient.care.com/peers/peer0.patient.care.com/tls/ca.crt
```
### Join peer0_Patient to the channel
```
peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block

peer channel list
```
### Anchor peer update
### ############## peer0_Hospital terminal ##############
```
peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
```
```
cd channel-artifacts
```
```
configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json

cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.HospitalMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.hospital.care.com","port": 7051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id ${CHANNEL_NAME} --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb
```
#
cd $~$..
```
peer channel update -f channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.care.com --tls --cafile $ORDERER_CA
```
### ############## peer0_Insurance terminal ##############
```
peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
```
```
cd channel-artifacts
```
```
cconfigtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.InsuranceMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.insurance.care.com","port": 9051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb
```
#
cd $~$..
```
peer channel update -f channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.care.com --tls --cafile $ORDERER_CA
```
### ############## peer0_Patient terminal ##############
```
peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
```
```
cd channel-artifacts
```
```
configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.PatientMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.patient.care.com","port": 11051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb
```
#
cd $~$..
```
peer channel update -f channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.care.com --tls --cafile $ORDERER_CA
```
```
peer channel getinfo -c $CHANNEL_NAME
```

### —-----------------Chaincode lifecycle—-------------------
###### ***Build the chaincode
### **************** peer0_Hospital terminal ******************
```
peer lifecycle chaincode package healthcare.tar.gz --path ../Chaincode/ --lang golang --label healthcare_1.0
```
```
peer lifecycle chaincode install healthcare.tar.gz
```
```
peer lifecycle chaincode queryinstalled
```
```
export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid healthcare.tar.gz)
```
### **************** peer0_Insurance terminal ******************
```
peer lifecycle chaincode install healthcare.tar.gz
```
```
export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid healthcare.tar.gz)
```
### **************** peer0_Patient terminal ******************
```
peer lifecycle chaincode install healthcare.tar.gz
```
```
export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid healthcare.tar.gz)
```
### **************** peer0_Hospital terminal ******************
```
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com --channelID $CHANNEL_NAME --name Healthcare --version 1.0 --collections-config ../Chaincode/collection-healthcare.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
```
### **************** peer0_Insurance terminal ******************
```
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com --channelID $CHANNEL_NAME --name Healthcare --version 1.0 --collections-config ../Chaincode/collection-healthcare.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
```
### **************** peer0_Patient terminal ******************
```
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com --channelID $CHANNEL_NAME --name Healthcare --version 1.0 --collections-config ../Chaincode/collection-healthcare.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
```
### **************** peer0_Hospital terminal ******************
```
peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name Healthcare --version 1.0 --sequence 1 --collections-config ../Chaincode/collection-healthcare.json --tls --cafile $ORDERER_CA --output json
```
```
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com --channelID $CHANNEL_NAME --name Healthcare --version 1.0 --sequence 1 --collections-config ../Chaincode/collection-healthcare.json --tls --cafile $ORDERER_CA --peerAddresses localhost:7051 --tlsRootCertFiles $HOSPITAL_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $INSURANCE_PEER_TLSROOTCERT
```
```
peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name Healthcare --cafile $ORDERER_CA
```
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n Healthcare --peerAddresses localhost:7051 --tlsRootCertFiles $HOSPITAL_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $INSURANCE_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $PATIENT_PEER_TLSROOTCERT -c '{"function":"AddPatient","Args":["Patient-01", "Amjad", "22", "Typhoid"]}'
```
```
peer chaincode query -C $CHANNEL_NAME -n Healthcare -c '{"Args":["GetPatient", "Patient-01"]}'
```
#### Functions for Patient Contract
#### 1. AddPatient
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n Healthcare --peerAddresses localhost:7051 --tlsRootCertFiles $HOSPITAL_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $INSURANCE_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $PATIENT_PEER_TLSROOTCERT -c '{"function":"AddPatient","Args":["Patient-01", "Amjad", "22", "Typhoid"]}'
```
#### 2. GetPatient
```
peer chaincode query -C $CHANNEL_NAME -n Healthcare -c '{"Args":["GetPatient", "Patient-01"]}'
```
#### 3. SubmitClaim
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n Healthcare --peerAddresses localhost:7051 --tlsRootCertFiles $HOSPITAL_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $INSURANCE_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $PATIENT_PEER_TLSROOTCERT -c '{"function":"SubmitClaim","Args":["Claim-01", "Patient-01", "Surgery", "200000"]}'
```
#### 4. GetClaim
```
peer chaincode query -C $CHANNEL_NAME -n Healthcare -c '{"Args":["GetClaim", "Claim-01"]}'
```
#### 5. RemovePatient
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n Healthcare --peerAddresses localhost:7051 --tlsRootCertFiles $HOSPITAL_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $INSURANCE_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $PATIENT_PEER_TLSROOTCERT -c '{"function":"RemovePatient","Args":["Patient-01"]}'
```
#### 6. GetPatientByRange
```
peer chaincode query -C $CHANNEL_NAME -n Healthcare -c '{"Args":["GetPatientByRange", "Patient-01", "Patient-03"]}'
```
#### 7. GetAllPatients
```
peer chaincode query -C $CHANNEL_NAME -n Healthcare -c '{"Args":["GetAllPatients"]}'
```
#### 8. GetPatientHistory
```
peer chaincode query -C $CHANNEL_NAME -n Healthcare -c '{"Args":["GetPatientHistory", "Patient-01"]}'
```
#### 9. GetPatientsWithPagination
```
peer chaincode query -C $CHANNEL_NAME -n Healthcare -c '{"Args":["GetPatientsWithPagination", "3", ""]}'
```
### --------Invoke Private Transaction----------
#### **************** peer0_Insurance terminal ******************
```
export POLICY=$(echo -n "HealthInsurance" | base64 | tr -d \\n)

export INSUREDAMT=$(echo -n "500000" | base64 | tr -d \\n)

export STATUS=$(echo -n "Approved" | base64 | tr -d \\n)
```
#### Functions for Insurance Contract
#### 1. CreateInsurance
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n Healthcare --peerAddresses localhost:7051 --tlsRootCertFiles $HOSPITAL_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $INSURANCE_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $PATIENT_PEER_TLSROOTCERT -c '{"Args":["InsuranceContract:CreateInsurance","Ins-01"]}' --transient "{\"policy\":\"$POLICY\",\"insuredAmt\":\"$INSUREDAMT\",\"status\":\"$STATUS\"}"
```
#### 2. SubmitClaimApproval
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n Healthcare --peerAddresses localhost:7051 --tlsRootCertFiles $HOSPITAL_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $INSURANCE_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $PATIENT_PEER_TLSROOTCERT -c '{"Args":["InsuranceContract:SubmitClaimApproval","Claim-01", "Approved"]}'
```
#### 3. ReadInsurance
```
peer chaincode query -C $CHANNEL_NAME -n Healthcare -c '{"Args":["InsuranceContract:ReadInsurance","Ins-01"]}'
```
#### 4. DeleteInsurance
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n Healthcare --peerAddresses localhost:7051 --tlsRootCertFiles $HOSPITAL_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $INSURANCE_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $PATIENT_PEER_TLSROOTCERT -c '{"Args":["InsuranceContract:DeleteInsurance","Ins-01"]}'
```
#### 5. GetInsuranceByRange
```
peer chaincode query -C $CHANNEL_NAME -n Healthcare -c '{"Args":["InsuranceContract:GetInsuranceByRange", "Ins-02", "Ins-04"]}'
```
#### 6. GetAllInsurance
```
peer chaincode query -C $CHANNEL_NAME -n Healthcare -c '{"Args":["InsuranceContract:GetAllInsurance"]}'
```

## Perform Transaction using Client App
### Start the network
To start the network execute this command 
```
./startHealthcareNetwork.sh
```
### Create the client folder
```
cd ..
```
```
mkdir Client
```
### Build the client application
```
cd Client

go mod init client

Create and build profile.go, connect.go, client.go, main.go

go get google.golang.org/grpc@v1.67.1

go mod tidy

go run .
```

## Events
### #Create a Events folder with in the MdAmjadAli_Healthcare directory
```
mkdir Events
```
### #Build the events code
```
cd Events

go mod init events

Create & Build profile.go, connect.go, events.go, main.go

go mod tidy
```
### Set a patient addition event in the AddPatient function in patient-contract.go
```
type EventData struct{
	Type string
	Name string
}
```
```
addPatientEventData := EventData{
	Type:  "Patient Addition",
	Name: name,
}

eventDataByte, _ := json.Marshal(addPatientEventData)

ctx.GetStub().SetEvent("Add Patient", eventDataByte)
```
### Start the network**
```
cd ..

cd Healthcare-network/
```
To start the automobile network execute this command
```
./startHealthcareNetwork.sh
```
### To Run Block Event
Open a terminal in the events folder & consider this terminal as to listening block events.
#
go run $~$.
#
Note: For checking newly created block do a transaction using client application.(open a terminal from Client folder and execute go run .)
### To Run Chaincode Event
open another terminal from Events folder and consider it as to listening chaincode events.
#
go run $~$.
#
Note: Do a add patient transaction using client application.( change AddPatient transaction arguments in main.go then execute go run . in the client terminal)
### To Run Private Block Event
Open a new terminal from Events folder and consider it to listening private blockevent.
#
go run $~$.
#
Note: Submit a CreateInsurance transaction using client application.(Edit main.go for CreateInsurance transaction and execute go run . in the client terminal)