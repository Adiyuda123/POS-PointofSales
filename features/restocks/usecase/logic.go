package usecase

import (
	"POS-PointofSales/features/restocks"
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
)

type restockLogic struct {
	u         restocks.Repository
	validator *validator.Validate
}

func New(r restocks.Repository) restocks.UseCase {
	return &restockLogic{
		u:         r,
		validator: validator.New(),
	}
}

// GetAllRestock implements restocks.UseCase.
func (pl *restockLogic) GetAllRestock(userID uint, limit, offset int, search string, fromDate, toDate time.Time) ([]restocks.Core, int, error) {
	result, totaldata, err := pl.u.SelectAllRestock(userID, limit, offset, search, fromDate, toDate)
	if err != nil {
		log.Error("failed to find all restock", err.Error())
		return []restocks.Core{}, totaldata, errors.New("internal server error")
	}

	return result, totaldata, nil
}

// AddRestock implements restocks.UseCase.
func (pl *restockLogic) AddRestock(userId uint, stockInput restocks.Core) error {
	err := pl.validator.Struct(stockInput)
	if err != nil {
		log.Error("validation error:", err.Error())
		return err
	}

	err = pl.u.AddRestock(userId, stockInput)
	if err != nil {
		log.Error("failed on calling restock query")
		if strings.Contains(err.Error(), "finding product") {
			log.Error("error on finding product (not found)")
			return errors.New("bad request, product not found")
		} else if strings.Contains(err.Error(), "affected") {
			log.Error("no rows affected on add product")
			return errors.New("data is up to date")
		}
		return err
	}
	return nil
}
