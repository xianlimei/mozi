package util

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestTraversalDirByExt(t *testing.T) {
	sofiles := []string{"test.so", "boo.so"}
	files := append(sofiles, "test.txt")
	for _, f := range files {
		fi, _ := os.Create(f)
		defer os.Remove(f)
		defer fi.Close()
	}

	filterRes := []string{}
	TraversalDirByExt(filepath.Join("."), ".so", func(f string) {
		filterRes = append(filterRes, f)
	})
	sort.Strings(filterRes)
	sort.Strings(sofiles)
	if len(sofiles) != len(filterRes) || !reflect.DeepEqual(filterRes, sofiles) {
		t.Errorf("should filter out all so file")
	}

}
