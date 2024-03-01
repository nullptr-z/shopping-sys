package model

type Stock struct {
	BaseModel

	Goods   int32 `gorm:"type:int;index"`
	Stocks  int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"` // 分布式锁-乐观锁
}

// type StockHistory struct {
// 	user   int32 // 哪个用户
// 	goods  int32 // 哪个商品
// 	nums   int32 // 占了多少库存
// 	order  int32 // 对应的订单号
// 	status int32 // 1 库存预先扣减，幂等。2 已支付
// }
