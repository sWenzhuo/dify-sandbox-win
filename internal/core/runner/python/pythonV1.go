package python

import (
	"crypto/rand"
	"dify-sandbox-win/internal/core/runner"
	"dify-sandbox-win/internal/core/runner/types"
	"dify-sandbox-win/internal/static"
	_ "embed"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// 脚本模板
//
//go:embed prescriptV1.py
var template_script []byte

func (p *PythonRunner) InitializeEnvironmentV1(code string, preload string, options *types.RunnerOptions) (string, string, error) {

	//  生成随机脚本名
	scriptName := strings.ReplaceAll(uuid.New().String(), "-", "_") + ".py"
	scriptPath := filepath.Join(LIB_PATH, scriptName)

	// 根据code进行XOR加密
	const keyLen = 64
	key := make([]byte, keyLen)
	if _, err := rand.Read(key); err != nil {
		return "", "", fmt.Errorf("生成密钥失败: %w", err)
	}
	encrypted_code := make([]byte, len(code))
	for i := 0; i < len(code); i++ {
		encrypted_code[i] = code[i] ^ key[i%keyLen]
	}
	// 根据code进行base64加密
	encodedCode := base64.StdEncoding.EncodeToString(encrypted_code)
	//根据脚本动态输入
	encodedKey := base64.StdEncoding.EncodeToString(key)

	// 填充模板
	finalScript := strings.ReplaceAll(string(template_script), "{{preload}}", preload)
	finalScript = strings.ReplaceAll(finalScript, "{{code}}", encodedCode)

	// 写入临时文件
	if err := os.WriteFile(scriptPath, []byte(finalScript), 0755); err != nil {
		return "", "", fmt.Errorf("写入脚本文件失败: %w", err)
	}

	return scriptPath, encodedKey, nil
}

func (p *PythonRunner) RunV1(
	code string,
	timeout time.Duration,
	stdin []byte,
	preload string,
	options *types.RunnerOptions,
) (chan []byte, chan []byte, chan bool, error) {
	configuration := static.GetDifySandboxGlobalConfigurations()

	// initialize the environment
	untrusted_code_path, key, err := p.InitializeEnvironmentV1(code, preload, options)
	if err != nil {
		return nil, nil, nil, err
	}
	// capture the output
	output_handler := runner.NewOutputCaptureRunner()
	output_handler.SetTimeout(timeout)
	output_handler.SetAfterExitHook(func() {
		// remove untrusted code
		os.Remove(untrusted_code_path)
	})

	// create a new process
	// path 是工作目录
	cmd := exec.Command(
		configuration.PythonPath,
		untrusted_code_path,
		LIB_PATH,
		key,
	)
	cmd.Env = []string{}
	cmd.Dir = LIB_PATH

	if configuration.Proxy.Socks5 != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("HTTPS_PROXY=%s", configuration.Proxy.Socks5))
		cmd.Env = append(cmd.Env, fmt.Sprintf("HTTP_PROXY=%s", configuration.Proxy.Socks5))
	} else if configuration.Proxy.Https != "" || configuration.Proxy.Http != "" {
		if configuration.Proxy.Https != "" {
			cmd.Env = append(cmd.Env, fmt.Sprintf("HTTPS_PROXY=%s", configuration.Proxy.Https))
		}
		if configuration.Proxy.Http != "" {
			cmd.Env = append(cmd.Env, fmt.Sprintf("HTTP_PROXY=%s", configuration.Proxy.Http))
		}
	}

	err = output_handler.CaptureOutput(cmd)
	if err != nil {
		return nil, nil, nil, err
	}

	return output_handler.GetStdout(), output_handler.GetStderr(), output_handler.GetDone(), nil
}
