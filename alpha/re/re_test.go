package re
import "testing"
func Test(t *testing.T){
    var tests = []struct{
        re string
        want bool
        inputs []string
    }{
        {
            "Marry|has|a|little|goat",
            true,
            []string{ "Marry", "has", "a", "little", "goat" },
        },
        {
            "Marry|has|a|little|goat",
            false,
            []string{ "Nothing", "lasts", "forever" },
        },
        {
            "[_a-zA-Z][_a-zA-Z0-9]*",
            true,
            []string{ "var", "foo", "base64" },
        },
        {
            "[_a-zA-Z][_a-zA-Z0-9]*",
            false,
            []string{ "one space", "1mystery" },
        },
    }
    for _, c := range tests{
        fa, err := RegExp(c.re)
        if err != nil{
            t.Error(err)
        }else{
            want := c.want
            for _, input := range c.inputs{
                got := fa.Test(input)
                if got != want{
                    t.Errorf("expect %v=RE/%v/.Test(%v), got %v instead",
                        want, c.re, input, got)
                }
            }
        }
    }
}
