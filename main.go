package main

import (
	_ "CRMsystemproject/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}

