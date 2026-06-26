package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"bookmark/internal/adapters/icon"
)

// ExitMessage renders a cancel/exit message with error styling when not plain text.
func ExitMessage(theme Theme, message string) string {
	if theme.PlainText {
		return message
	}

	title := fmt.Sprintf("%s %s", icon.Failure, message)
	titleStyled := lipgloss.NewStyle().Foreground(theme.Error).Bold(true).Render(title)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Error).
		Padding(1, 1).
		Margin(1, 1).
		Render(titleStyled)
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
