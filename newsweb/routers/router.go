package routers

import (
	"newsweb/controllers"
	"github.com/astaxie/beego"
    "github.com/astaxie/beego/context"
)

func init() {
    beego.InsertFilter("/article/*",beego.BeforeExec,filterFunc)
    beego.Router("/", &controllers.MainController{})
    beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleRegister")
    beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
    beego.Router("/article/index",&controllers.ArticleController{},"get:ShowIndex")
    beego.Router("/article/addarticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArticle")
    beego.Router("/article/content",&controllers.ArticleController{},"get:ShowContent")
    beego.Router("/article/editArticle",&controllers.ArticleController{},"get:ShowEditArticle;post:HandleEditArticle")
    beego.Router("/article/deleteArticle",&controllers.ArticleController{},"get:HandleDeleteArticle")
    beego.Router("/article/addType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
    beego.Router("/article/logout",&controllers.UserController{},"get:Showlogout")
    beego.Router("/article/deleteType",&controllers.ArticleController{},"get:DeleteArticleType")
    beego.Router("/redis",&controllers.RedisGit{},"get:ShowRedis")
}

func filterFunc(ctx*context.Context)  {
    userName:= ctx.Input.Session("userName")
    if userName==nil {
        ctx.Redirect(302,"/login")
    }

}