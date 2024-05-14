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
	tokenizeString(set, "{\"#METHOD\":\"GET\",\"#PATH\":\"/\",\"@FQON\":\"<ON>/?image_search=1&url=jfs/t20241008/237360/20/16226/70039/661737d0Fd9b16e33/e62800956c1486e6.jpg&is_first_query=0&trigger_query_type=5&is_user_rectangle=1&main_body_rectangle=289.632538,532.459778|388.045807,847.70575&gz=yes&aigc_type=1&merge_sku=yes&scene=0&log_id=867fc8a9d8a54112b5ec1b613e07a4d2&pvid=867fc8a9d8a54112b5ec1b613e07a4d2&uuid=b819c8c4a3f9b3df&qp_exclude=24&multi_suppliers=yes&shop_col=yes&area_ids=28,2495,2496,54832&user_pin=jd_jVJUzXmFdCpl&d_brand=OPPO&d_model=PFDM00&page=1&pagesize=10&ranktop=yes&jd_live=yes&brand_benefit=1&album_insert=yes&jd_jx_insert=yes&jd_hour_insert=yes&ffshop=yes&ico_col=yes&real_insert_valid=true&entity_merge_sku=yes&qp_disable=no&qar=no&client=1698305559501&collect_flower_store=no&videotab_valid=yes&filt_type=redisstore,1;productext2,b15v0;not_medical_service;lbs_exclude,19;productext3,b25v0;sku_state,1&expression_key=&delivertime=yes&delivery_time_col=yes&auction=yes&shop_multi=yes&pictag_requery_enable=yes&device_hash_value=10&deviceId=b819c8c4a3f9b3df&lightspeed=yes&jddj=2&platCode=android&jdblg=yes&appsource=h5&osVersion=13&resolution=2161*1080&top_store=yes&networkType=UNKNOWN&urlencode=yes&charset=utf8&uemps=0-0-0&brand_col=yes&price_col=yes&color_col=yes&size_col=yes&ext_attr=yes&ext_attr_sort=yes&new_brand_val=yes&oldware=yes&exist_col=yes&real_estate=yes&pers=1&channel=1&client_type=app&ignore_filt_type=;lbs_exclude,12&travel_col=yes&jdcs_wl=1&eid=eidAa03d812099senm5SeTvZQeW3+GJyIIOrJ55uVAReEP64pAZoUK4x3k5PubCYZR0DH1JmO1pYAGgO7xRG8aKsNWfJ0A1PHL7J6T7saR3AvPdeuEMw&clientChannel=2&securityToken=JD012145b9mnX1n0ztId171279762690405S6nhy48HhBmDREtkMkk4dzLazWF93aFEyxm_bYAJVe0Nlfdh36WAu23rqmVIGdgG_JRralBShlL5Ji93dk11DA1b8kxy1~BApXeVia5yetC9q4Q8kfutLLJaXe7fqgS1eG4tPBX9xJ1PdZfQozguC_xtCbcN6Vjf7t5e6Tn&clientIP=2408:8474:1700:31ab::1&search_differentiation=&bybt_traffic=1\",\"#CODE\":\"200\",\"#HTTP_RTO\":\"3000\",\"#HTTP_CTO\":\"3000\"}")
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
