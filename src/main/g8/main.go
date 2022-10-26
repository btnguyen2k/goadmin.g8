/*
Application Server bootstrapper.

@author Thanh Nguyen <btnguyen2k@gmail.com>
@since template-v0.4.r1
*/
package main

import (
	"math/rand"
	"time"

	"main/src/goadmin"
	"main/src/myapp"
)

func main() {
	// it is a good idea to initialize random seed
	rand.Seed(time.Now().UnixNano())

	// start Echo server with custom bootstrappers
	var bootstrappers = []goadmin.IBootstrapper{
		myapp.Bootstrapper,
	}
	goadmin.Start(bootstrappers...)
}
