package main

import(
    "fmt"
    "github.com/bombless/golangre/alpha/re"
    )
func testRE(reg string, expect bool, str ...string){
    fmt.Printf("tesing RE/%v/\n", reg)
    fa, err := re.RegExp(reg)
    if nil != err{
        fmt.Println(err)
        return
    }
    for _, val := range str{
        display := val
        if(val == ""){
            display = "<Empty String>"
        }
        te := fa.Test(val)
        if te == expect{
            fmt.Printf("passed: %v=RE/%v/.Test(%v)\n", expect, reg, display)
        }else{
            fmt.Printf("error: expect %v=RE/%v/.Test(%v), got %v instead\n", expect, reg, display, te)
        }
    }
}

func main(){
    testRE("a|b", true, "a", "b")
    testRE("a|b", false, "", "ab", "c", "aa")
    testRE("[ab]|b", true, "a", "b")
    testRE("[ab]|b", false, "", "ab", "c", "aa")
    testRE("[ac]|b", true, "a", "b", "c")
    testRE("[ac]|b", false, "", "ab", "aa")
    testRE("[ab]*", true, "", "a", "b", "ab", "aab", "aaaaa")
    testRE("[ab]*", false, "abc", "aac")
    testRE("ab", true, "ab")
    testRE("ab", false, "", "a", "b", "c")
    testRE("abc", true, "abc")
    testRE("abc", false, "", "a")
    testRE("(abc)*", true, "", "abc", "abcabcabc")
    testRE("(abc)*", false, "a", "c", "cba")
    testRE("[abc]*([de]|fg)*h", true, "acbdh", "dh", "fgfgh")
    testRE("[abc]*([de]|fg)*h", false, "acbgdh", "dfh")
    testRE("[^cat]*", true, "user", "feels", "good")
    testRE("[^cat]*", false, "little", "cat")
    testRE("dog|cat|fish", true, "dog", "cat", "fish")
    testRE("dog|cat|fish", false, "", "orange", "apple")
    testRE("[_a-zA-Z][_a-zA-Z0-9]*", true, "one", "_two")
    testRE("[_a-zA-Z][_a-zA-Z0-9]*", false, "1apple", "one space")
}
