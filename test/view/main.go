package main

import (
	"bassit/util"

	"github.com/gdamore/tcell/v2"
)

func testDrawTextLine(s tcell.Screen) {
	s.Clear()

	style := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)

	// Left aligned
	util.DrawTextLine(&s, 5, 30, 0, "Hello, World!", util.AlignLeft, style)
	// Left aligned, text too long
	util.DrawTextLine(&s, 5, 30, 1, "Hello, World! This is a longer text.", util.AlignLeft, style)

	// Center aligned
	util.DrawTextLine(&s, 5, 30, 3, "Hello, World!", util.AlignCenter, style)
	// Center aligned, text too long
	util.DrawTextLine(&s, 5, 30, 4, "Hello, World! This is a longer text.", util.AlignCenter, style)

	// Right aligned
	util.DrawTextLine(&s, 5, 30, 6, "Hello, World!", util.AlignRight, style)
	// Right aligned, text too long
	util.DrawTextLine(&s, 5, 30, 7, "Hello, World! This is a longer text.", util.AlignRight, style)

	for i := 0; i <= 30; i++ {
		s.SetContent(i, 8, rune(i%10+'0'), nil, tcell.StyleDefault)
	}

	util.DrawTextLine(&s, 0, 100, 10, "Tap [Enter] to exit", util.AlignLeft, tcell.StyleDefault)
	s.Show()
	waitForEnter(s)
}

func testDrawTextArea(s tcell.Screen) {
	s.Clear()

	style := tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite)
	textAreaX1, textAreaY1 := 5, 2
	textAreaX2, textAreaY2 := 35, 8

	// Draw a border for the text area for visualization
	for x := textAreaX1; x <= textAreaX2; x++ {
		s.SetContent(x, textAreaY1, '-', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
		s.SetContent(x, textAreaY2, '-', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
	}
	for y := textAreaY1 + 1; y < textAreaY2; y++ {
		s.SetContent(textAreaX1, y, '|', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
		s.SetContent(textAreaX2, y, '|', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
	}

	sampleText := "This is a test string.\nIt has multiple lines.\nAnd some of them might be a bit longer than the available width."

	// Test Case 1: Left, Top
	util.DrawTextArea(&s, textAreaX1+1, textAreaX2-1, textAreaY1+1, textAreaY2-1, sampleText, util.AlignLeft, util.AlignTop, style)

	util.DrawTextLine(&s, 0, 100, textAreaY2+2, "Left/Top. Tap [Enter] for next", util.AlignLeft, tcell.StyleDefault)
	s.Show()
	waitForEnter(s)
	s.Clear()

	// Draw border again
	for x := textAreaX1; x <= textAreaX2; x++ {
		s.SetContent(x, textAreaY1, '-', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
		s.SetContent(x, textAreaY2, '-', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
	}
	for y := textAreaY1 + 1; y < textAreaY2; y++ {
		s.SetContent(textAreaX1, y, '|', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
		s.SetContent(textAreaX2, y, '|', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
	}
	// Test Case 2: Center, Middle
	util.DrawTextArea(&s, textAreaX1+1, textAreaX2-1, textAreaY1+1, textAreaY2-1, sampleText, util.AlignCenter, util.AlignMiddle, style)
	util.DrawTextLine(&s, 0, 100, textAreaY2+2, "Center/Middle. Tap [Enter] for next", util.AlignLeft, tcell.StyleDefault)
	s.Show()
	waitForEnter(s)
	s.Clear()

	// Draw border again
	for x := textAreaX1; x <= textAreaX2; x++ {
		s.SetContent(x, textAreaY1, '-', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
		s.SetContent(x, textAreaY2, '-', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
	}
	for y := textAreaY1 + 1; y < textAreaY2; y++ {
		s.SetContent(textAreaX1, y, '|', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
		s.SetContent(textAreaX2, y, '|', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
	}
	// Test Case 3: Right, Bottom
	util.DrawTextArea(&s, textAreaX1+1, textAreaX2-1, textAreaY1+1, textAreaY2-1, sampleText, util.AlignRight, util.AlignBottom, style)
	util.DrawTextLine(&s, 0, 100, textAreaY2+2, "Right/Bottom. Tap [Enter] to exit", util.AlignLeft, tcell.StyleDefault)
	s.Show()
	waitForEnter(s)
}

func waitForEnter(s tcell.Screen) {
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
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	err = s.Init()
	if err != nil {
		panic(err)
	}
	defer s.Fini()

	testDrawTextLine(s)
	testDrawTextArea(s)
	util.DrawWelcome(&s)
	waitForEnter(s)
}
