package logstorage

import (
	"fmt"
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
			tokens = tokenizeStringsNew(tokens[:0], a)
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
			tokens = tokenizeStringsNew(tokens[:0], a)
		}
	})
}

//旧版本 BenchmarkTokenizeStrings-12    	  178954	      7396 ns/op	 456.44 MB/s	       0 B/op	       0 allocs/op
//新版本 BenchmarkTokenizeStrings-12    	  284407	      4331 ns/op	 779.42 MB/s	       0 B/op	       0 allocs/op
//      BenchmarkTokenizeStrings-12    	  214754	      5746 ns/op	 582.66 MB/s	     107 B/op	       4 allocs/op
//      BenchmarkTokenizeStrings-12    	  143713	      9408 ns/op	 355.87 MB/s	       1 B/op	       0 allocs/op

func TestTokenize(t *testing.T) {
	tokens := make(map[string]struct{}, 140)
	tokenizeString(tokens, "2023-11-27 12:18:05:392 INFO [<1503037114~87>] |\"com.jd.trade2.horizontal.biz.chinamain.domain.points.infrastructure.rpc.PointsRpcManagerImpl|getCalculatePoints-in\"|\"vielen123\"|\"#getCalculatePoints-积分规则计算入参in#\"|[{\"beanBalance\":13115,\"channelId\":10000,\"enjoyValue\":7698,\"orderAmount\":142.80,\"pin\":\"vielen123\",\"siteId\":301,\"skuInfoList\":[{\"properties\":{\"companyType\":\"0\",\"vender_bizid\":\"\",\"vender_col_type\":\"0\",\"categoryId\":\"11932\"},\"skuid\":\"10088665964982\",\"uuid\":\"1012_F2sb4w1285936148417642496\"},{\"properties\":{\"companyType\":\"0\",\"vender_bizid\":\"\",\"vender_col_type\":\"0\",\"categoryId\":\"11932\"},\"skuid\":\"10088665964984\",\"uuid\":\"1012_F2F1Vc1285937474951667712\"},{\"properties\":{\"companyType\":\"0\",\"vender_bizid\":\"\",\"vender_col_type\":\"0\",\"categoryId\":\"11932\"},\"skuid\":\"10088665964985\",\"uuid\":\"1012_F1b1h1v1285937504792027136\"},{\"properties\":{\"companyType\":\"0\",\"vender_bizid\":\"\",\"vender_col_type\":\"0\",\"categoryId\":\"11932\"},\"skuid\":\"10088665964981\",\"uuid\":\"1012_F1v2q3L1285937540369256448\"},{\"properties\":{\"companyType\":\"0\",\"vender_bizid\":\"\",\"vender_col_type\":\"0\",\"categoryId\":\"11932\"},\"skuid\":\"10088665964988\",\"uuid\":\"1012_F1U4AD1285936258924400640\"}],\"userInfoFlag\":\"0000000000000000201000000000010900100100000000003300201000000190010000000000000000000000000000000000\"}]")
	fmt.Printf("tokens:%d\n", len(tokens))
	for s, _ := range tokens {
		fmt.Printf("token:%s\n", s)
	}

}
