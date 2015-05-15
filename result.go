package benchmark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"time"
)

func (r Results) Save(file string) error {
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, b, 0644)
}

func LoadResultsFromFile(file string) (r Results, err error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &r)
	return
}

func (r Results) GroupByStatus() map[Status]Results {
	counter := make(map[Status]int)
	for _, v := range r {
		counter[v.Status] += 1
	}

	ret := make(map[Status]Results)
	for k, v := range counter {
		ret[k] = make(Results, 0, v)
	}

	for _, v := range r {
		ret[v.Status] = append(ret[v.Status], v)
	}

	return ret
}

func (r Results) GroupByInfo(key string) map[string]Results {
	counter := make(map[string]int)
	for _, v := range r {
		counter[v.Info[key]] += 1
	}

	ret := make(map[string]Results)
	for k, v := range counter {
		ret[k] = make(Results, 0, v)
	}

	for _, v := range r {
		ret[v.Info[key]] = append(ret[v.Info[key]], v)
	}

	return ret
}

type Stat struct {
	Mean, Th50, Th95, Th99, Max, Min int64
	Num                              int
	Throughput                       float64
}

func (s *Stat) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("count   %d\n", s.Num))
	buf.WriteString(fmt.Sprintf("latency[min]  %s\n", formatDuration(s.Min)))
	buf.WriteString(fmt.Sprintf("latency[max]  %s\n", formatDuration(s.Max)))
	buf.WriteString(fmt.Sprintf("latency[mean] %s\n", formatDuration(s.Mean)))
	buf.WriteString(fmt.Sprintf("latency[50%%]  %s\n", formatDuration(s.Th50)))
	buf.WriteString(fmt.Sprintf("latency[95%%]  %s\n", formatDuration(s.Th95)))
	buf.WriteString(fmt.Sprintf("latency[99%%]  %s\n", formatDuration(s.Th99)))
	buf.WriteString(fmt.Sprintf("throughput    %f\n", s.Throughput))

	return buf.String()
}

func (r Results) CalcStat() *Stat {
	num := len(r)
	if num == 0 {
		return &Stat{}
	}

	tmp := make([]int64, num)
	var total int64 = 0
	start := r[0].Start
	end := r[0].End
	for i, r := range r {
		d := r.End - r.Start
		tmp[i] = d
		total += d

		if start > r.Start {
			start = r.Start
		}
		if end < r.End {
			end = r.End
		}
	}

	sort.Sort(int64Slice(tmp))

	return &Stat{
		Mean:       total / int64(num),
		Th50:       tmp[num/2],
		Th95:       tmp[num*95/100],
		Th99:       tmp[num*99/100],
		Max:        tmp[num-1],
		Min:        tmp[0],
		Num:        num,
		Throughput: float64(num) * 1e9 / float64(end-start),
	}
}

func (r Results) FormatByStatus() string {
	var buf bytes.Buffer

	total := len(r)
	for status, r := range r.GroupByStatus() {
		s := r.CalcStat()
		buf.WriteString(fmt.Sprintf("status  %d\n", status))
		buf.WriteString(fmt.Sprintf("  count   %d/%d[%.2f%%]\n", s.Num, total, float64(s.Num)*100.0/float64(total)))
		buf.WriteString(fmt.Sprintf("  latency[min]  %s\n", formatDuration(s.Min)))
		buf.WriteString(fmt.Sprintf("  latency[max]  %s\n", formatDuration(s.Max)))
		buf.WriteString(fmt.Sprintf("  latency[mean] %s\n", formatDuration(s.Mean)))
		buf.WriteString(fmt.Sprintf("  latency[50%%]  %s\n", formatDuration(s.Th50)))
		buf.WriteString(fmt.Sprintf("  latency[95%%]  %s\n", formatDuration(s.Th95)))
		buf.WriteString(fmt.Sprintf("  latency[99%%]  %s\n", formatDuration(s.Th99)))
		buf.WriteString(fmt.Sprintf("  throughput    %f\n", s.Throughput))
	}

	return buf.String()
}

func formatDuration(d int64) string {
	str := time.Duration(d).String()
	pos := strings.IndexAny(str, ".smhnu")

	if pos < 5 {
		return strings.Repeat(" ", 5-pos) + str
	} else {
		return str
	}
}

type int64Slice []int64

func (a int64Slice) Len() int           { return len(a) }
func (a int64Slice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a int64Slice) Less(i, j int) bool { return a[i] < a[j] }
