package backend

import (
	"sync/atomic"

	"github.com/wailsapp/wails"
)

// 后端控制器
type Controller struct {
	runtime *wails.Runtime      // Wails Runtime Interface. 可以借此向前端发送信息
	logger  *wails.CustomLogger // Wails 输出流接口
	stop    int32               // 是否暂停所有 goroutine ，1为暂停，0则不暂停
}

// constructor of Controller
func NewController() *Controller {
	controller := new(Controller)
	controller.stop = 0
	return controller
}

// callback by wails
// To access this runtime, we create a struct method called "WailsInit"
func (contro *Controller) WailsInit(runtime *wails.Runtime) error {
	contro.runtime = runtime
	contro.logger = contro.runtime.Log.New("Controller")
	contro.logger.Info("I'm here")
	return nil
}

// 暂停所有 goroutine
func (contro *Controller) SetStop() {
	atomic.StoreInt32(&contro.stop, 1)
}
