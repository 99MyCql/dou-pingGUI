package ICMP

import (
    "bytes"
    "encoding/binary"
    "log"
    "math/rand"
    "time"
)

// the header of ICMP message
type ICMPHeader struct {
    Type        uint8   // 类型
    Code        uint8   // 代码
    CheckSum    uint16  // 校验和
    ID          uint16  // ID，用于匹配 ICMP 请求和响应报文
    Sequence    uint16  // 序号，用于匹配 ICMP 请求和响应报文
}

// constructor of ICMPHeader
func newICMPHeader (typ uint8, code uint8) *ICMPHeader {
    // 生成随机数种子
    // func Now() Time：返回当前时间
    // func (t Time) Unix() int64：返回从1970年1月1日到t的秒数
    rand.Seed(time.Now().Unix())

    // 构造 ICMP 报文首部
    var icmp_header *ICMPHeader = new(ICMPHeader)
    icmp_header.Type = typ
    icmp_header.Code = code
    icmp_header.CheckSum = 0
    icmp_header.ID = uint16(rand.Uint32())
    icmp_header.Sequence = uint16(rand.Uint32())

    return icmp_header
}

// convert header of ICMP message to bytes
func (icmp_header *ICMPHeader) toBytes () []byte {
    // 创建字节缓冲区
    buf := &bytes.Buffer{}

    // 写入 ICMP 首部数据，大端法
    if err := binary.Write(buf, binary.BigEndian, icmp_header); err != nil {
        log.Panic(err)
    }

    return buf.Bytes()
}