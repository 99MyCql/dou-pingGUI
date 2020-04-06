package main

import (
    "fmt"
    "github.com/leaanthony/mewn"
    "github.com/wailsapp/wails"
    "github.com/99MyCql/dou-pingGUI/backend"
    "net"
)

func basic() string {
  return "Hello World!"
}

func testBackend() {
    backend.Ping("192.168.0.100", []byte("hello world!"), 4)

	ip_addrs := backend.GetIPAddrs()
	fmt.Println(ip_addrs)

	for _, ip_addr := range ip_addrs {
		_, ip_net, err := net.ParseCIDR(ip_addr.String())
		if err != nil {
			panic(err)
		}
		fmt.Println(backend.PingIPNet(ip_net))
	}
}

func main() {

    js := mewn.String("./frontend/dist/app.js")
    css := mewn.String("./frontend/dist/app.css")

    app := wails.CreateApp(&wails.AppConfig{
        Width:  1024,
        Height: 768,
        Title:  "dou-pingGUI",
        JS:     js,
        CSS:    css,
        Colour: "#131313",
    })
    app.Bind(basic)
    app.Run()
}
