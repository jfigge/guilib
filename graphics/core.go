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

func Open(title string, width, height int32, handler GUIHandler) {
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(fmt.Errorf("failed to started sdl: %w", err))
	}

	window, err := sdl.CreateWindow(
		title,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		width, height,
		sdl.WINDOW_OPENGL,
	)
	if err != nil {
		panic(fmt.Errorf("failed to create window: %w", err))
	}
	c := &Canvas{
		Window: window,
		state:  initialized,
		done:   make(chan bool),
	}

	c.renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED) //|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		panic(fmt.Errorf("failed to create renderer: %w", err))
	}
	c.handler = handler
	c.handler.Init(c)
	c.start()
}

func (c *Canvas) start() {
	if c.state != initialized {
		return
	}
	c.state = running

	c.handler.Init(c)
	go c.gameLoop()

	fmt.Println("Event loop starting")
	for c.state == running {
		// Process window events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			c.lock.Lock()
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
			c.lock.Unlock()
		}
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

func (c *Canvas) gameLoop() {
	c.glwg.Add(1)
	go func() {
		defer func() {
			fmt.Println("Game loop Exited")
			c.glwg.Done()
		}()
		defer c.panicHandler("game loop")()
		fmt.Println("Game loop starting")
		c.frameRateTimer = time.NewTicker(time.Second / 60)
		for c.state == running {
			select {
			case <-c.done:
				fmt.Println("game Loop - Done")
				return
			case <-c.frameRateTimer.C:
				c.lock.Lock()
				// Update state
				c.handler.OnUpdate()

				// Handle draw canvas
				c.handler.OnDraw(c.renderer)

				// Render the image
				c.renderer.Present()
				c.lock.Unlock()
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
