package cartio

import (
	"bytes"
	"testing"
)

func TestWriteReadRoundTrip(t *testing.T) {
	m := Manifest{Title: "Hello", Author: "RF", Description: "d", Genre: "Action", Tags: []string{"a", "b"}, Entry: "main.lua"}
	assets := []Asset{{Name: "main.lua", Data: []byte("print('hi')")}, {Name: "sprites.png", Data: []byte{1, 2, 3}}}

	var buf bytes.Buffer
	if err := Write(&buf, m, assets, make(SFXMap), make(MusicMap), make(SpriteMap)); err != nil {
		t.Fatalf("write: %v", err)
	}

	result, err := Read(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if result.Manifest.Entry != "main.lua" || result.Manifest.Title != "Hello" {
		t.Fatalf("bad manifest: %+v", result.Manifest)
	}
	if len(result.Files) != 2 {
		t.Fatalf("expected 2 assets, got %d", len(result.Files))
	}
	if _, ok := result.Files["assets/main.lua"]; !ok {
		t.Fatalf("missing main.lua")
	}
}

func TestSortedAssetNames(t *testing.T) {
	m := map[string][]byte{
		"zebra":  []byte{1, 2, 3},
		"apple":  []byte{4, 5},
		"banana": []byte{6, 7, 8, 9},
		"cherry": []byte{10},
	}

	names := SortedAssetNames(m)

	expected := []string{"apple", "banana", "cherry", "zebra"}
	if len(names) != len(expected) {
		t.Fatalf("SortedAssetNames returned %d names, expected %d", len(names), len(expected))
	}

	for i, name := range names {
		if name != expected[i] {
			t.Errorf("SortedAssetNames[%d] = %q, expected %q", i, name, expected[i])
		}
	}

	// Test empty map
	empty := SortedAssetNames(map[string][]byte{})
	if len(empty) != 0 {
		t.Errorf("SortedAssetNames of empty map should return empty slice, got %d", len(empty))
	}

	// Test single element
	single := SortedAssetNames(map[string][]byte{"only": []byte{1}})
	if len(single) != 1 || single[0] != "only" {
		t.Errorf("SortedAssetNames of single element failed: %v", single)
	}
}
