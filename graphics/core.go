/*
 * Copyright (C) 2023 by Jason Figge
 */

package graphics

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type runState int

const (
	initialized runState = iota
	running
	terminating
	stopped
)

type Canvas struct {
	*sdl.Window
	state          runState
	renderer       *sdl.Renderer
	handler        GUIHandler
	glwg           sync.WaitGroup
	frameRateTimer *time.Ticker
	done           chan bool
	lock           sync.Mutex
}

type config struct {
	windowX       int32
	windowY       int32
	windowFlags   uint32
	framerate     int
	rendererFlags uint32
}

type ConfigOption func(c *config)

func Open(title string, width, height int32, handler GUIHandler, options ...ConfigOption) {
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(fmt.Errorf("failed to started sdl: %w", err))
	}

	cfg := &config{
		windowX:       sdl.WINDOWPOS_UNDEFINED,
		windowY:       sdl.WINDOWPOS_UNDEFINED,
		windowFlags:   sdl.WINDOW_OPENGL,
		framerate:     60,
		rendererFlags: sdl.RENDERER_ACCELERATED, //|sdl.RENDERER_PRESENTVSYNC
	}
	for _, option := range options {
		option(cfg)
	}

	window, err := sdl.CreateWindow(
		title,
		cfg.windowX,
		cfg.windowY,
		width, height,
		cfg.windowFlags,
	)
	if err != nil {
		panic(fmt.Errorf("failed to create window: %w", err))
	}
	c := &Canvas{
		state: initialized,
		done:  make(chan bool),
	}

	c.renderer, err = sdl.CreateRenderer(window, -1, cfg.rendererFlags)
	if err != nil {
		panic(fmt.Errorf("failed to create renderer: %w", err))
	}
	c.handler = handler
	c.handler.Init(c)
	c.start(cfg.framerate)
}

func (c *Canvas) start(framerate int) {
	if c.state != initialized {
		return
	}
	c.state = running

	c.handler.Init(c)
	go c.gameLoop(framerate)

	fmt.Println("Event loop starting")
	for c.state == running {
		// Process window events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			func() {
				c.lock.Lock()
				defer c.lock.Unlock()
				processed := false
				if c.handler != nil {
					processed = c.handler.Events(event)
				}
				if !processed {
					switch e := event.(type) {
					case *sdl.QuitEvent:
						fmt.Printf("Quit event: %+v\n", e)
						c.Quit()
					}
				}
			}()
		}
		time.Sleep(time.Millisecond * 10)
	}

	fmt.Println("Destroying handlers")
	c.handler.Destroy()
	fmt.Println("Destroying renderer")
	DeferError(c.renderer.Destroy)
	fmt.Println("Destroying canvas")
	DeferError(c.Destroy)
	fmt.Println("Waiting for event loop to quit")
	fmt.Println("Quitting sdl")
	sdl.Quit()
	fmt.Println("Stopped")

}

func (c *Canvas) gameLoop(framerate int) {
	c.glwg.Add(1)
	fr := time.Second / time.Duration(framerate)
	go func() {
		defer func() {
			fmt.Println("Game loop Exited")
			c.glwg.Done()
		}()
		defer c.panicHandler("game loop")()
		fmt.Println("Game loop starting")
		c.frameRateTimer = time.NewTicker(fr)
		for c.state == running {
			select {
			case <-c.done:
				fmt.Println("game Loop - Done")
				return
			case <-c.frameRateTimer.C:
				func() {
					c.lock.Lock()
					defer c.lock.Unlock()
					// Update state
					c.handler.OnUpdate()

					// Handle draw canvas
					c.handler.OnDraw(c.renderer)

					// Render the image
					c.renderer.Present()
				}()
			}
		}
	}()
}

func (c *Canvas) Quit() {
	if c.state == running {
		fmt.Println("terminating")
		c.state = terminating
		fmt.Println("stopping framerate timer")
		c.frameRateTimer.Stop()
		fmt.Println("Ending framerate timer")
		c.done <- true
		fmt.Println("waiting for game loop to exit")
		c.glwg.Wait()
		fmt.Println("Game loop complete.")
		c.state = stopped
	} else {
		fmt.Println("already quitting")
	}
}

func (c *Canvas) IsTerminated() bool {
	return c.state != running
}

func (c *Canvas) Renderer() *sdl.Renderer {
	return c.renderer
}

func (c *Canvas) panicHandler(name string) func() {
	return func() {
		if err := recover(); err != nil {
			fmt.Printf("Panic detected in %s: %v\n", name, err)
			for i := 3; i < 10; i++ {
				if _, file, line, ok := runtime.Caller(i); ok {
					fmt.Printf("%s:%d\n", file, line)
				}
			}
			c.Quit()
		}
		return
	}
}

func Framerate(fr int) ConfigOption {
	return func(c *config) {
		if fr > 0 {
			c.framerate = fr
		}
	}
}

func RendererFlags(flags uint32) ConfigOption {
	return func(c *config) {
		c.rendererFlags = flags &
			sdl.RENDERER_SOFTWARE &
			sdl.RENDERER_ACCELERATED &
			sdl.RENDERER_PRESENTVSYNC &
			sdl.RENDERER_TARGETTEXTURE
	}
}

func WindowPosition(x, y int32) ConfigOption {
	return func(c *config) {
		c.windowX = x
		c.windowY = y
	}
}

func WindowFlags(flags uint32) ConfigOption {
	return func(c *config) {
		c.windowFlags = flags &
			sdl.WINDOW_FULLSCREEN &
			sdl.WINDOW_FULLSCREEN_DESKTOP &
			sdl.WINDOW_OPENGL &
			sdl.WINDOW_VULKAN &
			sdl.WINDOW_HIDDEN &
			sdl.WINDOW_BORDERLESS &
			sdl.WINDOW_RESIZABLE &
			sdl.WINDOW_MINIMIZED &
			sdl.WINDOW_MAXIMIZED &
			sdl.WINDOW_INPUT_GRABBED &
			sdl.WINDOW_ALLOW_HIGHDPI
	}
}
