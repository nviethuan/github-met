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
	content.WriteString(`<style>
		@keyframes cardReveal {
			from { opacity: 0; transform: translateY(8px); }
			to { opacity: 1; transform: translateY(0); }
		}
		@keyframes titleReveal {
			0% { opacity: 0; transform: translateY(5px); }
			100% { opacity: 1; transform: translateY(0); }
		}
		@keyframes barGrow {
			from { transform: scaleX(0); }
			to { transform: scaleX(1); }
		}
		@media (prefers-reduced-motion: reduce) {
			*, *::before, *::after {
				animation-duration: 0.01ms !important;
				animation-delay: 0s !important;
				animation-iteration-count: 1 !important;
				transform: none !important;
			}
		}
	</style>`)
	content.WriteString(`<title id="title">Programming languages used by ` + html.EscapeString(username) + `</title>`)
	content.WriteString(`<desc id="description">Programming languages ranked by code bytes in owned, non-fork repositories.</desc>`)
	content.WriteString(`<rect width="495" height="` + strconv.Itoa(height) + `" rx="8" fill="` + background + `" stroke="` + borderColor + `" style="opacity: 0; animation: cardReveal 0.55s ease-out forwards"/>`)
	content.WriteString(`<text x="24" y="34" fill="` + textColor + `" font-family="Segoe UI, Ubuntu, sans-serif" font-size="18" font-weight="700" style="opacity: 0; animation: titleReveal 0.5s ease-out forwards 0.12s">Top Languages</text>`)
	content.WriteString(`<text x="24" y="57" fill="` + mutedColor + `" font-family="Segoe UI, Ubuntu, sans-serif" font-size="12" style="opacity: 0; animation: titleReveal 0.5s ease-out forwards 0.22s">@` + html.EscapeString(username) + `</text>`)
	content.WriteString(`<text x="471" y="35" text-anchor="end" fill="` + mutedColor + `" font-family="Segoe UI, Ubuntu, sans-serif" font-size="11" style="opacity: 0; animation: titleReveal 0.5s ease-out forwards 0.28s">RANKED BY CODE SIZE</text>`)

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
			delay := fmt.Sprintf("%.2fs", 0.32+float64(i)*0.1)
			content.WriteString(`<g font-family="Segoe UI, Ubuntu, sans-serif" style="opacity: 0; animation: titleReveal 0.45s ease-out forwards ` + delay + `">`)
			content.WriteString(`<text x="24" y="` + strconv.Itoa(y+16) + `" fill="` + mutedColor + `" font-size="11">` + strconv.Itoa(i+1) + `</text>`)
			content.WriteString(`<circle cx="48" cy="` + strconv.Itoa(y+12) + `" r="5" fill="` + color + `"/>`)
			content.WriteString(`<text x="62" y="` + strconv.Itoa(y+16) + `" fill="` + textColor + `" font-size="13" font-weight="600">` + label + `</text>`)
			content.WriteString(`<text x="471" y="` + strconv.Itoa(y+16) + `" text-anchor="end" fill="` + mutedColor + `" font-size="11">` + fmt.Sprintf("%.1f%%", percentage) + `</text>`)

			barWidth := percentage * 2.65
			if barWidth > 0 {
				content.WriteString(`<rect x="206" y="` + strconv.Itoa(y+23) + `" width="265" height="4" rx="2" fill="` + trackColor + `"/>`)
				content.WriteString(`<rect x="206" y="` + strconv.Itoa(y+23) + `" width="` + fmt.Sprintf("%.2f", barWidth) + `" height="4" rx="2" fill="` + color + `" style="transform-box: fill-box; transform-origin: left center; animation: barGrow 0.8s cubic-bezier(0.2, 0.8, 0.2, 1) both ` + delay + `"/>`)
			}
			content.WriteString(`</g>`)
		}
	}

	content.WriteString(`<text x="24" y="` + strconv.Itoa(height-12) + `" fill="` + mutedColor + `" font-family="Segoe UI, Ubuntu, sans-serif" font-size="10" style="opacity: 0; animation: titleReveal 0.45s ease-out forwards 0.7s">Based on code size in owned, non-fork repositories</text>`)
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
