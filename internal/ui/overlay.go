package ui

import (
	"strings"

	overlay "github.com/floatpane/bubble-overlay"
)

// SGR:
// SGR stands for Select Graphic Rendition. It is a parameter in ANSI/ECMA-48 escape sequences used to format terminal output text (controlling attributes like color, bold,
// underline, italic, blinking, or faint styling).

// ### What is an "SGR Reset"?
// An SGR reset sequence (specifically  \x1b[0m  or  \033[0m ) is the control code sent to the terminal to clear all applied text styles and colors, returning the terminal
// rendering behavior back to its default state.

// • Why it mattered here: The original  faintify  helper applied faint rendering ( \x1b[2m ) to the list background. However, any text in the list that reset its own styles using
// \x1b[0m  would automatically clear the faint formatting as well. To prevent this,  faintify  was designed to re-apply the faint sequence after every SGR reset sequence in the
// output string.

// Center renders the foreground string centered over the background string,
// automatically applying faint styling to the background and resetting styles on the foreground.
func Center(background, foreground string, width, height int) string {
	return overlay.Center(faintify(background), resetEachLine(foreground), width, height)
}

// faintify applies faint styling to all content in a string, preserving SGR resets
// but appending the faint style after each reset sequence.
func faintify(s string) string {
	const faintSeq = "\x1b[2m"
	const resetSeq = "\x1b[0m"
	return faintSeq + strings.ReplaceAll(s, resetSeq, resetSeq+faintSeq) + resetSeq
}

// resetEachLine forces a clean SGR reset at the start of every line, since
// overlay.Center's Line() resets SGR after pasting the popup but not before it,
// letting the faded background's trailing escape state bleed into the popup itself.
func resetEachLine(s string) string {
	const resetSeq = "\x1b[0m"
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = resetSeq + line
	}
	return strings.Join(lines, "\n")
}
