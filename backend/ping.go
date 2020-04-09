package backend

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync/atomic"

	"github.com/99MyCql/dou-pingGUI/backend/ICMP"
)

// 获取主机 IP 地址，可能有多个
func (contro *Controller) GetIPAddrs() []string {
	ip_addrs := make([]string, 0)

	// 获取所有网卡接口
	net_interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(net_interfaces); i++ {
		// 如果该接口处于活跃状态
		if (net_interfaces[i].Flags & net.FlagUp) != 0 {
			// 遍历该网卡接口的所有 IP 地址
			addrs, _ := net_interfaces[i].Addrs()
			for j := 0; j < len(addrs); j++ {
				// 按 CIDR 解析 IP 地址
				_, ip_net, err := net.ParseCIDR(addrs[j].String())
				// 如果 IP 不是环回地址（如：127.0.0.1），且是 IPv4 地址
				if err == nil && !ip_net.IP.IsLoopback() && ip_net.IP.To4() != nil {
					ip_addrs = append(ip_addrs, addrs[j].String())
				}
			}
		}
	}

	return ip_addrs
}

// ping 主机，主机名（域名或 IP 地址）错误会返回 false
func (contro *Controller) Ping(host string, data []byte, count uint) bool {
	// 判断域名或 IP 地址是否正确
	_, err := net.LookupHost(host)
	if err != nil {
		contro.logger.Info("ping err: host name error!")
		return false
	}

	contro.logger.Infof("ping %s with %d bytes \"%s\":\n", host, len(data), string(data))

	go func() {
		for i := uint(0); i < count; i++ {
			if contro.stop == 1 {
				atomic.StoreInt32(&contro.stop, 0)
				break
			}
			// 发送 ping 程序的 ICMP 报文
			suc, ip, ttl, duration, data_len := ICMP.SendICMP(host, 8, 0, data, 1)
			contro.logger.Infof("ping result: %t, ip:%s, ttl:%d, duration:%d, data len:%d",
				suc, ip, ttl, duration, data_len)
			// 向前端发送数据
			contro.runtime.Events.Emit("ping", suc, ip, ttl, duration, data_len)
		}
	}()

	return true
}

// ping 一个子网中的所有主机。
// 传入 CIDR 形式的 IP 地址，ping 该地址所在子网中的所有主机。
// 返回该子网中的主机数，如果不为正数，则说明传入的 IP 地址有问题
func (contro *Controller) PingIPNet(ip_CIDR string) int {
	// 先将 CIDR 形式的 IP 地址（如：192.168.0.103/24）解析为 IPNet 结构
	_, ip_net, err := net.ParseCIDR(ip_CIDR)
	if err != nil {
		contro.logger.Info("pingIPNet err: ip of CIDR error!")
		return 0
	}

	ip := binary.BigEndian.Uint32(ip_net.IP)
	mask := binary.BigEndian.Uint32(ip_net.Mask)

	go func() {
		for i := uint32(1); i < ^mask; i++ {
			if contro.stop == 1 {
				atomic.StoreInt32(&contro.stop, 0)
				break
			}
			host := fmt.Sprintf("%d.%d.%d.%d",
				byte((ip+i)>>24),
				byte((ip+i)>>16),
				byte((ip+i)>>8),
				byte(ip+i))
			// 发送 ping 程序的 ICMP 报文
			suc, ip, ttl, duration, data_len := ICMP.SendICMP(host, 8, 0, []byte{}, 1)
			contro.logger.Infof("pingNet result: %t, ip:%s, ttl:%d, duration:%d, data len:%d",
				suc, ip, ttl, duration, data_len)
			// 向前端发送数据
			contro.runtime.Events.Emit("pingIPNet", suc, ip, ttl, duration, data_len)
		}
	}()

	return int(^mask)-2
}
