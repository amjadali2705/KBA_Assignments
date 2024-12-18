# Fabric Network

### Create the folder structure
```
mkdir Fabric-network
```
```
cd Fabric-network
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
Build the docker-compose-2org.yaml in the docker folder
```
docker compose -f docker/docker-compose-2org.yaml up -d
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

export CHANNEL_NAME=mychannel
```
```
configtxgen -profile ChannelUsingRaft -outputBlock ./channel-artifacts/${CHANNEL_NAME}.block -channelID $CHANNEL_NAME
```
### Add the orderer to the channel
```
export ORDERER_CA=./organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

export ORDERER_ADMIN_TLS_SIGN_CERT=./organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt

export ORDERER_ADMIN_TLS_PRIVATE_KEY=./organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.key
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
Open another terminal with in Fabric-network folder, let's call this terminal as peer0_Org1 terminal.
### ############## peer0_Org1 terminal ##############
```
export FABRIC_CFG_PATH=./peercfg
export CHANNEL_NAME=mychannel
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export ORG1_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export ORG2_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
```
### Join peer0_org1 to the channel
```
peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block
```
```
peer channel list
```
#
Open another terminal with in Fabric-network folder, let's call this terminal as peer0_Org2 terminal.
### ############## peer0_Org2 terminal ##############
```
export FABRIC_CFG_PATH=./peercfg
export CHANNEL_NAME=mychannel 
export CORE_PEER_LOCALMSPID=Org2MSP 
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ADDRESS=localhost:9051 
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export ORG1_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export ORG2_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
```
### Join peer0_Org2 to the channel
```
peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block

peer channel list
```
### Anchor peer update
### ############## peer0_Org1 terminal ##############
```
peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
```
```
cd channel-artifacts
```
```
configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json

jq '.data.data[0].payload.data.config' config_block.json > config.json

cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.Org1MSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.org1.example.com","port": 7051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb

configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb

configtxlator compute_update --channel_id ${CHANNEL_NAME} --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json

echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json

configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb
```
#
cd ..
```
peer channel update -f channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA
```
### ############## peer0_Org2 terminal ##############
```
peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
```
```
cd channel-artifacts
```
```
configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json

jq '.data.data[0].payload.data.config' config_block.json > config.json

cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.Org2MSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.org2.example.com","port": 9051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb

configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb

configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json

echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json

configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb
```
#
cd ..
```
peer channel update -f channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA
```
```
peer channel getinfo -c $CHANNEL_NAME
```
### Chaincode lifecycle
Make sure that the chaincode is available in the Chaincode folder at the same location of Fabric-network folder. Edit the MSP IDs in the smart contract and create a new collection file.

### ############## peer0_Org1 terminal ##############
```
peer lifecycle chaincode package Insurance.tar.gz --path ../Chaincode/ --lang golang --label insurance_1.0
```
```
peer lifecycle chaincode install Insurance.tar.gz
```
```
peer lifecycle chaincode queryinstalled
```
### ############## peer0_Org2 terminal ##############
```
peer lifecycle chaincode install Insurance.tar.gz
```
```
peer lifecycle chaincode queryinstalled
```
### ############## peer0_Org1 terminal ##############
```
export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid Insurance.tar.gz)
```
```
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID $CHANNEL_NAME --name insurance --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --collections-config ../Chaincode/collection.json --tls --cafile $ORDERER_CA --waitForEvent
```
### ############## peer0_Org2 terminal ##############
```
export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid Insurance.tar.gz)
```
```
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID $CHANNEL_NAME --name insurance --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --collections-config ../Chaincode/collection.json --tls --cafile $ORDERER_CA --waitForEvent
```
### ############## peer0_Org1 terminal ##############
```
peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name insurance --version 1.0 --sequence 1 --collections-config ../Chaincode/collection.json --tls --cafile $ORDERER_CA --output json
```
```
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID $CHANNEL_NAME --name insurance --version 1.0 --sequence 1 --collections-config ../Chaincode/collection.json --tls --cafile $ORDERER_CA --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT
```
```
peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name insurance --cafile $ORDERER_CA
```
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n insurance --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT -c '{"function":"CreateAccidentReport","Args":["Report-11", "22/07/2024", "AccidentReport", "1234", "Car"]}'
```
```
peer chaincode query -C $CHANNEL_NAME -n insurance -c '{"Args":["GetAllAccidentReports"]}'
```
### ############## peer0_Org2 terminal ##############
```
peer chaincode query -C $CHANNEL_NAME -n insurance -c '{"Args":["GetAllAccidentReports"]}'
```