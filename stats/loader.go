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
	hs.stats = stats

	return nil
}
