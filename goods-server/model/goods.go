package model

type Category struct {
	BaseModel
	// gorm.Model
	Name             string `gorm:"type:varchar(20);not null"`
	ParentCategoryId int32
	ParentCategory   *Category
	Level            int32 `gorm:"type:int;not null;default:1"`
	IsTable          bool  `gorm:"not null;default:false"`
}

type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null"`
	Logo string `gorm:"type:varchar(20);default:'';not null"`
}

type GoodsCategoryBrand struct {
	BaseModel

	//  联合唯一索引
	CategoryId int32 `gorm:"type:int;index:index_category_brand;unique"`
	Category   Category
	BrandsId   int32 `gorm:"type:int;index:index_category_brand;unique"`
	Brands     Brands
}

// 重载默认表名
func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

// 轮播图
type Banner struct {
	BaseModel

	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null"`
	Index string `gorm:"type:varchar(20);not null"`
}

type Goods struct {
	BaseModel

	//  联合唯一索引
	CategoryId int32 `gorm:"type:int;not null"`
	Category   Category
	BrandsId   int32 `gorm:"type:int;not null"`
	Brands     Brands

	OnSale   bool `gorm:"default:false;not null"`
	ShipFree bool `gorm:"default:false;not null"`
	IsNew    bool `gorm:"default:false;not null"`
	IsHot    bool `gorm:"default:false;not null"`

	Name            string    `gorm:"type:varchar(50);not null"`
	GoodsSn         string    `gorm:"type:varchar(50);not null"`
	ClickNum        int32     `gorm:"type:int;default:0;not null"`
	SoldNum         int32     `gorm:"type:int;default:0;not null"`
	FavNum          int32     `gorm:"type:int;default:0;not null"`
	MarketPrice     float32   `gorm:"not null"`
	ShopPrice       float32   `gorm:"not null"`
	GoodsBrief      string    `gorm:"type:varchar(100);not null"`
	Images          GormSlice `gorm:"type:varchar(1000);not null"`
	DescImages      GormSlice `gorm:"type:varchar(1000);not null"`
	GoodsFrontImage string    `gorm:"type:varchar(200);not null"`
}
