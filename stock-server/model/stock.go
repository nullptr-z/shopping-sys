package model

import (
	"database/sql/driver"
	"encoding/json"
)

type Stock struct {
	BaseModel

	// 商品编号
	Goods int32 `gorm:"type:int;index"`
	// 库存数量
	Stocks  int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"` // 分布式锁-乐观锁
}

type GoodsDetail struct {
	Goods int32
	Num   int32
}

type GoodsDetailList []GoodsDetail

func (g GoodsDetailList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

func (g *GoodsDetailList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

type StockSellReback struct {
	OrderSn   string          `gorm:"type:varchar(200);index:idx_order_sn,unique"` // 对应的订单号,索引保证,库存回退幂等性
	Status    int32           `gorm:"type:int;"`                                   // 1 库存预先扣减，2 已归还
	GoodsList GoodsDetailList `gorm:"type:varchar(200);"`
	// User        int32           // 哪个用户
	// Goods       int32           // 哪个商品
	// Nums        int32           // 占了多少库存
}

func (s StockSellReback) TableName() string {
	return "stocksellreback"
}
