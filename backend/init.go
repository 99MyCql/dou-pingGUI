package backend

import "github.com/wailsapp/wails"

type Controller struct {
	runtime *wails.Runtime
	logger  *wails.CustomLogger
}

// constructor of Controller
func NewController() *Controller {
	return new(Controller)
}

// callback by wails
// To access this runtime, we create a struct method called "WailsInit"
func (contro *Controller) WailsInit(runtime *wails.Runtime) error {
	contro.runtime = runtime
	contro.logger = contro.runtime.Log.New("Controller")
	contro.logger.Info("I'm here")
	return nil
}
