package repository

import "context"

func (schema AllwiseInputSchema) FillMastercat(dst *Mastercat, ipix int64) {
	dst.ID = *schema.Source_id
	dst.Ra = *schema.Ra
	dst.Dec = *schema.Dec
	dst.Cat = "allwise"
	dst.Ipix = ipix
}

func (q *Queries) InsertMastercat(ctx context.Context, arg Mastercat) error {
	_, err := q.db.ExecContext(ctx, insertObject,
		arg.ID,
		arg.Ipix,
		arg.Ra,
		arg.Dec,
		arg.Cat,
	)
	return err
}
