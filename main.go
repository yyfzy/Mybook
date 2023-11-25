package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	ginRedis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github/yyfzy/mybook/config"
	"github/yyfzy/mybook/internal/repository"
	"github/yyfzy/mybook/internal/repository/dao"
	"github/yyfzy/mybook/internal/service"
	"github/yyfzy/mybook/internal/web"
	"github/yyfzy/mybook/internal/web/middleware"
	"github/yyfzy/mybook/pkg/ginx/ratelimit"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

func main() {
	db := initDB()
	server := initWebServer()
	u := initUser(db)
	u.RegisterRoutes(server)
	//server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello, yyf")
	})
	server.Run(":8080")

}

func initWebServer() *gin.Engine {
	server := gin.Default()
	redisClient := redis.NewClient(&redis.Options{Addr: "192.168.137.132:6379"})
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
	server.Use(cors.New(cors.Config{
		AllowAllOrigins: false,
		AllowOrigins:    nil,
		AllowOriginFunc: func(origin string) bool {
			return strings.Contains(origin, "abc.com")
		},
		AllowMethods:           nil,
		AllowHeaders:           nil,
		AllowCredentials:       false,
		ExposeHeaders:          nil,
		MaxAge:                 0,
		AllowWildcard:          false,
		AllowBrowserExtensions: false,
		AllowWebSockets:        false,
		AllowFiles:             false,
	}))

	//store := cookie.NewStore([]byte("secret"))
	store, err := ginRedis.NewStore(16, "tcp", config.Config.Redis.Addr, "",
		[]byte("3o4q6EshoibpRdTB6iPCayquqFmMQzkv"), []byte("naspBhPdXGTMOG9OoRaIukf48sf8WUXU"))
	if err != nil {
		panic(err)
	}

	server.Use(sessions.Sessions("mysession", store))
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/login", "/users/signup").
		Build())
	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	//db, err := gorm.Open(mysql.Open("root:root@tcp(webook-mysql:13309)/webook"))

	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
