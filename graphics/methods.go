/*
 * Copyright (C) 2023 by Jason Figge
 */

package graphics

import (
	"fmt"
	"time"

	"github.com/jfigge/guilib/graphics/fonts"

	"github.com/veandco/go-sdl2/sdl"
)

type CoreMethods struct {
	timer           time.Time
	frames          int
	frameRateWriter *fonts.Writer
}

func (cm *CoreMethods) WriteFrameRate(renderer *sdl.Renderer, x int32, y int32) error {
	if time.Now().After(cm.timer) {
		cm.timer = time.Now().Add(time.Second)
		if cm.frameRateWriter != nil {
			cm.frameRateWriter.Close()
		}
		var err error
		cm.frameRateWriter, err = fonts.Default.Writer(renderer, fmt.Sprintf("Frame rate: %d", cm.frames), 0xFFFFFF)
		if err != nil {
			return fmt.Errorf("failed to generate font writer: %w", err)
		}
		cm.frames = 0
	}
	cm.frames++
	if cm.frameRateWriter != nil {
		err := cm.frameRateWriter.Render(x, y)
		if err != nil {
			return fmt.Errorf("failed call render font writer: %w", err)
		}
	}

	return nil
}

func (cm *CoreMethods) Clear(renderer *sdl.Renderer, bgColor uint32) error {
	err := renderer.SetDrawColor(uint8(bgColor>>16), uint8(bgColor>>8), uint8(bgColor), 0xFF)
	if err != nil {
		return fmt.Errorf("failed to set renderer draw color: %w", err)
	}
	err = renderer.Clear()
	if err != nil {
		return fmt.Errorf("failed to set renderer draw color: %w", err)
	}
	return nil
}
