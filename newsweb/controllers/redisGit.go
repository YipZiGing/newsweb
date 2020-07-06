package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
)

type RedisGit struct {
	beego.Controller
}

func ShowRedis()  {
	//连接数据库
	coon,err:=redis.Dial("tcp",":6379")
	if err!=nil {
		beego.Error("连接失败",err)
		return
	}
	defer coon.Close()

	//操作函数
	resp,err:=coon.Do("mget","kk","vv","ll")

	//回复助手函数，类型转换
	//result,_:=redis.String(resp,err)
	result,_:=redis.Values(resp,err)

	var kk,vv string
	var ll int
	redis.Scan(result,&kk,&vv,&ll)


}