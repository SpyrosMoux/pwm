package secrets

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ListSecretsNumbered prints a tree with numeric indices for files.
// Indices correspond to the sorted list used by resolveSecretArg.
func (svc *Service) ListSecretsNumbered(root string, level int) error {
	files, err := svc.Store.GetFilesSortedByModTime(root)
	if err != nil {
		return err
	}

	indexMap := make(map[string]int)
	for i, f := range files {
		indexMap[f] = i + 1
	}

	// Build dirMin map: smallest index for files under each directory
	dirMin := make(map[string]int)
	const maxInt = int(^uint(0) >> 1)
	dirMin["."] = maxInt
	for rel, idx := range indexMap {
		if idx < dirMin["."] {
			dirMin["."] = idx
		}
		dir := filepath.Dir(rel)
		for dir != "." && dir != string(filepath.Separator) && dir != "" {
			cur, ok := dirMin[dir]
			if !ok || idx < cur {
				dirMin[dir] = idx
			}
			dir = filepath.Dir(dir)
		}
	}

	var printTree func(path string, level int) error
	printTree = func(path string, level int) error {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		type node struct {
			entry os.DirEntry
			key   int
		}
		nodes := make([]node, 0, len(entries))
		for _, entry := range entries {
			if entry.IsDir() {
				relDir, _ := filepath.Rel(root, filepath.Join(path, entry.Name()))
				key := dirMin[relDir]
				if key == 0 {
					key = maxInt
				}
				nodes = append(nodes, node{entry: entry, key: key})
			} else {
				relFile, _ := filepath.Rel(root, filepath.Join(path, entry.Name()))
				key := indexMap[relFile]
				if key == 0 {
					key = maxInt
				}
				nodes = append(nodes, node{entry: entry, key: key})
			}
		}

		sort.Slice(nodes, func(i, j int) bool {
			if nodes[i].key == nodes[j].key {
				return nodes[i].entry.Name() < nodes[j].entry.Name()
			}
			return nodes[i].key < nodes[j].key
		})

		for i, n := range nodes {
			entry := n.entry
			prefix := strings.Repeat("│   ", level)
			isLast := i == len(nodes)-1
			connector := "├── "
			if isLast {
				connector = "└── "
			}

			if entry.IsDir() {
				fmt.Printf("%s%s%s\n", prefix, connector, entry.Name())
				err := printTree(filepath.Join(path, entry.Name()), level+1)
				if err != nil {
					return err
				}
			} else {
				rel, err := filepath.Rel(root, filepath.Join(path, entry.Name()))
				if err != nil {
					return err
				}
				idx := indexMap[rel]
				if idx > 0 {
					fmt.Printf("%s%s{%d} %s\n", prefix, connector, idx, entry.Name())
				} else {
					fmt.Printf("%s%s%s\n", prefix, connector, entry.Name())
				}
			}
		}

		return nil
	}

	return printTree(root, level)
}
