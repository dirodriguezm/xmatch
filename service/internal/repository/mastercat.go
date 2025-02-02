package repository

type ParquetMastercat struct {
	ID   *string  `parquet:"name=id, type=BYTE_ARRAY"`
	Ipix *int64   `parquet:"name=ipix, type=INT64"`
	Ra   *float64 `parquet:"name=ra, type=DOUBLE"`
	Dec  *float64 `parquet:"name=dec, type=DOUBLE"`
	Cat  *string  `parquet:"name=cat, type=BYTE_ARRAY"`
}

func (m ParquetMastercat) ToInsertObjectParams() InsertObjectParams {
	return InsertObjectParams{
		ID:   *m.ID,
		Ipix: *m.Ipix,
		Ra:   *m.Ra,
		Dec:  *m.Dec,
		Cat:  *m.Cat,
	}
}

func (m ParquetMastercat) ToMastercat() Mastercat {
	return Mastercat{
		ID:   *m.ID,
		Ipix: *m.Ipix,
		Ra:   *m.Ra,
		Dec:  *m.Dec,
		Cat:  *m.Cat,
	}
}
