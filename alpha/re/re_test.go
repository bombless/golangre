package re
import "testing"
func Test(t *testing.T){
    var tests = []struct{
        re string
        matches, mismatches []string
    }{
        {
            `Marry|has|a|little|goat`,
            []string{ "Marry", "has", "a", "little", "goat" },
            []string{ "Nothing", "lasts", "forever" },
        },
        {
            `[_a-zA-Z][_a-zA-Z0-9]*`,
            []string{ "var", "foo", "base64" },
            []string{ "one space", "1mystery" },
        },
        {
            `we|need|parenthesis|and|brackets|\(*|\[*`,
            []string{ "(((", "[[" },
            []string{ ")", "()" },
        },
        {
            `#{,1}([a-f0-9]{6}|[a-f0-9]{3})`,
            []string{ "#a3c113" },
            []string{ "#4d82h4" },
        },
        {
            `[a-z0-9-]{1,}`,
            []string{ "my-title-here" },
            []string{ "my_title_here" },
        },
        {
            `[a-z0-9_.-]{1,}@[0-9a-z.-]{1,}.[a-z.]{2,6}`,
            []string{ "john@doe.com" },
            []string{ "john@doe.something" },
        },
        {
            `(https{,1}://|)[0-9a-z]([0-9a-z-]{1,}[0-9a-z].)*[a-z]{2,6}(/[/a-zA-Z.-]*|)`,
            []string{ "http://net.tutsplus.com/about" },
            []string{ "http://google.com/some/file!.html" },
        },
        {
            `((25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9]).){3}` +
            `(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])`,
            []string{ "73.60.124.136" },
            []string{ "256.7.3.1", "1.3.4" },
        },
    }
    for _, c := range tests{
        fa, err := RegExp(c.re)
        if err != nil{
            t.Error(err)
        }else{
            for want, inputs := range map[bool][]string{
                true: c.matches,
                false: c.mismatches,
            }{
                for _, input := range inputs{
                    got := fa.Test(input)
                    if got != want{
                        t.Errorf("expect %v=RE/%v/.Test(%v), got %v instead",
                            want, c.re, input, got)
                    }
                }
            }
        }
    }
}
func TestFailure(t *testing.T){
    tests := []string{
        "[a-b-c]", "*", "(", ")", "|*", "[|]", "[", "]", "[*]",
    }
    for _, c := range tests{
        _, got := RegExp(c)
        if got == nil{
            t.Error("result for " + c + " should cause failure")
        }
    }
}
