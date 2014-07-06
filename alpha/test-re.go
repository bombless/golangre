package main

import(
    "fmt"
    "./re"
    )
func testRE(reg string, expect bool, str ...string){
    fa, err := re.RegExp(reg)
    if nil != err{
        fmt.Println(err)
        return
    }
    for _, val := range str{
        if(val == ""){
            fmt.Print("testing <Empty String>...")
        }else{
            fmt.Print("testing ", val, "... ")
        }
        te := fa.Test(val)
        if te == expect{
            fmt.Println("passed")
        }else{
            fmt.Println("error: expect ", expect, ", got ", te, " instead")
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
}
