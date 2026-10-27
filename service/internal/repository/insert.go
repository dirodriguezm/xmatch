// Copyright 2024-2025 Diego Rodriguez Mancini
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package repository

import "context"

// Insertable is implemented by the models a catalog pipeline writes in bulk:
// each model knows the generated query that inserts it.
type Insertable interface {
	Insert(context.Context, *Queries) error
}

func (r Mastercat) Insert(ctx context.Context, q *Queries) error {
	return q.InsertObject(ctx, InsertObjectParams(r))
}

func (r Allwise) Insert(ctx context.Context, q *Queries) error {
	return q.InsertAllwise(ctx, InsertAllwiseParams(r))
}

func (r Erosita) Insert(ctx context.Context, q *Queries) error {
	return q.InsertErosita(ctx, InsertErositaParams(r))
}

func (r Gaia) Insert(ctx context.Context, q *Queries) error {
	return q.InsertGaia(ctx, InsertGaiaParams(r))
}

func (r Swift) Insert(ctx context.Context, q *Queries) error {
	return q.InsertSwift(ctx, r.ID)
}
