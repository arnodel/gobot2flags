package resources

//go:generate file2byteslice -package=resources -input=./greenflag.png -output=./greenflag.go -var=GreenFlagPng
//go:generate file2byteslice -package=resources -input=./greywalls.png -output=./greywalls.go -var=GreyWallsPng
//go:generate file2byteslice -package=resources -input=./floors.png -output=./floors.go -var=FloorsPng
//go:generate file2byteslice -package=resources -input=./robot.png -output=./robot.go -var=RobotPng
//go:generate file2byteslice -package=resources -input=./circuitboardtiles.png -output=./circuitboardtiles.go -var=CircuitBoarTiles

import (
	// Dummy imports for go.mod for some Go files with 'ignore' tags. For example, `go mod tidy` does not
	// recognize Go files with 'ignore' build tag.
	//
	// Note that this affects only importing this package, but not 'file2byteslice' commands in //go:generate.
	_ "github.com/hajimehoshi/file2byteslice"
)
