package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tcell "github.com/gdamore/tcell"
)

type coords struct {
	x int
	y int
}

const (
	UpDirection = iota
	DownDirection
	LeftDirection
	RightDirection
)

const snakeMaxLen = 1000
const gameSpeed = 2 // 1 to 10, 1= slow, 10=fast

var snakeBody [snakeMaxLen]coords
var headPos int
var tailPos int
var snakeLen int = 10
var head coords
var fruit coords

/*
Get an array, sufficiently big to hold the maximum snake. Establish two pointers, one for the head, one for the tail.
At the beginning, the tail would be in cell #1, the head in cell #3.
As the snake moves, move the head pointer to the right and write the new coordinate.
Then, if there's no food eaten, move the tail pointer to the right as well.
If either of the pointers tries to go beyond the rightmost end of the array, wrap them over to the beginning.
*/

var currentDirection int

func handleEvents(s tcell.Screen, quit chan (struct{})) {
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyEnter:
				close(quit)
			case tcell.KeyCtrlL:
				s.Sync()
			case tcell.KeyUp:
				currentDirection = UpDirection
			case tcell.KeyDown:
				currentDirection = DownDirection
			case tcell.KeyLeft:
				currentDirection = LeftDirection
			case tcell.KeyRight:
				currentDirection = RightDirection
			}
		case *tcell.EventResize:
			s.Sync()
		}
	}
}

func drawscreen(s tcell.Screen) {

	for {
		s.Clear()
		i := tailPos
		for i != headPos {
			s.SetCell(snakeBody[i].x, snakeBody[i].y, tcell.StyleDefault, '\U0001F311') // new moon
			i = (i + 1) % snakeMaxLen
		}
		s.SetCell(head.x, head.y, tcell.StyleDefault, '\U0001F315') // full moon unicode
		//draw fruit
		s.SetCell(fruit.x, fruit.y, tcell.StyleDefault, '\U0001F34E')
		s.Show()
	}

}

func update(s tcell.Screen) {
	for {
		time.Sleep(time.Duration(100/gameSpeed) * time.Millisecond)

		snakeBody[headPos] = head
		headPos = (headPos + 1) % snakeMaxLen
		tailPos = (tailPos + 1) % snakeMaxLen
		sizex, sizey := s.Size()
		switch currentDirection {
		case DownDirection:
			if head.y < sizey-1 {
				head.y++
			}
		case UpDirection:
			if head.y > 0 {
				head.y--
			}
		case LeftDirection:
			if head.x > 0 {
				head.x--
			}
		case RightDirection:
			if head.x < sizex-1 {
				head.x++
			}
		}
	}
}

func main() {

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	//s.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite))

	sizex, sizey := s.Size()
	head.x = sizex / 2
	head.y = sizey / 2
	currentDirection = RightDirection
	headPos = snakeLen
	tailPos = 0
	fruit = coords{rand.Int() % sizex, rand.Int() % sizey}
	quit := make(chan struct{})
	defer s.Fini()
	go handleEvents(s, quit)
	go update(s)
	go drawscreen(s)

	<-quit

}
