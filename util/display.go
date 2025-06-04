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

// FillArea fills a rectangular area on the screen with a specified rune and style.
// Parameters:
// - s: The screen to draw on
// - x1, x2: The horizontal boundary of the area (inclusive)
// - y1, y2: The vertical boundary of the area (inclusive)
// - bgOnly: If true, only the background color will be changed; if false, both character and style will be set
// - r: The rune to fill the area with (if bgOnly is false)
// - style: The style to apply to the area
// Note: If bgOnly is true, the rune parameter is ignored.
func FillArea(
	s *tcell.Screen,
	x1, x2 int,
	y1, y2 int,
	bgOnly bool,
	r rune,
	style tcell.Style,
) {
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			origR, _, origStyle, _ := (*s).GetContent(x, y)
			if bgOnly {
				// Only change the background color
				_, newBgColor, _ := style.Decompose()
				newStyle := origStyle.Background(newBgColor)
				(*s).SetContent(x, y, origR, nil, newStyle)
			} else {
				// Change both the character and the style
				(*s).SetContent(x, y, r, nil, style)
			}
		}
	}
}

func DrawWelcome(s *tcell.Screen) {
	content := `
              ╭◜‾‾◝╮        ╭◜‾‾◝╮        ╭◜‾‾◝╮        ╭◜‾‾◝╮             
            ◞◜      ◝◟    ◞◜      ◝◟    ◞◜      ◝◟    ◞◜      ◝◟           
           (          )  (          )  (          )  (          )          
            ◝‒◟ ╭╮ ◞‒◜    ◝‒◟ ╭╮ ◞‒◜    ◝‒◟ ╭╮ ◞‒◜    ◝‒◟ ╭╮ ◞‒◜           
               ╰┤├╯          ╰┤├╯          ╰┤├╯          ╰┤├╯              
         ◜‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾--_       
        /                                                           ‾-_    
       /   ◜‾‾╭◝         ◜‾|‾◝         ◜╮‾‾◝         ◜‾|‾◝             \   
      /    ▏ / ▕         ▏ | ▕         ▏ \ ▕         ▏ | ▕              \  
     /     ◟╯__◞         ◟_|_◞         ◟__╰◞         ◟_|_◞              ▕  
  __◞__--‾‾‾     ___---‾‾‾     ___---‾‾‾     ___---‾‾‾                  ▕  
  \\    ___---‾‾‾     /‾‾‾‾\‾‾‾     ___---‾‾‾                           /  
   \\‾‾‾     ___---‾‾‾▏(◜◞) ▏_---‾‾‾      ◟_                           /   
    \\_---‾‾‾     ___-\____/            _ ‾\◞  __--‾‾-_               /    
     \\  ___---‾‾‾                  (‾   \  _-‾        ‾◟_         _◞‾     
      \\‾                       (‾   ‾‾) _-‾              ‾------‾‾        
       \\__---__    ╭‾‾\   ◞-╮   ‾‾) ‾‾-‾                                  
                ◝_   \_◞-_  (‾\  ‾‾_-‾                                     
                  ◝_  \   \ ╰-╯‾_-‾                                        
                    ◝_ \_◞╯  _-‾                                           
                      \__--‾‾                                              
`
	w, h := (*s).Size()
	(*s).Clear()
	DrawTextArea(s, 0, w-1, 0, h-1, content, AlignCenter, AlignMiddle, tcell.StyleDefault)
	(*s).Show()
}
