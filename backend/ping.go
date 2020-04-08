package backend

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync/atomic"

	"github.com/99MyCql/dou-pingGUI/backend/ICMP"
)

// 获取主机 IP 地址，可能有多个
func (contro *Controller) GetIPAddrs() []net.Addr {
	ip_addrs := make([]net.Addr, 0)

	net_interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(net_interfaces); i++ {
		if (net_interfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := net_interfaces[i].Addrs()
			for j := 0; j < len(addrs); j++ {
				_, ip_net, err := net.ParseCIDR(addrs[j].String())
				if err == nil && !ip_net.IP.IsLoopback() && ip_net.IP.To4() != nil {
					ip_addrs = append(ip_addrs, addrs[j])
				}
			}
		}
	}

	return ip_addrs
}

// ping 主机，ping 通一次则返回 true
func (contro *Controller) Ping(host string, data []byte, count uint) bool {
	// 如果是主机名，则先解析主机名
	IP, err := net.LookupHost(host)
	if err != nil {
		contro.logger.Info("ping err: host name error!")
		return false
	}

	contro.logger.Infof("ping %s with %d bytes \"%s\":\n", IP[0], len(data), string(data))

	go func() {
		for i := uint(0); i < count; i++ {
			if contro.stop == 1 {
				atomic.StoreInt32(&contro.stop, 0)
				break
			}
			// 发送 ping 程序的 ICMP 报文
			suc, ttl, duration, data_len := ICMP.SendICMP(IP[0], 8, 0, data, 1)
			contro.logger.Infof("ping result: %t, ttl:%d, duration:%d, data len:%d", suc, ttl, duration, data_len)
			contro.runtime.Events.Emit("ping", suc, IP[0], ttl, duration, data_len)
		}
	}()

	return true
}

// ping 一个子网中的所有主机，返回可以 ping 通的 IP 地址
func (contro *Controller) PingIPNet(ip_net *net.IPNet) (ip_list []string) {
	ip := binary.BigEndian.Uint32(ip_net.IP)
	mask := binary.BigEndian.Uint32(ip_net.Mask)

	for i := uint32(1); i < ^mask; i++ {
		ip_string := fmt.Sprintf("%d.%d.%d.%d",
			byte((ip+i)>>24),
			byte((ip+i)>>16),
			byte((ip+i)>>8),
			byte(ip+i))
		if contro.Ping(ip_string, []byte{}, 1) {
			ip_list = append(ip_list, ip_string)
		}
	}

	return ip_list
}
