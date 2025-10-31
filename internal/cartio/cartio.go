package cartio

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"path"
	"sort"
)

// Manifest is the minimal metadata stored in an .rfs
type Manifest struct {
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	Genre       string   `json:"genre"`
	Tags        []string `json:"tags"`
	Entry       string   `json:"entry"`             // e.g. main.lua
	Palette     string   `json:"palette,omitempty"` // Optional palette name (e.g., "RetroForge 50")
	Scale       *int     `json:"scale,omitempty"`   // Optional default scale for cart display
}

// Asset represents a file to be packed.
type Asset struct {
	Name string
	Data []byte
}

// Write packs a manifest and assets into an .rfs (zip) archive.
func Write(w io.Writer, m Manifest, assets []Asset) error {
	zw := zip.NewWriter(w)
	defer zw.Close()

	// manifest.json
	mf, err := zw.Create("manifest.json")
	if err != nil {
		return err
	}
	enc := json.NewEncoder(mf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(&m); err != nil {
		return err
	}

	// assets/
	for _, a := range assets {
		name := path.Join("assets", a.Name)
		f, err := zw.Create(name)
		if err != nil {
			return err
		}
		if _, err := f.Write(a.Data); err != nil {
			return err
		}
	}
	return zw.Close()
}

// Read unpacks an .rfs archive into a manifest and asset map.
func Read(r io.ReaderAt, size int64) (Manifest, map[string][]byte, error) {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return Manifest{}, nil, err
	}
	var m Manifest
	files := make(map[string][]byte)
	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			return Manifest{}, nil, err
		}
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, rc); err != nil {
			rc.Close()
			return Manifest{}, nil, err
		}
		rc.Close()
		if f.Name == "manifest.json" {
			if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
				return Manifest{}, nil, err
			}
			continue
		}
		files[f.Name] = buf.Bytes()
	}
	return m, files, nil
}

// SortedAssetNames returns deterministic ordering for tests.
func SortedAssetNames(m map[string][]byte) []string {
	names := make([]string, 0, len(m))
	for k := range m {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
