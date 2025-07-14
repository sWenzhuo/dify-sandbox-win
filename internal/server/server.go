package server

import (
	"dify-sandbox-win/internal/controller"
	"dify-sandbox-win/internal/core/runner/python"
	"dify-sandbox-win/internal/static"
	"dify-sandbox-win/internal/utils/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func initConfig() {
	// auto migrate database
	err := static.InitConfig("conf/config.yaml")
	if err != nil {
		log.Panic("failed to init config: %v", err)
	}
	log.Info("config init success")

	err = static.SetupRunnerDependencies()
	if err != nil {
		log.Error("failed to setup runner dependencies: %v", err)
	}
	log.Info("runner dependencies init success")
}

func initServer() {
	config := static.GetDifySandboxGlobalConfigurations()
	if !config.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(gin.Recovery())
	if gin.Mode() == gin.DebugMode {
		r.Use(gin.Logger())
	}

	controller.Setup(r)

	r.Run(fmt.Sprintf(":%d", config.App.Port))
}

//func initServerV1() {
//	r := gin.Default()
//	r.Use(gin.Recovery())
//	if gin.Mode() == gin.DebugMode {
//		r.Use(gin.Logger())
//	}
//
//	controller.Setup(r)
//
//	r.Run(fmt.Sprintf(":%d", 8080))
//
//}

func initDependencies() {
	log.Info("installing python dependencies...")

	dependencies := static.GetRunnerDependencies() //如果没有就是空
	//安装依赖包
	//err := python.InstallDependencies(dependencies.PythonRequirements)

	config := static.GetDifySandboxGlobalConfigurations()
	err := python.InstallDependenciesV1(config.RequirementsFile)

	if err != nil {
		log.Panic("failed to install python dependencies: %v", err)
	}
	log.Info("python dependencies installed")

	log.Info("initializing python dependencies sandbox...")
	//查看环境是否正常
	err = python.PreparePythonDependenciesEnv_V1()

	//err = python.PreparePythonDependenciesEnv()
	if err != nil {
		log.Panic("failed to initialize python dependencies sandbox: %v", err)
	}
	log.Info("python dependencies sandbox initialized")

	// 定时异步更新依赖
	go func() {
		updateInterval := static.GetDifySandboxGlobalConfigurations().PythonDepsUpdateInterval
		tickerDuration, err := time.ParseDuration(updateInterval)
		if err != nil {
			log.Error("failed to parse python dependencies update interval, skip periodic updates: %v", err)
			return
		}
		ticker := time.NewTicker(tickerDuration)
		for range ticker.C {
			if err := updatePythonDependencies(dependencies); err != nil {
				log.Error("Failed to update Python dependencies: %v", err)
			}
		}
	}()
}

func updatePythonDependencies(dependencies static.RunnerDependencies) error {
	log.Info("Updating Python dependencies...")
	config := static.GetDifySandboxGlobalConfigurations()
	if err := python.InstallDependenciesV1(config.RequirementsFile); err != nil {
		log.Error("Failed to install Python dependencies: %v", err)
		return err
	}
	/*
		if err := python.InstallDependencies(dependencies.PythonRequirements); err != nil {
			log.Error("Failed to install Python dependencies: %v", err)
			return err
		}

	*/

	if err := python.PreparePythonDependenciesEnv_V1(); err != nil {
		log.Error("Failed to prepare Python dependencies environment: %v", err)
		return err
	}
	log.Info("Python dependencies updated successfully.")
	return nil
}

func Run() {

	// init config
	initConfig()
	// init dependencies, it will cost some times
	go initDependencies()
	initServer()
}

/*
func Run() {

	initServer()
}
*/
