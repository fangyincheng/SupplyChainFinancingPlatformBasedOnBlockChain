package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type SupplyChain struct {
}

// Init函数在链码实例化中初始化一些数据，同时，升级链码的时候也会调用这个函数
func (sc *SupplyChain) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke函数在链码调用的时候被调用到，会根据传入的funcName判断执行哪些链码逻辑
func (sc *SupplyChain) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	funcName, params := stub.GetFunctionAndParameters()
	switch funcName {

	case "insertInvoice": // 插入一条发票
		return sc.insertInvoice(stub, params)

	case "queryInvoiceBySomething": // 根据条件对couchdb进行rich query
		return sc.queryInvoiceBySomething(stub, params)

	case "updateChecked": // 确认字段修改
		return sc.updateChecked(stub, params)

	case "updateInto": // 融资字段修改
		return sc.updateInto(stub, params)

	case "summary": // 金额汇总
		return sc.summary(stub, params)

	default: // 如果不识别方法名，则返回error
		return shim.Error("Invalid funcName.")
	}
}

// main函数在初始化时启动容器中的链码
func main() {
	if err := shim.Start(new(SupplyChain)); err != nil {
		fmt.Printf("Error starting chaincode: %s.", err)
	}
}
