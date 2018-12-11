package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	/*"errors"
	"strconv"
	"string"
	"encoding/json"*/
	"github.com/hyperledger/fabric/core/chaincode/shim"
	/*"github.com/hyperledger/fabric/core/chaincode/lib/cid"*/
	pb "github.com/hyperledger/fabric/protos/peer"
)

type TradeWorkflowChaincode struct {
	testMode bool
}

func (t *TradeWorkflowChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Initializing Trade Workflow")
	_, args := stub.GetFunctionAndParameters()
	var err error
	// Upgrade Mode 1: Leaves ledger state as-is
	if len(args) == 0 {
		return shim.Sunccess(nil)
	}

	// Upgrade mode 2: change all the names and account balances
	if len(args) != 8 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 8: {"+
			"Exporter, "+
			"Exporter's Bank, "+
			"Exporter's Account Balance, "+
			"Importer, "+
			"Importer's Bank, "+
			"Importer's Account Balance, "+
			"Carrier, "+
			"Regulatory Authority"+
			"}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	// TYpe checks
	_, err = strconv.Atoi(string(args[2]))
	if err != nil {
		fmt.Printf("Exporter's account balance must be an integer.Found %s\n", args[2])
		return shim.Error(err.Error())
	}

	_, err = strconv.Atoi(string(args[5]))
	if err != nil {
		fmt.Printf("Importer's account balance must be an integer.Found %s\n", args[5])
		return shim.Error(err.Error())
	}

	roleKeys := []string{expKey, ebKey, expBalKey, impKey, ibKey, impBalKey, carKey, raKey}
	for i, roleKey := range roleKeys {
		err = stub.PutState(roleKey, []byte(args[i]))
		if err != nil {
			fmt.Errorf("Error recording key %s: %s\n", roleKey, err.Error())
			return shim.Error(err.Error())
		}
	}
	return shim.Success(nil)
}

func (t *TradeWorkflowChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("TradeWorkFlow Invoked")
	// Step#1: get function name and argunments
	function, args := stub.GetFunctionAndParameters()

	// Step#2: Check function name and route to the relevant function accordingly
	if function == "requestTrade" {
		// Importer requests a trade
		return t.requestTrade(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "acceptTrade" {
		// Exporter accepts a trade
		// return t.acceptTrade(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "requestLC" {
		// Importer requests an L/C
		// return t.requestLC(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "issueLC" {
		// Importer's Bank issues an L/C
		// return t.issueLC(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "acceptLC" {
		// Exporter's Bank accepts an L/C
		// return t.acceptLC(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "requestEL" {
		// Exporter requests an E/L
		// return t.requestEL(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "issueEL" {
		// Regulatory Authority issues an E/L
		// return t.issueEL(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "prepareShipment" {
		// Exporter prepares a shipment
		// return t.prepareShipment(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "acceptShipmentAndIssueBL" {
		// Carrier validates the shipment and issues a B/L
		// return t.acceptShipmentAndIssueBL(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "requestPayment" {
		// Exporter's Bank requests a payment
		// return t.requestPayment(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "makePayment" {
		// Importer's Bank makes a payment
		// return t.makePayment(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "updateShipmentLocation" {
		// Carrier updates the shipment location
		// return t.updateShipmentLocation(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "getTradeStatus" {
		// Get status of trade agreement
		// return t.getTradeStatus(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "getLCStatus" {
		// Get the L/C status
		// return t.getLCStatus(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "getELStatus" {
		// Get the E/L status
		// return t.getELStatus(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "getShipmentLocation" {
		// Get the shipment location
		// return t.getShipmentLocation(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "getBillOfLading" {
		// Get the bill of lading
		// return t.getBillOfLading(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "getAccountBalance" {
		// Get account balance: Exporter/Importer
		// return t.getAccountBalance(stub, creatorOrg, creatorCertIssuer, args)
		/*} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, creatorOrg, creatorCertIssuer, args)*/
	}
	return shim.Error("Invalid invoke function name")
}

// Request a trade agreement
func (t *TradeWorkflowChaincode) requestTrade(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var tradeKey string
	var tradeAgreement *TradeAgreement
	var tradeAgreementBytes []byte
	var amount int
	var err error

	// ADD TRADELIMIT RETRIEVAL HERE

	// Access control: Only an Importer Org member can invoke this transaction
	if !t.testMode && !authenticateImporterOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Importer Org. Access denied.")
	}

	if len(args) != 3 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 3: {ID, Amount, Description of Goods}. Found %d", len(args)))
		return shim.Error(err.Error())
	}

	amount, err = strconv.Atoi(string(args[1]))
	if err != nil {
		return shim.Error(err.Error())
	}

	// ADD TRADE LIMIT CHECK HERE

	tradeAgreement = &TradeAgreement{amount, args[2], REQUESTED, 0}
	tradeAgreementBytes, err = json.Marshal(tradeAgreement)
	if err != nil {
		return shim.Error("Error marshaling trade agreement structure")
	}

	// Write the state to the ledger
	tradeKey, err = getTradeKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(tradeKey, tradeAgreementBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Trade %s request recorded\n", args[0])

	return shim.Success(nil)
}

func (t *TradeWorkflowChaincode) acceptTrade(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {

}
