package logstorage

import (
	"fmt"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/logstorage/hashset"
	"strings"
	"testing"
)

func Benchmark_TokenizeStrings(b *testing.B) {
	a := strings.Split(benchLogs, "\n")

	b.ReportAllocs()
	b.SetBytes(int64(len(benchLogs)))
	b.RunParallel(func(pb *testing.PB) {
		var tokens []string
		for pb.Next() {
			tokens = tokenizeStrings(tokens[:0], a)
		}
	})
}

func Benchmark_TokenizeStringsNew(b *testing.B) {
	a := strings.Split(benchLogs, "\n")

	b.ReportAllocs()
	b.SetBytes(int64(len(benchLogs)))
	b.RunParallel(func(pb *testing.PB) {
		var tokens []string
		for pb.Next() {
			tokens = tokenizeStrings(tokens[:0], a)
		}
	})
}

func Benchmark_TokenizeSmall(b *testing.B) {
	a := strings.Split(benchLogs, "\n")[:1]

	b.ReportAllocs()
	b.SetBytes(int64(len(a[0])))
	b.RunParallel(func(pb *testing.PB) {
		var tokens []string
		for pb.Next() {
			tokens = tokenizeStrings(tokens[:0], a)
		}
	})
}

func Benchmark_TokenizeSmall_new(b *testing.B) {
	a := strings.Split(benchLogs, "\n")[:1]

	b.ReportAllocs()
	b.SetBytes(int64(len(a[0])))
	b.RunParallel(func(pb *testing.PB) {
		var tokens []string
		for pb.Next() {
			tokens = tokenizeStrings(tokens[:0], a)
		}
	})
}

//旧版本 BenchmarkTokenizeStrings-12    	  178954	      7396 ns/op	 456.44 MB/s	       0 B/op	       0 allocs/op
//新版本 BenchmarkTokenizeStrings-12    	  284407	      4331 ns/op	 779.42 MB/s	       0 B/op	       0 allocs/op
//      BenchmarkTokenizeStrings-12    	  214754	      5746 ns/op	 582.66 MB/s	     107 B/op	       4 allocs/op
//      BenchmarkTokenizeStrings-12    	  143713	      9408 ns/op	 355.87 MB/s	       1 B/op	       0 allocs/op

func TestTokenize(t *testing.T) {
	//tokens := make(map[string]struct{}, 140)
	set := hashset.GetHashSet()
	tokenizeString(set, "jd_7495f77e623dc 流程编排配置号")
	fmt.Printf("tokens:%d\n", set.Size())
	iter := set.Iterator()
	for iter.Next() {
		fmt.Printf("token:%s\n", iter.At())
	}
	//for s, _ := range set {
	//
	//}
}

func TestLetter(t *testing.T) {
	var result string

	// 拼接小写字母
	for i := 'a'; i <= 'z'; i++ {
		result += string(i)
	}

	// 拼接大写字母
	for i := 'A'; i <= 'Z'; i++ {
		result += string(i)
	}

	// 拼接数字
	for i := '0'; i <= '9'; i++ {
		result += string(i)
	}

	finalResult := fmt.Sprintf("%s", result)
	fmt.Println(finalResult)
}
