//go:build !release
// +build !release

package main

import (
	"game_bp/starter"

	"github.com/Liphium/magic/v2"
)

func main() {
	magic.Start(starter.BuildMagicConfig())
}
