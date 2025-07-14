package python

import (
	"dify-sandbox-win/internal/static"
	"testing"
)

//安装依赖测试

func TestInstallDependencies(t *testing.T) {
	static.InitConfig("D:\\myproject\\dify-sandbox-win\\conf\\config.yaml")
	InstallDependenciesV1("D:\\myproject\\dify-sandbox-win\\conf\\requirements.txt")
}
