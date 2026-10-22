package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/dirodriguezm/healpix"
)

// pixelRangeCTE builds the common table expression that exposes pixel ranges
// as rows with start and stop columns, plus its bind parameters. The returned
// bool is false when no usable range remains (empty input or every range
// degenerate), in which case callers should return no results without
// querying. Degenerate ranges (Start >= Stop) are skipped.
func pixelRangeCTE(ranges []healpix.PixelRange) (string, []any, bool) {
	values := make([]string, 0, len(ranges))
	params := make([]any, 0, len(ranges)*2)
	for _, r := range ranges {
		if r.Start >= r.Stop {
			continue
		}
		values = append(values, "(?, ?)")
		params = append(params, r.Start, r.Stop)
	}
	if len(values) == 0 {
		return "", nil, false
	}
	return fmt.Sprintf("WITH range(start, stop) AS (VALUES %s)", strings.Join(values, ",")), params, true
}

// QueryMetadataFromPixelRanges runs a metadata query against table filtered by
// the given pixel ranges and maps every matched row with mapRow. The query
// selects every column of table followed by mastercat.ra and mastercat.dec,
// joining through the mastercat.ipix column:
//
//	SELECT DISTINCT <table>.*, mastercat.ra, mastercat.dec
//	FROM range
//	JOIN mastercat
//	  ON mastercat.ipix >= range.start
//	 AND mastercat.ipix < range.stop
//	JOIN <table> ON <table>.id = mastercat.id
//
// The catalog adapters own the table name and the row scanning, so each
// mapRow must scan in the column order above. table must be a trusted catalog
// table name, never user input. Returns nil without querying when no usable
// range remains after skipping degenerate (Start >= Stop) ranges.
//
// Each usable range contributes two bind parameters, so the query is bounded
// by SQLite's SQLITE_MAX_VARIABLE_NUMBER (32,766 with mattn/go-sqlite3). At
// the current resolver order 18, a 1 degree radius already needs about 31k
// bind parameters, and radii beyond roughly 1.05 degrees (3,780 arcseconds)
// routinely exceed the limit and fail with "too many SQL variables". Cone
// searches are intended for arcsecond-scale radii, where this ceiling is
// never reached; supporting larger radii requires chunking the ranges into
// multiple queries.
func (q *Queries) QueryMetadataFromPixelRanges(
	ctx context.Context,
	table string,
	ranges []healpix.PixelRange,
	mapRow func(rows *sql.Rows) (Metadata, error),
) ([]Metadata, error) {
	cte, params, ok := pixelRangeCTE(ranges)
	if !ok {
		return nil, nil
	}

	query := cte + fmt.Sprintf(`
SELECT DISTINCT %[1]s.*, mastercat.ra, mastercat.dec
FROM range
JOIN mastercat
  ON mastercat.ipix >= range.start
 AND mastercat.ipix < range.stop
JOIN %[1]s ON %[1]s.id = mastercat.id`, table)

	rows, err := q.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	items := make([]Metadata, 0)
	for rows.Next() {
		item, err := mapRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
