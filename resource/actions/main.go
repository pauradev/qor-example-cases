package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/pauradev/gorm"
	_ "github.com/pauradev/gorm/dialects/postgres"
	"github.com/pauradev/admin"
	"github.com/pauradev/qor-website-cases/config"
	"github.com/pauradev/roles"
	appkitlog "github.com/theplant/appkit/log"
	"github.com/theplant/appkit/server"
)

type Order struct {
	gorm.Model
	Num   string
	Price float64
}

// run with dummy data
// MODE=data go run main.go

func main() {
	db := config.DB
	if os.Getenv("MODE") == "data" {
		db.DropTable(&Order{})
		db.AutoMigrate(&Order{})
		db.Create(&Order{Num: "T00001", Price: 1000})
		db.Create(&Order{Num: "T00002", Price: 1500})
	} else {
		db.AutoMigrate(&Order{})
	}

	adm := config.Admin
	orderR := adm.AddResource(&Order{})
	orderR.Action(&admin.Action{
		Name: "Action with error",
		Handler: func(argument *admin.ActionArgument) error {
			argument.Context.AddError(fmt.Errorf("This is a error"))
			return fmt.Errorf("This is a error in return")
		},
		Modes:      []string{"edit", "index", "show"},
		Permission: roles.Allow(roles.CRUD, roles.Anyone),
	})

	type CollectionRes struct {
		Name string
	}
	orderR.Action(&admin.Action{
		Name: "Action runs collection",
		Handler: func(argument *admin.ActionArgument) error {
			return nil
		},
		Resource:   adm.NewResource(&CollectionRes{}),
		Modes:      []string{"edit", "collection"},
		Permission: roles.Allow(roles.CRUD, roles.Anyone),
	})

	orderR.Action(&admin.Action{
		Name: "Action runs OK",
		Handler: func(argument *admin.ActionArgument) error {
			return nil
		},
		Modes:      []string{"edit", "index", "show"},
		Permission: roles.Allow(roles.CRUD, roles.Anyone),
	})
	orderR.Action(&admin.Action{
		Name: "Action with download",
		Handler: func(argument *admin.ActionArgument) error {
			w := argument.Context.Writer
			w.Header().Set("Content-Disposition", "attachment; filename=file.csv")
			w.Write([]byte("Hello world"))
			return nil
		},
		Modes:      []string{"index", "menu_item"},
		Permission: roles.Allow(roles.CRUD, roles.Anyone),
	})

	mux := http.NewServeMux()
	adm.MountTo("/admin", mux)
	color.Green("URL: %v", "http://localhost:3000/admin/orders")
	server.ListenAndServe(server.Config{Addr: ":3000"}, appkitlog.Default(), mux)
}
