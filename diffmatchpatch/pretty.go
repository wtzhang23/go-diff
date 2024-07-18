package diffmatchpatch

import (
	"bytes"
	"html"
	"strings"
)

// DiffPrettyHtml converts a []Diff into a pretty HTML report.
// It is intended as an example from which to write one's own display functions.
func (dmp *DiffMatchPatch) DiffPrettyHtml(diffs []Diff) string {
	var buff bytes.Buffer
	for _, diff := range diffs {
		text := strings.Replace(html.EscapeString(diff.Text), "\n", "&para;<br>", -1)
		switch diff.Type {
		case DiffInsert:
			_, _ = buff.WriteString("<ins style=\"background:#e6ffe6;\">")
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("</ins>")
		case DiffDelete:
			_, _ = buff.WriteString("<del style=\"background:#ffe6e6;\">")
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("</del>")
		case DiffEqual:
			_, _ = buff.WriteString("<span>")
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("</span>")
		}
	}
	return buff.String()
}

// DiffPrettyText converts a []Diff into a colored text report.
func (dmp *DiffMatchPatch) DiffPrettyText(diffs []Diff) string {
	var buff bytes.Buffer
	for _, diff := range diffs {
		text := diff.Text

		switch diff.Type {
		case DiffInsert:
			_, _ = buff.WriteString("\x1b[32m")
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("\x1b[0m")
		case DiffDelete:
			_, _ = buff.WriteString("\x1b[31m")
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("\x1b[0m")
		case DiffEqual:
			_, _ = buff.WriteString(text)
		}
	}

	return buff.String()
}

type LinesPrettyConfig struct {
	// Color toggles adding terminal color codes to the output.
	Color bool
	// Spacing is the space between the diff marker (e.g. "+", "-", " ") and the text.
	Spacing string
	// Context is the number of lines of context to show around the diff.
	Context int
}

// DiffLinesPrettyText converts two strings into a pretty text report of the diffs by line.
func (dmp *DiffMatchPatch) DiffLinesPrettyText(config LinesPrettyConfig, from, to string) string {
	linesDiff := dmp.DiffLines(from, to)
	var buff bytes.Buffer
	firstLine := true
	pushLine := func(marker rune, line string) {
		if !firstLine {
			buff.WriteRune('\n')
		}
		firstLine = false
		buff.WriteRune(marker)
		buff.WriteString(config.Spacing)
		buff.WriteString(line)
	}
	for diffIdx, diff := range linesDiff {
		switch diff.Type {
		case DiffInsert:
			if config.Color {
				_, _ = buff.WriteString("\x1b[32m")
			}
			split := strings.Split(diff.Text, "\n")
			for li, line := range split {
				if li == len(split)-1 && diffIdx < len(linesDiff)-1 {
					continue
				}
				pushLine('+', line)
			}
			if config.Color {
				_, _ = buff.WriteString("\x1b[0m")
			}
		case DiffDelete:
			if config.Color {
				_, _ = buff.WriteString("\x1b[31m")
			}
			split := strings.Split(diff.Text, "\n")
			for li, line := range split {
				if li == len(split)-1 && diffIdx < len(linesDiff)-1 {
					continue
				}
				pushLine('-', line)
			}
			if config.Color {
				_, _ = buff.WriteString("\x1b[0m")
			}
		case DiffEqual:
			split := strings.Split(diff.Text, "\n")
			for li, line := range split {
				if li == len(split)-1 && diffIdx < len(linesDiff)-1 {
					continue
				}

				// add context lines before the diff
				if diffIdx > 0 && li < config.Context {
					pushLine(' ', line)
					continue
				}
				// add context lines after the diff
				if diffIdx < len(linesDiff)-1 && li >= len(split)-config.Context {
					pushLine(' ', line)
					continue
				}
			}
		}
	}
	return buff.String()
}
