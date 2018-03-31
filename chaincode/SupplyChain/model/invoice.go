package model

/*
 * 数据类型，存储的发票格式
 */
type Invoice struct {
	// InvoiceId    string  `json:"invoice_id"`    //发票id，一个唯一标识
	UserName     string  `json:"user_name"`     //组织用户名
	SupplierName string  `json:"supplier_name"` //供应商名称
	SupplierTel  string  `json:"supplier_tel"`  //供应商联系方式
	CoreName     string  `json:"core_name"`     //核心企业名称
	CoreTel      string  `json:"core_tel"`      //核心企业联系方式
	BankName     string  `json:"bank_name"`     //银行保理商名称
	BankTel      string  `json:"bank_tel"`      //银行保理商联系方式
	Money        float64 `json:"money"`         //融资金额
	TrueMoney    float64 `json:"true_money"`    //放款金额
	Checked      string  `json:"checked"`       //核心企业检查(ready、yes、no)
	Into         string  `json:"into"`          //银行或保理商融资(ready、in、yes、no)
	Term         string  `json:"term"`          //融资还款期限
	Date         string  `json:"date"`          //融资放款日期
}
