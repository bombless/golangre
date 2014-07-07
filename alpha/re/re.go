package re
import(
    "fmt"
    "errors"
    "reflect"
    )
type TransitionPair struct{
    First CanMatch
    Second int
}
type FiniteAutomachine struct{
    StatusMap map[int][]TransitionPair
    ClosureList map[int][]int
    Final int
}
func singleTransition(c CanMatch)FiniteAutomachine{
    return FiniteAutomachine{map[int][]TransitionPair{0:{{c, 1}}}, map[int][]int{}, 1}
}
func voidTransition()FiniteAutomachine{
    return FiniteAutomachine{map[int][]TransitionPair{}, map[int][]int{}, 0}
}

func(fa FiniteAutomachine)Kleene()FiniteAutomachine{
    if !inArray(fa.ClosureList[fa.Final], 0){
        fa.ClosureList[fa.Final] = append(fa.ClosureList[fa.Final], 0)
    }
    if !inArray(fa.ClosureList[0], fa.Final){
        fa.ClosureList[0] = append(fa.ClosureList[0], fa.Final)
    }
    return fa
}
func(lhs FiniteAutomachine)Pipe(rhs FiniteAutomachine)FiniteAutomachine{
    shiftLhs := 1
    shiftRhs := lhs.Final + 2
    finalLhs := lhs.Final + shiftLhs
    finalRhs := rhs.Final + shiftRhs
    final := rhs.Final + shiftRhs + 1
    statusMap := map[int][]TransitionPair{}
    closureList := map[int][]int{0:{shiftLhs, shiftRhs}, finalLhs:{final}, finalRhs:{final}}
    for from, transits := range lhs.StatusMap{
        pairs := []TransitionPair{}
        for _, pair := range transits{
            pairs = append(pairs, TransitionPair{pair.First, pair.Second + shiftLhs})
        }
        statusMap[from + shiftLhs] = pairs
    }
    for from, transits := range rhs.StatusMap{
        pairs := []TransitionPair{}
        for _, pair := range transits{
            pairs = append(pairs, TransitionPair{pair.First, pair.Second + shiftRhs})
        }
        statusMap[from + shiftRhs] = pairs
    }
    for from, list := range lhs.ClosureList{
        arr := closureList[from + shiftLhs]
        for _, to := range list{
            arr = append(arr, to + shiftLhs)
        }
        closureList[from + shiftLhs] = arr
    }
    for from, list := range rhs.ClosureList{
        arr := closureList[from + shiftRhs]
        for _, to := range list{
            arr = append(arr, to + shiftRhs)
        }
        closureList[from + shiftRhs] = arr
    }
    return FiniteAutomachine{statusMap, closureList, final}
}
func(lhs FiniteAutomachine)Concat(rhs FiniteAutomachine)FiniteAutomachine{
    shift := lhs.Final + 1
    final := rhs.Final + shift
    closureList := lhs.ClosureList
    statusMap := lhs.StatusMap
    for from, list := range rhs.ClosureList{
        arr := []int{}
        for _, to := range list{
            arr = append(arr, to + shift)
        }
        closureList[from + shift] = arr
    }
    closureList[lhs.Final] = append(closureList[lhs.Final], shift)
    for from, mapping := range rhs.StatusMap{
        record := []TransitionPair{}
        for _, pair := range mapping{
            record = append(record, TransitionPair{pair.First, pair.Second + shift})
        }
        statusMap[from + shift] = record
    }
    return FiniteAutomachine{statusMap, closureList, final}
}
func(fa FiniteAutomachine)GetClosures(id int, acc []int)[]int{
    ret := []int{}
    if !inArray(acc, id){
        acc = append(acc, id)
    }
    for _, val := range fa.ClosureList[id]{
        if !inArray(acc, val){
            ret = append(ret, append(fa.GetClosures(val, acc), val)...)
        }
    }   
    return ret
}
func inArray(arr []int, item int)bool{
    for _, val := range arr{
        if val == item{
            return true
        }
    }
    return false
}
func(fa FiniteAutomachine)Test(str string)bool{
    streams := append([]int{0}, fa.GetClosures(0, []int{})...)
    waitForAttach := []int{}
    for _, char := range str{        
        for index := 0; index < len(streams); index += 1{
            match := false
            for _, pair := range fa.StatusMap[streams[index]]{
                if pair.First.Match(char){
                    match = true
                    streams[index] = pair.Second
                    waitForAttach = append(waitForAttach, fa.GetClosures(pair.Second, []int{})...)
                }
            }
            if !match{
            //remove index
                streams = append(streams[:index], streams[index + 1:]...)
                index -= 1
            }
        }
        streams = append(streams, waitForAttach...)
        waitForAttach = []int{}
    }
    return inArray(streams, fa.Final)
}
type CanMatch interface{
    Match(rune)bool
}
type Rune struct{
    Value rune
}
func(this Rune)Match(r rune)bool{
    return this.Value == r
}
func makeRune(r rune)Rune{
    return Rune{r}
}
type EnumClass struct{
    Set map[rune]struct{}
}
func(c EnumClass)Match(r rune)bool{
    _, ret := c.Set[r]
    return ret
}
func makeEnumClass(arr []Rune)EnumClass{
    s := EnumClass{map[rune]struct{}{}}
    for _, v := range arr{
        s.Set[v.Value] = struct{}{}
    }
    return s
}
type RangeClass struct{
    Min rune
    Max rune
}
func(c RangeClass)Match(r rune)bool{
    return r >= c.Min && r <= c.Max
}
func makeRangeClass(min rune, max rune)RangeClass{
    return RangeClass{min, max}
}
type MixedClass struct{
    Collection []CanMatch
}
func (c MixedClass)Match(r rune)bool{
    for _, v := range c.Collection{
        if v.Match(r){
            return true
        }
    }
    return false
}
func makeMixedClass(collection ...CanMatch)MixedClass{
    return MixedClass{collection}
}
type NegativeClass struct{
    Value MixedClass
}
func(c NegativeClass)Match(r rune)bool{
    return !c.Value.Match(r)
}
func makeNegativeClass(c MixedClass)NegativeClass{
    return NegativeClass{c}
}
type Group struct{
    Content []interface{}
}
type handleFunction func(interface{}, []interface{})([]interface{}, error)
type Pipe struct{}
type Kleene struct{}
type GroupStart struct{}
type GroupEnd struct{}
type ClassStart struct{}
type ClassEnd struct{}
func filterEscaping(r rune)interface{}{
    switch r{
    case 't': return '\t'
    case 'n': return '\n'
    case 'r': return '\r'
    case '|': return '|'
    case '*': return '*'
    case '(': return '('
    case ')': return ')'
    case '[': return '['
    case ']': return ']'
    case '\\': return '\\'
    }
    return errors.New("unexpected " + string([]rune{r}))
}
func filterNormal(r rune)interface{}{
    switch r{
    case '*': return Kleene{}
    case '|': return Pipe{}
    case '(': return GroupStart{}
    case ')': return GroupEnd{}
    case '[': return ClassStart{}
    case ']': return ClassEnd{}
    }
    return r
}
func lexing(reg string)([]interface{}, error){
    escaping := false
    ret := []interface{}{}
    for _, v := range reg{
        if escaping{
            escaped := filterEscaping(v)
            if typeName(escaped) == "error"{
                return ret, escaped.(error)
            }
            ret = append(ret, escaped)
            escaping = false
        }else{
            if v == '\\'{
                escaping = true
                continue
            }
            ret = append(ret, filterNormal(v))
        }
    }
    return ret, nil
}
func typeName(i interface{})string{
    switch reflect.TypeOf(i).String(){
    case "re.GroupStart": return "GroupStart"
    case "re.GroupEnd": return "GroupEnd"
    case "re.ClassStart": return "ClassStart"
    case "re.ClassEnd": return "ClassEnd"
    case "re.Group": return "Group"
    case "re.Class": return "Class"
    case "re.Kleene": return "Kleene"
    case "re.Pipe": return "Pipe"
    case "int32": return "rune"
    case "re.Rune": return "Rune"
    case "re.EnumClass": return "EnumClass"
    case "re.RangeClass": return "RangeClass"
    case "re.MixedClass": return "MixedClass"
    case "re.NegativeClass": return "NegativeClass"
    case "*errors.errorString": return "error"
    }
    return "other"
}
func chooseFunction(item interface{})handleFunction{
    switch typeName(item){
    case "GroupStart": return funcGroupStart
    case "GroupEnd": return funcGroupEnd
    case "ClassStart": return funcClassStart
    case "ClassEnd": return funcClassEnd
    case "Kleene": return funcKleene
    case "Pipe": return funcPipe
    case "rune": return funcRune
    }
    return funcError
}
func funcError(item interface{}, stack []interface{})([]interface{}, error){
    return stack, errors.New("unexpected input")
}
func funcGroupStart(item interface{}, stack []interface{})([]interface{}, error){
    return append(stack, GroupStart{}), nil
}
func funcGroupEnd(item interface{}, stack []interface{})([]interface{}, error){
    i := len(stack) - 1
    for i >= 0 && typeName(stack[i]) != "GroupStart"{
        i -= 1
    }
    if i < 0{
        return stack, errors.New("unexpected GroupEnd")
    }
    pack := []interface{}{}
    for j := i + 1; j < len(stack); j += 1{
        name := typeName(stack[j])
        if name == "ClassStart" || name == "ClassEnd" || name == "GroupEnd"{
            return stack, errors.New(fmt.Sprintf("unexpected %v", name))
        }
        pack = append(pack, stack[j])
    }
    stack = stack[:i]
    if len(pack) == 0{
        return stack, errors.New("Empty group not allowed")
    }
    return append(stack, Group{pack}), nil
}
func funcClassStart(item interface{}, stack []interface{})([]interface{}, error){
    return append(stack, ClassStart{}), nil
}
func funcClassEnd(item interface{}, stack []interface{})([]interface{}, error){
    i := len(stack) - 1
    for i >= 0 && typeName(stack[i]) != "ClassStart"{
        i -= 1
    }
    if i < 0{
        return stack, errors.New("unexpected ClassEnd")
    }
    pack := []CanMatch{}
    for j := i + 1; j < len(stack); j += 1{
        name := typeName(stack[j])
        if name != "Rune"{
            return stack, errors.New(fmt.Sprintf("unexpected %v", name))
        }
        r := stack[j].(Rune)
        if r.Value == '-' && j > i + 1 && j < len(stack) - 1{
            if j - 2 > i && typeName(stack[j - 2]) == "Rune" && stack[j - 2].(Rune).Value == '-'{
                notice := string([]rune{stack[j - 2].(Rune).Value, stack[j - 1].(Rune).Value})
                return stack,
                    errors.New("unexpected `-` after " + notice)
            }
            min, max := stack[j - 1].(Rune), stack[j + 1].(Rune)
            if min.Value > max.Value{
                return stack,
                    errors.New(fmt.Sprintf("value of %c bigger than %c", min.Value, max.Value))
            }
            pack[len(pack) - 1] = makeRangeClass(min.Value, max.Value)
            j += 1
        }else{
            pack = append(pack, r)
        }
    }
    stack = stack[:i]
    needNegative := typeName(pack[0]) == "Rune" && pack[0].(Rune).Value == '^'
    if needNegative{
        pack = pack[1:]
    }
    rangeCollection := []CanMatch{}
    runeCollection := []CanMatch{}
    for _, v := range pack{
        if typeName(v) == "Rune"{
            runeCollection = append(runeCollection, v)
        }else{
            rangeCollection = append(rangeCollection, v)
        }
    }
    mixedClass := makeMixedClass(append(rangeCollection, runeCollection...)...)
    if needNegative{
        return append(stack, makeNegativeClass(mixedClass)), nil
    }
    return append(stack, mixedClass), nil    
}
func funcKleene(item interface{}, stack []interface{})([]interface{}, error){
    return append(stack, Kleene{}), nil
}
func funcPipe(item interface{}, stack []interface{})([]interface{}, error){
    return append(stack, Pipe{}), nil
}
func funcRune(item interface{}, stack []interface{})([]interface{}, error){
    return append(stack, makeRune(item.(rune))), nil
}
func compile(seq []interface{})([]interface{}, error){
    if len(seq) == 0{
        return []interface{}{}, nil
    }
    stack := []interface{}{}
    for _, item := range seq{
        fun := chooseFunction(item)
        var err error
        stack, err = fun(item, stack)
        if err != nil{
            return []interface{}{}, err
        }
    }
    for _, v := range stack{
        name := typeName(v)
        if name == "ClassStart" || name == "ClassEnd" ||
            name == "GroupStart" || name == "GroupEnd"{
            return []interface{}{}, errors.New(fmt.Sprintf("unexpected %v", name))
        }
    }
    return stack, nil
}
func(p Pipe)String()string{
    return "|"
}
func(k Kleene)String()string{
    return "*"
}
func(s GroupStart)String()string{
    return "GroupStart"
}
func(e GroupEnd)String()string{
    return "GroupEnd"
}
func(s ClassStart)String()string{
    return "ClassStart"
}
func(e ClassEnd)String()string{
    return "ClassEnd"
}
func construct(seq []interface{})(FiniteAutomachine, error){
    if len(seq) == 0{
        return voidTransition(), nil
    }else if len(seq) == 1{
        switch typeName(seq[0]){
            case "Group": return construct(seq[0].(Group).Content)
            case "MixedClass": return singleTransition(seq[0].(MixedClass)), nil
            case "NegativeClass": return singleTransition(seq[0].(NegativeClass)), nil
            case "Rune": return singleTransition(seq[0].(Rune)), nil
        }
        return FiniteAutomachine{}, errors.New(fmt.Sprintf("unexpected %v", typeName(seq[0])))
    }
    pipe := -1
    for i, v := range seq{
        if typeName(v) == "Pipe"{
            pipe = i
            break
        }
    }
    if pipe > -1{
        left, err := construct(seq[:pipe])
        if err != nil{
            return FiniteAutomachine{}, err
        }
        right, err := construct(seq[(pipe + 1):])
        if err != nil{
            return FiniteAutomachine{}, err
        }
        return left.Pipe(right), nil
    }
    left, err := construct(seq[:1])
    if err != nil{
        return FiniteAutomachine{}, err
    }
    var right FiniteAutomachine
    if typeName(seq[1]) == "Kleene"{
        left = left.Kleene()
        right, err = construct(seq[2:])
    }else{
        right, err = construct(seq[1:])
    }
    if err != nil{
        return FiniteAutomachine{}, err
    }
    return left.Concat(right), nil
}
func RegExp(reg string)(FiniteAutomachine, error){
    seq, err := lexing(reg)
    if err != nil{
        return FiniteAutomachine{}, err
    }
    seq, err = compile(seq)
    if err != nil{
        return FiniteAutomachine{}, err
    }
    return construct(seq)
}
