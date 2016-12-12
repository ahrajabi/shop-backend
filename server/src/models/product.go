package models

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/media_library"

	"../constants"
	"fmt"
)

type Product struct {
	gorm.Model
  	Image       media_library.MediaLibraryStorage	`sql:"size:4294967295;" media_library:"url:/system/{{class}}/{{primary_key}}/{{column}}.{{extension}};path:./private"`
	Name        string 	`sql:"size:50"`
	Description string
	ImgUrl      string
	Price       int
	CatID       uint
	Likes       uint
	Categories  []Category `gorm:"many2many:product_category;"`
}

func CreateProductTable() {
	db := getConnectionDB()
	defer db.Close()

	if !db.HasTable(&Product{}) {
		db.CreateTable(&Product{})
	}
}

func SaveProduct(name, description string) {
	db := getConnectionDB()
	defer db.Close()
	db.Save(
		&Product{
			Name:        name,
			Description: description,
		})
}

func GetProductByID(id uint) (*Product, error) {
	result := new(Product)

	db := getConnectionDB()
	defer db.Close()
	db.First(result, id)

	if result.Name != "" {
		return result, nil
	}
	return result, constants.PRODUCT_ID_ERR
}

func GetProductsByCatID(id uint) ([]Product, error) {
	var productsArr []Product
	cat, err := GetCategoryByID(id)

	if err != nil {
		return productsArr, err
	}

	db := getConnectionDB()
	defer db.Close()
	//cat := new(Category)
	//db.Where("categories = ?", cat).where().Find(&productsArr)
	//product := new(Product)
	//db.Model(&productsArr).
	//	Association("Categories").Find(&productsArr)
	//db.Find(&productsArr, "product_category = ?", cat.ID)
	//db.Model(&cat).Related(&productsArr, "Categories")
	var result []uint
	db.Raw("Select * from products where deleted_at IS NULL AND ID IN (Select product_id from product_category where category_id = ?)", cat.ID).Scan(&productsArr)
	fmt.Println(result)

	return productsArr, nil
}

func GetProductsByPage(pageNum uint) ([]Product, error) {
	if pageNum > GetProductsPageCount() {
		return nil, constants.END_PAGE_ERR
	}

	var productsArr []Product
	offset := (pageNum - 1) * 10

	db := getConnectionDB()
	defer db.Close()
	db.Limit(10).
		Offset(offset).Find(&productsArr)

	return productsArr, nil
}

func GetProductsByName(name string) ([]Product, error) {
	var productsArr []Product

	db := getConnectionDB()
	defer db.Close()
	db.Where("name = ?", name).Find(&productsArr)

	if len(productsArr) > 0 {
		return productsArr, nil
	} else {
		return productsArr, constants.PRODUCT_NAME_ERR
	}
}

func GetProductsPageCount() uint {
	var productsCount uint

	db := getConnectionDB()
	defer db.Close()
	db.Model(&Product{}).Count(&productsCount)

	pageCountInt :=
		(productsCount / 10)
	pageCountFloat :=
		(float32(productsCount) / 10.0)

	if float32(pageCountInt) < pageCountFloat {
		return (pageCountInt + 1)
	}
	return pageCountInt
}
