//go:build !js && !wasm

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
	"github.com/AndrewDonelson/retroforge-engine/internal/engine"
	"github.com/AndrewDonelson/retroforge-engine/internal/sdlrun"
)

func packDir(dir, out string) error {
	// read manifest.json
	mfBytes, err := os.ReadFile(filepath.Join(dir, "manifest.json"))
	if err != nil {
		return err
	}
	var m cartio.Manifest
	if err := json.Unmarshal(mfBytes, &m); err != nil {
		return err
	}
	var assets []cartio.Asset
	// walk assets/
	assetsDir := filepath.Join(dir, "assets")
	_ = filepath.WalkDir(assetsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(assetsDir, path)
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		assets = append(assets, cartio.Asset{Name: rel, Data: b})
		return nil
	})
	var buf bytes.Buffer
	if err := cartio.Write(&buf, m, assets); err != nil {
		return err
	}
	return os.WriteFile(out, buf.Bytes(), 0644)
}

func savePNG(path string, w, h int, rgba []uint8) error {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	copy(img.Pix, rgba)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func main() {
	pack := flag.String("pack", "", "pack cart directory into .rfs (specify input dir)")
	cart := flag.String("cart", "", "run .rfs cart (specify file path)")
	frames := flag.Int("frames", 1, "frames to run when executing a cart (headless)")
	out := flag.String("out", "", "output PNG path (headless). Omit to disable.")
	window := flag.Bool("window", false, "open window and run until ESC/Close")
	scale := flag.Int("scale", 2, "window scale (integer)")
	flag.Parse()

	if *pack != "" {
		outFile := *pack + ".rf"
		if err := packDir(*pack, outFile); err != nil {
			panic(err)
		}
		println("packed:", outFile)
		return
	}

	if *cart != "" {
		e := engine.New(60)
		defer e.Close()
		if err := e.LoadCartFile(*cart); err != nil {
			panic(err)
		}
		if *window {
			if err := sdlrun.RunWindow(e, *scale); err != nil {
				panic(err)
			}
			return
		}
		// headless
		e.RunFrames(*frames)
		if *out != "" {
			if err := savePNG(*out, e.Ren.Width(), e.Ren.Height(), e.Ren.Pixels()); err != nil {
				panic(err)
			}
			println("wrote:", *out)
		} else {
			pix := e.Ren.Pixels()
			_ = color.RGBA{R: pix[0], G: pix[1], B: pix[2], A: pix[3]}
		}
		return
	}

	flag.Usage()
}
