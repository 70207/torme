package bencode

import (

    "io"
    "io/ioutil"
    "log"
    "strconv"


    "torme/torrent"
    

)


type decState int

const (
    decStateNew decState = iota
    decStateHalf 
)




type Decoder struct {
	data       []byte
    pos          int
    size         int

    level        int

    bdict       *BDict
    
    state        decState

    info         torrent.TorrentInfo

}


 func (d Decoder)PrintK(){
     str := d.bdict.GetK()
     log.Printf(str)
 }



func (d *Decoder) Decode(r io.Reader) (*torrent.TorrentInfo, error) {
    contents, err := ioutil.ReadAll(r)
    if err != nil{
        log.Println("bdecode failed, read all failed, err:", err)
        return nil, err
    }

    d.data = contents

    err = d.parse()
    if err != nil{
        log.Println("bdecode failed, parse failed, err:", err)
        return nil, err
    }

    d.handleDict()

    return &d.info, nil
}

const ht = "0123456789abcdef"

func (d *Decoder) handleDict(){
    str := d.bdict.Pairs["encoding"].(*BString)
    k := string(str.V)
    d.info.Encoding = k
    info := d.bdict.Pairs["info"].(*BDict)
    pieces := info.Pairs["pieces"].(*BString)

    pl := len(pieces.V)/20

    

    for i:=0; i < pl; i++{
        hex := make([]byte, 40)
        for i, v := range pieces.V[i*20:(i+1) * 20] {
		    hex[i*2] = ht[v>>4]
		    hex[i*2+1] = ht[v&0x0f]
	    }

        d.info.Pieces = append(d.info.Pieces, hex)
    }

    files := info.Pairs["files"].(*BList)

    fel := len(files.Eles)
    d.info.Files = make([]torrent.FileInfo, fel)

    for i:= 0; i < fel; i++{
        fd := files.Eles[i].(*BDict)
        fdl := fd.Pairs["length"].(*BInt).V
        d.info.Files[i].Length = fdl
        fdll := fd.Pairs["path"].(*BList)
        for _, fp := range fdll.Eles{
            path := fp.GetV()
            d.info.Files[i].Path = append(d.info.Files[i].Path, *path)
        }

    }
    
    name := info.Pairs["name"].GetV()
    d.info.Name = *name

    publisher := info.Pairs["publisher"].GetV()
    d.info.Publisher = *publisher

    publisherUrl := info.Pairs["publisher-url"].GetV()
    d.info.PublisherUrl = *publisherUrl

    announceList := d.bdict.Pairs["announce-list"].(*BList)
    anlSize := len(announceList.Eles)
    d.info.Announces = make([]torrent.AnnounceList, anlSize)

    for i:=0; i < anlSize; i++{
        anl := announceList.Eles[i].(*BList)
        for _, an := range anl.Eles{
            anv := an.GetV()
            d.info.Announces[i].Announce = append(d.info.Announces[i].Announce, *anv)
        }
    }

    comment := d.bdict.Pairs["comment"].(*BString).GetV()
    d.info.Comment = *comment

}




func (d *Decoder)parse()(err error){
    d.size = len(d.data)

    log.Println("decode parse size:", d.size)

    var r bool
    d.bdict, r = d.parseDict()

    if !r{
        log.Println("decoder parse failed")
    }

    log.Println("decode end, pos:", d.pos)
    log.Println("pair size:", len(d.bdict.Pairs))

    return nil
}


func (d *Decoder)parseDict()(*BDict, bool){
    
    d.level++
    dict := &BDict{
        Pairs : make(map[string]BData),
    }
    d.pos++   
    for d.pos < d.size {
        if d.data[d.pos] == 'e'{
            d.pos++
            break
        }
        key, r := d.parseSingleString()
        if !r {
            log.Println("parse dict failed, get key failed")
            return nil, false
        }

        skey := string(key)

        r = d.parseKV(dict, &skey)
        if !r{
            log.Println("parse dict failed, get item failed")
            return nil, false
        }
    
    }


    if dict.Pairs == nil{
        log.Fatalln("parse dict error, pairs is nil")
    }

    d.level--

    return dict, true
}

func (d *Decoder)parseKV(dict *BDict, key *string)(bool){
    c := d.data[d.pos]

    if c >= '0' && c <= '9'{
        str, r := d.parseSingleString()
        if !r {
            log.Println("parse dict item failed, key:", *key)
            return false
        }

        dict.Pairs[*key] = &BString{
            V : str,
        }
        return true
    }

    if c == 'l' {
        list, r := d.parseList()
        if !r {
            log.Println("parse list failed")
            return false
        }


        dict.Pairs[*key] = list

        return true
    }

    if c == 'd' {
        dc, r := d.parseDict()
        if !r {
            log.Println("parse dict failed")
            return false
        }
        dict.Pairs[*key] = dc

        return true
    }

    if c == 'i' {
        v, r := d.parseInt()
        if !r {
            log.Println("parse list failed")
            return false
        }
        dict.Pairs[*key] = &BInt{
            V : v,
        }

        return true
    }



    return true
}

func (d *Decoder)parseList()(*BList, bool){
    d.pos++
    d.level++

    var ls = &BList{}
    

    for d.pos < d.size {
        c := d.data[d.pos]

        if c == 'd' {
            dc, r := d.parseDict()
            if !r {
                log.Println("parse dict failed")
                return nil, false
            }

            ls.Eles = append(ls.Eles, dc)

        }else if c == 'i' {
            i, r := d.parseInt()
            if !r {
                log.Println("parse dict failed")
                return nil, false
            }

            iv := &BInt{
                V : i,
            }

            ls.Eles = append(ls.Eles, iv)
        }else if c == 'l' {
            l, r := d.parseList()
            if !r {
                log.Println("parse list failed")
                return nil, false
            }

            ls.Eles = append(ls.Eles, l)
        }else if c >= '0' && c <= '9' {
            s, r := d.parseSingleString()
            if !r {
                log.Println("parse string failed")
                return nil, false
            }

            vs := &BString{
                V : s,
            }

            ls.Eles = append(ls.Eles, vs)
        }else if c == 'e' {
            d.pos++
            break
        }
    }

    d.level--
    return ls, true
}


func (d *Decoder)parseSingleString()([]byte, bool){
    var bp = d.pos
    var keyl int

    for d.pos < d.size {
        
        if d.data[d.pos] == ':'{
            _lens := d.data[bp:d.pos]
            _str := string(_lens)
            _keyl, err := strconv.Atoi(_str)
            if err != nil{
                log.Printf("parse key failed, str:%s, bp:%d, pos:%d, bpc:%c\n", _str, bp, d.pos, d.data[bp])
                return nil, false
            }

            keyl = _keyl
            break
        }

        d.pos++
    
    }

    d.pos++

    if d.pos + keyl > d.size {
        log.Printf("parse single string failed, pos:%d, keyl:%d, size:%d\n", d.pos, keyl, d.size)
        return nil, false
    }

    keys := d.data[d.pos : d.pos + keyl]
    d.pos += keyl

   


    return keys, true
}

func (d *Decoder)parseInt()(int64, bool){
    d.pos++
    bp := d.pos
    for d.pos < d.size{
        c := d.data[d.pos]
        if c == 'e'{
            break
        }

        d.pos++
    }

    if d.data[d.pos] != 'e'{
        log.Println("parse int failed")
        return 0, false
    }

    dls := d.data[bp:d.pos]
    dlst := string(dls)
    dl, err := strconv.Atoi(dlst)

    if err != nil{
        log.Println("parse int failed, err:", err)
        return 0, false
    }
    d.pos++

    return int64(dl), true
    
}