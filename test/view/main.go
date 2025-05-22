package main

import (
	"bassit/view"

	"github.com/gdamore/tcell/v2"
)

func testDrawTextLine() {
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	err = s.Init()
	if err != nil {
		panic(err)
	}
	defer s.Fini()

	s.Clear()

	style := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)

	// Left aligned
	view.DrawTextLine(&s, 5, 30, 0, "Hello, World!", view.AlignLeft, style)
	// Left aligned, text too long
	view.DrawTextLine(&s, 5, 30, 1, "Hello, World! This is a longer text.", view.AlignLeft, style)

	// Center aligned
	view.DrawTextLine(&s, 5, 30, 3, "Hello, World!", view.AlignCenter, style)
	// Center aligned, text too long
	view.DrawTextLine(&s, 5, 30, 4, "Hello, World! This is a longer text.", view.AlignCenter, style)

	// Right aligned
	view.DrawTextLine(&s, 5, 30, 6, "Hello, World!", view.AlignRight, style)
	// Right aligned, text too long
	view.DrawTextLine(&s, 5, 30, 7, "Hello, World! This is a longer text.", view.AlignRight, style)

	for i := 0; i <= 30; i++ {
		s.SetContent(i, 8, rune(i%10+'0'), nil, tcell.StyleDefault)
	}

	view.DrawTextLine(&s, 0, 100, 10, "Tap [Enter] to exit", view.AlignLeft, tcell.StyleDefault)

	// Wait for enter key to exit
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEnter {
				return
			}
		}
	}
}

func main() {
	testDrawTextLine()
}
