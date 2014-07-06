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
type PositiveClass struct{
    Set map[rune]struct{}
}
func(c PositiveClass)Match(r rune)bool{
    _, ret := c.Set[r]
    return ret
}
func makePositiveClass(arr []Rune)PositiveClass{
    s := PositiveClass{map[rune]struct{}{}}
    for _, v := range arr{
        s.Set[v.Value] = struct{}{}
    }
    return s
}
type NegativeClass struct{
    Set map[rune]struct{}
}
func(c NegativeClass)Match(r rune)bool{
    _, ret := c.Set[r]
    return !ret
}
func makeNegativeClass(arr []Rune)NegativeClass{
    s := NegativeClass{map[rune]struct{}{}}
    for _, v := range arr{
        s.Set[v.Value] = struct{}{}
    }
    return s
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
func filtering(reg string)([]interface{}, error){
    escaping := false
    ret := []interface{}{}
    for _, v := range reg{
        if escaping{
            escaped := filterEscaping(v)
            if reflect.TypeOf(escaped).String() == "*errors.errorString"{
                return []interface{}{}, escaped.(error)
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
    case "re.PositiveClass": return "PositiveClass"
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
    pack := []Rune{}
    for j := i + 1; j < len(stack); j += 1{
        name := typeName(stack[j])
        if name != "Rune"{
            fmt.Printf("#line300, %v\n", stack)
            return stack, errors.New(fmt.Sprintf("unexpected %v", name))
        }
        pack = append(pack, stack[j].(Rune))
    }
    stack = stack[:i]
    var c CanMatch
    if len(pack) == 0{
        return stack, errors.New("Empty positive class not allowed")
    }else if pack[0].Value != '^'{
            c = makePositiveClass(pack)
    }else if len(pack) == 1{
        return stack, errors.New("Empty  negative class not allowed")
    }else{
        c = makeNegativeClass(pack[1:])
    }
    return append(stack, c), nil    
}
func funcKleene(item interface{}, stack []interface{})([]interface{}, error){
    for _, v := range stack{
        name := typeName(v)
        if name == "ClassStart" || name == "ClassEnd"{
            return stack, errors.New("can not have kleene in class")
        }
    }
    return append(stack, Kleene{}), nil
}
func funcPipe(item interface{}, stack []interface{})([]interface{}, error){
    for _, v := range stack{
        name := typeName(v)
        if name == "ClassStart" || name == "ClassEnd"{
            return stack, errors.New("can not have pipe in class")
        }
    }
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
func handle(reg string)([]interface{}, error){
    seq, err := filtering(reg)
    if err != nil{
        return []interface{}{}, err
    }
    return compile(seq)
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
            case "PositiveClass": return singleTransition(seq[0].(PositiveClass)), nil
            case "NegativeClass": return singleTransition(seq[0].(NegativeClass)), nil
            case "Rune": return singleTransition(seq[0].(Rune)), nil
        }
        return FiniteAutomachine{}, errors.New(fmt.Sprintf("unexpected %v", typeName(seq[0])))
    }
    pipe := -1
    for i, v := range seq{
        if typeName(v) == "Pipe"{
            pipe = i
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
    seq, err := handle(reg)
    if err != nil{
        return FiniteAutomachine{}, err
    }
    return construct(seq)
}
