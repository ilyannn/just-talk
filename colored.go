package main

import (
	"github.com/charmbracelet/lipgloss"
	_log "github.com/charmbracelet/log"
	"os"
)

var _BoldStyle = lipgloss.NewStyle().Bold(true)
var JustBold = _BoldStyle.Render("just")

var PromptStyle = lipgloss.NewStyle().Bold(true)
var LLMOutputStyle = lipgloss.NewStyle().Italic(true)

func _createLogger() *_log.Logger {
	levelString := os.Getenv("JUST_TALK_LOG_LEVEL")

	if levelString == "" {
		levelString = "info"
	}

	level, levelParseErr := _log.ParseLevel(levelString)

	if levelParseErr != nil {
		level = _log.InfoLevel
	}

	// Create a new logger
	logger := _log.NewWithOptions(os.Stderr, _log.Options{
		ReportTimestamp: true,
		Level:           level,
	})

	// Define the styles
	var styles = _log.DefaultStyles()
	styles.Levels[_log.InfoLevel] = lipgloss.NewStyle().
		SetString("INFO").
		Bold(true).
		Foreground(lipgloss.Color("42"))
	styles.Values["stderr"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	styles.Values["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	styles.Values["command"] = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
	styles.Timestamp = lipgloss.NewStyle().Faint(true)

	// Set the styles on the default logger
	logger.SetStyles(styles)

	if levelParseErr != nil {
		logger.With("err", levelParseErr).Errorf("Failed to parse log level")
	}
	return logger
}

var log = _createLogger()
