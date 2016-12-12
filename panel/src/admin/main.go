package admin

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/media_library"

	"../../../server/src/models"
)



func MakeAdminPanel() {
	db := getConnectionDB()
	db.AutoMigrate(
		&models.Product{}, &models.Category{}, &models.Collection{})

	// Initalize
	myAdmin := admin.New(&qor.Config{DB: db})

	// Create resources from GORM-backend model
	categories := myAdmin.AddResource(&models.Category{})

	// PRODUCT
	product := myAdmin.AddResource(&models.Product{})

	product.Filter(&admin.Filter{
		Name:   "Categories",
		Config: &admin.SelectOneConfig{RemoteDataResource: categories},
	})
	product.UseTheme("grid")


	myAdmin.AddResource(&models.Collection{})
	assetManager := myAdmin.AddResource(&media_library.AssetManager{}, &admin.Config{Invisible: true})
	product.Meta(&admin.Meta{Name: "Description", Config: &admin.RichEditorConfig{AssetManager: assetManager, Plugins: []admin.RedactorPlugin{
		{Name: "medialibrary", Source: "/admin/assets/javascripts/qor_redactor_medialibrary.js"},
		{Name: "table", Source: "/javascripts/redactor_table.js"},
	},
		Settings: map[string]interface{}{
			"medialibraryUrl": "/admin/product_images",
		},
	}})
	// Register route
	mux := http.NewServeMux()
	// amount to /admin, so visit `/admin` to view the admin interface
	myAdmin.MountTo("/admin", mux)

	myAdmin.SetSiteName("SHOP admin panel")

	log.Println("Admin panel listening on: 9000")
	http.ListenAndServe(":9000", mux)
}

func getConnectionDB() *gorm.DB {
	db, err := gorm.Open(
		"mysql", "iris:iris@/iris_db?charset=utf8&parseTime=True&loc=Local")
	checkErr(err)
	return db
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
