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
}

// DiffLinesPrettyText converts two strings into a pretty text report of the diffs by line.
func (dmp *DiffMatchPatch) DiffLinesPrettyText(config LinesPrettyConfig, text1, text2 string) string {
	fromRunes, toRunes, linesMap := dmp.DiffLinesToRunes(text1, text2)
	runeDiffs := dmp.DiffMainRunes(fromRunes, toRunes, false)
	diffs := dmp.DiffCharsToLines(runeDiffs, linesMap)
	patches := dmp.PatchMake(diffs)

	var buff bytes.Buffer
	writeBlock := func(marker rune, block string, color string, lastDiff bool) {
		if config.Color && color != "" {
			_, _ = buff.WriteString(color)
		}

		startIdx := 0
		for startIdx < len(block) {
			endIdx := strings.Index(block[startIdx:], "\n")
			if endIdx == -1 {
				buff.WriteRune(marker)
				buff.WriteString(config.Spacing)
				buff.WriteString(block[startIdx:])
				if !lastDiff {
					buff.WriteRune('\n')
				}
				break
			}
			buff.WriteRune(marker)
			buff.WriteString(config.Spacing)
			buff.WriteString(block[startIdx : startIdx+endIdx+1])
			startIdx += endIdx + 1
		}
		if config.Color && color != "" {
			_, _ = buff.WriteString("\x1b[0m")
		}
	}
	for _, patch := range patches {
		patch.addCoordsToBuffer(&buff)
		buff.WriteRune('\n')
		for di, diff := range patch.diffs {
			lastDiff := di == len(patch.diffs)-1
			switch diff.Type {
			case DiffInsert:
				writeBlock('+', diff.Text, "\x1b[32m", lastDiff)
			case DiffDelete:
				writeBlock('-', diff.Text, "\x1b[31m", lastDiff)
			case DiffEqual:
				writeBlock(' ', diff.Text, "", lastDiff)
			}
		}
	}
	return buff.String()
}
