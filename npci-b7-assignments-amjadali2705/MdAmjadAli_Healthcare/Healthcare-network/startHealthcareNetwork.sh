#!/bin/bash

# Register the CA admin for each organization
echo "------------Register the ca admin for each organization----------------"
docker compose -f docker/docker-compose-ca.yaml up -d
sleep 3

# Set permissions for the organizations directory
sudo chmod -R 777 organizations/

echo "------------Register and enroll the users for each organization-----------"
chmod +x registerEnroll.sh
./registerEnroll.sh
sleep 3

# Build the infrastructure
echo "------------Build the infrastructure-----------------"
docker compose -f docker/docker-compose-3org.yaml up -d
sleep 3

# Generate the genesis block for the channel
echo "------------Generate the genesis block-------------------------------"
export FABRIC_CFG_PATH=${PWD}/config
export CHANNEL_NAME=carechannel
configtxgen -profile ThreeOrgsChannel -outputBlock ${PWD}/channel-artifacts/${CHANNEL_NAME}.block -channelID $CHANNEL_NAME
sleep 2

# Create the application channel
echo "------------Create the application channel----------------------------"
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/care.com/orderers/orderer.care.com/msp/tlscacerts/tlsca.care.com-cert.pem
export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/care.com/orderers/orderer.care.com/tls/server.crt
export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/care.com/orderers/orderer.care.com/tls/server.key

osnadmin channel join --channelID $CHANNEL_NAME --config-block ${PWD}/channel-artifacts/$CHANNEL_NAME.block -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
sleep 2

# List the channels
osnadmin channel list -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
sleep 2

# Set environment variables for the Hospital peer
export FABRIC_CFG_PATH=${PWD}/peercfg
export CORE_PEER_LOCALMSPID=HospitalMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/hospital.care.com/peers/peer0.hospital.care.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/hospital.care.com/users/Admin@hospital.care.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export HOSPITAL_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/hospital.care.com/peers/peer0.hospital.care.com/tls/ca.crt
export INSURANCE_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/insurance.care.com/peers/peer0.insurance.care.com/tls/ca.crt
export PATIENT_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/patient.care.com/peers/peer0.patient.care.com/tls/ca.crt
sleep 2

# Join Hospital peer to the channel
echo "------------Join Hospital peer to the channel------------"
peer channel join -b ${PWD}/channel-artifacts/${CHANNEL_NAME}.block
sleep 3

# List the channels
echo "------------Channel List------------"
peer channel list

# Update Hospital anchor peer
echo "------------Hospital anchor peer update----------------"
peer channel fetch config ${PWD}/channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
sleep 1

cd channel-artifacts
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
cd ..

# Update the anchor peer in the channel
peer channel update -f ${PWD}/channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com --tls --cafile $ORDERER_CA
sleep 1

# Package chaincode
echo "------------Package chaincode----------------"
peer lifecycle chaincode package healthcare.tar.gz --path ${PWD}/../Chaincode/ --lang golang --label healthcare_1.0
sleep 1

# Install chaincode on Hospital peer
echo "------------Install chaincode on Hospital peer----------------"
peer lifecycle chaincode install healthcare.tar.gz
sleep 3

# Query installed chaincode
peer lifecycle chaincode queryinstalled
sleep 1

# Get chaincode package ID
export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid healthcare.tar.gz)

# Approve chaincode on Hospital peer
echo "------------Approve chaincode on Hospital peer----------------"
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com --channelID $CHANNEL_NAME --name Healthcare --version 1.0 --collections-config ../Chaincode/collection-healthcare.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
sleep 2

# Set environment variables for Insurance peer
export CORE_PEER_LOCALMSPID=InsuranceMSP
export CORE_PEER_ADDRESS=localhost:9051
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/insurance.care.com/peers/peer0.insurance.care.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/insurance.care.com/users/Admin@insurance.care.com/msp

# Join Insurance peer to the channel
echo "------------Join Insurance peer to the channel------------"
peer channel join -b ${PWD}/channel-artifacts/$CHANNEL_NAME.block
sleep 1

# Update Insurance anchor peer
echo "------------Insurance anchor peer update----------------"
peer channel fetch config ${PWD}/channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
sleep 1

cd channel-artifacts
configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json config_copy.json
jq '.channel_group.groups.Application.groups.InsuranceMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.insurance.care.com","port": 9051}]},"version": "0"}}' config_copy.json > modified_config.json
configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output config_update.pb
configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb
cd ..

# Update the anchor peer in the channel
peer channel update -f ${PWD}/channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com --tls --cafile $ORDERER_CA
sleep 1

# Install chaincode on Insurance peer
echo "------------Install chaincode on Insurance peer----------------"
peer lifecycle chaincode install healthcare.tar.gz
sleep 3

# Query installed chaincode
peer lifecycle chaincode queryinstalled
sleep 1

# Approve chaincode on Insurance peer
echo "------------Approve chaincode on Insurance peer----------------"
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com --channelID $CHANNEL_NAME --name Healthcare --version 1.0 --collections-config ../Chaincode/collection-healthcare.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
sleep 1

# Set environment variables for Patient peer
export CORE_PEER_LOCALMSPID=PatientMSP
export CORE_PEER_ADDRESS=localhost:11051
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/patient.care.com/peers/peer0.patient.care.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/patient.care.com/users/Admin@patient.care.com/msp

# Join Patient peer to the channel
echo "------------Join Patient peer to the channel------------"
peer channel join -b ${PWD}/channel-artifacts/$CHANNEL_NAME.block
sleep 1

# Update Patient anchor peer
echo "------------Patient anchor peer update----------------"
peer channel fetch config ${PWD}/channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
sleep 1

cd channel-artifacts
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
cd ..

# Update the anchor peer in the channel
peer channel update -f ${PWD}/channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com --tls --cafile $ORDERER_CA
sleep 1

# Install chaincode on Patient peer
echo "------------Install chaincode on Patient peer----------------"
peer lifecycle chaincode install healthcare.tar.gz
sleep 3

# Query installed chaincode
peer lifecycle chaincode queryinstalled
sleep 1

# Approve chaincode on Patient peer
echo "------------Approve chaincode on Patient peer----------------"
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com --channelID $CHANNEL_NAME --name Healthcare --version 1.0 --collections-config ../Chaincode/collection-healthcare.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
sleep 1

# Commit chaincode
echo "------------Commit chaincode----------------"
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.care.com --channelID $CHANNEL_NAME --name Healthcare --version 1.0 --sequence 1 --collections-config ../Chaincode/collection-healthcare.json --tls --cafile $ORDERER_CA --peerAddresses localhost:7051 --tlsRootCertFiles $HOSPITAL_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $INSURANCE_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $PATIENT_PEER_TLSROOTCERT
sleep 1

# Query the committed chaincode
peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name Healthcare --cafile $ORDERER_CA
