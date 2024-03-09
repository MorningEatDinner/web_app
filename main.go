package main

// 使用:: MailHog -smtp-bind-addr 0.0.0.0:1028 -api-bind-addr 127.0.0.1:8026 -ui-bind-addr 127.0.0.1:8026 我这样就行了

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xiaorui/web_app/cmd"
	"github.com/xiaorui/web_app/controller"
	"github.com/xiaorui/web_app/pkg/console"
	"github.com/xiaorui/web_app/pkg/snowflake"

	"github.com/xiaorui/web_app/dao/mysql"
	"github.com/xiaorui/web_app/dao/redis"
	"github.com/xiaorui/web_app/logger"
	"github.com/xiaorui/web_app/settings"
	"go.uber.org/zap"
)

// @title 热点论坛
// @version 1.0
// @description 这是一个热点论坛项目， 能够根据当下热点来向用户展示论坛帖子
// @termsOfService http://swagger.io/terms/
// @contact.name xiaorui zheng
// @contact.url http://www.swagger.io/support
// @contact.email 1298453249@qq.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:9000
// @BasePath /api/v1
func main() {

	var rootCmd = &cobra.Command{
		Use:   "bluebell",
		Short: "[Start] bluebell...",
		Long:  `Default will run "serve" command, you can use "-h" flag to see all subcommands`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(os.Args) < 2 {
				os.Args = append(os.Args, "conf/config.yaml")
				// fmt.Println("Please set your config file!")
			}

			//1. 加载配置文件
			if err := settings.Init(os.Args[1]); err != nil {
				fmt.Printf("init settings failed, err:%v", err)
				return
			}
			//2 初始化日志文件
			if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
				fmt.Printf("init logger failed, err:%v", err)
				return
			}
			zap.L().Sync()
			zap.L().Debug("logger init success...")
			//3. 初始化mysql
			if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
				fmt.Printf("init mysql failed, err:%v", err)
				return
			}
			//4 初始化redis
			if err := redis.Init(settings.Conf.RedisConfig); err != nil {
				fmt.Printf("init redis failed, err:%v", err)
				return
			}

			//初始化雪花算法， 用于创建用户id
			if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
				fmt.Printf("snowflake.Init err:%v", err)
				return
			}

			//注册gin中的validator校验器
			if err := controller.InitTrans("zh"); err != nil {
				fmt.Printf("controller.InitTrans err:%v", err)
				return
			}
		},
	}
	defer mysql.Close()
	defer redis.Close()

	// 注册子命令
	rootCmd.AddCommand(
		cmd.CmdServe,
	)

	cmd.RegisterDefaultCmd(rootCmd, cmd.CmdServe)

	// 配置默认运行 Web 服务
	cmd.RegisterDefaultCmd(rootCmd, cmd.CmdServe)

	// 注册全局参数，--env
	cmd.RegisterGlobalFlags(rootCmd)

	// 执行主命令
	if err := rootCmd.Execute(); err != nil {
		console.Exit(fmt.Sprintf("Failed to run app with %v: %s", os.Args, err.Error()))
	}
}
