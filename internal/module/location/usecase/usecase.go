package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"murakali/config"
	"murakali/internal/model"
	"murakali/internal/module/location"
	"murakali/internal/module/location/delivery/body"
	"murakali/pkg/postgre"
	"net/http"
	"strings"
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

func (u *locationUC) GetCity(ctx context.Context, provinceID int) (*body.CityResponse, error) {
	city := body.CityResponse{Rows: make([]model.City, 0)}
	cityRedis, err := u.locationRepo.GetCityRedis(ctx, provinceID)
	if err != nil {
		res, err := u.GetCityRajaOngkir(provinceID)
		if err != nil {
			return nil, err
		}

		for _, result := range res.Rajaongkir.Results {
			city.Rows = append(city.Rows, model.City{
				CityID: result.CityID,
				City:   fmt.Sprintf("%s %s", result.Type, result.CityName),
			})
		}

		redisValue, err := json.Marshal(city)
		if err != nil {
			return nil, err
		}

		if err := u.locationRepo.InsertCityRedis(ctx, provinceID, string(redisValue)); err != nil {
			return nil, err
		}

		return &city, nil
	}

	if err := json.Unmarshal([]byte(cityRedis), &city); err != nil {
		return nil, err
	}

	return &city, nil
}

func (u *locationUC) GetSubDistrict(ctx context.Context, province, city string) (*body.SubDistrictResponse, error) {
	subDistrict := body.SubDistrictResponse{Rows: make([]model.SubDistrict, 0)}
	subDistrictRedis, err := u.locationRepo.GetSubDistrictRedis(ctx, province, city)
	if err != nil {
		res, err := u.GetDataFromKodePos(province, city, "")
		if err != nil {
			return nil, err
		}

		subDistrictMap := make(map[string]model.SubDistrict)
		for _, value := range res.Data {
			_, ok := subDistrictMap[value.Subdistrict]
			if strings.Contains(value.Province, province) && strings.Contains(value.City, city) && !ok {
				subDistrictMap[value.Subdistrict] = model.SubDistrict{SubDistrict: value.Subdistrict}
				subDistrict.Rows = append(subDistrict.Rows, subDistrictMap[value.Subdistrict])
			}
		}

		redisValue, err := json.Marshal(subDistrict)
		if err != nil {
			return nil, err
		}

		if err := u.locationRepo.InsertSubDistrictRedis(ctx, province, city, string(redisValue)); err != nil {
			return nil, err
		}

		return &subDistrict, nil
	}

	if err := json.Unmarshal([]byte(subDistrictRedis), &subDistrict); err != nil {
		return nil, err
	}

	return &subDistrict, nil
}

func (u *locationUC) GetUrban(ctx context.Context, province, city, subDistrict string) (*body.UrbanResponse, error) {
	urban := body.UrbanResponse{Rows: make([]model.Urban, 0)}
	urbanRedis, err := u.locationRepo.GetUrbanRedis(ctx, province, city, subDistrict)
	if err != nil {
		res, err := u.GetDataFromKodePos(province, city, subDistrict)
		if err != nil {
			return nil, err
		}

		urbanMap := make(map[string]model.Urban)
		for _, value := range res.Data {
			_, ok := urbanMap[value.Urban]
			if strings.Contains(value.Province, province) && strings.Contains(value.City, city) && strings.Contains(value.Subdistrict, subDistrict) && !ok {
				urbanMap[value.Urban] = model.Urban{Urban: value.Urban, PostalCode: value.Postalcode}
				urban.Rows = append(urban.Rows, urbanMap[value.Urban])
			}
		}

		redisValue, err := json.Marshal(urban)
		if err != nil {
			return nil, err
		}

		if errInsert := u.locationRepo.InsertUrbanRedis(ctx, province, city, subDistrict, string(redisValue)); errInsert != nil {
			return nil, errInsert
		}

		return &urban, err
	}

	if err := json.Unmarshal([]byte(urbanRedis), &urban); err != nil {
		return nil, err
	}

	return &urban, nil
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

func (u *locationUC) GetCityRajaOngkir(provinceID int) (*body.RajaOngkirCityResponse, error) {
	var responseOngkir body.RajaOngkirCityResponse
	url := fmt.Sprintf("%s/city?province=%d", u.cfg.External.OngkirAPIURL, provinceID)

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

func (u *locationUC) GetDataFromKodePos(province, city, subdistrict string) (*body.KodePosResponse, error) {
	var responseKodePos body.KodePosResponse
	url := fmt.Sprintf("%s/search/?q=%s %s %s", u.cfg.External.KodePosURL, province, city, subdistrict)
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	readErr := json.NewDecoder(res.Body).Decode(&responseKodePos)
	if readErr != nil {
		return nil, err
	}

	return &responseKodePos, nil
}
