package restful

import (
	"fmt"
	"github.com/ironbang/proxypool/database/struct_"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

func RESTFul() {
	fmt.Println("启动RESTFul模块...")
	app := iris.Default()
	app.Get("/", func(context context.Context) {
		context.JSON(iris.Map{
			"/v1/get-ip": "批量获取ip",
		})
	})
	v1 := app.Party("/v1")
	{
		v1.Get("/getip", func(ctx context.Context) {
			count, err := ctx.URLParamInt("count")
			if err != nil {
				count = 30
			}
			limit, err := ctx.URLParamFloat64("limit")
			if err != nil {
				limit = 0.2
			}
			ips, err := struct_.GetReliability(count, limit)
			if err == nil {
				r := iris.Map{
					"total": len(ips),
					"info":  ips,
				}
				ctx.JSON(r)
			} else {
				r := iris.Map{
					"error": err.Error(),
				}
				ctx.JSON(r)
			}
		})
	}
	app.Run(iris.Addr(":8080"))
}
