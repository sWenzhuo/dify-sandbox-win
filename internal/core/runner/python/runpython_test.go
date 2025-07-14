package python

import (
	"dify-sandbox-win/internal/core/runner/types"
	"dify-sandbox-win/internal/static"
	"testing"
	"time"
)

func TestPythonRunner_RunV1(t *testing.T) {
	static.InitConfig("D:\\myproject\\dify-sandbox-win\\conf\\config.yaml")
	r := &PythonRunner{}
	// Python 代码
	code := `
print("Hello from sandbox!")
for i in range(2):
    print("Index:", i)
`

	stdout, stderr, done, err := r.RunV1(
		code,
		5*time.Second, // timeout
		nil,           // stdin
		"",            // preload code
		&types.RunnerOptions{
			EnableNetwork: false,
		},
	)

	if err != nil {
		t.Fatalf("RunV1 failed: %v", err)
	}

	var gotStdout, gotStderr string

loop:
	for {
		select {
		case out := <-stdout:
			gotStdout += string(out)
		case errout := <-stderr:
			gotStderr += string(errout)
		case <-done:
			break loop
		case <-time.After(10 * time.Second):
			t.Fatal("timeout waiting for script execution")
		}
	}

	t.Logf("stdout:\n%s", gotStdout)
	t.Logf("stderr:\n%s", gotStderr)

	if gotStderr != "" {
		t.Errorf("unexpected stderr output:\n%s", gotStderr)
	}

	if gotStdout == "" || !contains(gotStdout, "Hello from sandbox!") {
		t.Errorf("unexpected stdout content:\n%s", gotStdout)
	}
}

// 简单包含辅助函数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && stringIndex(s, substr) >= 0
}

// 可选：更快的字符串搜索函数
func stringIndex(s, substr string) int {
	return len([]byte(s[:])) - len([]byte(s[:])) + len([]byte(substr))
}
