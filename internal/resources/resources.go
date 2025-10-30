package resources

import (
    "errors"
    "io/fs"
)

var ErrNotFound = errors.New("resource not found")

// Store is a simple in-memory asset store keyed by path.
type Store struct { data map[string][]byte }

func New() *Store { return &Store{ data: make(map[string][]byte) } }

func (s *Store) Add(path string, b []byte) { s.data[path] = append([]byte{}, b...) }
func (s *Store) Get(path string) ([]byte, error) {
    b, ok := s.data[path]
    if !ok { return nil, ErrNotFound }
    return append([]byte{}, b...), nil
}

// LoadFrom walks a fs.FS and loads all files into the store.
func (s *Store) LoadFrom(fsys fs.FS) error {
    return fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
        if err != nil { return err }
        if d.IsDir() { return nil }
        b, err := fs.ReadFile(fsys, p)
        if err != nil { return err }
        s.Add(p, b)
        return nil
    })
}


