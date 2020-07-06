package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"newsweb/models"
)

type UserController struct {
	beego.Controller
}

//展示注册页面
func (this *UserController) ShowRegister() {
	this.TplName = "register.html"
}

//处理注册数据
func (this *UserController) HandleRegister() {
	//获取前端的数据
	username := this.GetString("userName")
	pwd := this.GetString("password")
	//判断数据
	if username == "" || pwd == "" {
		beego.Error("数据为空")
		this.TplName = "register.html"
		return
	}
	//处理数据
	//创建ORM对象
	o := orm.NewOrm()
	//创建插入对象
	var user models.User
	//对象赋值
	user.Name = username
	user.Pwd = pwd
	//插入
	_, err := o.Insert(&user)
	if err != nil {
		beego.Error(err, "注册失败")
		return
	}
	//this.Ctx.WriteString("恭喜你注册成功")
	this.Redirect("/login", 302)

}

//展示登陆页面
func (this *UserController) ShowLogin() {
	userName := this.Ctx.GetCookie("userName")
	if userName == "" {
		this.Data["userName"] = ""
		this.Data["checked"] = ""
	} else {
		this.Data["userName"] = userName
		this.Data["checked"] = "checked"
	}

	this.TplName = "login.html"
}

//处理登录数据
func (this *UserController) HandleLogin() {
	//接收登录数据
	userName := this.GetString("userName")
	pwd := this.GetString("password")
	//校验登录数据
	if userName == "" || pwd == "" {
		beego.Error("登录账号或密码不能为空")
		this.TplName = "login.html"
		return
	}
	//处理数据
	//与数据库数据进行校验
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil {
		beego.Error("用户名不存在")
		this.TplName = "login.html"
		return
	}
	if user.Pwd != pwd {
		beego.Error("密码错误")
		this.TplName = "login.html"
		return
	}

	//获取用户名存cookie
	remember := this.GetString("remember")
	if remember == "on" {
		this.Ctx.SetCookie("userName", userName, 60 * 60 * 24)
	} else {
		this.Ctx.SetCookie("userName", userName, -1)
	}

	this.SetSession("userName", userName)

	//返回数据
	//this.Ctx.WriteString("登录成功")
	//跳转页面
	this.Redirect("/article/index", 302)

}

//退出登录
func (this *UserController) Showlogout() {
	this.DelSession("userName")
	this.Redirect("/login", 302)
}
