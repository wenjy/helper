package helper

import "testing"

var (
	v1    = "2.0.1"
	v1_1  = "2.0.1.1"
	v11   = "v2.0.1"
	v12   = "v2.0.2"
	v12_1 = "v2.0.1.1"
	v2    = "2.1.1"
	v3    = "2.10.22.1"
	v4    = "2.09.22.2"
	v5    = "2.9.22.2"
)

// go test -test.bench=".*"
func BenchmarkVersion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		VersionCompare("2.3.1", "2.1.3.4", ">=")
	}
}

// 该函数比较两个版本号是否相等，是否大于或小于的关系
// 返回值：0表示v1与v2相等；1表示v1大于v2；2表示v1小于v2
func TestCompare(t *testing.T) {
	if 0 != Compare(v1, v1) {
		t.Errorf("v1 %s == v1 %s ", v1, v1)
	}

	if 1 != Compare(v1_1, v1) {
		t.Errorf("v1_1 %s == v1 %s ", v1_1, v1)
	}

	if 0 != Compare(v1, v11) {
		t.Errorf("v1 %s == v11 %s ", v1, v11)
	}

	if 1 != Compare(v2, v1) {
		t.Errorf("v2 %s > v1 %s", v2, v1)
	}

	if 1 != Compare(v12, v12_1) {
		t.Errorf("v12 %s > v12_1 %s", v12, v12_1)
	}

	if 2 != Compare(v11, v12) {
		t.Errorf("v11 %s > v2 %s", v11, v12)
	}

	if 2 != Compare(v1, v2) {
		t.Errorf("v1 %s < v2 %s", v1, v2)
	}

}

func TestVersionCompare(t *testing.T) {
	if !VersionCompare(v1, v1_1, "<") {
		t.Errorf("v1 %s < v1_1 %s ", v1, v1_1)
	}

	if !VersionCompare(v2, v1_1, ">") {
		t.Errorf("v2 %s < v1_1 %s ", v2, v1_1)
	}

	if !VersionCompare(v2, v1, ">=") {
		t.Errorf("v2 %s >= v1 %s ", v2, v1)
	}

	if !VersionCompare(v1, v2, "<=") {
		t.Errorf("v1 %s >= v2 %s ", v1, v2)
	}

	if !VersionCompare(v1, v1, "==") {
		t.Errorf("v1 %s == v1 %s ", v1, v1)
	}
	if !VersionCompare(v3, v4, ">") {
		t.Errorf("v3 %s > v4 %s ", v3, v4)
	}
	if !VersionCompare(v3, v5, ">") {
		t.Errorf("v3 %s > v5 %s ", v3, v5)
	}

	if !VersionCompare(v4, v5, "==") {
		t.Errorf("v4 %s == v5 %s ", v4, v5)
	}

}
