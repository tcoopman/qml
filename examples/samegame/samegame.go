package main

import (
	"fmt"
	"github.com/niemeyer/qml"
	"os"
    "math/rand"
)

const (
	MAX_COL   = 10
	MAX_ROW   = 15
	MAX_INDEX = MAX_COL * MAX_ROW
)

//TODO unique seed
var r = rand.New(rand.NewSource(0))

type Game struct {
	MaxColumn  int
	MaxRow     int
	MaxIndex   int
	Board      []qml.Object
	Block      *Block
	fillFound  int
	floorBoard []int
	parent     qml.Object
}

func (g *Game) index(col, row int) int {
	return col + (row * g.MaxColumn)
}

func (g *Game) StartNewGame(parent qml.Object) {
	for _, b := range g.Board {
		b.Destroy()
	}

	g.parent = parent
	w := parent.Int("width")
	h := parent.Int("height")
	blockSize := parent.Int("blockSize")
	g.MaxColumn = w / blockSize
	g.MaxRow = h / blockSize
	g.MaxIndex = g.MaxColumn * g.MaxRow

	g.Block.BlockSize = blockSize

	g.Board = make([]qml.Object, g.MaxIndex, g.MaxIndex)
	for col := 0; col < g.MaxColumn; col++ {
		for row := 0; row < g.MaxRow; row++ {
			g.Board[g.index(col, row)] = g.Block.createBlock(col, row, parent)
		}
	}
}

func (g *Game) HandleClick(xPos, yPos int) {
	fmt.Println(xPos)
	fmt.Println(yPos)
	col := xPos / g.Block.BlockSize
	row := yPos / g.Block.BlockSize
	fmt.Printf("Clicking on col: %d, row: %d\n", col, row)

	if col >= g.MaxColumn || col < 0 || row >= g.MaxRow || row < 0 {
		return
	}
	if g.Board[g.index(col, row)] == nil {
		fmt.Println("it is nil")
		return
	}
	g.floodFill(col, row, -1)
	if g.fillFound <= 0 {
		fmt.Println("fillFound <= 0")
		return
	}

	// Set the score
	score := g.parent.Int("score")
	score += (g.fillFound - 1) * (g.fillFound - 1)
	g.parent.Set("score", score)

	g.shuffleDown()
	g.victoryCheck()
}

func (g *Game) floodFill(col, row, typ int) {
	if g.Board[g.index(col, row)] == nil {
		fmt.Println("is it nil???")
		return
	}
	first := false
	if typ == -1 {
		first = true
		typ = g.Board[g.index(col, row)].Int("type")

		g.fillFound = 0
		g.floorBoard = make([]int, g.MaxIndex, g.MaxIndex)
	}

	if col >= g.MaxColumn || col < 0 || row >= g.MaxRow || row < 0 {
		return
	}
	if g.floorBoard[g.index(col, row)] == 1 || (!first && typ != g.Board[g.index(col, row)].Int("type")) {
		return
	}

	g.floorBoard[g.index(col, row)] = 1
	g.floodFill(col+1, row, typ)
	g.floodFill(col-1, row, typ)
	g.floodFill(col, row+1, typ)
	g.floodFill(col, row-1, typ)
	if first && g.fillFound == 0 {
		return //Can't remove single blocks
	}
	g.Board[g.index(col, row)].Set("opacity", 0)
	g.Board[g.index(col, row)] = nil
	g.fillFound += 1
}

func (g *Game) shuffleDown() {
    // Fall down
    for col := 0; col < g.MaxColumn; col++ {
        fallDist := 0
        for row := g.MaxRow -1; row >= 0; row-- {
            if g.Board[g.index(row, col)] == nil {
                fallDist += 1
            } else {
                if fallDist > 0 {
                    obj := g.Board[g.index(col, row)]
                    y := obj.Int("y")
                    fmt.Printf("y was: %d\n", y)
                    y += fallDist + g.Block.BlockSize
                    fmt.Printf("Setting y to: %d\n", y)
                    obj.Set("y", y)
                    g.Board[g.index(col, row+fallDist)] = obj
                    g.Board[g.index(col, row)] = nil
                }
            }
        }
    }

    // Fall to the left
    fallDist := 0
    for col := 0; col < g.MaxColumn; col++ {
        if g.Board[g.index(col, g.MaxRow - 1)] == nil {
            fallDist += 1
        } else {
            if fallDist > 0 {
                for row := 0; row < g.MaxRow; row++ {
                    obj := g.Board[g.index(col, row)]
                    if obj == nil {
                        continue
                    }

                    x := obj.Int("x")
                    x -= fallDist * g.Block.BlockSize
                    obj.Set("x", x)
                    g.Board[g.index(col - fallDist, row)] = obj
                    g.Board[g.index(col, row)] = nil
                }
            }
        }
    }
}

func (g *Game) victoryCheck() {
	//TODO
}

type Block struct {
	Component qml.Object
	BlockSize int
}

func (b *Block) createBlock(col, row int, parent qml.Object) qml.Object {
	dynamicBlock := b.Component.Create(nil)
	dynamicBlock.Set("parent", parent)

    dynamicBlock.Set("type", r.Intn(3))
	dynamicBlock.Set("x", col*b.BlockSize)
	dynamicBlock.Set("y", row*b.BlockSize)
	dynamicBlock.Set("width", b.BlockSize)
	dynamicBlock.Set("height", b.BlockSize)

	return dynamicBlock
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	qml.Init(nil)
	engine := qml.NewEngine()

	component, err := engine.LoadFile("samegame.qml")
	if err != nil {
		return err
	}

	game := Game{
		MaxColumn: MAX_COL,
		MaxRow:    MAX_ROW,
		MaxIndex:  MAX_COL * MAX_ROW,
	}

	context := engine.Context()
	context.SetVar("game", &game)

	win := component.CreateWindow(nil)

	blockComponent, err := engine.LoadFile("Block.qml")
	if err != nil {
		return err
	}

	block := &Block{Component: blockComponent}
	game.Block = block

	win.Show()
	win.Wait()

	return nil
}