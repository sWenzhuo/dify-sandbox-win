package python

import (
	"dify-sandbox-win/internal/static"
	"testing"
)

// 测试必要的包安装检测代码
func TestValidateRequiredPackages(t *testing.T) {
	static.InitConfig("D:\\myproject\\dify-sandbox-win\\conf\\config.yaml")
	validateRequiredPackages(static.GetDifySandboxGlobalConfigurations())
}

// 测试额外的包安装检测代码
func TestVaildateRequirementsFile(t *testing.T) {
	static.InitConfig("D:\\myproject\\dify-sandbox-win\\conf\\config.yaml")
	validateRequirementsFile(static.GetDifySandboxGlobalConfigurations())
}
