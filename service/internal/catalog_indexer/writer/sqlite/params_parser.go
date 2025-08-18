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

package sqlite_writer

import "github.com/dirodriguezm/xmatch/service/internal/repository"

type ParamsParser[I, O any] interface {
	Parse(I) O
}

type MastercatParser struct{}

func (p MastercatParser) Parse(obj repository.Mastercat) repository.InsertObjectParams {
	return repository.InsertObjectParams{
		ID:   obj.ID,
		Ipix: obj.Ipix,
		Ra:   obj.Ra,
		Dec:  obj.Dec,
		Cat:  obj.Cat,
	}
}

type AllwiseParser struct{}

func (p AllwiseParser) Parse(obj repository.Metadata) repository.Metadata {
	return repository.InsertAllwiseParams{
		ID:         obj.(repository.Allwise).ID,
		W1mpro:     obj.(repository.Allwise).W1mpro,
		W1sigmpro:  obj.(repository.Allwise).W1sigmpro,
		W2mpro:     obj.(repository.Allwise).W2mpro,
		W2sigmpro:  obj.(repository.Allwise).W2sigmpro,
		W3mpro:     obj.(repository.Allwise).W3mpro,
		W3sigmpro:  obj.(repository.Allwise).W3sigmpro,
		W4mpro:     obj.(repository.Allwise).W4mpro,
		W4sigmpro:  obj.(repository.Allwise).W4sigmpro,
		JM2mass:    obj.(repository.Allwise).JM2mass,
		JMsig2mass: obj.(repository.Allwise).JMsig2mass,
		HM2mass:    obj.(repository.Allwise).HM2mass,
		HMsig2mass: obj.(repository.Allwise).HMsig2mass,
		KM2mass:    obj.(repository.Allwise).KM2mass,
		KMsig2mass: obj.(repository.Allwise).KMsig2mass,
	}
}
