package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"

	"./model"
)

// params [name, tel]
// func (sc *SupplyChain) createAccount(stub shim.ChaincodeStubInterface, params []string) peer.Response {
// 	if len(params) != 2 {
// 		return shim.Error("Params's size is not two.")
// 	}

// 	creator, _ := stub.GetCreator() // 获取证书的byte数据
// 	uname, err := getCN(creator)    // 根据byte获取证书CN
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}

// 	return shim.Success(nil)
// }

// {SupplierName SupplierTel CoreName CoreTel Money}
func (sc *SupplyChain) insertInvoice(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	if len(params) != 1 {
		return shim.Error("Params's size is not one.")
	}

	invoice := model.Invoice{}
	err := json.Unmarshal([]byte(params[0]), &invoice)
	if err != nil {
		return shim.Error("Unmarshal failed for params[0]. Err: " + err.Error())
	}

	creator, _ := stub.GetCreator()
	uname, _ := getCN(creator)
	if uname == "" {
		return shim.Error("Get CN in cert failed.")
	}
	invoice.UserName = uname

	invoice.Checked = "ready"
	invoice.Into = "ready"

	number, err := stub.GetState(uname)
	if err != nil {
		err = stub.PutState(uname, []byte("0"))
		if err != nil {
			return shim.Error("PutState for " + uname + " failed. Err: " + err.Error())
		}
		number = []byte("0")
	}

	// 写入状态数据
	data, _ := json.Marshal(invoice)
	err = stub.PutState(uname+string(number), data)
	if err != nil {
		return shim.Error("PutState failed. Err: " + err.Error())
	}
	num, err := strconv.Atoi(string(number))
	if err != nil {
		return shim.Error("Atoi failed. Err: " + err.Error())
	}
	err = stub.PutState(uname, []byte(strconv.Itoa(num+1)))
	if err != nil {
		return shim.Error("PutState for " + uname + " failed. Err: " + err.Error())
	}
	return shim.Success(nil)
}

// Checked Into Date(start) Date(end) ident
// checked into date(start) date(end) ident(s 供应商,c 核心企业,b 银行)
func (sc *SupplyChain) queryInvoiceBySomething(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	if len(params) != 5 {
		return shim.Error("Params's size is not five.")
	}
	if params[4] == "" {
		return shim.Error("No ident.")
	}

	// 获取用户组织
	creator, _ := stub.GetCreator()
	uname, _ := getCN(creator)
	if uname == "" {
		return shim.Error("Get CN in cert failed.")
	}
	orgName, err := getOrgNameFromCN(uname)
	if err != nil {
		return shim.Error("Get orgName failed. Err: " + err.Error())
	}

	// 检查参数并构造查询语句
	queryStringObj := &model.RichQuery{}
	queryStringObj.Selector = &model.InvoiceForQuery{}
	for index, value := range params {
		if value != "" {
			switch index {
			case 0:
				queryStringObj.Selector.Checked = value
			case 1:
				queryStringObj.Selector.Into = value
			case 2:
				queryStringObj.Selector.Date.Gte = value
			case 3:
				queryStringObj.Selector.Date.Lte = value
			case 4:
				if value == "s" {
					queryStringObj.Selector.SupplierName = orgName
				} else if value == "c" {
					queryStringObj.Selector.CoreName = orgName
				} else if value == "b" {
					queryStringObj.Selector.BankName = orgName
				}
			}
		}
	}
	queryStringByte, err := json.Marshal(queryStringObj)
	if err != nil {
		return shim.Error("Structed querystring failed. Err: " + err.Error())
	}

	// 查询
	stateIterator, err := stub.GetQueryResult(string(queryStringByte))
	if err != nil {
		return shim.Error("Query failed. Err: " + err.Error())
	}
	// 根据查询结果构造jso array
	var buffer bytes.Buffer
	buffer.WriteString("[")
	isFirst := true
	for stateIterator.HasNext() {
		queryResult, _ := stateIterator.Next()
		if !isFirst {
			buffer.WriteString(",")
		}
		buffer.Write([]byte(queryResult.String()))
		isFirst = false
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

// key Checked
// key checked
func (sc *SupplyChain) updateChecked(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	if len(params) != 2 {
		return shim.Error("Params's size is not two.")
	}
	if params[1] != "yes" && params[1] != "no" {
		return shim.Error("Params[1] is invalid.")
	}

	// 根据key查询并反序列化
	data, err := stub.GetState(params[0])
	if err != nil {
		return shim.Error("No data for " + params[0] + ". Err: " + err.Error())
	}
	invoice := model.Invoice{}
	err = json.Unmarshal(data, &invoice)
	if err != nil {
		return shim.Error("Unmarshal failed for data. Err: " + err.Error())
	}
	if invoice.Checked != "ready" {
		return shim.Error("Checked is not ready.")
	}

	// 获取用户组织
	creator, _ := stub.GetCreator()
	uname, _ := getCN(creator)
	if uname == "" {
		return shim.Error("Get CN in cert failed.")
	}
	orgName, err := getOrgNameFromCN(uname)
	if err != nil {
		return shim.Error("Get orgName failed. Err: " + err.Error())
	}

	// 修改checked并提交
	if orgName != invoice.CoreName {
		return shim.Error("Permission denied for org.")
	}
	invoice.Checked = params[1]
	newData, _ := json.Marshal(invoice)
	err = stub.PutState(params[0], newData)
	if err != nil {
		return shim.Error("PutState failed. Err: " + err.Error())
	}
	return shim.Success(nil)
}

// key Into BankName BankTel
// key into bankName bankTel
func (sc *SupplyChain) updateInto(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	if len(params) != 2 && len(params) != 4 {
		return shim.Error("Params's size is not two or four.")
	}
	if params[1] != "in" && params[1] != "yes" && params[1] != "no" {
		return shim.Error("Params[1] is invalid.")
	}

	// 根据key查询并反序列化
	data, err := stub.GetState(params[0])
	if err != nil {
		return shim.Error("No data for " + params[0] + ". Err: " + err.Error())
	}
	invoice := model.Invoice{}
	err = json.Unmarshal(data, &invoice)
	if err != nil {
		return shim.Error("Unmarshal failed for data. Err: " + err.Error())
	}

	// 获取用户组织
	creator, _ := stub.GetCreator()
	uname, _ := getCN(creator)
	if uname == "" {
		return shim.Error("Get CN in cert failed.")
	}
	orgName, err := getOrgNameFromCN(uname)
	if err != nil {
		return shim.Error("Get orgName failed. Err: " + err.Error())
	}

	// 修改checked并提交
	if invoice.Into == "ready" {
		if len(params) != 4 {
			return shim.Error("Params's size is not four.")
		}
		if params[1] != "in" {
			return shim.Error("Params[1] is not 'in'.")
		}
		if orgName != invoice.BankName {
			return shim.Error("Permission denied for org.")
		}
	} else if invoice.Into == "in" {
		if len(params) != 2 {
			return shim.Error("Params's size is not two.")
		}
		if params[1] != "yes" && params[1] != "no" {
			return shim.Error("Params[1] is not 'yes' or 'no'.")
		}
		// if orgName != invoice.SupplierName {
		// 	return shim.Error("Org is invalid.")
		// }
		if uname != invoice.UserName {
			return shim.Error("Permission denied for user.")
		}
	} else {
		return shim.Error("Into is not 'ready' or 'in'")
	}
	invoice.Into = params[1]
	newData, _ := json.Marshal(invoice)
	err = stub.PutState(params[0], newData)
	if err != nil {
		return shim.Error("PutState failed. Err: " + err.Error())
	}
	return shim.Success(nil)
}

// ident
// ident(s 供应商,c 核心企业,b 银行)
func (sc *SupplyChain) summary(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	if len(params) != 1 {
		return shim.Error("Params's size is not one.")
	}

	// 获取用户组织
	creator, _ := stub.GetCreator()
	uname, _ := getCN(creator)
	if uname == "" {
		return shim.Error("Get CN in cert failed.")
	}
	orgName, err := getOrgNameFromCN(uname)
	if err != nil {
		return shim.Error("Get orgName failed. Err: " + err.Error())
	}

	// 检查参数并构造查询语句
	queryStringObj := &model.RichQuery{}
	queryStringObj.Selector = &model.InvoiceForQuery{}
	switch params[0] {
	case "s":
		queryStringObj.Selector.SupplierName = orgName
	case "c":
		queryStringObj.Selector.CoreName = orgName
	case "b":
		queryStringObj.Selector.BankName = orgName
	}

	queryStringByte, err := json.Marshal(queryStringObj)
	if err != nil {
		return shim.Error("Structed querystring failed. Err: " + err.Error())
	}

	// 查询
	stateIterator, err := stub.GetQueryResult(string(queryStringByte))
	if err != nil {
		return shim.Error("Query failed. Err: " + err.Error())
	}

	// 循环相加
	var money, trueMoney float64
	money = 0
	trueMoney = 0
	for stateIterator.HasNext() {
		queryResult, _ := stateIterator.Next()
		invoice := model.Invoice{}
		err = json.Unmarshal(queryResult.GetValue(), &invoice)
		if err != nil {
			return shim.Error("Unmarshal failed for queryResult.Value. Err: " + err.Error())
		}
		money += invoice.Money
		trueMoney += invoice.TrueMoney
	}
	moneyStr := strconv.FormatFloat(money, 'f', -1, 64)
	trueMoneyStr := strconv.FormatFloat(trueMoney, 'f', -1, 64)
	res := "{\"money\":" + moneyStr + ",\"trueMoney\":" + trueMoneyStr + "}"
	return shim.Success([]byte(res))
}

/*
 * 从user的证书CN字段中截取组织名
 */
func getOrgNameFromCN(cn string) (string, error) {

	if cn == "" {
		return "", fmt.Errorf("Arg is nil.")
	}
	start := strings.Index(cn, "@")
	end := strings.Index(cn, ".")
	if start == -1 || end == -1 {
		return "", fmt.Errorf("Arg is invalid.")
	}
	orgName := cn[start+1 : end]
	return orgName, nil
}

/*
 * 根据byte获取证书CN
 */
func getCN(certByte []byte) (string, error) {
	certStart := bytes.IndexAny(certByte, "-----BEGIN") // 从certByte中找到cert
	if certStart == -1 {
		return "", fmt.Errorf("No cert found.")
	}
	certByte = certByte[certStart:]
	block, _ := pem.Decode(certByte) // 解码
	if block == nil {
		return "", fmt.Errorf("Could not decode the PEM structure.")
	}
	cert, err := x509.ParseCertificate(block.Bytes) // 转为cert
	if err != nil {
		return "", fmt.Errorf("ParseCertificate failed: %s.", err)
	}
	return cert.Subject.CommonName, nil
}
