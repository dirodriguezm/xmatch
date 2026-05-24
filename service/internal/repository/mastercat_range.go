package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/dirodriguezm/healpix"
)

func (q *Queries) FindObjectsInPixelRanges(ctx context.Context, ranges []healpix.PixelRange) ([]Mastercat, error) {
	if len(ranges) == 0 {
		return nil, nil
	}
	values := make([]string, 0, len(ranges))
	params := make([]interface{}, 0, len(ranges)*2)
	for _, r := range ranges {
		if r.Start >= r.Stop {
			continue
		}
		values = append(values, "(?, ?)")
		params = append(params, r.Start, r.Stop)
	}
	if len(values) == 0 {
		return nil, nil
	}

	query := fmt.Sprintf(`WITH range(start, stop) AS (VALUES %s)
SELECT DISTINCT mastercat.id, mastercat.ipix, mastercat.ra, mastercat.dec, mastercat.cat
FROM range
JOIN mastercat
  ON mastercat.ipix >= range.start
 AND mastercat.ipix < range.stop`, strings.Join(values, ","))

	rows, err := q.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Mastercat, 0)
	for rows.Next() {
		var i Mastercat
		if err := rows.Scan(
			&i.ID,
			&i.Ipix,
			&i.Ra,
			&i.Dec,
			&i.Cat,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
