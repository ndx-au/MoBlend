package resources_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ndx-au/MoBlend/internal/resources"
)

func TestExtractBasePack(t *testing.T) {
	dir := t.TempDir()
	path, err := resources.ExtractBasePack(dir)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(path) != "base.blend" {
		t.Fatalf("unexpected path base: %s", path)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() < 1024 {
		t.Fatalf("pack too small: %d bytes", info.Size())
	}
}

func TestBasePackSHA256Stable(t *testing.T) {
	a, err := resources.BasePackSHA256()
	if err != nil {
		t.Fatal(err)
	}
	b, err := resources.BasePackSHA256()
	if err != nil {
		t.Fatal(err)
	}
	if a != b || len(a) != 64 {
		t.Fatalf("unstable or bad sha256: %q / %q", a, b)
	}
}
