package model

/*
 * 用于couchdb rich query的模板
 */

type ContrastForFloat struct {
	Lt  float64 `json:"$lt,omitempty"`
	Lte float64 `json:"$lte,omitempty"`
	Gt  float64 `json:"$gt,omitempty"`
	Gte float64 `json:"$gte,omitempty"`
	Eq  float64 `json:"$eq,omitempty"`
	Ne  float64 `json:"$ne,omitempty"`
}

type ContrastForString struct {
	Lt  string `json:"$lt,omitempty"`
	Lte string `json:"$lte,omitempty"`
	Gt  string `json:"$gt,omitempty"`
	Gte string `json:"$gte,omitempty"`
	Eq  string `json:"$eq,omitempty"`
	Ne  string `json:"$ne,omitempty"`
}

type InvoiceForQuery struct {
	UserName     string             `json:"user_name,omitempty"`     //组织用户名
	SupplierName string             `json:"supplier_name,omitempty"` //供应商名称
	SupplierTel  string             `json:"supplier_tel,omitempty"`  //供应商联系方式
	CoreName     string             `json:"core_name,omitempty"`     //核心企业名称
	CoreTel      string             `json:"core_tel,omitempty"`      //核心企业联系方式
	BankName     string             `json:"bank_name,omitempty"`     //银行保理商名称
	BankTel      string             `json:"bank_tel,omitempty"`      //银行保理商联系方式
	Money        *ContrastForFloat  `json:"money,omitempty"`         //融资金额
	TrueMoney    *ContrastForFloat  `json:"true_money,omitempty"`    //放款金额
	Checked      string             `json:"checked,omitempty"`       //核心企业检查(ready、yes、no)
	Into         string             `json:"into,omitempty"`          //银行或保理商融资(ready、in、yes、no)
	Term         string             `json:"term,omitempty"`          //融资还款期限
	Date         *ContrastForString `json:"date,omitempty"`          //融资放款日期(0000-00-00)
}

type RichQuery struct {
	Selector *InvoiceForQuery `json:"selector,omitempty"`
	Limit    float64          `json:"limit,omitempty"`
	Skip     float64          `json:"skip,omitempty"`
}
