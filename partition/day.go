package partition

import (
	"os"
	"time"

	"github.com/seefan/goerr"
	"github.com/seefan/ssdbproxy/conf"
)

const DayFormat = "20060102"

//按天进行分区
type DayPartition struct {
	config *conf.ProxyConf
}

func NewDayPartition(c *conf.ProxyConf) *DayPartition {
	return &DayPartition{
		config: c,
	}
}

func (d *DayPartition) Get(day string) (string, error) {
	date, err := time.Parse(d.config.Partition.Pattern, day)
	if err != nil {
		return "", goerr.NewError(err, "key is inconsistent date format")
	}
	return date.Format(DayFormat), nil
}

func (d *DayPartition) getKey(day time.Time) string {
	return day.Format(DayFormat)
}
func (d *DayPartition) getAllKey(end time.Time, limit int) (re []string) {
	for i := 0; i < limit; i++ {
		now := end.AddDate(0, 0, -i)
		re = append(re, d.getKey(now))
	}
	return
}
func (d *DayPartition) Clean(limit int, dayKeys []os.FileInfo) (re []string) {
	dayRange := make(map[string]interface{})
	allKey := d.getAllKey(time.Now().AddDate(0, 0, 1), limit)
	for _, k := range allKey {
		dayRange[k] = nil
	}
	for _, day := range dayKeys {
		if _, ok := dayRange[day.Name()]; !ok && day.IsDir() && day.Name() != "." && day.Name() != ".." {
			re = append(re, day.Name())
		}
	}
	return
}
