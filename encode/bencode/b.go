package bencode



type BType int

const (
    BTypeList BType = iota
    BTypeDict 
    BTypeString
    BTypeInt
    BTypeEnd
)


type BData interface{
    GetVType()( tp BType)
    GetK() string
    GetV() *string
}

type BDict struct{
    Pairs       map[string]BData
}

type BDictPair struct {
    Key         *string
    V            BData
}

type BString struct {
    V          []byte
}

type BInt struct{
    V           int64
}

type BList struct{
    Eles    []BData
}

func (bb BDict)GetVType()(tp BType){
    return BTypeDict
}


func (bb BList)GetVType()(tp BType){
    return BTypeList
}


func (bb BString)GetVType()(tp BType){
    return BTypeString
}

func (bb BInt)GetVType()(tp BType){
    return BTypeInt
}

func (bb BDict)GetK() string{
    str := "\n{\n"

    b := false
    for k, v := range bb.Pairs{
        if b {
            str += ",\n"
        }
        b = true
        str += k
        str += ":"
        str += v.GetK()
    
    }

    str += "\n}\n"

    return str
}



func (bb BDictPair)GetK() string{
    str := *bb.Key
    str += ":"
    str += bb.V.GetK()

    return str
}

func (bb BList)GetK() string{
    return "[LIST]"
}

func (bb BString)GetK() string{
    return "[STRING]"
}

func (bb BInt)GetK() string{
    return "[INT]"
}

func (dict BDict)Get(key string)BData{
    v, ok := dict.Pairs[key]; if !ok{
        return nil
    }

    return v
}

func (dict BDict)GetStr(key string)*string{
    v, ok := dict.Pairs[key]; if !ok{
        return nil
    }

    return v.GetV()
}

func (bb BDict)GetV() *string{
    str := string("dict v")
    return &str
}

func (bb BString)GetV() *string{
    str := "\""
    str += string(bb.V)
    str += "\""
    return &str
}




func (bb BList)GetV() *string{
    str := "["
    bf := false
    for _, bl := range bb.Eles{
        if bf{
            str += ","
        }

        v := bl.GetV()
        str += *v
    }

    str += "]"
    return &str
}



func (bb BInt)GetV() *string{
    str := string("int v")
    return &str
}


