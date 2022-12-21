package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"murakali/config"
	"murakali/internal/module/location"
	"murakali/internal/module/location/delivery/body"
	"murakali/pkg/postgre"
	"net/http"
)

type locationUC struct {
	cfg          *config.Config
	txRepo       *postgre.TxRepo
	locationRepo location.Repository
}

func NewLocationUseCase(cfg *config.Config, txRepo *postgre.TxRepo, locationRepo location.Repository) location.UseCase {
	return &locationUC{cfg: cfg, txRepo: txRepo, locationRepo: locationRepo}
}

func (u *locationUC) GetProvince(ctx context.Context) (*body.ProvinceResponse, error) {
	province := body.ProvinceResponse{}
	provinceRedis, err := u.locationRepo.GetProvinceRedis(ctx)
	if err != nil {
		res, err := u.GetProvinceRajaOngkir()
		if err != nil {
			return nil, err
		}

		province.Rows = res.RajaOngkir.Results
		redisValue, err := json.Marshal(province)
		if err != nil {
			return nil, err
		}

		if err := u.locationRepo.InsertProvinceRedis(ctx, string(redisValue)); err != nil {
			return nil, err
		}

		return &province, nil
	}

	if err := json.Unmarshal([]byte(provinceRedis), &province); err != nil {
		return nil, err
	}

	return &province, nil
}

func (u *locationUC) GetProvinceRajaOngkir() (*body.RajaOngkirProvinceResponse, error) {
	var responseOngkir body.RajaOngkirProvinceResponse
	url := fmt.Sprintf("%s/province", u.cfg.External.OngkirAPIURL)
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("key", u.cfg.External.OngkirAPIKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	readErr := json.NewDecoder(res.Body).Decode(&responseOngkir)
	if readErr != nil {
		return nil, err
	}

	return &responseOngkir, nil
}
