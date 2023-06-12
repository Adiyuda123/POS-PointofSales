package usecase

import (
	"POS-PointofSales/features/restocks"
	"errors"
	"strings"

	"github.com/labstack/gommon/log"
)

type restockLogic struct {
	u restocks.Repository
}

func New(r restocks.Repository) restocks.UseCase {
	return &restockLogic{
		u: r,
	}
}

// AddRestock implements restocks.UseCase.
func (pl *restockLogic) AddRestock(userId uint, stockInput restocks.Core) error {
	err := pl.u.AddRestock(userId, stockInput)
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
