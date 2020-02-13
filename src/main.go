package main

import (
	_ "github.com/GoAdminGroup/go-admin/adapter/gin"
	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/examples/datamodel"
	"github.com/GoAdminGroup/go-admin/modules/config"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/mysql"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/plugins/admin"
	"github.com/GoAdminGroup/go-admin/plugins/example"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/GoAdminGroup/go-admin/template/login"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/themes/adminlte"
	"github.com/gin-gonic/gin"
	tt "github.com/xiaowei520/go-admin-x/src/datamodel"
)

func main() {
	r := gin.Default()

	eng := engine.Default()

	// global config
	cfg := config.Config{
		Databases: config.DatabaseList{
			"default": {
				Host: "127.0.0.1",
				Port: "3306",
				User: "root",
				//Pwd:          "123456",
				Pwd:        "",
				Name:       "goadmin",
				MaxIdleCon: 50,
				MaxOpenCon: 150,
				Driver:     "mysql",
			},
		},
		UrlPrefix: "admin",
		// STORE 必须设置且保证有写权限，否则增加不了新的管理员用户
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Language: language.CN,
		// 开发模式
		Debug: true,
		// 日志文件位置，需为绝对路径
		InfoLogPath:   "/Users/admin/Proj/Go/src/github.com/xiaowei520/go-admin-x/logs/info.log",
		AccessLogPath: "/Users/admin/Proj/Go/src/github.com/xiaowei520/go-admin-x/logs/access.log",
		ErrorLogPath:  "/Users/admin/Proj/Go/src/github.com/xiaowei520/go-admin-x/logs/error.log",
		ColorScheme:   adminlte.ColorschemeSkinBlack,
	}

	// Generators： 详见 https://github.com/GoAdminGroup/go-admin/blob/master/examples/datamodel/tables.go
	adminPlugin := admin.NewAdmin(datamodel.Generators)

	// 增加 chartjs 组件
	template.AddComp(chartjs.NewChart())

	// 增加 generator, 第一个参数是对应的访问路由前缀
	// 例子:
	//
	// "user" => http://localhost:9033/admin/info/user
	//
	adminPlugin.AddGenerator("user", datamodel.GetUserTable)

	// 增加登录组件
	template.AddLoginComp(login.GetLoginComponent())

	// 自定义首页

	r.GET("/admin", func(ctx *gin.Context) {
		eng.Content(ctx, func(ctx interface{}) (types.Panel, error) {
			return datamodel.GetContent()
		})
	})
	r.GET("/push", func(ctx *gin.Context) {
		eng.Content(ctx, func(ctx interface{}) (types.Panel, error) {
			return tt.GetContent()
			//return datamodel.GetContent()
		})
	})

	//r.GET("/custom", func(ctx *gin.Context) {
	//	engine.Content(ctx, func(ctx interface{}) (types.Panel, error) {
	//		return datamodel.GetContent()
	//	})
	//})

	// you can custom a plugin like:
	examplePlugin := example.NewExample()

	_ = eng.AddConfig(cfg).AddPlugins(adminPlugin, examplePlugin).Use(r)

	_ = r.Run(":9033")
}
