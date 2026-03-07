package store

import (
	"os"
	"path/filepath"
	"sort"
	"time"
)

type Storer interface {
	StoreFile(fileName string, hex []byte) (string, error)
	GetFilesSortedByModTime(root string) ([]string, error)
	ReadFile(fileName string) ([]byte, error)
	CheckFileExists(fileName string) (bool, error)
	RemoveFile(fileName string) error
	GetFilePath(fileName string) (string, error)
}

type FileStore struct {
	StorageLocation string
}

// StoreFile stores a secret in a default or user-defined directory
// provided by the --location flag
func (fStore *FileStore) StoreFile(secretName string, hex []byte) (string, error) {
	_, err := os.Stat(fStore.StorageLocation)
	if err != nil {
		err := os.Mkdir(fStore.StorageLocation, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	dstPath := fStore.StorageLocation + "/" + secretName

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}

	_, err = dstFile.Write([]byte(hex))
	if err != nil {
		return "", err
	}

	defer dstFile.Close()

	return dstPath, nil
}

// GetFilesSortedByModTime returns a list of files (relative paths) under root
// sorted by modification time ascending (oldest first). If mod times are equal
// it falls back to lexical order.
func (fStore *FileStore) GetFilesSortedByModTime(root string) ([]string, error) {
	type fileEntry struct {
		rel     string
		modTime time.Time
	}

	var entries []fileEntry
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		entries = append(entries, fileEntry{rel: rel, modTime: info.ModTime()})
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].modTime.Equal(entries[j].modTime) {
			return entries[i].rel < entries[j].rel
		}
		return entries[i].modTime.Before(entries[j].modTime)
	})

	var files []string
	for _, e := range entries {
		files = append(files, e.rel)
	}
	return files, nil
}

func (fStore *FileStore) ReadFile(fileName string) ([]byte, error) {
	hex, err := os.ReadFile(fStore.StorageLocation + "/" + fileName)
	if err != nil {
		return nil, err
	}

	return hex, nil
}

func (fStore *FileStore) CheckFileExists(fileName string) (bool, error) {
	_, err := os.Stat(fStore.StorageLocation + "/" + fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, err
}

func (fStore *FileStore) RemoveFile(fileName string) error {
	err := os.Remove(fStore.StorageLocation + "/" + fileName)
	if err != nil {
		return err
	}
	return nil
}

func (fStore *FileStore) GetFilePath(fileName string) (string, error) {
	exists, err := fStore.CheckFileExists(fileName)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", os.ErrNotExist
	}

	return fStore.StorageLocation + "/" + fileName, nil
}
