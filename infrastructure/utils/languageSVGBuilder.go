package utils

import (
	"fmt"
	"html"
	"strconv"
	"strings"
	"unicode/utf8"

	"github-met/types"
)

func LanguagesSVG(username string, languages []types.LanguageData, theme string) string {
	const headerHeight = 82
	const rowHeight = 38
	const footerHeight = 34

	background := "#FFFEFE"
	textColor := "#151515"
	mutedColor := "#464646"
	borderColor := "#E4E2E2"
	trackColor := "#F1F3F5"
	if theme == "dark" {
		background = "#151515"
		textColor = "#FFFEFE"
		mutedColor = "#E4E2E2"
		borderColor = "#464646"
		trackColor = "#292929"
	}

	height := headerHeight + len(languages)*rowHeight + footerHeight
	var content strings.Builder
	content.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 495 ` + strconv.Itoa(height) + `" width="495" height="` + strconv.Itoa(height) + `" role="img" aria-labelledby="title description">`)
	content.WriteString(`<title id="title">Programming languages used by ` + html.EscapeString(username) + `</title>`)
	content.WriteString(`<desc id="description">Programming languages ranked by code bytes in owned, non-fork repositories.</desc>`)
	content.WriteString(`<rect width="495" height="` + strconv.Itoa(height) + `" rx="8" fill="` + background + `" stroke="` + borderColor + `"/>`)
	content.WriteString(`<text x="24" y="34" fill="` + textColor + `" font-family="Segoe UI, Ubuntu, sans-serif" font-size="18" font-weight="700">Top Languages</text>`)
	content.WriteString(`<text x="24" y="57" fill="` + mutedColor + `" font-family="Segoe UI, Ubuntu, sans-serif" font-size="12">@` + html.EscapeString(username) + `</text>`)
	content.WriteString(`<text x="471" y="35" text-anchor="end" fill="` + mutedColor + `" font-family="Segoe UI, Ubuntu, sans-serif" font-size="11">RANKED BY CODE SIZE</text>`)

	var totalBytes int64
	for _, language := range languages {
		if language.Bytes > 0 {
			totalBytes += language.Bytes
		}
	}

	if len(languages) == 0 {
		content.WriteString(`<text x="24" y="108" fill="` + mutedColor + `" font-family="Segoe UI, Ubuntu, sans-serif" font-size="13">No language data found.</text>`)
	} else {
		for i, language := range languages {
			y := headerHeight + i*rowHeight
			percentage := 0.0
			if totalBytes > 0 {
				percentage = float64(language.Bytes) * 100 / float64(totalBytes)
			}

			color := language.Color
			if !isHexColor(color) {
				color = "#8b949e"
			}
			label := html.EscapeString(truncateLanguageName(language.Name, 24))
			content.WriteString(`<g font-family="Segoe UI, Ubuntu, sans-serif">`)
			content.WriteString(`<text x="24" y="` + strconv.Itoa(y+16) + `" fill="` + mutedColor + `" font-size="11">` + strconv.Itoa(i+1) + `</text>`)
			content.WriteString(`<circle cx="48" cy="` + strconv.Itoa(y+12) + `" r="5" fill="` + color + `"/>`)
			content.WriteString(`<text x="62" y="` + strconv.Itoa(y+16) + `" fill="` + textColor + `" font-size="13" font-weight="600">` + label + `</text>`)
			content.WriteString(`<text x="471" y="` + strconv.Itoa(y+16) + `" text-anchor="end" fill="` + mutedColor + `" font-size="11">` + fmt.Sprintf("%.1f%%", percentage) + `</text>`)

			barWidth := percentage * 2.65
			if barWidth > 0 {
				content.WriteString(`<rect x="206" y="` + strconv.Itoa(y+23) + `" width="265" height="4" rx="2" fill="` + trackColor + `"/>`)
				content.WriteString(`<rect x="206" y="` + strconv.Itoa(y+23) + `" width="` + fmt.Sprintf("%.2f", barWidth) + `" height="4" rx="2" fill="` + color + `"/>`)
			}
			content.WriteString(`</g>`)
		}
	}

	content.WriteString(`<text x="24" y="` + strconv.Itoa(height-12) + `" fill="` + mutedColor + `" font-family="Segoe UI, Ubuntu, sans-serif" font-size="10">Based on code size in owned, non-fork repositories</text>`)
	content.WriteString(`</svg>`)
	return content.String()
}

func truncateLanguageName(name string, limit int) string {
	if utf8.RuneCountInString(name) <= limit {
		return name
	}

	runes := []rune(name)
	return string(runes[:limit-1]) + "…"
}

func isHexColor(color string) bool {
	if len(color) != 7 || color[0] != '#' {
		return false
	}
	for _, char := range color[1:] {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
			return false
		}
	}
	return true
}
