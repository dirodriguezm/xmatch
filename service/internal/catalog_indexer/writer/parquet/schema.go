package parquet_writer

type IndexerSchema struct {
	Id   string  `parquet:"name=id, type=BYTE_ARRAY"`
	Ra   float64 `parquet:"name=ra, type=DOUBLE"`
	Dec  float64 `parquet:"name=dec, type=DOUBLE"`
	Ipix int     `parquet:"name=ipix, type=INT64"`
	Cat  string  `parquet:"name=cat, type=BYTE_ARRAY"`
}
