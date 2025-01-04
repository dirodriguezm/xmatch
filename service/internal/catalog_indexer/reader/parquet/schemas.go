package parquet_reader

type AllWiseSchema struct {
	Designation *string  `parquet:"name=designation, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	Ra          *float64 `parquet:"name=ra, type=DOUBLE"`
	Dec         *float64 `parquet:"name=dec, type=DOUBLE"`
}

type VlassSchema struct{}

type ZtfSchema struct{}
