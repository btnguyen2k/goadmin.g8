/*
Application Server bootstrapper.

@author Thanh Nguyen <btnguyen2k@gmail.com>
@since template-v0.4.r1
*/
package main

import (
	"main/src/goadmin"
	"main/src/myapp"
	"math/rand"
	"time"
)

func main() {
	// it is a good idea to initialize random seed
	rand.Seed(time.Now().UnixNano())

	// start Echo server with custom bootstrappers
	// bootstrapper routine is passed the echo.Echo instance as argument, and also has access to global variables:
	// - Application configurations via goadmin.AppConfig
	// - Echo server via goadmin.EchoServer
	var bootstrappers = []goadmin.IBootstrapper{
		myapp.Bootstrapper,
	}
	goadmin.Start(bootstrappers...)
}
