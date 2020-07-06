package main

import (
	_ "newsweb/routers"
	"github.com/astaxie/beego"
	_ "newsweb/models"
)

func main() {
	beego.AddFuncMap("prePage",ShowPrePage)
	beego.AddFuncMap("nextPage",ShowNextPage)
	beego.Run()
}

func ShowPrePage(pageindex int)int  {
	if pageindex <= 1{
		return  1
	}
	return pageindex - 1
}

func ShowNextPage(pageindex int,pagecount float64)int  {
	if pageindex >= int(pagecount)  {
		return int(pagecount)
	}
	return pageindex + 1
}