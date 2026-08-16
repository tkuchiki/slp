package stats

import (
	"io"

	"gopkg.in/yaml.v2"
)

func (hs *QueryStats) LoadStats(r io.Reader) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	var stats []*QueryStat
	if err := yaml.Unmarshal(buf, &stats); err != nil {
		return err
	}
	for _, stat := range stats {
		restoreLegacySampleCounts(stat)
	}

	hs.stats = stats

	return nil
}

func restoreLegacySampleCounts(stat *QueryStat) {
	restoreLegacyTimeSampleCount(stat.QueryTime, stat.Cnt)
	restoreLegacyTimeSampleCount(stat.LockTime, stat.Cnt)
	restoreLegacyNumberSampleCount(stat.RowsSent, stat.Cnt)
	restoreLegacyNumberSampleCount(stat.RowsExamined, stat.Cnt)
	restoreLegacyNumberSampleCount(stat.RowsAffected, stat.Cnt)
	restoreLegacyNumberSampleCount(stat.BytesSent, stat.Cnt)
}

func restoreLegacyTimeSampleCount(stats *timeStats, fallback int) {
	if stats == nil || stats.SampleCount != nil {
		return
	}

	stats.SampleCount = &fallback
}

func restoreLegacyNumberSampleCount(stats *numberStats, fallback int) {
	if stats == nil || stats.SampleCount != nil {
		return
	}

	stats.SampleCount = &fallback
}
