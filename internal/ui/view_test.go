package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/kreuger97/hy2tool/internal/config"
)

func TestWelcomeView(t *testing.T) {
	lipgloss.SetHasDarkBackground(true)

	m := model{
		state:   stateWelcome,
		cfg:     config.Default(),
		width:   80,
		height:  24,
		ready:   true,
		taskIdx: -1,
	}

	view := m.View()
	if view == "" {
		t.Fatal("welcome view is empty")
	}
	if !strings.Contains(view, "Hysteria 2") {
		t.Error("welcome view missing title")
	}
	if !strings.Contains(view, "Enter") {
		t.Error("welcome view missing prompt")
	}
	t.Logf("welcome view:\n%s", view)
}

func TestProcessingView(t *testing.T) {
	lipgloss.SetHasDarkBackground(true)

	m := model{
		state:   stateProcessing,
		cfg:     config.Default(),
		width:   80,
		height:  24,
		ready:   true,
		taskIdx: 0,
		tasks: []task{
			{name: "Install Hysteria 2", status: taskRunning},
			{name: "Generate SSL Certificate", status: taskPending},
			{name: "Write Server Config", status: taskPending},
		},
	}

	view := m.View()
	if view == "" {
		t.Fatal("processing view is empty")
	}
	t.Logf("processing view:\n%s", view)
}

func TestSummaryView(t *testing.T) {
	lipgloss.SetHasDarkBackground(true)

	m := model{
		state:    stateSummary,
		cfg:      config.Config{Port: "443", Password: "test123"},
		width:    80,
		height:   24,
		ready:    true,
		serverIP: "1.2.3.4",
	}

	view := m.View()
	if view == "" {
		t.Fatal("summary view is empty")
	}
	if !strings.Contains(view, "1.2.3.4") {
		t.Error("summary view missing server IP")
	}
	if !strings.Contains(view, "test123") {
		t.Error("summary view missing password")
	}
	if !strings.Contains(view, "Installation Complete") {
		t.Error("summary view missing title")
	}
	t.Logf("summary view:\n%s", view)
}
