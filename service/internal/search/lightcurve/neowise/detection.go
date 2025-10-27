package neowise

import "strconv"

type Detection struct {
	Mjd       float64 `json:"mjd" parquet:"name=mjd, type=DOUBLE"`
	Ra        float64 `json:"ra" parquet:"name=ra, type=DOUBLE"`
	Dec       float64 `json:"dec" parquet:"name=dec, type=DOUBLE"`
	W1mpro    float64 `json:"w1mpro" parquet:"name=w1mpro, type=DOUBLE"`
	W1sigmpro float32 `json:"w1sigmpro" parquet:"name=w1sigmpro, type=FLOAT"`
	W2mpro    float64 `json:"w2mpro" parquet:"name=w2mpro, type=DOUBLE"`
	W2sigmpro float32 `json:"w2sigmpro" parquet:"name=w2sigmpro, type=FLOAT"`
	Cntr      int64   `json:"cntr" parquet:"name=cntr, type=INT64"`
	Source_id string  `json:"source_id" parquet:"name=source_id, type=BYTE_ARRAY"`
}

func (d Detection) GetId() string {
	return d.Source_id
}

func (d Detection) GetObjectId() string {
	return strconv.FormatInt(d.Cntr, 10)
}

func (d Detection) GetBrightness() float64 {
	return d.W1mpro
}

func (d Detection) GetBrightnessError() float32 {
	return d.W1sigmpro
}

func (d Detection) GetMjd() float64 {
	return d.Mjd
}
