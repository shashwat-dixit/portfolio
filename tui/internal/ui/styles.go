package ui

import "charm.land/lipgloss/v2"

var (
	colorAccent   = lipgloss.Color("14")
	colorTitle    = lipgloss.Color("15")
	colorMuted    = lipgloss.Color("8")
	colorDim      = lipgloss.Color("245")
	colorActive   = lipgloss.Color("0")
	colorActiveBg = lipgloss.Color("14")
	colorBorder   = lipgloss.Color("240")
	colorError    = lipgloss.Color("9")
	colorChip     = lipgloss.Color("6")
)

func titleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(colorTitle).Bold(true)
}

func mutedStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(colorMuted)
}

func dimStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(colorDim)
}

func accentStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(colorAccent)
}

func errorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(colorError)
}

func navItemStyle(active bool) lipgloss.Style {
	s := lipgloss.NewStyle().Padding(0, 1)
	if active {
		return s.Foreground(colorActive).Background(colorActiveBg).Bold(true)
	}
	return s.Foreground(colorDim)
}

func chipStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(colorChip).Padding(0, 1).Border(lipgloss.NormalBorder()).BorderForeground(colorBorder)
}

func headerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(colorBorder).Padding(0, 1)
}

func footerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(colorMuted)
}
