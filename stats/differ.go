package stats

import "fmt"

type Differ struct {
	From *QueryStat
	To   *QueryStat
}

func NewDiffer(from, to *QueryStat) *Differ {
	return &Differ{
		From: from,
		To:   to,
	}
}

func (d *Differ) DiffCnt() string {
	v := d.To.Cnt - d.From.Cnt
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

// query time
func (d *Differ) DiffMaxQueryTime() string {
	v := d.To.MaxQueryTime() - d.From.MaxQueryTime()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffMinQueryTime() string {
	v := d.To.MinQueryTime() - d.From.MinQueryTime()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffSumQueryTime() string {
	v := d.To.SumQueryTime() - d.From.SumQueryTime()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffAvgQueryTime() string {
	v := d.To.AvgQueryTime() - d.From.AvgQueryTime()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffPNQueryTime(n int) string {
	v := d.To.PNQueryTime(n) - d.From.PNQueryTime(n)
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

// lock time
func (d *Differ) DiffMaxLockTime() string {
	v := d.To.MaxLockTime() - d.From.MaxLockTime()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffMinLockTime() string {
	v := d.To.MinLockTime() - d.From.MinLockTime()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffSumLockTime() string {
	v := d.To.SumLockTime() - d.From.SumLockTime()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffAvgLockTime() string {
	v := d.To.AvgLockTime() - d.From.AvgLockTime()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffPNLockTime(n int) string {
	v := d.To.PNLockTime(n) - d.From.PNLockTime(n)
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

// rows sent
func (d *Differ) DiffMaxRowsSent() string {
	v := d.To.MaxRowsSent() - d.From.MaxRowsSent()
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

func (d *Differ) DiffMinRowsSent() string {
	v := d.To.MinRowsSent() - d.From.MinRowsSent()
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

func (d *Differ) DiffSumRowsSent() string {
	v := d.To.SumRowsSent() - d.From.SumRowsSent()
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

func (d *Differ) DiffAvgRowsSent() string {
	v := d.To.AvgRowsSent() - d.From.AvgRowsSent()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffPNRowsSent(n int) string {
	v := d.To.PNRowsSent(n) - d.From.PNRowsSent(n)
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

// rows examined
func (d *Differ) DiffMaxRowsExamined() string {
	v := d.To.MaxRowsExamined() - d.From.MaxRowsExamined()
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

func (d *Differ) DiffMinRowsExamined() string {
	v := d.To.MinRowsExamined() - d.From.MinRowsExamined()
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

func (d *Differ) DiffSumRowsExamined() string {
	v := d.To.SumRowsExamined() - d.From.SumRowsExamined()
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

func (d *Differ) DiffAvgRowsExamined() string {
	v := d.To.AvgRowsExamined() - d.From.AvgRowsExamined()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffPNRowsExamined(n int) string {
	v := d.To.PNRowsExamined(n) - d.From.PNRowsExamined(n)
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

// rows affected
func (d *Differ) DiffMaxRowsAffected() string {
	v := d.To.MaxRowsAffected() - d.From.MaxRowsAffected()
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

func (d *Differ) DiffMinRowsAffected() string {
	v := d.To.MinRowsAffected() - d.From.MinRowsAffected()
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

func (d *Differ) DiffSumRowsAffected() string {
	v := d.To.SumRowsAffected() - d.From.SumRowsAffected()
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

func (d *Differ) DiffAvgRowsAffected() string {
	v := d.To.AvgRowsAffected() - d.From.AvgRowsAffected()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffPNRowsAffected(n int) string {
	v := d.To.PNRowsAffected(n) - d.From.PNRowsAffected(n)
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

// bytes sent
func (d *Differ) DiffMaxBytesSent() string {
	v := d.To.MaxBytesSent() - d.From.MaxBytesSent()
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

func (d *Differ) DiffMinBytesSent() string {
	v := d.To.MinBytesSent() - d.From.MinBytesSent()
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

func (d *Differ) DiffSumBytesSent() string {
	v := d.To.SumBytesSent() - d.From.SumBytesSent()
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

func (d *Differ) DiffAvgBytesSent() string {
	v := d.To.AvgBytesSent() - d.From.AvgBytesSent()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffPNBytesSent(n int) string {
	v := d.To.PNBytesSent(n) - d.From.PNBytesSent(n)
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

func DiffCountAll(from, to map[string]int) map[string]string {
	counts := make(map[string]string, 1)
	keys := []string{"count"}

	for _, key := range keys {
		v := to[key] - from[key]
		if v >= 0 {
			counts[key] = fmt.Sprintf("+%d", v)
		} else {
			counts[key] = fmt.Sprintf("-%d", v)
		}
	}

	return counts
}
