package main

import (
	"POS-PointofSales/app/config"
	"POS-PointofSales/app/database"
	"POS-PointofSales/app/routes"
	authHandler "POS-PointofSales/features/auth/handler"
	authRepo "POS-PointofSales/features/auth/repository"
	authLogic "POS-PointofSales/features/auth/usecase"
	pHandler "POS-PointofSales/features/products/handler"
	pRepo "POS-PointofSales/features/products/repository"
	pLogic "POS-PointofSales/features/products/usecase"
	rHandler "POS-PointofSales/features/restocks/handler"
	rRepo "POS-PointofSales/features/restocks/repository"
	rLogic "POS-PointofSales/features/restocks/usecase"
	tHandler "POS-PointofSales/features/transactions/handler"
	tRepo "POS-PointofSales/features/transactions/repository"
	tLogic "POS-PointofSales/features/transactions/usecase"
	uHandler "POS-PointofSales/features/users/handler"
	uRepo "POS-PointofSales/features/users/repository"
	uLogic "POS-PointofSales/features/users/usecase"

	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	cfg := config.ReadEnv()
	if cfg == nil {
		log.Fatal("Failed to read configuration")
	}
	db := database.InitDBMySql(*cfg)
	if db == nil {
		log.Fatal("Failed to connect to database")
	}
	database.Migrate(db)

	aMdl := authRepo.New(db)
	aSrv := authLogic.New(aMdl)
	aCtl := authHandler.New(aSrv)

	uMdl := uRepo.New(db)
	uSrv := uLogic.New(uMdl)
	uCtl := uHandler.New(uSrv)

	pMdl := pRepo.New(db)
	pSrv := pLogic.New(pMdl)
	pCtl := pHandler.New(pSrv, uMdl)

	rMdl := rRepo.New(db)
	rSrv := rLogic.New(rMdl)
	rCtl := rHandler.New(rSrv)

	tMdl := tRepo.New(db)
	tSrv := tLogic.New(tMdl)
	tCtl := tHandler.New(tSrv, pMdl, uMdl)

	routes.AuthRoutes(e, aCtl)
	routes.UserRoutes(e, uCtl)
	routes.ProductRoutes(e, pCtl)
	routes.RestockRoutes(e, rCtl)
	routes.TransactionRoutes(e, tCtl)

	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
