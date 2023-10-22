package openapi

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestA(t *testing.T) {
	res := filepath.Join("/", defaultYsId, "subdir1")
	fmt.Println(res)
}
