package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	EXPORTER   = "LumberInc"
	EXPBANK    = "LumberBank"
	EXPBALANCE = 100000
	IMPORTER   = "WoodenToys"
	IMPBANK    = "ToyBank"
	IMPBALANCE = 200000
	CARRIER    = "UniversalFrieght"
	REGAUTH    = "ForestryDepartment"
)

func getInitArguments() [][]byte {
	return [][]byte{[]byte("init"),
		[]byte("LumberInc"),
		[]byte("LumberBank"),
		[]byte("100000"),
		[]byte("WoodenToys"),
		[]byte("ToyBank"),
		[]byte("200000"),
		[]byte("UniversalFrieght"),
		[]byte("ForestryDepartment")}
}

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	fmt.Println("Initiating unit testing...")
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	} else {

	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was", string(bytes), "and not", value, "as expected")
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, function string, name string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte(function), []byte(name)})
	if res.Status != shim.OK {
		fmt.Println("Query", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", name, "failed to get value")
		t.FailNow()
	}
	payload := string(res.Payload)
	if payload != value {
		fmt.Println("Query value", name, "was", payload, "and not", value, "as expected")
		t.FailNow()
	}
}

func checkBadInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status == shim.OK {
		fmt.Println("Invoke", args, "unexpectedly succeeded")
		t.FailNow()
	}
}

func TestTradeWorkflow_Init(t *testing.T) {
	scc := new(TradeWorkflowChaincode)
	scc.testMode = true
	stub := shim.NewMockStub("Trade Workflow", scc)

	// Init
	checkInit(t, stub, getInitArguments())

	checkState(t, stub, "Exporter", EXPORTER)
	checkState(t, stub, "ExportersBank", EXPBANK)
}

func TestTradeWorkflow_Agreement(t *testing.T) {
	fmt.Println("**********Initiating TradeWorkflow Agreement Test Cases**********")
	scc := new(TradeWorkflowChaincode)
	scc.testMode = true
	stub := shim.NewMockStub("Trade Workflow", scc)

	checkInit(t, stub, getInitArguments())
	// Invoking requestTrade()
	tradeID := "2ks89j9"
	amount := 50000
	descGoods := "Wood for Toys"
	checkInvoke(t, stub, [][]byte{[]byte("requestTrade"), []byte(tradeID), []byte(strconv.Itoa(amount)), []byte(descGoods)})

	tradeAgreement := &TradeAgreement{amount, descGoods, REQUESTED, 0} // Check file asset.go
	tradeAgreementBytes, _ := json.Marshal(tradeAgreement)
	tradeKey, _ := stub.CreateCompositeKey("Trade", []string{tradeID})
	checkState(t, stub, tradeKey, string(tradeAgreementBytes))

	expectedResp := "{\"Status\":\"REQUESTED\"}"
	checkQuery(t, stub, "getTradeStatus", tradeID, expectedResp)

	// Invoke bad 'acceptTrade' and verify unchanged state
	checkBadInvoke(t, stub, [][]byte{[]byte("acceptTrade")})
	badTradeID := "abcd"
	checkBadInvoke(t, stub, [][]byte{[]byte("acceptTrade"), []byte(badTradeID)})
	checkState(t, stub, tradeKey, string(tradeAgreementBytes))
	checkQuery(t, stub, "getTradeStatus", tradeID, expectedResp)
}
