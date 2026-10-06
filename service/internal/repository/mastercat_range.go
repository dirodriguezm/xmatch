package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/dirodriguezm/healpix"
)

// FindObjectsInPixelRanges returns the mastercat objects whose ipix falls in
// any of the given ranges, without materializing a pixel list. Each usable
// range contributes two bind parameters, so it shares the SQLite
// bind-variable ceiling documented on QueryMetadataFromPixelRanges: at the
// current resolver order 18, radii up to about 1 degree stay under the limit
// and radii beyond roughly 1.05 degrees fail with "too many SQL variables".
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
	defer func() { _ = rows.Close() }()

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
