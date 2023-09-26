/*
 * Copyright (C) 2023 by Jason Figge
 */

package fonts

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Font int

const (
	Default Font = iota
)

var (
	fonts    = []Font{Default}
	fontSrc  = []string{"Tahoma.ttf"}
	ttfFonts = make([]*ttf.Font, len(fonts))
)

type Writer struct {
	texture  *sdl.Texture
	renderer *sdl.Renderer
	srcRect  *sdl.Rect
	destRect *sdl.Rect
}

func (w *Writer) Render(x int32, y int32) error {
	w.destRect.X = x
	w.destRect.Y = y
	return w.renderer.Copy(w.texture, w.srcRect, w.destRect)
}

func (w *Writer) Close() {
	w.texture.Destroy()
}

func Fonts() []Font {
	return fonts
}

func LoadFonts(renderer *sdl.Renderer) error {
	err := ttf.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize fonts package: %v", err)
	}

	var dir string
	dir, err = os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to get working directory: %w", err)
	}
	base := filepath.Join(dir, "resources", "fonts")
	for i, src := range fontSrc {
		ttfFonts[i], err = ttf.OpenFont(filepath.Join(base, fontSrc[i]), 12)
		if err != nil {
			return fmt.Errorf("failed to load font [%s]: %v", src, err)
		}
	}
	return nil
}

func FreeFonts() {
	if ttfFonts != nil {
		for _, ttfFont := range ttfFonts {
			if ttfFont != nil {
				ttfFont.Close()
			}
		}
	}
}

func (f Font) Font() *ttf.Font {
	return ttfFonts[f]
}

func (f Font) Size(text string) (int, int, error) {
	return ttfFonts[f].SizeUTF8(text)
}

func (f Font) Writer(renderer *sdl.Renderer, text string, fgColor uint32) (*Writer, error) {
	surface, err := ttfFonts[f].RenderUTF8Blended(text, sdl.Color(color.RGBA{
		R: uint8(fgColor >> 24),
		G: uint8(fgColor >> 16),
		B: uint8(fgColor >> 8),
		A: uint8(fgColor >> 24),
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to render text: %w", err)
	}
	defer surface.Free()

	var texture *sdl.Texture
	texture, err = renderer.CreateTextureFromSurface(surface)
	if err != nil {
		surface.Free()
		return nil, fmt.Errorf("failed to create texture from surface: %w", err)
	}

	srcRect := &sdl.Rect{X: 0, Y: 0, W: surface.W, H: surface.H}
	destRect := &sdl.Rect{X: 0, Y: 0, W: surface.W, H: surface.H}

	return &Writer{
		texture:  texture,
		renderer: renderer,
		srcRect:  srcRect,
		destRect: destRect,
	}, nil
}

func (f Font) PrintfAt(renderer *sdl.Renderer, x, y int32, fgColor uint32, format string, a ...any) {
	text := fmt.Sprintf(format, a...)
	writer, _ := f.Writer(renderer, text, fgColor)
	writer.Render(x, y)
	writer.Close()
}
