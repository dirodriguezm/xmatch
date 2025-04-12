package repository

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAllwiseInputSchemaToMastercat(t *testing.T) {
	source_id := "source_id"
	ra := 0.0
	dec := 0.0
	catalog := "allwise"
	ipix := int64(0)
	a := &AllwiseInputSchema{
		Source_id: &source_id,
		Ra:        &ra,
		Dec:       &dec,
	}
	require.Implements(t, (*InputSchema)(nil), a)
	expected := ParquetMastercat{
		ID:   &source_id,
		Ra:   &ra,
		Dec:  &dec,
		Cat:  &catalog,
		Ipix: &ipix,
	}
	actual := a.ToMastercat(0)
	require.Equal(t, expected, actual)
}

func TestAllwiseInputSchema_SetField(t *testing.T) {
	source := "source_id"
	ra := float64(0.0)
	dec := float64(0.0)
	tests := []struct {
		// Named input parameters for target function.
		name string
		val  any
	}{
		{"Source_id", &source},
		{"Ra", &ra},
		{"Dec", &dec},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AllwiseInputSchema{}
			a.SetField(tt.name, tt.val)
			val := reflect.ValueOf(a)
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			field := val.FieldByName(tt.name)
			if !field.IsValid() {
				t.Fatalf("field %s is not valid", tt.name)
			}
			require.Equal(t, tt.val, field.Interface())
		})
	}
}
