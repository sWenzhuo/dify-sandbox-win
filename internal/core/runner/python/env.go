package python

import (
	"dify-sandbox-win/internal/core/runner"
	"dify-sandbox-win/internal/static"
	"dify-sandbox-win/internal/types"
	"dify-sandbox-win/internal/utils/log"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path"
)

//go:embed env.sh
var env_script string

//go:embed checkNess.py
var checkNess_script_path string

//go:embed checkReq.py
var checkReq_script_path string

func validatePythonInterpreter(pythonPath string) error {
	cmd := exec.Command(pythonPath, "--version")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("python interpreter check failed: %v, output: %s", err, string(output))
	}
	return nil
}

func validatePythonLibPaths(libPaths []string) error {
	for _, libPath := range libPaths {
		if _, err := os.Stat(libPath); err != nil {
			log.Warn("python lib path %s is not available: %v", libPath, err)
			continue
		}
	}
	return nil
}
func validateRequiredPackages(config types.DifySandboxGlobalConfigurations) error {
	checkNess_script_path = "D:\\myproject\\dify-sandbox-win\\internal\\core\\runner\\python\\checkNess.py"

	args := append([]string{checkNess_script_path}, config.RequirementsPages...)

	cmd := exec.Command(config.PythonPath, args...)
	output, err := cmd.CombinedOutput() //查看脚本输出，并查看是否运行成功
	if err != nil {
		return fmt.Errorf("python interpreter check failed: %v, output: %s", err, string(output))
	}
	return nil
}

func validateRequirementsFile(config types.DifySandboxGlobalConfigurations) error {
	requirementsPath := config.RequirementsFile
	scriptPath := "D:\\myproject\\dify-sandbox-win\\internal\\core\\runner\\python\\checkReq.py" // 脚本路径

	cmd := exec.Command(config.PythonPath, scriptPath, requirementsPath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("requirements.txt 校验失败：\n%s", string(output))
	}
	fmt.Println("requirements.txt 校验通过：\n" + string(output))
	return nil
}

func PreparePythonDependenciesEnv_V1() error {
	config := static.GetDifySandboxGlobalConfigurations()
	// 1. 验证 Python 解释器
	if err := validatePythonInterpreter(config.PythonPath); err != nil {
		return fmt.Errorf("python interpreter validation failed: %v", err)
	}
	// 2. 验证 Python 库路径
	if err := validatePythonLibPaths(config.PythonLibPaths); err != nil {
		return fmt.Errorf("python lib paths validation failed: %v", err)
	}
	// 3. 验证必要的包是否已经安装
	if err := validateRequiredPackages(config); err != nil {
		return fmt.Errorf("required packages validation failed: %v", err)
	}
	// 4，验证requirements里面的包是否安装
	if err := validateRequirementsFile(config); err != nil {
		return fmt.Errorf("requirements file packages validation failed: %v", err)
	}
	return nil
}

func PreparePythonDependenciesEnv() error {
	config := static.GetDifySandboxGlobalConfigurations()
	runner := runner.TempDirRunner{}
	err := runner.WithTempDir("/", []string{}, func(root_path string) error {
		err := os.WriteFile(path.Join(root_path, "env.sh"), []byte(env_script), 0755)
		if err != nil {
			return err
		}

		for _, lib_path := range config.PythonLibPaths {
			// check if the lib path is available
			if _, err := os.Stat(lib_path); err != nil {
				log.Warn("python lib path %s is not available", lib_path)
				continue
			}
			exec_cmd := exec.Command(
				"bash",
				path.Join(root_path, "env.sh"),
				lib_path,
				LIB_PATH,
			)
			exec_cmd.Stderr = os.Stderr

			if err := exec_cmd.Run(); err != nil {
				return err
			}
		}

		os.RemoveAll(root_path)
		return nil
	})

	return err
}
