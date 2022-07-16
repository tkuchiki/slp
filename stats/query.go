package stats

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"sync"

	mlog "github.com/percona/go-mysql/log"
	"github.com/tkuchiki/slp/errors"
	"github.com/tkuchiki/slp/helper"
	"github.com/tkuchiki/slp/options"
)

type hints struct {
	values map[string]int
	len    int
	mu     sync.RWMutex
}

func newHints() *hints {
	return &hints{
		values: make(map[string]int),
	}
}

func (h *hints) loadOrStore(key string) int {
	h.mu.Lock()
	defer h.mu.Unlock()
	_, ok := h.values[key]
	if !ok {
		h.values[key] = h.len
		h.len++
	}

	return h.values[key]
}

type QueryStats struct {
	hints                     *hints
	stats                     queryStats
	useQueryTimePercentile    bool
	useLockTimePercentile     bool
	useRowsSentPercentile     bool
	useRowsExaminedPercentile bool
	useRowsAffectedPercentile bool
	useBytesSentPercentile    bool
	filter                    *Filter
	options                   *options.Options
	sortOptions               *SortOptions
	queryMatchingGroups       []*regexp.Regexp
}

func NewQueryStats(useQueryTimePercentile, useLockTimePercentile, useRowsSentPercentile, useRowsExaminedPercentile, useRowsAffectedPercentile, useBytesSent bool) *QueryStats {
	return &QueryStats{
		hints:                     newHints(),
		stats:                     make([]*QueryStat, 0),
		useQueryTimePercentile:    useQueryTimePercentile,
		useLockTimePercentile:     useLockTimePercentile,
		useRowsSentPercentile:     useRowsSentPercentile,
		useRowsExaminedPercentile: useRowsExaminedPercentile,
		useRowsAffectedPercentile: useRowsAffectedPercentile,
		useBytesSentPercentile:    useBytesSent,
	}
}

func (qs *QueryStats) Set(query string, querytime, lockTime float64, rowsSent, rowsExamined, rowsAffected, bytesSent uint64) {
	if len(qs.queryMatchingGroups) > 0 {
		for _, re := range qs.queryMatchingGroups {
			if ok := re.Match([]byte(query)); ok {
				pattern := re.String()
				query = pattern
				break
			}
		}
	}

	idx := qs.hints.loadOrStore(query)

	if idx >= len(qs.stats) {
		qs.stats = append(qs.stats, newQueryStat(query, qs.useQueryTimePercentile, qs.useLockTimePercentile, qs.useRowsSentPercentile, qs.useRowsExaminedPercentile, qs.useRowsAffectedPercentile, qs.useBytesSentPercentile))
	}

	qs.stats[idx].Set(querytime, lockTime, rowsSent, rowsExamined, rowsAffected, bytesSent)
}

func (qs *QueryStats) Stats() []*QueryStat {
	return qs.stats
}

func (qs *QueryStats) CountUris() int {
	return qs.hints.len
}

func (qs *QueryStats) SetOptions(options *options.Options) {
	qs.options = options
}

func (qs *QueryStats) SetSortOptions(options *SortOptions) {
	qs.sortOptions = options
}

func (qs *QueryStats) SetQueryMatchingGroups(groups []string) error {
	queryGroups, err := helper.CompileQueryMatchingGroups(groups)
	if err != nil {
		return err
	}

	qs.queryMatchingGroups = queryGroups

	return nil
}

func (qs *QueryStats) InitFilter(options *options.Options) error {
	qs.filter = NewFilter(options)
	return qs.filter.Init()
}

func (qs *QueryStats) DoFilter(qstat *mlog.Event) (bool, error) {
	err := qs.filter.Do(qstat)
	if err == errors.SkipReadLineErr {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (qs *QueryStats) CountAll() map[string]int {
	counts := make(map[string]int, 1)

	for _, s := range qs.stats {
		counts["count"] += s.Cnt
	}

	return counts
}

func (qs *QueryStats) SortWithOptions() {
	qs.Sort(qs.sortOptions, qs.options.Reverse)
}

type QueryStat struct {
	Query        string       `yaml:"query"`
	Cnt          int          `yaml:"count"`
	QueryTime    *timeStats   `yaml:"query_time"`
	LockTime     *timeStats   `yaml:"lock_time"`
	RowsSent     *numberStats `yaml:"rows_sent"`
	RowsExamined *numberStats `yaml:"rows_examined"`
	RowsAffected *numberStats `yaml:"rows_affected"`
	BytesSent    *numberStats `yaml:"bytes_sent"`
}

type queryStats []*QueryStat

func newQueryStat(query string, useQueryTimePercentile, useLockTimePercentile, useRowsSentPercentile, useRowsExaminedPercentile, useRowsAffectedPercentile, useBytesSent bool) *QueryStat {
	return &QueryStat{
		Query:        query,
		QueryTime:    newTimeStats(useQueryTimePercentile),
		LockTime:     newTimeStats(useLockTimePercentile),
		RowsSent:     newNumberStats(useRowsSentPercentile),
		RowsExamined: newNumberStats(useRowsExaminedPercentile),
		RowsAffected: newNumberStats(useRowsAffectedPercentile),
		BytesSent:    newNumberStats(useBytesSent),
	}
}

func (qs *QueryStat) Set(queryTime, lockTime float64, rowsSent, rowsExamined, rowsAffected, bytesSent uint64) {
	qs.Cnt++
	qs.QueryTime.Set(queryTime)
	qs.LockTime.Set(lockTime)
	qs.RowsSent.Set(rowsSent)
	qs.RowsExamined.Set(rowsExamined)
	qs.RowsAffected.Set(rowsAffected)
	qs.BytesSent.Set(bytesSent)
}

func (qs *QueryStat) Count() int {
	return qs.Cnt
}

func (qs *QueryStat) StrCount() string {
	return fmt.Sprint(qs.Cnt)
}

// query_time
func (qs *QueryStat) MaxQueryTime() float64 {
	return qs.QueryTime.Max
}

func (qs *QueryStat) StrMaxQueryTime() string {
	return fmt.Sprintf("%.6f", qs.QueryTime.Max)
}

func (qs *QueryStat) MinQueryTime() float64 {
	return qs.QueryTime.Min
}

func (qs *QueryStat) StrMinQueryTime() string {
	return fmt.Sprintf("%.6f", qs.QueryTime.Min)
}

func (qs *QueryStat) SumQueryTime() float64 {
	return qs.QueryTime.Sum
}

func (qs *QueryStat) StrSumQueryTime() string {
	return fmt.Sprintf("%.6f", qs.QueryTime.Sum)
}

func (qs *QueryStat) AvgQueryTime() float64 {
	return qs.QueryTime.Avg(qs.Cnt)
}

func (qs *QueryStat) StrAvgQueryTime() string {
	return fmt.Sprintf("%.6f", qs.QueryTime.Avg(qs.Cnt))
}

func (qs *QueryStat) PNQueryTime(n int) float64 {
	return qs.QueryTime.PN(qs.Cnt, n)
}

func (qs *QueryStat) StrPNQueryTime(n int) string {
	return fmt.Sprintf("%.6f", qs.QueryTime.PN(qs.Cnt, n))
}

func (qs *QueryStat) StddevQueryTime() float64 {
	return qs.QueryTime.Stddev(qs.Cnt)
}

func (qs *QueryStat) StrStddevQueryTime() string {
	return fmt.Sprintf("%.6f", qs.QueryTime.Stddev(qs.Cnt))
}

// lock_time
func (qs *QueryStat) MaxLockTime() float64 {
	return qs.LockTime.Max
}

func (qs *QueryStat) StrMaxLockTime() string {
	return fmt.Sprintf("%.6f", qs.LockTime.Max)
}

func (qs *QueryStat) MinLockTime() float64 {
	return qs.LockTime.Min
}

func (qs *QueryStat) StrMinLockTime() string {
	return fmt.Sprintf("%.6f", qs.LockTime.Min)
}

func (qs *QueryStat) SumLockTime() float64 {
	return qs.LockTime.Sum
}

func (qs *QueryStat) StrSumLockTime() string {
	return fmt.Sprintf("%.6f", qs.LockTime.Sum)
}

func (qs *QueryStat) AvgLockTime() float64 {
	return qs.LockTime.Avg(qs.Cnt)
}

func (qs *QueryStat) StrAvgLockTime() string {
	return fmt.Sprintf("%.6f", qs.LockTime.Avg(qs.Cnt))
}

func (qs *QueryStat) PNLockTime(n int) float64 {
	return qs.LockTime.PN(qs.Cnt, n)
}

func (qs *QueryStat) StrPNLockTime(n int) string {
	return fmt.Sprintf("%.6f", qs.LockTime.PN(qs.Cnt, n))
}

func (qs *QueryStat) StddevLockTime() float64 {
	return qs.LockTime.Stddev(qs.Cnt)
}

func (qs *QueryStat) StrStddevLockTime() string {
	return fmt.Sprintf("%.6f", qs.LockTime.Stddev(qs.Cnt))
}

// rows_sent
func (qs *QueryStat) MaxRowsSent() uint64 {
	return qs.RowsSent.Max
}

func (qs *QueryStat) StrMaxRowsSent() string {
	return fmt.Sprint(qs.RowsSent.Max)
}

func (qs *QueryStat) MinRowsSent() uint64 {
	return qs.RowsSent.Min
}

func (qs *QueryStat) StrMinRowsSent() string {
	return fmt.Sprint(qs.RowsSent.Min)
}

func (qs *QueryStat) SumRowsSent() uint64 {
	return qs.RowsSent.Sum
}

func (qs *QueryStat) StrSumRowsSent() string {
	return fmt.Sprint(qs.RowsSent.Sum)
}

func (qs *QueryStat) AvgRowsSent() float64 {
	return qs.RowsSent.Avg(qs.Cnt)
}

func (qs *QueryStat) StrAvgRowsSent() string {
	return fmt.Sprintf("%.6f", qs.RowsSent.Avg(qs.Cnt))
}

func (qs *QueryStat) PNRowsSent(n int) uint64 {
	return qs.RowsSent.PN(qs.Cnt, n)
}

func (qs *QueryStat) StrPNRowsSent(n int) string {
	return fmt.Sprint(qs.RowsSent.PN(qs.Cnt, n))
}

func (qs *QueryStat) StddevRowsSent() float64 {
	return qs.RowsSent.Stddev(qs.Cnt)
}

func (qs *QueryStat) StrStddevRowsSent() string {
	return fmt.Sprintf("%.6f", qs.RowsSent.Stddev(qs.Cnt))
}

// rows_examined
func (qs *QueryStat) MaxRowsExamined() uint64 {
	return qs.RowsExamined.Max
}

func (qs *QueryStat) StrMaxRowsExamined() string {
	return fmt.Sprint(qs.RowsExamined.Max)
}

func (qs *QueryStat) MinRowsExamined() uint64 {
	return qs.RowsExamined.Min
}

func (qs *QueryStat) StrMinRowsExamined() string {
	return fmt.Sprint(qs.RowsExamined.Min)
}

func (qs *QueryStat) SumRowsExamined() uint64 {
	return qs.RowsExamined.Sum
}

func (qs *QueryStat) StrSumRowsExamined() string {
	return fmt.Sprint(qs.RowsExamined.Sum)
}

func (qs *QueryStat) AvgRowsExamined() float64 {
	return qs.RowsExamined.Avg(qs.Cnt)
}

func (qs *QueryStat) StrAvgRowsExamined() string {
	return fmt.Sprintf("%.6f", qs.RowsExamined.Avg(qs.Cnt))
}

func (qs *QueryStat) PNRowsExamined(n int) uint64 {
	return qs.RowsExamined.PN(qs.Cnt, n)
}

func (qs *QueryStat) StrPNRowsExamined(n int) string {
	return fmt.Sprint(qs.RowsExamined.PN(qs.Cnt, n))
}

func (qs *QueryStat) StddevRowsExamined() float64 {
	return qs.RowsExamined.Stddev(qs.Cnt)
}

func (qs *QueryStat) StrStddevRowsExamined() string {
	return fmt.Sprintf("%.6f", qs.RowsExamined.Stddev(qs.Cnt))
}

//	rows_affected
func (qs *QueryStat) MaxRowsAffected() uint64 {
	return qs.RowsAffected.Max
}

func (qs *QueryStat) StrMaxRowsAffected() string {
	return fmt.Sprint(qs.RowsAffected.Max)
}

func (qs *QueryStat) MinRowsAffected() uint64 {
	return qs.RowsAffected.Min
}

func (qs *QueryStat) StrMinRowsAffected() string {
	return fmt.Sprint(qs.RowsAffected.Min)
}

func (qs *QueryStat) SumRowsAffected() uint64 {
	return qs.RowsAffected.Sum
}

func (qs *QueryStat) StrSumRowsAffected() string {
	return fmt.Sprint(qs.RowsAffected.Sum)
}

func (qs *QueryStat) AvgRowsAffected() float64 {
	return qs.RowsAffected.Avg(qs.Cnt)
}

func (qs *QueryStat) StrAvgRowsAffected() string {
	return fmt.Sprintf("%.6f", qs.RowsAffected.Avg(qs.Cnt))
}

func (qs *QueryStat) PNRowsAffected(n int) uint64 {
	return qs.RowsAffected.PN(qs.Cnt, n)
}

func (qs *QueryStat) StrPNRowsAffected(n int) string {
	return fmt.Sprint(qs.RowsAffected.PN(qs.Cnt, n))
}

func (qs *QueryStat) StddevRowsAffected() float64 {
	return qs.RowsAffected.Stddev(qs.Cnt)
}

func (qs *QueryStat) StrStddevRowsAffected() string {
	return fmt.Sprintf("%.6f", qs.RowsAffected.Stddev(qs.Cnt))
}

//	bytes_sent"
func (qs *QueryStat) MaxBytesSent() uint64 {
	return qs.BytesSent.Max
}

func (qs *QueryStat) StrMaxBytesSent() string {
	return fmt.Sprint(qs.BytesSent.Max)
}

func (qs *QueryStat) MinBytesSent() uint64 {
	return qs.BytesSent.Min
}

func (qs *QueryStat) StrMinBytesSent() string {
	return fmt.Sprint(qs.BytesSent.Min)
}

func (qs *QueryStat) SumBytesSent() uint64 {
	return qs.BytesSent.Sum
}

func (qs *QueryStat) StrSumBytesSent() string {
	return fmt.Sprint(qs.BytesSent.Sum)
}

func (qs *QueryStat) AvgBytesSent() float64 {
	return qs.BytesSent.Avg(qs.Cnt)
}

func (qs *QueryStat) StrAvgBytesSent() string {
	return fmt.Sprintf("%.6f", qs.BytesSent.Avg(qs.Cnt))
}

func (qs *QueryStat) PNBytesSent(n int) uint64 {
	return qs.BytesSent.PN(qs.Cnt, n)
}

func (qs *QueryStat) StrPNBytesSent(n int) string {
	return fmt.Sprint(qs.BytesSent.PN(qs.Cnt, n))
}

func (qs *QueryStat) StddevBytesSent() float64 {
	return qs.BytesSent.Stddev(qs.Cnt)
}

func (qs *QueryStat) StrStddevBytesSent() string {
	return fmt.Sprintf("%.6f", qs.BytesSent.Stddev(qs.Cnt))
}

func percentRank(n int, pi int) int {
	if pi == 0 {
		return 0
	} else if pi == 100 {
		return n - 1
	}

	p := float64(pi) / 100.0
	pos := int(float64(n+1) * p)
	if pos < 0 {
		pos = 0
	}

	return pos - 1
}

type timeStats struct {
	Max           float64 `yaml:"max"`
	Min           float64 `yaml:"min"`
	Sum           float64 `yaml:"sum"`
	UsePercentile bool
	Percentiles   []float64 `yaml:"percentiles"`
}

func newTimeStats(usePercentile bool) *timeStats {
	return &timeStats{
		UsePercentile: usePercentile,
		Percentiles:   make([]float64, 0),
	}
}

func (ts *timeStats) Set(val float64) {
	if ts.Max < val {
		ts.Max = val
	}

	if ts.Min >= val || ts.Min == 0 {
		ts.Min = val
	}

	ts.Sum += val

	if ts.UsePercentile {
		ts.Percentiles = append(ts.Percentiles, val)
	}
}

func (ts *timeStats) Avg(cnt int) float64 {
	return ts.Sum / float64(cnt)
}

func (ts *timeStats) PN(cnt, n int) float64 {
	if !ts.UsePercentile {
		return 0.0
	}

	plen := percentRank(cnt, n)
	ts.Sort()
	return ts.Percentiles[plen]
}

func (ts *timeStats) Stddev(cnt int) float64 {
	if !ts.UsePercentile {
		return 0.0
	}

	var stdd float64
	avg := ts.Avg(cnt)
	n := float64(cnt)

	for _, v := range ts.Percentiles {
		stdd += (v - avg) * (v - avg)
	}

	return math.Sqrt(stdd / n)
}

func (ts *timeStats) Sort() {
	sort.Slice(ts.Percentiles, func(i, j int) bool {
		return ts.Percentiles[i] < ts.Percentiles[j]
	})
}

type numberStats struct {
	Max           uint64 `yaml:"max"`
	Min           uint64 `yaml:"min"`
	Sum           uint64 `yaml:"sum"`
	UsePercentile bool
	Percentiles   []uint64 `yaml:"percentiles"`
}

func newNumberStats(usePercentile bool) *numberStats {
	return &numberStats{
		UsePercentile: usePercentile,
		Percentiles:   make([]uint64, 0),
	}
}

func (ns *numberStats) Set(val uint64) {
	if ns.Max < val {
		ns.Max = val
	}

	if ns.Min >= val || ns.Min == 0.0 {
		ns.Min = val
	}

	ns.Sum += val

	if ns.UsePercentile {
		ns.Percentiles = append(ns.Percentiles, val)
	}
}

func (ns *numberStats) Avg(cnt int) float64 {
	return float64(ns.Sum) / float64(cnt)
}

func (ns *numberStats) PN(cnt, n int) uint64 {
	if !ns.UsePercentile {
		return 0.0
	}

	plen := percentRank(cnt, n)
	ns.Sort()
	return ns.Percentiles[plen]
}

func (ns *numberStats) Stddev(cnt int) float64 {
	if !ns.UsePercentile {
		return 0.0
	}

	var stdd float64
	avg := ns.Avg(cnt)
	n := float64(cnt)

	for _, v := range ns.Percentiles {
		stdd += (float64(v) - avg) * (float64(v) - avg)
	}

	return math.Sqrt(stdd / n)
}

func (ns *numberStats) Sort() {
	sort.Slice(ns.Percentiles, func(i, j int) bool {
		return ns.Percentiles[i] < ns.Percentiles[j]
	})
}
