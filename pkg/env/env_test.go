package env

import "testing"

func TestMissing(t *testing.T) {
	t.Setenv("FRONTEND_ROOT_URL", "http://example.com")
	if hasMissing("FRONTEND_ROOT_URL") {
		t.Fatal("設定済みの環境変数が未設定として報告されました")
	}

	t.Setenv("FRONTEND_ROOT_URL", "")
	if !hasMissing("FRONTEND_ROOT_URL") {
		t.Fatal("空文字の環境変数が未設定として報告されませんでした")
	}
}

func hasMissing(name string) bool {
	for _, n := range Missing() {
		if n == name {
			return true
		}
	}
	return false
}
