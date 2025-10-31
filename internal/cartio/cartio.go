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
// Optionally includes sfx.json, music.json, and sprites.json if provided.
func Write(w io.Writer, m Manifest, assets []Asset, sfx SFXMap, music MusicMap, sprites SpriteMap) error {
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

	// Always write sfx.json, music.json, and sprites.json (even if empty {})
	// This ensures all carts have these files as per user requirement
	sfxFile, err := zw.Create("assets/sfx.json")
	if err != nil {
		return err
	}
	enc = json.NewEncoder(sfxFile)
	enc.SetIndent("", "  ")
	if err = enc.Encode(sfx); err != nil {
		return err
	}

	musicFile, err := zw.Create("assets/music.json")
	if err != nil {
		return err
	}
	enc = json.NewEncoder(musicFile)
	enc.SetIndent("", "  ")
	if err = enc.Encode(music); err != nil {
		return err
	}

	spritesFile, err := zw.Create("assets/sprites.json")
	if err != nil {
		return err
	}
	enc = json.NewEncoder(spritesFile)
	enc.SetIndent("", "  ")
	if err = enc.Encode(sprites); err != nil {
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

// ReadResult contains all data read from a cart
type ReadResult struct {
	Manifest Manifest
	SFX      SFXMap
	Music    MusicMap
	Sprites  SpriteMap
	Files    map[string][]byte
}

// Read unpacks an .rfs archive into a manifest, sfx, music, and asset map.
func Read(r io.ReaderAt, size int64) (ReadResult, error) {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return ReadResult{}, err
	}
	var m Manifest
	var sfxMap SFXMap
	var musicMap MusicMap
	var spriteMap SpriteMap
	files := make(map[string][]byte)

	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			return ReadResult{}, err
		}
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, rc); err != nil {
			rc.Close()
			return ReadResult{}, err
		}
		rc.Close()

		switch f.Name {
		case "manifest.json":
			if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
				return ReadResult{}, err
			}
			continue
		case "assets/sfx.json":
			if err := json.Unmarshal(buf.Bytes(), &sfxMap); err != nil {
				// If sfx.json is invalid, use empty map
				sfxMap = make(SFXMap)
			}
			continue
		case "assets/music.json":
			if err := json.Unmarshal(buf.Bytes(), &musicMap); err != nil {
				// If music.json is invalid, use empty map
				musicMap = make(MusicMap)
			}
			continue
		case "assets/sprites.json":
			if err := json.Unmarshal(buf.Bytes(), &spriteMap); err != nil {
				// If sprites.json is invalid, use empty map
				spriteMap = make(SpriteMap)
			}
			continue
		}
		files[f.Name] = buf.Bytes()
	}

	// Initialize empty maps if files weren't found
	if sfxMap == nil {
		sfxMap = make(SFXMap)
	}
	if musicMap == nil {
		musicMap = make(MusicMap)
	}
	if spriteMap == nil {
		spriteMap = make(SpriteMap)
	}

	return ReadResult{
		Manifest: m,
		SFX:      sfxMap,
		Music:    musicMap,
		Sprites:  spriteMap,
		Files:    files,
	}, nil
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
