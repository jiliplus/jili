package binancecollector

// v1 data struct
type trade struct {
	ID           int
	Price        string
	Quantity     string
	Time         int
	IsBuyerMaker bool
	IsBestMatch  bool
	Symbol       string `gorm:"-"` // 本字段不会保存到数据库
}

func newTrade(symbol string) *trade {
	return &trade{
		Symbol: symbol,
	}
}

func (t *trade) TableName() string {
	return t.Symbol
}
