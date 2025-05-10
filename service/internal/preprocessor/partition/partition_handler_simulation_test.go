package partition

import (
	"testing"
)

func FuzzGetPartition(f *testing.F) {
	hash := make(map[string]string)
	handler := PartitionHandler{NumPartitions: 16, PartitionLevels: 2}

	f.Add("ZTF20aaelulu")
	f.Add("0438p015_ac51-015430")
	f.Add("VLASS2QLCIR J074401.08+503011.3")

	f.Fuzz(func(t *testing.T, oid string) {
		p, err := handler.GetPartition(oid)
		if err != nil {
			t.Fatal(err)
		}

		dir, err := p.LevelsToString(handler.NumPartitions)
		if err != nil {
			t.Fatal(err)
		}

		if _, ok := hash[oid]; !ok {
			hash[oid] = dir
		} else {
			if dir != hash[oid] {
				t.Fatalf("Expected %s, got %s", hash[oid], dir)
			}
		}
	})

}
