package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileStoreStoreReadCheckRemoveAndPath(t *testing.T) {
	base := t.TempDir()
	vault := filepath.Join(base, "vault")
	fs := &FileStore{StorageLocation: vault}

	path, err := fs.StoreFile("github", []byte("payload"))
	if err != nil {
		t.Fatalf("StoreFile() error = %v", err)
	}
	if path != filepath.Join(vault, "github") {
		t.Fatalf("StoreFile() path = %q, want %q", path, filepath.Join(vault, "github"))
	}

	got, err := fs.ReadFile("github")
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(got) != "payload" {
		t.Fatalf("ReadFile() = %q, want %q", string(got), "payload")
	}

	exists, err := fs.CheckFileExists("github")
	if err != nil {
		t.Fatalf("CheckFileExists() error = %v", err)
	}
	if !exists {
		t.Fatalf("CheckFileExists() = false, want true")
	}

	fullPath, err := fs.GetFilePath("github")
	if err != nil {
		t.Fatalf("GetFilePath() error = %v", err)
	}
	if fullPath != filepath.Join(vault, "github") {
		t.Fatalf("GetFilePath() = %q, want %q", fullPath, filepath.Join(vault, "github"))
	}

	if err := fs.RemoveFile("github"); err != nil {
		t.Fatalf("RemoveFile() error = %v", err)
	}

	exists, err = fs.CheckFileExists("github")
	if err != nil {
		t.Fatalf("CheckFileExists() after remove error = %v", err)
	}
	if exists {
		t.Fatalf("CheckFileExists() after remove = true, want false")
	}

	_, err = fs.GetFilePath("github")
	if !os.IsNotExist(err) {
		t.Fatalf("GetFilePath() missing error = %v, want os.ErrNotExist", err)
	}
}

func TestFileStoreGetFilesSortedByModTime(t *testing.T) {
	root := t.TempDir()
	fs := &FileStore{StorageLocation: root}

	files := []string{"c", "b", "a"}
	for _, name := range files {
		if err := os.WriteFile(filepath.Join(root, name), []byte(name), 0600); err != nil {
			t.Fatalf("WriteFile(%s) error = %v", name, err)
		}
	}

	old := time.Now().Add(-2 * time.Hour)
	newer := time.Now().Add(-1 * time.Hour)

	if err := os.Chtimes(filepath.Join(root, "a"), old, old); err != nil {
		t.Fatalf("Chtimes(a) error = %v", err)
	}
	if err := os.Chtimes(filepath.Join(root, "b"), old, old); err != nil {
		t.Fatalf("Chtimes(b) error = %v", err)
	}
	if err := os.Chtimes(filepath.Join(root, "c"), newer, newer); err != nil {
		t.Fatalf("Chtimes(c) error = %v", err)
	}

	got, err := fs.GetFilesSortedByModTime(root)
	if err != nil {
		t.Fatalf("GetFilesSortedByModTime() error = %v", err)
	}

	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("GetFilesSortedByModTime() len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("GetFilesSortedByModTime()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
