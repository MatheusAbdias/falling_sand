package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 400
	screenHeight = 400
	cellSize     = 1
)

type Game struct {
	pixels         []byte
	grid           [][]bool
	updateInterval time.Duration
	lastUpdateTime time.Time
}

func NewGame() *Game {
	gridWidth := screenWidth / cellSize
	gridHeight := screenHeight / cellSize
	grid := make([][]bool, gridWidth)
	for i := range grid {
		grid[i] = make([]bool, gridHeight)
	}

	return &Game{
		pixels:         make([]byte, screenWidth*screenHeight*4),
		grid:           grid,
		updateInterval: time.Millisecond * 5,
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	for x, vx := range g.grid {
		for y, vy := range vx {
			px := float32(cellSize * x)
			py := float32(cellSize * y)

			if vy {
				g.drawFilled(screen, px, py)
			}
		}
	}
}

func (g *Game) drawFilled(screen *ebiten.Image, px, py float32) {
	white := color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	vector.DrawFilledRect(
		screen,
		px,
		py,
		cellSize,
		cellSize,
		white,
		true,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Update() error {
	now := time.Now()
	if now.Sub(g.lastUpdateTime) >= g.updateInterval {
		g.lastUpdateTime = now

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			mouseX, mouseY := ebiten.CursorPosition()
			mouseCol := mouseX / cellSize
			mouseRow := mouseY / cellSize
			matrix := 10
			extend := matrix / 2

			for x := -extend; x < extend; x++ {
				for y := -extend; y < extend; y++ {
					col := mouseCol + x
					row := mouseRow + y
					if col >= 0 && col < len(g.grid) && row >= 0 &&
						row < len(g.grid[col]) {
						g.grid[col][row] = true
					}
				}
			}
		}

		for x := 0; x < len(g.grid); x++ {
			for y := len(g.grid[x]) - 1; y >= 0; y-- {
				currentCell := g.grid[x][y]
				if currentCell {
					if y < len(g.grid[x])-1 {
						bellowCell := g.grid[x][y+1]
						if !bellowCell {
							g.grid[x][y] = false
							g.grid[x][y+1] = true
						} else {
							pos := 1
							if rand.Float64() < 0.5 {
								pos = pos * -1
							}
							var bellowFCell bool
							var bellowSCell bool
							if x < len(g.grid)-1 && x > 0 {
								bellowFCell = g.grid[x+pos][y+1]
								bellowSCell = g.grid[x-pos][y+1]
								if !bellowFCell {
									g.grid[x][y] = false
									g.grid[x+pos][y+1] = true
								} else if !bellowSCell {
									g.grid[x][y] = false
									g.grid[x-pos][y+1] = true
								}
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func main() {
	g := NewGame()

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("## Landing Sand ##")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
