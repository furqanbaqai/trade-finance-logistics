# Chaper 4: Running chain code

## Start chaincode
```
$ docker exec -it chaincode bash
# In chaincode:
$ CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=tw:0 ./trade_workflow_v1
```
## Install chaincode on the network
```
$ docker exec -it cli bash
# Install Chaincode
$ peer chaincode install -p chaincodedev/chaincode/trade_workflow_v1 -n tw -v 0
# Instantiate the chain code
$ peer chaincode instantiate -n tw -v 0 -c '{"Args":["init","LumberInc","LumberBank","100000","WoodenToys","ToyBank","200000","UniversalFreight","ForestryDepartment"]}' -C tradechannel
# Invoke Chaincode
$ peer chaincode invoke -n tw -c '{"Args":["requestTrade", "50000","15000","Wood for Toys"]}' -C tradechannel
$ peer chaincode invoke -n tw -c '{"Args":["getTradeStatus","50000"]}' -C tradechannel
```
## Section: ABAC (pg. no. 105)

Run following commands on dev network:
```
$ docker exec -it cli bash
root@af4edc280f83:~# fabric-ca-client enroll -u http://admin:adminpw@ca:7054
2018/12/03 16:02:03 [INFO] Created a default configuration file at /root/.fabric-ca-client/fabric-ca-client-config.yaml
2018/12/03 16:02:03 [INFO] generating key: &{A:ecdsa S:256}
2018/12/03 16:02:03 [INFO] encoded CSR
2018/12/03 16:02:03 [INFO] Stored client certificate at /root/.fabric-ca-client/msp/signcerts/cert.pem
2018/12/03 16:02:03 [INFO] Stored root CA certificate at /root/.fabric-ca-client/msp/cacerts/ca-7054.pem

root@af4edc280f83:~# fabric-ca-client register --id.name user1 --id.secret pwd1 --id.type user --id.attrs 'importer=true:ecert' -u http://ca:7054
2018/12/03 16:09:01 [INFO] Configuration file location: /root/.fabric-ca-client/fabric-ca-client-config.yaml
Password: pwd1
```