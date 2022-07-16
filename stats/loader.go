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
	err = yaml.Unmarshal(buf, &stats)
	hs.stats = stats

	return err
}
