package main

import (
	"fmt"
	"net"

	"github.com/99MyCql/dou-pingGUI/backend"
	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"
)

func basic() string {
	return "Hello World!"
}

func testBackend(controller backend.Controller) {
	controller.Ping("192.168.0.100", []byte("hello world!"), 4)

	ip_addrs := controller.GetIPAddrs()
	fmt.Println(ip_addrs)

	for _, ip_addr := range ip_addrs {
		_, ip_net, err := net.ParseCIDR(ip_addr.String())
		if err != nil {
			panic(err)
		}
		fmt.Println(controller.PingIPNet(ip_net))
	}
}

func main() {

	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:     1024,
		Height:    600,
		Title:     "dou-pingGUI",
		JS:        js,
		CSS:       css,
		Colour:    "#131313",
		Resizable: true,
	})
	app.Bind(backend.NewController())
	app.Run()
}
