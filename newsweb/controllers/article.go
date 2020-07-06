package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"newsweb/models"
	"math"
	"github.com/gomodule/redigo/redis"
	"bytes"
	"encoding/gob"
)

type ArticleController struct {
	beego.Controller
}

//封装函数处理校验上传图片的业务
func CheckImg(this *ArticleController, key string) string {
	//获取图片信息,用getfile,返回三个值，文件指针，文件头，err
	file, head, err := this.GetFile(key)

	//检验数据
	if err != nil {
		beego.Error("获取用户添加文章数据失败", err)
		this.TplName = "update.html"
		return ""
	}
	defer file.Close()

	//判断图片的大小
	if head.Size > 5000000 {
		beego.Error("图片过大")
		this.TplName = "update.html"
		return ""
	}

	//判断图片类型
	//获取文件的后缀名
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		beego.Error("图片格式错误")
		this.TplName = "update.html"
		return ""
	}

	//防止重名
	filename := time.Now().Format("2006-01-02 15:04:05")

	//操作数据
	this.SaveToFile("uploadname", "./static/img/"+filename+ext)

	return "/static/img/" + filename + ext
}

//展示首页
func (this *ArticleController) ShowIndex() {
	//获取数据库数据
	//创建ORM对象
	o := orm.NewOrm()
	//定义对象储存获取数据,多行数据所以用切片
	var articles []models.Article
	//制定获取数据表,返回queryseter
	qs := o.QueryTable("article")
	//获取所有数据
	//qs.All(&articles)

	//获取所有文章类型
	var articleTypes []models.ArticleType


	//连接redis数据库
	conn,err:=redis.Dial("tcp",":6379")
	if err!=nil {
		beego.Error("数据库连接失败",err)
		return
	}
	defer conn.Close()

	//操作数据库获取数据
	data,err:=redis.Bytes(conn.Do("get","articleTypes"))
	if len(data)==0 {
		//第一次没有获取到数据才从mysql中获取数据
		o.QueryTable("ArticleType").All(&articleTypes)
		var buffer bytes.Buffer
		enco:=gob.NewEncoder(&buffer)
		enco.Encode(&articleTypes)

		conn.Do("set","articleTypes",buffer.Bytes())
	}else {
		deco:=gob.NewDecoder(bytes.NewReader(data))
		deco.Decode(&articleTypes)
	}


	this.Data["articleTypes"] = articleTypes

	typeName := this.GetString("select")

	//定义每页展示的条数
	pagesize := 2

	//获取前端的pageindex页数
	pageindex, err := this.GetInt("pageindex")
	if err != nil {
		pageindex = 1
	}

	var count int64
	if typeName == "" {
		count, _ = qs.RelatedSel("ArticleType").Count()
		qs.Limit(pagesize, pagesize*(pageindex-1)).RelatedSel("ArticleType").All(&articles)

	} else {
		count, _ = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).Count()
		qs.Limit(pagesize, pagesize*(pageindex-1)).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articles)

	}
	//定义总条数变量
	this.Data["count"] = count

	//总页数
	pagecount := math.Ceil(float64(count) / float64(pagesize))
	this.Data["pagecount"] = pagecount

	this.Data["typeName"] = typeName

	//传递数据给前端
	this.Data["pageindex"] = pageindex
	this.Data["articles"] = articles

	this.Layout = "layout.html"

	this.TplName = "index.html"
}

//展示添加文章页面
func (this *ArticleController) ShowAddArticle() {
	//获取数据
	o := orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)

	//传递数据
	this.Data["articleTypes"] = articleTypes

	this.Layout = "layout.html"
	this.TplName = "add.html"
}

//处理添加文章页面数据
func (this *ArticleController) HandleAddArticle() {
	//获取数据
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	//获取图片信息,用getfile,返回三个值，文件指针，文件头，err
	file, head, err := this.GetFile("uploadname")

	//检验数据
	if articleName == "" || content == "" || err != nil {
		beego.Error("获取用户添加文章数据失败", err)
		this.TplName = "add.html"
		return
	}
	defer file.Close()

	//判断图片的大小
	if head.Size > 5000000 {
		beego.Error("图片过大")
		this.TplName = "add.html"
		return
	}

	//判断图片类型
	//获取文件的后缀名
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		beego.Error("图片格式错误")
		this.TplName = "add.html"
		return
	}

	//防止重名
	filename := time.Now().Format("2006-01-02 15:04:05")

	//操作数据
	this.SaveToFile("uploadname", "./static/img/"+filename+ext)

	//创建orm对象
	o := orm.NewOrm()
	//创建插入对象
	var article models.Article
	//给插入对象赋值
	article.Content = content
	article.Title = articleName
	article.Img = "/static/img/" + filename + ext

	//获取类型数据
	typeName := this.GetString("select")
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Read(&articleType, "TypeName")

	article.ArticleType = &articleType

	//插入
	_, err = o.Insert(&article)
	if err != nil {
		beego.Error("插入失败", err)
		this.TplName = "add.html"
		return
	}

	//返回数据
	this.Redirect("/article/index", 302)
}

//展示首页的查看详情页面
func (this *ArticleController) ShowContent() {
	//获取前端传过来的数据
	id, err := this.GetInt("id")
	//校验数据
	if err != nil {
		beego.Error("获取失败", err)
		this.TplName = "index.html"
		return
	}

	//处理数据
	//创建ORM对象
	o := orm.NewOrm()
	//实例化对象
	var article models.Article
	//查询条件赋值
	article.Id2 = id
	//查询
	o.Read(&article)

	//增加阅读次数
	article.ReadCount += 1
	o.Update(&article)

	//多对多插入
	m2m := o.QueryM2M(&article, "Users")
	userName := this.GetSession("userName").(string)
	var user models.User
	user.Name = userName
	o.Read(&user, "Name")
	m2m.Add(user)

	//加载关系,多对多查询
	//o.LoadRelated(&article,"Users")
	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id2", article.Id2).Distinct().All(&users)

	this.Data["users"] = users

	//返回数据
	this.Data["article"] = article

	this.Layout = "layout.html"
	//this.LayoutSections = make(map[string]string)
	//this.LayoutSections["jsFile"] = "index.js"

	this.TplName = "content.html"
}

//展示首页的编辑文章里的页面
func (this *ArticleController) ShowEditArticle() {
	//获取数据
	id, err := this.GetInt("id")
	if err != nil {
		beego.Error(err)
		this.TplName = "index.html"
		return
	}
	o := orm.NewOrm()
	var article models.Article
	article.Id2 = id
	o.Read(&article)

	//传递数据
	this.Data["article"] = article

	this.Layout = "layout.html"
	this.TplName = "update.html"

}

//处理编辑文章页面
func (this *ArticleController) HandleEditArticle() {
	//获取数据
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	//调用方法获取和校验图片文件,返回文件路径
	fliename := CheckImg(this, "uploadname")
	id, err := this.GetInt("id")
	//校验数据
	if articleName == "" || content == "" || fliename == "" || err != nil {
		beego.Error("更改的数据不能为空", err)
		this.TplName = "update.html"
		return
	}

	//处理数据
	o := orm.NewOrm()
	var article models.Article
	article.Id2 = id
	err = o.Read(&article)
	if err != nil {
		beego.Error("更新的数据不存在")
		this.TplName = "update.html"
		return
	}
	article.Title = articleName
	article.Content = content
	article.Img = fliename
	o.Update(&article)

	//返回数据
	this.Redirect("/article/index", 302)

}

//首页的删除文章操作
func (this *ArticleController) HandleDeleteArticle() {
	//获取数据
	id, err := this.GetInt("id")
	//校验数据
	if err != nil {
		beego.Error("获取数据失败", err)
		this.TplName = "index.html"
		return
	}
	//处理数据
	o := orm.NewOrm()
	var article models.Article
	article.Id2 = id
	_, err = o.Delete(&article)
	if err != nil {
		beego.Error("删除数据失败", err)
		this.TplName = "index.html"
		return
	}
	//返回数据
	this.Redirect("/article/index", 302)
}

//展示首页的添加分类页面
func (this *ArticleController) ShowAddType() {
	//获取数据
	o := orm.NewOrm()
	var ArticleTypes []models.ArticleType
	qs := o.QueryTable("ArticleType")
	qs.All(&ArticleTypes)

	//返回数据
	this.Data["ArticleTypes"] = ArticleTypes
	this.Layout = "layout.html"
	this.TplName = "addType.html"
}

//操作添加分类页面数据
func (this *ArticleController) HandleAddType() {
	//获取数据
	typeName := this.GetString("typeName")
	//校验数据
	if typeName == "" {
		beego.Error("填入数据不能为空")
		this.TplName = "addType.html"
		return
	}
	//处理数据
	o := orm.NewOrm()
	var ArticleType models.ArticleType
	ArticleType.TypeName = typeName
	_, err := o.Insert(&ArticleType)
	if err != nil {
		beego.Error("charu shibai ", err)
		this.TplName = "addType.html"
		return
	}

	//返回数据
	this.Redirect("/article/addType", 302)
}

//删除文章类型
func (this *ArticleController) DeleteArticleType() {
	id, err := this.GetInt("id")
	if err != nil {
		beego.Error("获取数据失败")
		this.TplName = "addType.html"
		return
	}
	o := orm.NewOrm()
	var articletype models.ArticleType
	articletype.Id = id
	o.Delete(&articletype)

	this.Redirect("/article/addType", 302)

}
