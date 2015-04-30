package benchmark

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestReadWriteFile(t *testing.T) {
	now := time.Now().UnixNano()
	results := Results{
		Result{Status: 0, Start: now, End: now + 1, Info: "a"},
		Result{Status: 1, Start: now, End: now + 2, Info: "b"},
		Result{Status: 0, Start: now, End: now + 3, Info: "a"},
		Result{Status: 1, Start: now, End: now + 4, Info: "b"},
		Result{Status: 0, Start: now, End: now + 5, Info: "a"},
	}

	file, _ := ioutil.TempFile(os.TempDir(), "bench_")
	defer os.Remove(file.Name())

	results.Save(file.Name())
	dst, _ := LoadResultsFromFile(file.Name())

	if !reflect.DeepEqual(results, dst) {
		t.Error("load err: %v", dst)
	}

}

func TestGroupingByStatus(t *testing.T) {
	now := time.Now().UnixNano()

	results := Results{
		Result{Status: 0, Start: now, End: now + 1},
		Result{Status: 1, Start: now, End: now + 2},
		Result{Status: 0, Start: now, End: now + 3},
		Result{Status: 1, Start: now, End: now + 4},
		Result{Status: 0, Start: now, End: now + 5},
	}
	ret := results.GroupByStatus()

	if !reflect.DeepEqual(ret[0], Results{results[0], results[2], results[4]}) {
		t.Error("slice: %v", ret[0])
	}

	if !reflect.DeepEqual(ret[1], Results{results[1], results[3]}) {
		t.Error("slice: %v", ret[1])
	}
}

func TestGroupingByInfo(t *testing.T) {
	now := time.Now().UnixNano()

	results := Results{
		Result{Status: 0, Start: now, End: now + 1, Info: "a"},
		Result{Status: 1, Start: now, End: now + 2, Info: "b"},
		Result{Status: 0, Start: now, End: now + 3, Info: "a"},
		Result{Status: 1, Start: now, End: now + 4, Info: "b"},
		Result{Status: 0, Start: now, End: now + 5, Info: "a"},
	}
	ret := results.GroupByInfo()

	if !reflect.DeepEqual(ret["a"], Results{results[0], results[2], results[4]}) {
		t.Error("slice: %v", ret["a"])
	}

	if !reflect.DeepEqual(ret["b"], Results{results[1], results[3]}) {
		t.Error("slice: %v", ret["b"])
	}
}

func TestStat(t *testing.T) {
	data := make(Results, 1000)
	for i := 0; i < 1000; i++ {
		now := time.Now().UnixNano()
		data[i] = Result{Status: 0, Start: now, End: now + (int64(i)*37)%1000*1000}
	}

	s := data.CalcStat()
	if s.Max != 999*1000 {
		t.Errorf("max: %d", s.Max)
	}
	if s.Min != 0 {
		t.Errorf("min: %d", s.Min)
	}
	if s.Num != 1000 {
		t.Errorf("num: %d", s.Num)
	}
	if s.Mean != 999*1000/2 {
		t.Errorf("mean: %d", s.Mean)
	}
	if s.Th50 != 500*1000 {
		t.Errorf("th50: %d", s.Th50)
	}
	if s.Th95 != 950*1000 {
		t.Errorf("th95: %d", s.Th95)
	}
	if s.Th99 != 990*1000 {
		t.Errorf("th99: %d", s.Th99)
	}

}
