package main

import (
	U "github.com/Golevka2001/bassit/util"

	"github.com/gdamore/tcell/v2"
)

func testDrawTextLine(s tcell.Screen) {
	s.Clear()

	style := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)

	// Left aligned
	U.DrawTextLine(&s, 5, 30, 0, "Hello, World!", U.AlignLeft, style)
	// Left aligned, text too long
	U.DrawTextLine(&s, 5, 30, 1, "Hello, World! This is a longer text.", U.AlignLeft, style)

	// Center aligned
	U.DrawTextLine(&s, 5, 30, 3, "Hello, World!", U.AlignCenter, style)
	// Center aligned, text too long
	U.DrawTextLine(&s, 5, 30, 4, "Hello, World! This is a longer text.", U.AlignCenter, style)

	// Right aligned
	U.DrawTextLine(&s, 5, 30, 6, "Hello, World!", U.AlignRight, style)
	// Right aligned, text too long
	U.DrawTextLine(&s, 5, 30, 7, "Hello, World! This is a longer text.", U.AlignRight, style)

	for i := 0; i <= 30; i++ {
		s.SetContent(i, 8, rune(i%10+'0'), nil, tcell.StyleDefault)
	}

	U.DrawTextLine(&s, 0, 100, 10, "Tap [Enter] to exit", U.AlignLeft, tcell.StyleDefault)
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
	U.DrawTextArea(&s, textAreaX1+1, textAreaX2-1, textAreaY1+1, textAreaY2-1, sampleText, U.AlignLeft, U.AlignTop, style)

	U.DrawTextLine(&s, 0, 100, textAreaY2+2, "Left/Top. Tap [Enter] for next", U.AlignLeft, tcell.StyleDefault)
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
	U.DrawTextArea(&s, textAreaX1+1, textAreaX2-1, textAreaY1+1, textAreaY2-1, sampleText, U.AlignCenter, U.AlignMiddle, style)
	U.DrawTextLine(&s, 0, 100, textAreaY2+2, "Center/Middle. Tap [Enter] for next", U.AlignLeft, tcell.StyleDefault)
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
	U.DrawTextArea(&s, textAreaX1+1, textAreaX2-1, textAreaY1+1, textAreaY2-1, sampleText, U.AlignRight, U.AlignBottom, style)
	U.DrawTextLine(&s, 0, 100, textAreaY2+2, "Right/Bottom. Tap [Enter] to exit", U.AlignLeft, tcell.StyleDefault)
	s.Show()
	waitForEnter(s)
}

func testFillArea(s tcell.Screen) {
	s.Clear()

	// Test Case 1: Fill area with character and style
	fillAreaX1, fillAreaY1 := 5, 2
	fillAreaX2, fillAreaY2 := 25, 6

	// Draw a border for visualization
	for x := fillAreaX1 - 1; x <= fillAreaX2+1; x++ {
		s.SetContent(x, fillAreaY1-1, '-', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
		s.SetContent(x, fillAreaY2+1, '-', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
	}
	for y := fillAreaY1; y <= fillAreaY2; y++ {
		s.SetContent(fillAreaX1-1, y, '|', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
		s.SetContent(fillAreaX2+1, y, '|', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
	}

	// Fill area with character and style (bgOnly = false)
	fillStyle := tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorYellow)
	U.FillArea(&s, fillAreaX1, fillAreaX2, fillAreaY1, fillAreaY2, false, '█', fillStyle)

	U.DrawTextLine(&s, 0, 100, fillAreaY2+3, "Filled with '█' character and blue background. Tap [Enter] for background only", U.AlignLeft, tcell.StyleDefault)
	s.Show()
	waitForEnter(s)

	// Clear the screen for the next test
	s.Clear()
	for x := fillAreaX1 - 1; x <= fillAreaX2+1; x++ {
		s.SetContent(x, fillAreaY1-1, '-', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
		s.SetContent(x, fillAreaY2+1, '-', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
	}
	for y := fillAreaY1; y <= fillAreaY2; y++ {
		s.SetContent(fillAreaX1-1, y, '|', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
		s.SetContent(fillAreaX2+1, y, '|', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
	}

	// Fill area with background only (bgOnly = true)
	fillStyle = tcell.StyleDefault.Background(tcell.ColorGreen)
	U.FillArea(&s, fillAreaX1, fillAreaX2, fillAreaY1, fillAreaY2, true, ' ', fillStyle)

	U.DrawTextLine(&s, 0, 100, fillAreaY2+3, "Filled with green background only. Tap [Enter] to exit", U.AlignLeft, tcell.StyleDefault)
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
	testFillArea(s)
	U.DrawWelcome(&s)
	waitForEnter(s)
}
