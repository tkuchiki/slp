package stats

import (
	"fmt"
	"sort"
)

const (
	SortCount           = "Count"
	SortQuery           = "Query"
	SortMaxQueryTime    = "MaxQueryTime"
	SortMinQueryTime    = "MinQueryTime"
	SortAvgQueryTime    = "AvgQueryTime"
	SortSumQueryTime    = "SumQueryTime"
	SortPNQueryTime     = "PNQueryTime"
	SortMaxLockTime     = "MaxLockTime"
	SortMinLockTime     = "MinLockTime"
	SortAvgLockTime     = "AvgLockTime"
	SortSumLockTime     = "SumLockTime"
	SortPNLockTime      = "PNLockTime"
	SortMaxRowsSent     = "MaxRowsSent"
	SortMinRowsSent     = "MinRowsSent"
	SortAvgRowsSent     = "AvgRowsSent"
	SortSumRowsSent     = "SumRowsSent"
	SortPNRowsSent      = "PNRowsSent"
	SortMaxRowsExamined = "MaxRowsExamined"
	SortMinRowsExamined = "MinRowsExamined"
	SortAvgRowsExamined = "AvgRowsExamined"
	SortSumRowsExamined = "SumRowsExamined"
	SortPNRowsExamined  = "PNRowsExamined"
	SortMaxRowsAffected = "MaxRowsAffected"
	SortMinRowsAffected = "MinRowsAffected"
	SortAvgRowsAffected = "AvgRowsAffected"
	SortSumRowsAffected = "SumRowsAffected"
	SortPNRowsAffected  = "PNRowsAffected"
	SortMaxBytesSent    = "MaxBytesSent"
	SortMinBytesSent    = "MinBytesSent"
	SortAvgBytesSent    = "AvgBytesSent"
	SortSumBytesSent    = "SumBytesSent"
	SortPNBytesSent     = "PNBytesSent"
)

type SortOptions struct {
	options    map[string]string
	sortType   string
	percentile int
}

func NewSortOptions() *SortOptions {
	options := map[string]string{
		"count":             SortCount,
		"query":             SortQuery,
		"max-query-time":    SortMaxQueryTime,
		"min-query-time":    SortMinQueryTime,
		"avg-query-time":    SortAvgQueryTime,
		"sum-query-time":    SortSumQueryTime,
		"pn-query-time":     SortPNQueryTime,
		"max-lock-time":     SortMaxLockTime,
		"min-lock-time":     SortMinLockTime,
		"avg-lock-time":     SortAvgLockTime,
		"sum-lock-time":     SortSumLockTime,
		"pn-lock-time":      SortPNLockTime,
		"max-rows-sent":     SortMaxRowsSent,
		"min-rows-sent":     SortMinRowsSent,
		"avg-rows-sent":     SortAvgRowsSent,
		"sum-rows-sent":     SortSumRowsSent,
		"pn-rows-sent":      SortPNRowsSent,
		"max-rows-examined": SortMaxRowsExamined,
		"min-rows-examined": SortMinRowsExamined,
		"avg-rows-examined": SortAvgRowsExamined,
		"sum-rows-examined": SortSumRowsExamined,
		"pn-rows-examined":  SortPNRowsExamined,
		"max-rows-affected": SortMaxRowsAffected,
		"min-rows-affected": SortMinRowsAffected,
		"avg-rows-affected": SortAvgRowsAffected,
		"sum-rows-affected": SortSumRowsAffected,
		"pn-rows-affected":  SortPNRowsAffected,
		"max-bytes-sent":    SortMaxBytesSent,
		"min-bytes-sent":    SortMinBytesSent,
		"avg-bytes-sent":    SortAvgBytesSent,
		"sum-bytes-sent":    SortSumBytesSent,
		"pn-bytes-sent":     SortPNBytesSent,
	}

	return &SortOptions{
		options: options,
	}
}

func (so *SortOptions) SetAndValidate(opt string) error {
	_, ok := so.options[opt]
	if ok {
		so.sortType = so.options[opt]
		return nil
	}

	var n int
	_, err := fmt.Sscanf(opt, "p%d", &n)
	if err != nil {
		return err
	}

	if n < 0 && n > 100 {
		return fmt.Errorf("enum value must be one of count,query,max-(query-time|lock-time|rows-sent|rows-examined|rows-affected|bytes-sent),min-(query-time|lock-time|rows-sent|rows-examined|rows-affected|bytes-sent),avg-(query-time|lock-time|rows-sent|rows-examined|rows-affected|bytes-sent),sum-(query-time|lock-time|rows-sent|rows-examined|rows-affected|bytes-sent),pN(N = 0 ~ 100)-(query-time|lock-time|rows-sent|rows-examined|rows-affected|bytes-sent),stddev, got '%s'", opt)
	}

	so.sortType = so.options["pn"]
	so.percentile = n

	return nil
}

func (so *SortOptions) SortType() string {
	return so.sortType
}

func (so *SortOptions) Percentile() int {
	return so.percentile
}

func (qs *QueryStats) Sort(sortOptions *SortOptions, reverse bool) {
	switch sortOptions.sortType {
	case SortCount:
		qs.SortCount(reverse)
	case SortQuery:
		qs.SortQuery(reverse)
	// query time
	case SortMaxQueryTime:
		qs.SortMaxQueryTime(reverse)
	case SortMinQueryTime:
		qs.SortMinQueryTime(reverse)
	case SortSumQueryTime:
		qs.SortSumQueryTime(reverse)
	case SortAvgQueryTime:
		qs.SortAvgQueryTime(reverse)
	case SortPNQueryTime:
		qs.SortPNQueryTime(reverse)
	// lock time
	case SortMaxLockTime:
		qs.SortMaxLockTime(reverse)
	case SortMinLockTime:
		qs.SortMinLockTime(reverse)
	case SortSumLockTime:
		qs.SortSumLockTime(reverse)
	case SortAvgLockTime:
		qs.SortAvgLockTime(reverse)
	case SortPNLockTime:
		qs.SortPNLockTime(reverse)
		// rows sent
	case SortMaxRowsSent:
		qs.SortMaxRowsSent(reverse)
	case SortMinRowsSent:
		qs.SortMinRowsSent(reverse)
	case SortSumRowsSent:
		qs.SortSumRowsSent(reverse)
	case SortAvgRowsSent:
		qs.SortAvgRowsSent(reverse)
	case SortPNRowsSent:
		qs.SortPNRowsSent(reverse)
	// rows examined
	case SortMaxRowsExamined:
		qs.SortMaxRowsExamined(reverse)
	case SortMinRowsExamined:
		qs.SortMinRowsExamined(reverse)
	case SortSumRowsExamined:
		qs.SortSumRowsExamined(reverse)
	case SortAvgRowsExamined:
		qs.SortAvgRowsExamined(reverse)
	case SortPNRowsExamined:
		qs.SortPNRowsExamined(reverse)
	// rows affected
	case SortMaxRowsAffected:
		qs.SortMaxRowsAffected(reverse)
	case SortMinRowsAffected:
		qs.SortMinRowsAffected(reverse)
	case SortSumRowsAffected:
		qs.SortSumRowsAffected(reverse)
	case SortAvgRowsAffected:
		qs.SortAvgRowsAffected(reverse)
	case SortPNRowsAffected:
		qs.SortPNRowsAffected(reverse)
	// bytes sent
	case SortMaxBytesSent:
		qs.SortMaxBytesSent(reverse)
	case SortMinBytesSent:
		qs.SortMinBytesSent(reverse)
	case SortSumBytesSent:
		qs.SortSumBytesSent(reverse)
	case SortAvgBytesSent:
		qs.SortAvgBytesSent(reverse)
	case SortPNBytesSent:
		qs.SortPNBytesSent(reverse)
	default:
		qs.SortCount(reverse)
	}
}

func (qs *QueryStats) SortCount(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].Count() > qs.stats[j].Count()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].Count() < qs.stats[j].Count()
		})
	}
}

func (qs *QueryStats) SortQuery(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].Query > qs.stats[j].Query
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].Query < qs.stats[j].Query
		})
	}
}

// query time
func (qs *QueryStats) SortMaxQueryTime(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MaxQueryTime() > qs.stats[j].MaxQueryTime()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MaxQueryTime() < qs.stats[j].MaxQueryTime()
		})
	}
}

func (qs *QueryStats) SortMinQueryTime(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MinQueryTime() > qs.stats[j].MinQueryTime()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MinQueryTime() < qs.stats[j].MinQueryTime()
		})
	}
}

func (qs *QueryStats) SortAvgQueryTime(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].AvgQueryTime() > qs.stats[j].AvgQueryTime()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].AvgQueryTime() < qs.stats[j].AvgQueryTime()
		})
	}
}

func (qs *QueryStats) SortSumQueryTime(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].SumQueryTime() > qs.stats[j].SumQueryTime()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].SumQueryTime() < qs.stats[j].SumQueryTime()
		})
	}
}

func (qs *QueryStats) SortPNQueryTime(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].PNQueryTime(qs.sortOptions.percentile) > qs.stats[j].PNQueryTime(qs.sortOptions.percentile)
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].PNQueryTime(qs.sortOptions.percentile) < qs.stats[j].PNQueryTime(qs.sortOptions.percentile)
		})
	}
}

// lock time
func (qs *QueryStats) SortMaxLockTime(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MaxLockTime() > qs.stats[j].MaxLockTime()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MaxLockTime() < qs.stats[j].MaxLockTime()
		})
	}
}

func (qs *QueryStats) SortMinLockTime(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MinLockTime() > qs.stats[j].MinLockTime()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MinLockTime() < qs.stats[j].MinLockTime()
		})
	}
}

func (qs *QueryStats) SortAvgLockTime(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].AvgLockTime() > qs.stats[j].AvgLockTime()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].AvgLockTime() < qs.stats[j].AvgLockTime()
		})
	}
}

func (qs *QueryStats) SortSumLockTime(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].SumLockTime() > qs.stats[j].SumLockTime()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].SumLockTime() < qs.stats[j].SumLockTime()
		})
	}
}

func (qs *QueryStats) SortPNLockTime(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].PNLockTime(qs.sortOptions.percentile) > qs.stats[j].PNLockTime(qs.sortOptions.percentile)
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].PNLockTime(qs.sortOptions.percentile) < qs.stats[j].PNLockTime(qs.sortOptions.percentile)
		})
	}
}

// rows sent
func (qs *QueryStats) SortMaxRowsSent(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MaxRowsSent() > qs.stats[j].MaxRowsSent()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MaxRowsSent() < qs.stats[j].MaxRowsSent()
		})
	}
}

func (qs *QueryStats) SortMinRowsSent(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MinRowsSent() > qs.stats[j].MinRowsSent()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MinRowsSent() < qs.stats[j].MinRowsSent()
		})
	}
}

func (qs *QueryStats) SortAvgRowsSent(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].AvgRowsSent() > qs.stats[j].AvgRowsSent()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].AvgRowsSent() < qs.stats[j].AvgRowsSent()
		})
	}
}

func (qs *QueryStats) SortSumRowsSent(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].SumRowsSent() > qs.stats[j].SumRowsSent()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].SumRowsSent() < qs.stats[j].SumRowsSent()
		})
	}
}

func (qs *QueryStats) SortPNRowsSent(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].PNRowsSent(qs.sortOptions.percentile) > qs.stats[j].PNRowsSent(qs.sortOptions.percentile)
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].PNRowsSent(qs.sortOptions.percentile) < qs.stats[j].PNRowsSent(qs.sortOptions.percentile)
		})
	}
}

// rows examined
func (qs *QueryStats) SortMaxRowsExamined(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MaxRowsExamined() > qs.stats[j].MaxRowsExamined()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MaxRowsExamined() < qs.stats[j].MaxRowsExamined()
		})
	}
}

func (qs *QueryStats) SortMinRowsExamined(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MinRowsExamined() > qs.stats[j].MinRowsExamined()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MinRowsExamined() < qs.stats[j].MinRowsExamined()
		})
	}
}

func (qs *QueryStats) SortAvgRowsExamined(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].AvgRowsExamined() > qs.stats[j].AvgRowsExamined()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].AvgRowsExamined() < qs.stats[j].AvgRowsExamined()
		})
	}
}

func (qs *QueryStats) SortSumRowsExamined(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].SumRowsExamined() > qs.stats[j].SumRowsExamined()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].SumRowsExamined() < qs.stats[j].SumRowsExamined()
		})
	}
}

func (qs *QueryStats) SortPNRowsExamined(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].PNRowsExamined(qs.sortOptions.percentile) > qs.stats[j].PNRowsExamined(qs.sortOptions.percentile)
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].PNRowsExamined(qs.sortOptions.percentile) < qs.stats[j].PNRowsExamined(qs.sortOptions.percentile)
		})
	}
}

// rows affected
func (qs *QueryStats) SortMaxRowsAffected(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MaxRowsAffected() > qs.stats[j].MaxRowsAffected()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MaxRowsAffected() < qs.stats[j].MaxRowsAffected()
		})
	}
}

func (qs *QueryStats) SortMinRowsAffected(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MinRowsAffected() > qs.stats[j].MinRowsAffected()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MinRowsAffected() < qs.stats[j].MinRowsAffected()
		})
	}
}

func (qs *QueryStats) SortAvgRowsAffected(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].AvgRowsAffected() > qs.stats[j].AvgRowsAffected()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].AvgRowsAffected() < qs.stats[j].AvgRowsAffected()
		})
	}
}

func (qs *QueryStats) SortSumRowsAffected(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].SumRowsAffected() > qs.stats[j].SumRowsAffected()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].SumRowsAffected() < qs.stats[j].SumRowsAffected()
		})
	}
}

func (qs *QueryStats) SortPNRowsAffected(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].PNRowsAffected(qs.sortOptions.percentile) > qs.stats[j].PNRowsAffected(qs.sortOptions.percentile)
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].PNRowsAffected(qs.sortOptions.percentile) < qs.stats[j].PNRowsAffected(qs.sortOptions.percentile)
		})
	}
}

// bytes sent
func (qs *QueryStats) SortMaxBytesSent(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MaxBytesSent() > qs.stats[j].MaxBytesSent()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MaxBytesSent() < qs.stats[j].MaxBytesSent()
		})
	}
}

func (qs *QueryStats) SortMinBytesSent(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MinBytesSent() > qs.stats[j].MinBytesSent()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].MinBytesSent() < qs.stats[j].MinBytesSent()
		})
	}
}

func (qs *QueryStats) SortAvgBytesSent(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].AvgBytesSent() > qs.stats[j].AvgBytesSent()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].AvgBytesSent() < qs.stats[j].AvgBytesSent()
		})
	}
}

func (qs *QueryStats) SortSumBytesSent(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].SumBytesSent() > qs.stats[j].SumBytesSent()
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].SumBytesSent() < qs.stats[j].SumBytesSent()
		})
	}
}

func (qs *QueryStats) SortPNBytesSent(reverse bool) {
	if reverse {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].PNBytesSent(qs.sortOptions.percentile) > qs.stats[j].PNBytesSent(qs.sortOptions.percentile)
		})
	} else {
		sort.Slice(qs.stats, func(i, j int) bool {
			return qs.stats[i].PNBytesSent(qs.sortOptions.percentile) < qs.stats[j].PNBytesSent(qs.sortOptions.percentile)
		})
	}
}
