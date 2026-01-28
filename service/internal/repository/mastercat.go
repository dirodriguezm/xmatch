package repository

import "context"

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
