package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

func loadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

const MapWidth int = 80 * 2
const MapHeight int = 24 * 2

type Tile struct {
	id           int
	display      byte
	displayColor color.Color
	name         string
	collides     bool
}

// color
// https://lospec.com/palette-list/aap-64

var aapDarkGrey = color.RGBA{20, 16, 19, 0xFF}
var aapPink = color.RGBA{245, 160, 151, 0xFF}
var aapGrey = color.RGBA{50, 43, 40, 0xFF}
var aapMuddy = color.RGBA{121, 103, 85, 0xFF}

func lookupTile(id int) Tile {
	Tiles := map[int]Tile{
		0x00: Tile{0x00, ' ', aapDarkGrey, "Empty", false},
		0x01: Tile{0x01, '@', aapPink, "Player", false},
		0x02: Tile{0x02, '-', aapMuddy, "Ground", false},
		0x03: Tile{0x03, '#', aapGrey, "Wall", true},
		0x04: Tile{},
		0x05: Tile{},
		0x06: Tile{},
		0x07: Tile{},
		0x08: Tile{},
		0x09: Tile{},
		0x0A: Tile{},
		0x0B: Tile{},
		0x0C: Tile{},
		0x0D: Tile{},
		0x0E: Tile{},
		0x0F: Tile{},
		0x10: Tile{},
		0x11: Tile{},
		0x12: Tile{},
	}
	return Tiles[id]
}

func generateMap(sepr int) [MapWidth][MapHeight]Tile {

	// var n = opensimplex.New(rand.Int63())

	var a = [MapWidth][MapHeight]Tile{}

	var i, j int
	for i = 0; i < MapWidth; i++ {
		for j = 0; j < MapHeight; j++ {
			a[i][j] = lookupTile(0x03)
			//var index = n.Eval2(((float64)(i)*0.1)+0.5, ((float64)(j)*0.1)+0.5)
			//if index > 0.3 {
			//	a[i][j] = lookupTile(0x02)
			//} else {
			//	a[i][j] = lookupTile(0x03)
			//}
		}
	}

	var roomCount = rand.Intn(40) + 5
	var roomSize = 10
	var minRoomSize = 3
	var c int

	var roomMap = make([]Rect, 0)
	for c = 0; c < roomCount; c++ {
		var p1 = MapWidth/3 + rand.Intn(MapWidth/5)
		var p2 = MapHeight/3 + rand.Intn(MapHeight/5)
		var p3 = rand.Intn(roomSize) + minRoomSize
		var p4 = rand.Intn(roomSize) + minRoomSize
		m := Rect{pos: Vec2{x: p1, y: p2}, size: Vec2{x: p3 + 5, y: p4}}
		roomMap = append(roomMap, m)
	}

	var movedRooms [50]bool
	for true {
		var p int
		for p = 0; p < roomCount; p++ {
			movedRooms[p] = false
		}

		var c int
		for c = 0; c < roomCount; c++ {
			var j int
			for j = 0; j < roomCount; j++ {
				if c == j {
					continue
				}
				if !movedRooms[j] {
					if roomMap[c].expand(1).intersects(roomMap[j].expand(1)) {
						movedRooms[j] = true
						var sepVec = Vec2Math.add(roomMap[c].center(), roomMap[j].center().invert()).normalise()
						if sepVec == Vec2Zero {
							sepVec = Vec2Math.random(Vec2Zero)
						}
						roomMap[c].pos = Vec2{
							roomMap[c].pos.x + sepVec.x, roomMap[c].pos.y + sepVec.y,
						}
					}
				}
			}
		}

		var finished = true
		for p = 0; p < roomCount; p++ {
			if movedRooms[p] {
				finished = false
				break
			}
		}
		if finished {
			break
		}
	}

	var p = 0
	for p = 0; p < roomCount; p++ {
		buildRoom(&a, roomMap[p].pos.x, roomMap[p].pos.y, roomMap[p].size.x, roomMap[p].size.y)
	}

	return a
}

func buildRoom(m *[MapWidth][MapHeight]Tile, mx int, my int, width int, height int) {
	var x, y int
	for x = mx; float64(x) < math.Min(float64(MapWidth), float64(mx+width)); x++ {
		for y = my; float64(y) < math.Min(float64(MapHeight), float64(my+height)); y++ {
			if x < 0 {
				x = 0
			}
			if y < 0 {
				y = 0
			}
			m[x][y] = lookupTile(0x02)
		}
	}
}

var pos = Vec2{0, 0}
var tileMap = [MapWidth][MapHeight]Tile{}

func isCollision(m [MapWidth][MapHeight]Tile, x int, y int) bool {
	return m[x][y].collides
}

func move(m [MapWidth][MapHeight]Tile, x int, y int) bool {
	var nx = pos.x + x
	var ny = pos.y + y
	if nx < 0 || ny < 0 || nx >= MapWidth || ny >= MapHeight {
		return false
	}
	if isCollision(m, nx, ny) {
		return false
	}
	pos.x = nx
	pos.y = ny
	return true
}

func run() {
	var title = "Iveria"
	cfg := pixelgl.WindowConfig{
		Title:  title,
		Bounds: pixel.R(0, 0, 1300, 768),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	face, err := loadTTF("FiraCode-Regular.ttf", 16)
	if err != nil {
		panic(err)
	}

	atlas := text.NewAtlas(face, text.ASCII)
	txt := text.New(pixel.V(50, 500), atlas)

	txt.Color = colornames.Lightgrey

	fps := time.Tick(time.Second / 120)
	r := time.Now().UnixMicro()
	rand.Seed(r)

	tileMap = generateMap(0)
	c_sepr := 0

	for !win.Closed() {

		txt.Clear()
		var i, j int
		for i = 0; i < MapHeight; i++ {
			for j = 0; j < MapWidth; j++ {
				// fmt.Printf("a[%d][%d] = %d\n", i, j, getIdentity()[i][j])
				var Tile = tileMap[j][i]
				if (j) == pos.x && (i) == pos.y {
					Tile = lookupTile(0x02)
				}
				txt.Color = Tile.displayColor
				_, err := txt.WriteString(string(Tile.display))
				if err != nil {
					fmt.Printf("Invalid write\n")
					return
				}
			}
			txt.WriteRune('\n')
		}

		if win.JustPressed(pixelgl.KeyEnter) || win.Repeated(pixelgl.KeyEnter) {
			txt.WriteRune('\n')
		}

		if win.JustPressed(pixelgl.KeyRight) || win.Repeated(pixelgl.KeyRight) {
			move(tileMap, 1, 0)
		}
		if win.JustPressed(pixelgl.KeyLeft) || win.Repeated(pixelgl.KeyLeft) {
			move(tileMap, -1, 0)
		}
		if win.JustPressed(pixelgl.KeyUp) || win.Repeated(pixelgl.KeyUp) {
			move(tileMap, 0, -1)
		}
		if win.JustPressed(pixelgl.KeyDown) || win.Repeated(pixelgl.KeyDown) {
			move(tileMap, 0, 1)
		}
		if win.JustPressed(pixelgl.KeySpace) || win.Repeated(pixelgl.KeySpace) {
			c_sepr = c_sepr + 1
			rand.Seed(r)
			tileMap = generateMap(c_sepr)
		}

		win.Clear(aapDarkGrey)
		txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))
		win.Update() // Swap buffers and poll events

		<-fps
	}
}

func main() {
	pixelgl.Run(run)
}
