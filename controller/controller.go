package controller

import (
	"blog/dao"
	"blog/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
	"html/template"
	"net/http"
	"strconv"
)

// 基于Cookie实现用户登录认证的中间件

func Middleware(c *gin.Context) {
	// 在返回页面之前要先校验是否存在username的Cookie
	// 获取Cookie
	username, err := c.Cookie("username")
	if err != nil {
		// 直接让跳转到登录页面

		toPath := fmt.Sprintf("%s?next=%s", "/login", c.Request.URL.Path)
		c.Redirect(http.StatusMovedPermanently, toPath)
		return
	}
	// 用户已经登录了
	c.Set("username", username) // 在上下文中设置一个键值对
	c.Next()                    // 继续后续的处理函数
}
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	u := dao.Mgr.Login(username)
	if username == "" || password == "" {
		c.HTML(200, "login.html", "用户名密码不能为空")
	} else {

		if u.Username == username && u.Password == password {
			c.SetCookie("username", u.Username, 20, "/", "127.0.0.1", false, true)
			fmt.Println("登录成功", username)
			c.Redirect(301, "post_index")
		} else {
			c.HTML(200, "login.html", "用户名密码错误")
		}
	}

}
func RegisterUser(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	user := model.User{
		Username: username,
		Password: password,
	}
	//执行数据库添加
	dao.Mgr.RegisterUser(&user)

	c.Redirect(301, "/")

}

func GoRegister(c *gin.Context) {
	c.HTML(200, "register.html", nil)
}

func Index(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func GoLogin(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}

// GetPostIndex 博客列表
func GetPostIndex(c *gin.Context) {

	posts := dao.Mgr.GetAllPost()
	c.HTML(200, "postIndex.html", posts)
}

// AddPost 添加博客
func AddPost(c *gin.Context) {
	tmpUsername, ok := c.Get("username")
	if !ok {
		// 如果取不到值，说明前面中间件出问题了
		c.Redirect(http.StatusMovedPermanently, "/login")
		return
	}
	username, ok := tmpUsername.(string)
	if !ok {
		// 类型断言失败
		c.Redirect(http.StatusMovedPermanently, "/login")
		return
	}
	fmt.Println("555", username)

	title := c.PostForm("title")
	tag := c.PostForm("tag")
	context := c.PostForm("context")

	post := model.Post{
		Title:   title,
		Tag:     tag,
		Content: context,
	}

	dao.Mgr.AddPost(&post)

	c.Redirect(302, "/post_index")
}

// GoAddPost 跳转到添加博客
func GoAddPost(c *gin.Context) {
	c.HTML(200, "post.html", nil)
}

// PostDetail 查看详细博客
func PostDetail(c *gin.Context) {
	//获取传递的id
	s := c.Query("id")
	pid, _ := strconv.Atoi(s)
	p := dao.Mgr.GetPost(pid)
	content := blackfriday.Run([]byte(p.Content))
	c.HTML(200, "detail.html", gin.H{
		"Pid":     pid,
		"Title":   p.Title,
		"Content": template.HTML(content),
	})
}

//UpDate 修改博客
func UpDate(c *gin.Context) {
	tmpUsername, ok := c.Get("username")
	if !ok {
		// 如果取不到值，说明前面中间件出问题了
		c.Redirect(http.StatusMovedPermanently, "/login")
		return
	}
	username, ok := tmpUsername.(string)
	if !ok {
		// 类型断言失败
		c.Redirect(http.StatusMovedPermanently, "/login")
		return
	}
	fmt.Println("6666", username)

	s := c.PostForm("id")
	pid, _ := strconv.Atoi(s)
	fmt.Println("qq", pid)

	title := c.PostForm("title")
	tag := c.PostForm("tag")
	content := c.PostForm("context")

	post := model.Post{
		Title:   title,
		Tag:     tag,
		Content: content,
	}

	dao.Mgr.Update(pid, &post)
	c.Redirect(302, "/post_index")
}

// GoUpdate 跳转到修改博客
func GoUpdate(c *gin.Context) {
	//获取传递的id
	s := c.Query("id")
	pid, _ := strconv.Atoi(s)
	fmt.Println("rr", pid)
	p := dao.Mgr.GetPost(pid)

	c.HTML(200, "updateIndex.html", gin.H{
		"Pid":     p.ID,
		"Title":   p.Title,
		"Content": p.Content,
		"Tag":     p.Tag,
	})

}
