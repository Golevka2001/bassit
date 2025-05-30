package util

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type HorizontalAlign int

const (
	AlignLeft HorizontalAlign = iota
	AlignCenter
	AlignRight
)

type VerticalAlign int

const (
	AlignTop VerticalAlign = iota
	AlignMiddle
	AlignBottom
)

// DrawTextLine draws a text line with a given alignment at the specified position.
// Parameters:
// - s: The screen to draw on
// - x1, x2: The horizontal boundary of the text line (inclusive)
// - y: The y-coordinate of the text line
// - text: The text to draw
// - align: AlignLeft|AlignCenter|AlignRight
// - style: The style of the text
func DrawTextLine(
	s *tcell.Screen,
	x1, x2 int,
	y int,
	text string,
	align HorizontalAlign,
	style tcell.Style,
) {
	if x1 > x2 || y < 0 || text == "" {
		return
	}

	runes := []rune(text)
	textVisualWidth := runewidth.StringWidth(string(runes))
	if textVisualWidth == 0 {
		return
	}

	// Calculate the ideal starting X position for the text
	idealStartX := x1
	availableWidth := x2 - x1 + 1
	switch align {
	case AlignCenter:
		idealStartX = x1 + (availableWidth-textVisualWidth)/2
	case AlignRight:
		idealStartX = x2 - textVisualWidth + 1
		// AlignLeft is the default (idealStartX = x1)
	}

	textCursor := idealStartX

	for _, r := range runes {
		runeW := runewidth.RuneWidth(r)
		// Skip zero-width runes (like combining characters)
		if runeW == 0 {
			continue
		}

		// Boundary check
		if textCursor+runeW-1 < x1 {
			textCursor += runeW
			continue
		}
		if textCursor > x2 {
			break
		}

		// Determine the actual screen column to draw this rune
		actualDrawX := max(textCursor, x1)

		// Check if the rune, starting at actualDrawX, fits completely within [x1, x2]
		if actualDrawX+runeW-1 <= x2 {
			(*s).SetContent(actualDrawX, y, r, nil, style)
		} else {
			break
		}

		textCursor += runeW
	}
}

// DrawTextArea draws a text area with a given alignment at the specified position.
// A new line should be added after each line of text.
// Parameters:
// - s: The screen to draw on
// - x1, x2: The horizontal boundary of the text area (inclusive)
// - y1, y2: The vertical boundary of the text area (inclusive)
// - text: The text to draw
// - horizontalAlign: AlignLeft|AlignCenter|AlignRight
// - verticalAlign: AlignTop|AlignMiddle|AlignBottom
// - style: The style of the text
func DrawTextArea(
	s *tcell.Screen,
	x1, x2 int,
	y1, y2 int,
	text string,
	horizontalAlign HorizontalAlign,
	verticalAlign VerticalAlign,
	style tcell.Style,
) {
	lines := strings.Split(text, "\n")
	numLines := len(lines)
	if numLines == 0 {
		return
	}

	startY := y1
	availableHeight := y2 - y1 + 1
	switch verticalAlign {
	case AlignMiddle:
		startY = y1 + (availableHeight-numLines)/2
	case AlignBottom:
		startY = y2 - numLines + 1
	}

	for i, line := range lines {
		currentY := startY + i
		if currentY > y2 || currentY < y1 { // Ensure drawing within y1 and y2
			continue
		}
		DrawTextLine(s, x1, x2, currentY, line, horizontalAlign, style)
	}
}

func DrawWelcome(s *tcell.Screen) {
	content := `
              ╭◜‾‾◝╮        ╭◜‾‾◝╮        ╭◜‾‾◝╮        ╭◜‾‾◝╮             
            ◞◜      ◝◟    ◞◜      ◝◟    ◞◜      ◝◟    ◞◜      ◝◟           
           (          )  (          )  (          )  (          )          
            ◝‒◟ ╭╮ ◞‒◜    ◝‒◟ ╭╮ ◞‒◜    ◝‒◟ ╭╮ ◞‒◜    ◝‒◟ ╭╮ ◞‒◜           
               ╰┤├╯          ╰┤├╯          ╰┤├╯          ╰┤├╯              
         ◜‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾--ˍ       
        /                                                           ‾-ˍ    
       /   ◜‾‾/◝         ◜‾|‾◝         ◜\‾‾◝         ◜‾|‾◝             \   
      /    ▏ / ▕         ▏ | ▕         ▏ \ ▕         ▏ | ▕              \  
     /     ◟/ˍˍ◞         ◟ˍ|ˍ◞         ◟ˍˍ\◞         ◟ˍ|ˍ◞              ▕  
  ˍˍ◞ˍˍ--‾‾‾     ˍˍˍ---‾‾‾     ˍˍˍ---‾‾‾     ˍˍˍ---‾‾‾                  ▕  
  \\    ˍˍˍ---‾‾‾     ˍˍˍ---‾‾‾    ˍ ˍˍ---‾‾‾                           /  
   \\‾‾‾      ˍˍˍ---‾‾‾    ˍˍˍ---‾‾‾      ◟ˍ                           /   
    \\ˍˍ---‾‾‾    ˍˍˍ---‾‾‾             ˍ ‾◟◞  ˍˍ--‾‾‾-ˍ              /    
     \\  ˍˍˍ---‾‾‾                  /‾   \  ˍ-‾         ‾ˍ         ˍˍ‾     
      \\‾           ˍˍˍ         /‾  ‾‾‾) ˍ-‾              ‾------‾‾        
       \\ˍˍ---ˍˍ    \  \    ◜\  ‾‾‾) ‾ˍ-‾                                  
                ⌝ˍ   \-◠-ˍ  (‾\  ‾‾ˍ-‾                                     
                  ⌝ˍ  \   \ ╰-╯‾ˍ-‾                                        
                    ⌝ˍ \ˍ◞╯  ˍ-‾                                           
                      \ˍˍ--‾‾                                              
`
	w, h := (*s).Size()
	(*s).Clear()
	DrawTextArea(s, 0, w-1, 0, h-1, content, AlignCenter, AlignMiddle, tcell.StyleDefault)
	(*s).Show()
}
