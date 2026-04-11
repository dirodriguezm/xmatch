package ztfdr

import (
	"fmt"
	"strconv"
)

type Detection struct {
	Oid      int64   `json:"oid"`
	FilterId int     `json:"filterid"`
	FieldId  int     `json:"fieldid"`
	Rcid     int     `json:"rcid"`
	Nepochs  int     `json:"nepochs"`
	ObjRa    float64 `json:"objra"`
	ObjDec   float64 `json:"objdec"`
	Hmjd     float64 `json:"hmjd"`
	Mag      float64 `json:"mag"`
	Magerr   float64 `json:"magerr"`
}

func (d Detection) GetId() string {
	return fmt.Sprintf("%d_%s", d.Oid, strconv.FormatFloat(d.Hmjd, 'f', -1, 64))
}

func (d Detection) GetObjectId() string {
	return strconv.FormatInt(d.Oid, 10)
}

func (d Detection) GetBrightness() float64 {
	return d.Mag
}

func (d Detection) GetBrightnessError() float32 {
	return float32(d.Magerr)
}

func (d Detection) GetMjd() float64 {
	return d.Hmjd
}
