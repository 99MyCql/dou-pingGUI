/*
ICMP's Echo or Echo Reply Message:
    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |     Type      |     Code      |          Checksum             |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |           Identifier          |        Sequence Number        |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |     Data ...
    +-+-+-+-+-

Type:
    8 for echo message;
    0 for echo reply message.

Code:
    0

Checksum:
    The checksum is the 16-bit ones’s complement of the one’s
    complement sum of the ICMP message starting with the ICMP Type.
    For computing the checksum , the checksum field should be zero.
    If the total length is odd, the received data is padded with one
    octet of zeros for computing the checksum. This checksum may be
    replaced in the future.

Identifier:
    If code = 0, an identifier to aid in matching echos and replies,
    may be zero.

Sequence Number:
    If code = 0, a sequence number to aid in matching echos and
    replies, may be zero.

Description:
    The data received in the echo message must be returned in the echo
    reply message.

    The identifier and sequence number may be used by the echo sender
    to aid in matching the replies with the echo requests. For
    example, the identifier might be used like a port in TCP or UDP to
    identify a session, and the sequence number might be incremented
    on each echo request sent. The echoer returns these same values
    in the echo reply.

    Code 0 may be received from a gateway or a host.
 */

package ICMP

import (
    "encoding/binary"
    "net"
    "time"
)

// computing the checksum
func checkSum(data []byte) uint16 {
    var length int = len(data)

    // If the total length is odd, the received data is padded with one
    //    octet of zeros for computing the checksum.
    if length % 2 == 1 {
        data = append(data, 0)
        length++
    }

    // 16个比特相加
    var sum uint32
    for i:=0; i < length; i+=2 {
        sum += uint32(data[i]) << 8 + uint32(data[i+1])
    }

    sum += (sum >> 16)

    return uint16(^sum)
}

// send ICMP messgae
func SendICMP(host string, typ uint8, code uint8, icmp_data []byte, timeout uint) (
    suc bool, ip string, ttl uint8, duration int64, data_len uint) {
    defer (func() {
        if err := recover(); err != nil {
            suc = false
        }
    })()

    // 创建 ICMP 报文首部
    icmp_header := newICMPHeader(typ, code)

    // 将ICMP首部和数据部分拼接，从而计算校验和
    icmp_header.CheckSum = checkSum(append(icmp_header.toBytes(), icmp_data...))

    // 转为字节数组
    icmp := append(icmp_header.toBytes(), icmp_data...)

    // 连接对应IP地址。但，ICMP并不需要建立连接。所以，此处并未建立连接，而只是走个流程
    conn, err := net.Dial("ip4:icmp", host)
    if err != nil {
        g_logger.Panic(err)
    }

    // 记录目的 IP 地址，防止 panic 被触发，导致 ip 未赋值
    ip = conn.RemoteAddr().String()

    defer conn.Close()

    // 设置 conn 读写超限时间，即发送和接收报文的时限
    if err := conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
        g_logger.Panic(err)
    }

    // 开始时间
    start_time := time.Now()

    // 向连接中发送 ICMP 报文，此处才真正发送 ICMP 报文
    if _, err := conn.Write(icmp); err != nil {
        g_logger.Panic(err)
    }

    // 接收 ICMP 响应报文，实际接收到的为IP报文，前20字节为IP报文首部
    recv_buf := make([]byte, 65535) // 2^16-1 = 65535，IP报文最大长度
    if _, err := conn.Read(recv_buf); err != nil {
        g_logger.Panic(err)
    }

    // 往返时间
    duration = time.Now().Sub(start_time).Milliseconds()

    length := uint(binary.BigEndian.Uint16(recv_buf[2:4]))

    // 前20个字节为IP报文首部，第8个字节为TTL，第20~28个字节为ICMP首部
    return true, conn.RemoteAddr().String(), recv_buf[8], duration, length-20-8
}
