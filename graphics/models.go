/*
 * Copyright (C) 2023 by Jason Figge
 */

package graphics

import (
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

type GUIHandler interface {
	Init(canvas *Canvas)
	Events(event sdl.Event) bool
	OnUpdate()
	OnDraw(renderer *sdl.Renderer)
	Destroy()
}

type BaseHandler struct {
	lock       sync.Mutex
	destroyers []func()
}

func (b *BaseHandler) AddDestroyer(destroyer func()) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.destroyers == nil {
		b.destroyers = make([]func(), 0)
	}
	b.destroyers = append(b.destroyers, destroyer)
}

func (b *BaseHandler) Quit() {
	sdl.PushEvent(&sdl.QuitEvent{
		Type:      256,
		Timestamp: 2000,
	})
}

func (b *BaseHandler) Init(canvas *Canvas) {
}

func (b *BaseHandler) Events(event sdl.Event) bool {
	return false
}

func (b *BaseHandler) OnUpdate() {
}

func (b *BaseHandler) OnDraw(renderer *sdl.Renderer) {
}

func (b *BaseHandler) Destroy() {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.destroyers != nil {
		for _, destroyer := range b.destroyers {
			destroyer()
		}
		b.destroyers = nil
	}
}
