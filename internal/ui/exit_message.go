package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"bookmark/internal/adapters/icon"
)

// ExitMessage renders a standard framed exit message.
func ExitMessage(theme Theme, message string, mutedText bool) string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Text)).
		Margin(0, 2, 0, 2)

	if mutedText {
		return style.Render(message)
	}

	return style.Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Muted)).
		Foreground(lipgloss.Color(theme.Secondary)).
		Bold(true).
		Margin(1, 1).
		Padding(1, 2).Render(message)
}

// SuccessMessage renders a confirmation message for bookmark add/edit/delete.
// action should be "created", "updated", or "deleted".
// body is the pre-formatted content to display below the title.
// inline renders title and body on one line without a border.
func SuccessMessage(theme Theme, action, body string, inline bool) string {
	if theme.PlainText {
		return fmt.Sprintf("Bookmark %s: %s", action, body)
	}

	title := fmt.Sprintf("%s Bookmark %s:", icon.Success, action)

	titleStyled := lipgloss.NewStyle().Foreground(theme.Success).Bold(true).Render(title)
	bodyStyled := lipgloss.NewStyle().Foreground(theme.Text).Render(body)

	sep := "\n"
	if inline {
		sep = " "
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border).
		Padding(1, 1).
		Margin(1, 1).
		Render(titleStyled + sep + bodyStyled)
}

// ErrorMessage renders a styled error message with failure icon.
func ErrorMessage(theme Theme, msg string) string {
	if theme.PlainText {
		return fmt.Sprintf("Error: %s", msg)
	}

	title := fmt.Sprintf("%s Error:", icon.Failure)

	content := lipgloss.NewStyle().Foreground(theme.Error).Bold(true).Render(title) +
		"\n" +
		lipgloss.NewStyle().Foreground(theme.Text).Render(msg)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Error).
		Padding(1, 1).
		Margin(1, 1).
		Render(content)
}
