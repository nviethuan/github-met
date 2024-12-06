package utils

import types "github-met/types"

func ChunkWeeks(weeks []types.ContributionWeek, chunkSize int) [][]types.ContributionWeek {
	var chunkedWeeks [][]types.ContributionWeek
	for i := 0; i < len(weeks); i += chunkSize {
		end := i + chunkSize
		if end > len(weeks) {
			end = len(weeks)
		}
		chunkedWeeks = append(chunkedWeeks, weeks[i:end])
	}
	return chunkedWeeks
}
