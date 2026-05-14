package allwise

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

const displayName = "AllWISE"

type InputSchema struct {
	Source_id    *string  `parquet:"name=source_id, type=BYTE_ARRAY"`
	Cntr         *int64   `parquet:"name=cntr, type=INT64"`
	Ra           *float64 `parquet:"name=ra, type=DOUBLE"`
	Dec          *float64 `parquet:"name=dec, type=DOUBLE"`
	W1mpro       *float64 `parquet:"name=w1mpro, type=DOUBLE"`
	W1sigmpro    *float64 `parquet:"name=w1sigmpro, type=DOUBLE"`
	W2mpro       *float64 `parquet:"name=w2mpro, type=DOUBLE"`
	W2sigmpro    *float64 `parquet:"name=w2sigmpro, type=DOUBLE"`
	W3mpro       *float64 `parquet:"name=w3mpro, type=DOUBLE"`
	W3sigmpro    *float64 `parquet:"name=w3sigmpro, type=DOUBLE"`
	W4mpro       *float64 `parquet:"name=w4mpro, type=DOUBLE"`
	W4sigmpro    *float64 `parquet:"name=w4sigmpro, type=DOUBLE"`
	J_m_2mass    *float64 `parquet:"name=j_m_2mass, type=DOUBLE"`
	H_m_2mass    *float64 `parquet:"name=h_m_2mass, type=DOUBLE"`
	K_m_2mass    *float64 `parquet:"name=k_m_2mass, type=DOUBLE"`
	J_msig_2mass *float64 `parquet:"name=j_msig_2mass, type=DOUBLE"`
	H_msig_2mass *float64 `parquet:"name=h_msig_2mass, type=DOUBLE"`
	K_msig_2mass *float64 `parquet:"name=k_msig_2mass, type=DOUBLE"`
}

type Adapter struct {
	repo *repository.Queries
}

func init() {
	catalog.Register("allwise", func(repo *repository.Queries) (catalog.CatalogAdapter, error) {
		return &Adapter{repo: repo}, nil
	})
}

func (a Adapter) Name() string {
	return "allwise"
}

func (a Adapter) NewRawRecord() any {
	return InputSchema{}
}

func (a Adapter) NewMetadataRecord() any {
	return repository.Allwise{}
}

func (a Adapter) BulkInsertMetadata(ctx context.Context, rows []any) error {
	if a.repo == nil {
		return fmt.Errorf("allwise adapter has no repository")
	}
	return a.repo.BulkInsertAllwise(ctx, rows)
}

func (a Adapter) GetByID(ctx context.Context, id string) (any, error) {
	if a.repo == nil {
		return nil, fmt.Errorf("allwise adapter has no repository")
	}
	return a.repo.GetAllwise(ctx, id)
}

func (a Adapter) BulkGetByID(ctx context.Context, ids []string) (any, error) {
	if a.repo == nil {
		return nil, fmt.Errorf("allwise adapter has no repository")
	}
	return a.repo.BulkGetAllwise(ctx, ids)
}

func (a Adapter) GetFromPixels(ctx context.Context, pixels []int64) ([]repository.Metadata, error) {
	if a.repo == nil {
		return nil, fmt.Errorf("allwise adapter has no repository")
	}
	rows, err := a.repo.GetAllwiseFromPixels(ctx, pixels)
	if err != nil {
		return nil, err
	}
	result := make([]repository.Metadata, len(rows))
	for i, r := range rows {
		result[i] = convertAllwiseFromPixelsRowToMetadata(r)
	}
	return result, nil
}

func (a Adapter) GetCoordinates(raw any) (float64, float64, error) {
	schema, ok := raw.(InputSchema)
	if !ok {
		return 0, 0, fmt.Errorf("expected allwise.InputSchema, got %T", raw)
	}
	ra := 0.0
	dec := 0.0
	if schema.Ra != nil {
		ra = *schema.Ra
	}
	if schema.Dec != nil {
		dec = *schema.Dec
	}
	return ra, dec, nil
}

func (a Adapter) ConvertToMastercat(raw any, ipix int64) (repository.Mastercat, error) {
	schema, ok := raw.(InputSchema)
	if !ok {
		return repository.Mastercat{}, fmt.Errorf("expected allwise.InputSchema, got %T", raw)
	}
	mc := repository.Mastercat{
		Ipix: ipix,
		Cat:  "allwise",
	}
	if schema.Source_id != nil {
		mc.ID = *schema.Source_id
	}
	if schema.Ra != nil {
		mc.Ra = *schema.Ra
	}
	if schema.Dec != nil {
		mc.Dec = *schema.Dec
	}
	return mc, nil
}

func (a Adapter) ConvertToMetadataFromRaw(raw any) (any, error) {
	schema := raw.(InputSchema)
	return convertAllwiseInputToMetadata(schema), nil
}

func convertAllwiseFromPixelsRowToMetadata(row repository.GetAllwiseFromPixelsRow) repository.Metadata {
	return repository.Metadata{
		ID:      row.ID,
		Catalog: displayName,
		Ra:      row.Ra,
		Dec:     row.Dec,
		Object: repository.Allwise{
			ID:         row.ID,
			Cntr:       row.Cntr,
			W1mpro:     row.W1mpro,
			W1sigmpro:  row.W1sigmpro,
			W2mpro:     row.W2mpro,
			W2sigmpro:  row.W2sigmpro,
			W3mpro:     row.W3mpro,
			W3sigmpro:  row.W3sigmpro,
			W4mpro:     row.W4mpro,
			W4sigmpro:  row.W4sigmpro,
			JM2mass:    row.JM2mass,
			JMsig2mass: row.JMsig2mass,
			HM2mass:    row.HM2mass,
			HMsig2mass: row.HMsig2mass,
			KM2mass:    row.KM2mass,
			KMsig2mass: row.KMsig2mass,
		},
	}
}

func convertAllwiseInputToMetadata(schema InputSchema) repository.Allwise {
	allwise := repository.Allwise{}
	if schema.Source_id != nil {
		allwise.ID = *schema.Source_id
	}
	if schema.Cntr != nil {
		allwise.Cntr = *schema.Cntr
	}
	if schema.W1mpro != nil {
		allwise.W1mpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.W1mpro, Valid: true}}
	} else {
		allwise.W1mpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.W1sigmpro != nil {
		allwise.W1sigmpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.W1sigmpro, Valid: true}}
	} else {
		allwise.W1sigmpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.W2mpro != nil {
		allwise.W2mpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.W2mpro, Valid: true}}
	} else {
		allwise.W2mpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.W2sigmpro != nil {
		allwise.W2sigmpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.W2sigmpro, Valid: true}}
	} else {
		allwise.W2sigmpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.W3mpro != nil {
		allwise.W3mpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.W3mpro, Valid: true}}
	} else {
		allwise.W3mpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.W3sigmpro != nil {
		allwise.W3sigmpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.W3sigmpro, Valid: true}}
	} else {
		allwise.W3sigmpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.W4mpro != nil {
		allwise.W4mpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.W4mpro, Valid: true}}
	} else {
		allwise.W4mpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.W4sigmpro != nil {
		allwise.W4sigmpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.W4sigmpro, Valid: true}}
	} else {
		allwise.W4sigmpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.J_m_2mass != nil {
		allwise.JM2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.J_m_2mass, Valid: true}}
	} else {
		allwise.JM2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.J_msig_2mass != nil {
		allwise.JMsig2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.J_msig_2mass, Valid: true}}
	} else {
		allwise.JMsig2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.H_m_2mass != nil {
		allwise.HM2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.H_m_2mass, Valid: true}}
	} else {
		allwise.HM2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.H_msig_2mass != nil {
		allwise.HMsig2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.H_msig_2mass, Valid: true}}
	} else {
		allwise.HMsig2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.K_m_2mass != nil {
		allwise.KM2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.K_m_2mass, Valid: true}}
	} else {
		allwise.KM2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.K_msig_2mass != nil {
		allwise.KMsig2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.K_msig_2mass, Valid: true}}
	} else {
		allwise.KMsig2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	return allwise
}
