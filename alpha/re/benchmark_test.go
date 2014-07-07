package re
import "testing"
func BenchmarkBasicCompiling(b *testing.B){
    re := "[_a-zA-Z][_a-zA-Z0-9]*"
    for i := 0; i < b.N; i += 1{
            RegExp(re)
    }
}
func BenchmarkBasicMatching(b *testing.B){
    fa, _ := RegExp("[_a-zA-Z][_a-zA-Z0-9]*")
    tests := []string{
        "dogs and cats",
        "should_not_pass",
    }
    for _, v := range tests{
        for i := 0; i < b.N; i += 1{
            _ = fa.Test(v)
        }
    }
}
