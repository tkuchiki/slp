package stats

import (
	"fmt"
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/tkuchiki/slp/helper"
	"github.com/tkuchiki/slp/html"
)

type keywordMode string

const (
	simpleMode   keywordMode = "simple"
	standardMode keywordMode = "standard"
	allMode      keywordMode = "all"
)

var outputKeywords = map[keywordMode][]string{
	simpleMode:   {"count", "query", "query-time"},
	standardMode: {"count", "query", "query-time", "lock-time", "rows-sent", "rows-examined"},
	allMode:      {"count", "query", "query-time", "lock-time", "rows-sent", "rows-examined", "rows-affected", "bytes-sent"},
}

var statsKeyPrefixes = []string{"min", "max", "sum", "avg"}

func OutputKeywords(percentiles []int) []string {
	var s []string
	for _, key := range outputKeywords[allMode] {
		if key == "count" || key == "query" {
			s = append(s, key)
			continue
		}

		for _, prefix := range statsKeyPrefixes {
			s = append(s, fmt.Sprintf("%s-%s", prefix, key))
		}

		for _, p := range percentiles {
			s = append(s, fmt.Sprintf("p%d-%s", p, key))
		}
	}

	return s
}

func keywords(percentiles []int, keys []string) []string {
	var s []string

	for _, key := range keys {
		if key == "count" || key == "query" {
			s = append(s, key)
			continue
		}

		for _, prefix := range statsKeyPrefixes {
			s = append(s, fmt.Sprintf("%s-%s", prefix, key))
		}
		for _, p := range percentiles {
			s = append(s, fmt.Sprintf("p%d-%s", p, key))
		}
	}

	return s
}

func defaultHeaders(percentiles []int, keys []string) []string {
	headers := headersMap(percentiles)

	var s []string

	for _, key := range keys {
		if key == "count" || key == "query" {
			s = append(s, headers[key])
			continue
		}

		for _, prefix := range statsKeyPrefixes {
			s = append(s, headers[fmt.Sprintf("%s-%s", prefix, key)])
		}
		for _, p := range percentiles {
			s = append(s, headers[fmt.Sprintf("p%d-%s", p, key)])
		}
	}

	return s
}

var defalutHeaders = map[string]string{
	"count":             "Count",
	"query":             "Query",
	"min-query-time":    "Min(QueryTime)",
	"max-query-time":    "Max(QueryTime)",
	"sum-query-time":    "Sum(QueryTime)",
	"avg-query-time":    "Avg(QueryTime)",
	"min-lock-time":     "Min(LockTime)",
	"max-lock-time":     "Max(LockTime)",
	"sum-lock-time":     "Sum(LockTime)",
	"avg-lock-time":     "Avg(LockTime)",
	"min-rows-sent":     "Min(RowsSent)",
	"max-rows-sent":     "Max(RowsSent)",
	"sum-rows-sent":     "Sum(RowsSent)",
	"avg-rows-sent":     "Avg(RowsSent)",
	"min-rows-examined": "Min(RowsExamined)",
	"max-rows-examined": "Max(RowsExamined)",
	"sum-rows-examined": "Sum(RowsExamined)",
	"avg-rows-examined": "Avg(RowsExamined)",
	"min-rows-affected": "Min(RowsAffected)",
	"max-rows-affected": "Max(RowsAffected)",
	"sum-rows-affected": "Sum(RowsAffected)",
	"avg-rows-affected": "Avg(RowsAffected)",
	"min-bytes-sent":    "Min(BytesSent)",
	"max-bytes-sent":    "Max(BytesSent)",
	"sum-bytes-sent":    "Sum(BytesSent)",
	"avg-bytes-sent":    "Avg(BytesSent)",
}

var percentileMap = map[string]string{
	"query-time":    "QueryTime",
	"lock-time":     "LockTime",
	"rows-sent":     "RowsSent",
	"rows-examined": "RowsExamined",
	"rows-affected": "RowsAffected",
	"bytes-sent":    "BytesSent",
}

func headersMap(percentiles []int) map[string]string {
	headers := defalutHeaders

	for _, p := range percentiles {
		for pkey, pval := range percentileMap {
			key := fmt.Sprintf("p%d-%s", p, pkey)
			val := fmt.Sprintf("P%d(%s)", p, pval)
			headers[key] = val
		}
	}

	return headers
}

type PrintOptions struct {
	noHeaders       bool
	showFooters     bool
	paginationLimit int
}

func NewPrintOptions(noHeaders, showFooters bool, paginationLimit int) *PrintOptions {
	return &PrintOptions{
		noHeaders:       noHeaders,
		showFooters:     showFooters,
		paginationLimit: paginationLimit,
	}
}

type Printer struct {
	keywords     []string
	format       string
	percentiles  []int
	printOptions *PrintOptions
	headers      []string
	headersMap   map[string]string
	writer       io.Writer
	all          bool
}

func NewPrinter(w io.Writer, val, format string, percentiles []int, printOptions *PrintOptions) *Printer {
	p := &Printer{
		format:       format,
		percentiles:  percentiles,
		headersMap:   headersMap(percentiles),
		writer:       w,
		printOptions: printOptions,
	}

	switch keywordMode(val) {
	case simpleMode:
		p.keywords = keywords(percentiles, outputKeywords[simpleMode])
		p.headers = defaultHeaders(percentiles, outputKeywords[simpleMode])
	case standardMode:
		p.keywords = keywords(percentiles, outputKeywords[standardMode])
		p.headers = defaultHeaders(percentiles, outputKeywords[standardMode])
	case allMode:
		p.keywords = keywords(percentiles, outputKeywords[allMode])
		p.headers = defaultHeaders(percentiles, outputKeywords[allMode])
		p.all = true
	default:
		p.keywords = helper.SplitCSV(val)
		for _, key := range p.keywords {
			p.headers = append(p.headers, p.headersMap[key])
			if key == "all" {
				p.keywords = keywords(percentiles, outputKeywords[allMode])
				p.headers = defaultHeaders(percentiles, outputKeywords[allMode])
				p.all = true
				break
			}
		}
	}

	return p
}

func (p *Printer) Validate() error {
	if p.all {
		return nil
	}

	invalids := make([]string, 0, len(p.keywords))
	for _, key := range p.keywords {
		if _, ok := p.headersMap[key]; !ok {
			invalids = append(invalids, key)
		}
	}

	if len(invalids) > 0 {
		return fmt.Errorf("invalid keywords: %s", strings.Join(invalids, ","))
	}

	return nil
}

func (p *Printer) GenerateLine(s *QueryStat) []string {
	keyLen := len(p.keywords)
	line := make([]string, 0, keyLen)

	for i := 0; i < keyLen; i++ {
		switch p.keywords[i] {
		case "count":
			line = append(line, s.StrCount())
		case "query":
			line = append(line, s.Query)
		case "min-query-time":
			line = append(line, s.StrMinQueryTime())
		case "max-query-time":
			line = append(line, s.StrMaxQueryTime())
		case "sum-query-time":
			line = append(line, s.StrSumQueryTime())
		case "avg-query-time":
			line = append(line, s.StrAvgQueryTime())
		case "min-lock-time":
			line = append(line, s.StrMinLockTime())
		case "max-lock-time":
			line = append(line, s.StrMaxLockTime())
		case "sum-lock-time":
			line = append(line, s.StrSumLockTime())
		case "avg-lock-time":
			line = append(line, s.StrAvgLockTime())
		case "min-rows-sent":
			line = append(line, s.StrMinRowsSent())
		case "max-rows-sent":
			line = append(line, s.StrMaxRowsSent())
		case "sum-rows-sent":
			line = append(line, s.StrSumRowsSent())
		case "avg-rows-sent":
			line = append(line, s.StrAvgRowsSent())
		case "min-rows-examined":
			line = append(line, s.StrMinRowsExamined())
		case "max-rows-examined":
			line = append(line, s.StrMaxRowsExamined())
		case "sum-rows-examined":
			line = append(line, s.StrSumRowsExamined())
		case "avg-rows-examined":
			line = append(line, s.StrAvgRowsExamined())
		case "min-rows-affected":
			line = append(line, s.StrMinRowsAffected())
		case "max-rows-affected":
			line = append(line, s.StrMaxRowsAffected())
		case "sum-rows-affected":
			line = append(line, s.StrSumRowsAffected())
		case "avg-rows-affected":
			line = append(line, s.StrAvgRowsAffected())
		case "min-bytes-sent":
			line = append(line, s.StrMinBytesSent())
		case "max-bytes-sent":
			line = append(line, s.StrMaxBytesSent())
		case "sum-bytes-sent":
			line = append(line, s.StrSumBytesSent())
		case "avg-bytes-sent":
			line = append(line, s.StrAvgBytesSent())
		default: // percentile
			var n int
			var key string
			_, err := fmt.Sscanf(p.keywords[i], "p%d-%s", &n, &key)
			if err != nil {
				continue
			}
			switch key {
			case "query-time":
				line = append(line, s.StrPNQueryTime(n))
			case "lock-time":
				line = append(line, s.StrPNLockTime(n))
			case "rows-sent":
				line = append(line, s.StrPNRowsSent(n))
			case "rows-examined":
				line = append(line, s.StrPNRowsExamined(n))
			case "rows-affected":
				line = append(line, s.StrPNRowsAffected(n))
			case "bytes-sent":
				line = append(line, s.StrPNBytesSent(n))
			}
		}
	}

	return line
}

func formattedLineWithDiff(val, diff string) string {
	if diff == "+0" || diff == "+0.000" {
		return val
	}
	return fmt.Sprintf("%s (%s)", val, diff)
}

func (p *Printer) GenerateLineWithDiff(from, to *QueryStat) []string {
	keyLen := len(p.keywords)
	line := make([]string, 0, keyLen)

	differ := NewDiffer(from, to)

	for i := 0; i < keyLen; i++ {
		switch p.keywords[i] {
		case "count":
			line = append(line, formattedLineWithDiff(to.StrCount(), differ.DiffCnt()))
		case "query":
			line = append(line, to.Query)
		case "min-query-time":
			line = append(line, formattedLineWithDiff(to.StrMinQueryTime(), differ.DiffMinQueryTime()))
		case "max-query-time":
			line = append(line, formattedLineWithDiff(to.StrMaxQueryTime(), differ.DiffMaxQueryTime()))
		case "sum-query-time":
			line = append(line, formattedLineWithDiff(to.StrSumQueryTime(), differ.DiffSumQueryTime()))
		case "avg-query-time":
			line = append(line, formattedLineWithDiff(to.StrAvgQueryTime(), differ.DiffAvgQueryTime()))
		case "min-lock-time":
			line = append(line, formattedLineWithDiff(to.StrMinLockTime(), differ.DiffMinLockTime()))
		case "max-lock-time":
			line = append(line, formattedLineWithDiff(to.StrMaxLockTime(), differ.DiffMaxLockTime()))
		case "sum-lock-time":
			line = append(line, formattedLineWithDiff(to.StrSumLockTime(), differ.DiffSumLockTime()))
		case "avg-lock-time":
			line = append(line, formattedLineWithDiff(to.StrAvgLockTime(), differ.DiffAvgLockTime()))
		case "min-rows-sent":
			line = append(line, formattedLineWithDiff(to.StrMinRowsSent(), differ.DiffMinRowsSent()))
		case "max-rows-sent":
			line = append(line, formattedLineWithDiff(to.StrMaxRowsSent(), differ.DiffMaxRowsSent()))
		case "sum-rows-sent":
			line = append(line, formattedLineWithDiff(to.StrSumRowsSent(), differ.DiffSumRowsSent()))
		case "avg-rows-sent":
			line = append(line, formattedLineWithDiff(to.StrAvgRowsSent(), differ.DiffAvgRowsSent()))
		case "min-rows-examined":
			line = append(line, formattedLineWithDiff(to.StrMinRowsExamined(), differ.DiffMinRowsExamined()))
		case "max-rows-examined":
			line = append(line, formattedLineWithDiff(to.StrMaxRowsExamined(), differ.DiffMaxRowsExamined()))
		case "sum-rows-examined":
			line = append(line, formattedLineWithDiff(to.StrSumRowsExamined(), differ.DiffSumRowsExamined()))
		case "avg-rows-examined":
			line = append(line, formattedLineWithDiff(to.StrAvgRowsExamined(), differ.DiffAvgRowsExamined()))
		case "min-rows-affected":
			line = append(line, formattedLineWithDiff(to.StrMinRowsAffected(), differ.DiffMinRowsAffected()))
		case "max-rows-affected":
			line = append(line, formattedLineWithDiff(to.StrMaxRowsAffected(), differ.DiffMaxRowsAffected()))
		case "sum-rows-affected":
			line = append(line, formattedLineWithDiff(to.StrSumRowsAffected(), differ.DiffSumRowsAffected()))
		case "avg-rows-affected":
			line = append(line, formattedLineWithDiff(to.StrAvgRowsAffected(), differ.DiffAvgRowsAffected()))
		case "min-bytes-sent":
			line = append(line, formattedLineWithDiff(to.StrMinBytesSent(), differ.DiffMinBytesSent()))
		case "max-bytes-sent":
			line = append(line, formattedLineWithDiff(to.StrMaxBytesSent(), differ.DiffMaxBytesSent()))
		case "sum-bytes-sent":
			line = append(line, formattedLineWithDiff(to.StrSumBytesSent(), differ.DiffSumBytesSent()))
		case "avg-bytes-sent":
			line = append(line, formattedLineWithDiff(to.StrAvgBytesSent(), differ.DiffAvgBytesSent()))
		default: // percentile
			var n int
			var key string
			_, err := fmt.Sscanf(p.keywords[i], "p%d-%s", &n, &key)
			if err != nil {
				continue
			}

			switch key {
			case "query-time":
				line = append(line, to.StrPNQueryTime(n), differ.DiffPNQueryTime(n))
			case "lock-time":
				line = append(line, to.StrPNLockTime(n), differ.DiffPNLockTime(n))
			case "rows-sent":
				line = append(line, to.StrPNRowsSent(n), differ.DiffPNRowsSent(n))
			case "rows-examined":
				line = append(line, to.StrPNRowsExamined(n), differ.DiffPNRowsExamined(n))
			case "rows-affected":
				line = append(line, to.StrPNRowsAffected(n), differ.DiffPNRowsAffected(n))
			case "bytes-sent":
				line = append(line, to.StrPNBytesSent(n), differ.DiffPNBytesSent(n))
			}
		}
	}

	return line
}

func (p *Printer) GenerateFooter(counts map[string]int) []string {
	keyLen := len(p.keywords)
	line := make([]string, 0, keyLen)

	for i := 0; i < keyLen; i++ {
		switch p.keywords[i] {
		case "count":
			line = append(line, fmt.Sprint(counts["count"]))
		default:
			line = append(line, "")
		}
	}

	return line
}

func (p *Printer) GenerateFooterWithDiff(countsFrom, countsTo map[string]int) []string {
	keyLen := len(p.keywords)
	line := make([]string, 0, keyLen)
	counts := DiffCountAll(countsFrom, countsTo)

	for i := 0; i < keyLen; i++ {
		switch p.keywords[i] {
		case "count":
			line = append(line, formattedLineWithDiff(fmt.Sprint(countsTo["count"]), counts["count"]))
		default:
			line = append(line, "")
		}
	}

	return line
}

func (p *Printer) SetFormat(format string) {
	p.format = format
}

func (p *Printer) SetHeaders(headers []string) {
	p.headers = headers
}

func (p *Printer) SetWriter(w io.Writer) {
	p.writer = w
}

func (p *Printer) Print(qs, qsTo *QueryStats) {
	switch p.format {
	case "table":
		p.printTable(qs, qsTo)
	case "md", "markdown":
		p.printMarkdown(qs, qsTo)
	case "tsv":
		p.printTSV(qs, qsTo)
	case "csv":
		p.printCSV(qs, qsTo)
	case "html":
		p.printHTML(qs, qsTo)
	}
}

func findQueryStatFrom(qsFrom *QueryStats, qsTo *QueryStat) *QueryStat {
	for _, sFrom := range qsFrom.stats {
		if sFrom.Query == qsTo.Query {
			return sFrom
		}
	}
	return nil
}

func (p *Printer) printTable(qsFrom, qsTo *QueryStats) {
	table := tablewriter.NewWriter(p.writer)
	table.SetHeader(p.headers)
	if qsTo == nil {
		for _, s := range qsFrom.stats {
			data := p.GenerateLine(s)
			table.Append(data)
		}
	} else {
		for _, to := range qsTo.stats {
			from := findQueryStatFrom(qsFrom, to)

			var data []string
			if from == nil {
				data = p.GenerateLine(to)
			} else {
				data = p.GenerateLineWithDiff(from, to)
			}
			table.Append(data)
		}
	}

	if p.printOptions.showFooters {
		var footer []string
		if qsTo == nil {
			footer = p.GenerateFooter(qsFrom.CountAll())
		} else {
			footer = p.GenerateFooterWithDiff(qsFrom.CountAll(), qsTo.CountAll())
		}
		table.SetFooter(footer)
		table.SetFooterAlignment(tablewriter.ALIGN_LEFT)
	}

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}

func (p *Printer) printMarkdown(qsFrom, qsTo *QueryStats) {
	table := tablewriter.NewWriter(p.writer)
	table.SetHeader(p.headers)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	if qsTo == nil {
		for _, s := range qsFrom.stats {
			data := p.GenerateLine(s)
			table.Append(data)
		}
	} else {
		for _, to := range qsTo.stats {
			from := findQueryStatFrom(qsFrom, to)

			var data []string
			if from == nil {
				data = p.GenerateLine(to)
			} else {
				data = p.GenerateLineWithDiff(from, to)
			}
			table.Append(data)
		}
	}

	if p.printOptions.showFooters {
		var footer []string
		if qsTo == nil {
			footer = p.GenerateFooter(qsFrom.CountAll())
		} else {
			footer = p.GenerateFooterWithDiff(qsFrom.CountAll(), qsTo.CountAll())
		}
		table.Append(footer)
	}

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}

func (p *Printer) printTSV(qsFrom, qsTo *QueryStats) {
	if !p.printOptions.noHeaders {
		fmt.Println(strings.Join(p.headers, "\t"))
	}

	var data []string
	if qsTo == nil {
		for _, s := range qsFrom.stats {
			data = p.GenerateLine(s)
			fmt.Println(strings.Join(data, "\t"))
		}
	} else {
		for _, to := range qsTo.stats {
			from := findQueryStatFrom(qsFrom, to)

			if from == nil {
				data = p.GenerateLine(to)
			} else {
				data = p.GenerateLineWithDiff(from, to)
			}
			fmt.Println(strings.Join(data, "\t"))
		}
	}
}

func (p *Printer) printCSV(qsFrom, qsTo *QueryStats) {
	if !p.printOptions.noHeaders {
		fmt.Println(strings.Join(p.headers, ","))
	}

	var data []string
	if qsTo == nil {
		for _, s := range qsFrom.stats {
			data = p.GenerateLine(s)
			fmt.Println(strings.Join(data, ","))
		}
	} else {
		for _, to := range qsTo.stats {
			from := findQueryStatFrom(qsFrom, to)

			if from == nil {
				data = p.GenerateLine(to)
			} else {
				data = p.GenerateLineWithDiff(from, to)
			}
			fmt.Println(strings.Join(data, ","))
		}
	}
}

func (p *Printer) printHTML(qsFrom, qsTo *QueryStats) {
	var data [][]string

	if qsTo == nil {
		for _, s := range qsFrom.stats {
			data = append(data, p.GenerateLine(s))
		}
	} else {
		for _, to := range qsTo.stats {
			from := findQueryStatFrom(qsFrom, to)

			if from == nil {
				data = append(data, p.GenerateLine(to))
			} else {
				data = append(data, p.GenerateLineWithDiff(from, to))
			}
		}
	}
	content, _ := html.RenderTableWithGridJS("slp", p.headers, data, p.printOptions.paginationLimit)
	fmt.Println(content)
}
