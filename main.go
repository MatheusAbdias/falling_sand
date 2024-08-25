package main

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 400
	screenHeight = 400
	cellSize     = 10
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
		updateInterval: time.Millisecond * 100,
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	for x, vx := range g.grid {
		for y, vy := range vx {
			px := float32(cellSize * x)
			py := float32(cellSize * y)

			if vy {
				g.drawFilled(screen, px, py)
			} else {
				g.drawStoke(screen, px, py)
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

func (g *Game) drawStoke(screen *ebiten.Image, px, py float32) {
	white := color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	vector.StrokeRect(
		screen,
		px,
		py,
		cellSize,
		cellSize,
		1,
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
			cellX := mouseX / cellSize
			cellY := mouseY / cellSize
			if cellX >= 0 && cellX < len(g.grid) && cellY >= 0 && cellY < len(g.grid[cellX]) {
				g.grid[cellX][cellY] = true
			}
		}

		for x := 0; x < len(g.grid); x++ {
			for y := len(g.grid[x]) - 1; y >= 0; y-- {
				if g.grid[x][y] && y < len(g.grid[x])-1 && !g.grid[x][y+1] {
					g.grid[x][y+1] = true
					g.grid[x][y] = false
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
