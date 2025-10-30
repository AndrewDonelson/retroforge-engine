package resources

import (
    "testing"
    "testing/fstest"
)

func TestStoreAddGet(t *testing.T) {
    s := New()
    s.Add("a.txt", []byte("hello"))
    b, err := s.Get("a.txt")
    if err != nil || string(b) != "hello" { t.Fatalf("got %q err %v", string(b), err) }
    // ensure copy
    b[0] = 'X'
    b2, _ := s.Get("a.txt")
    if string(b2) != "hello" { t.Fatalf("mutability leak") }
}

func TestLoadFromFS(t *testing.T) {
    fsys := fstest.MapFS{
        "dir/file.txt":  &fstest.MapFile{Data: []byte("data")},
        "root.json":     &fstest.MapFile{Data: []byte("{}")},
    }
    s := New()
    if err := s.LoadFrom(fsys); err != nil { t.Fatalf("load err: %v", err) }
    if _, err := s.Get("root.json"); err != nil { t.Fatalf("missing root.json") }
    if _, err := s.Get("dir/file.txt"); err != nil { t.Fatalf("missing dir/file.txt") }
    if _, err := s.Get("nope.txt"); !errorsIs(err, ErrNotFound) { t.Fatalf("expected not found") }
}

func errorsIs(err, target error) bool { if err == nil { return target == nil }; return err.Error() == target.Error() }


