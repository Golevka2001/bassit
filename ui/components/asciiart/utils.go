package asciiart

import "github.com/charmbracelet/lipgloss/v2"

// cropLineToWidth will be called when the line is wider than the target width
func cropLineToWidth(line string, targetWidth int) string {
	var trimmed string
	var currentWidth int
	excess := lipgloss.Width(line) - targetWidth
	startOffset := excess / 2

	for _, r := range line {
		rw := lipgloss.Width(string(r))
		if currentWidth >= startOffset && lipgloss.Width(trimmed) < targetWidth {
			trimmed += string(r)
		}
		currentWidth += rw
		if lipgloss.Width(trimmed) >= targetWidth {
			break
		}
	}
	return trimmed
}
