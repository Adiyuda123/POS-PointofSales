package main

import (
	"POS-PointofSales/app/config"
	"POS-PointofSales/app/database"
	"POS-PointofSales/app/routes"

	authHandler "POS-PointofSales/features/auth/handler"
	authRepo "POS-PointofSales/features/auth/repository"
	authLogic "POS-PointofSales/features/auth/usecase"

	uHandler "POS-PointofSales/features/users/handler"
	uRepo "POS-PointofSales/features/users/repository"
	uLogic "POS-PointofSales/features/users/usecase"

	productsHandler "POS-PointofSales/features/products/handler"
	productsRepo "POS-PointofSales/features/products/repository"
	productsLogic "POS-PointofSales/features/products/usecase"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	cfg := config.InitConfig()
	db := database.InitDBMySql(*cfg)
	database.Migrate(db)

	aMdl := authRepo.New(db)
	aSrv := authLogic.New(aMdl)
	aCtl := authHandler.New(aSrv)

	uMdl := uRepo.New(db)
	uSrv := uLogic.New(uMdl)
	uCtl := uHandler.New(uSrv)

	productsMdl := productsRepo.New(db)
	productsSrv := productsLogic.New(productsMdl)
	productsCtl := productsHandler.New(productsSrv)

	routes.AuthRoutes(e, aCtl)
	routes.UserRoutes(e, uCtl)
	routes.ProductRoutes(e, productsCtl)

	if err := e.Start(":8080"); err != nil {
		e.Logger.Fatal("cannot start server", err.Error())
	}
}
