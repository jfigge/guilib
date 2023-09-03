/*
 * Copyright (C) 2023 by Jason Figge
 */

package graphics

import (
	"log"
	"runtime"
)

func DeferError(f func() error) {
	err := f()
	if err != nil {
		if _, file, line, ok := runtime.Caller(0); ok {
			log.Printf("trapped unexpected error from line %s:%d. Error: %v\n", file, line, err)
		} else {
			log.Printf("trapped unexpected error: %v\n", err)
		}
	}
}

func ErrorTrap(err error) {
	if err != nil {
		if _, file, line, ok := runtime.Caller(0); ok {
			log.Printf("trapped unexpected error from line %s:%d. Error: %v\n", file, line, err)
		} else {
			log.Printf("trapped unexpected error: %v\n", err)
		}
	}
}

func FMap(value, srcMin, srcMax, destMin, destMax float32) float32 {
	return (value - srcMin) * (destMax - destMin) / (srcMax - srcMin) + destMin;
}