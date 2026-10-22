package env

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestGetWarnsWhenNotSet(t *testing.T) {
	hook := test.NewGlobal()

	t.Setenv("DIGICORE_ENV_TEST_UNSET", "")
	Get("DIGICORE_ENV_TEST_UNSET")
	entry := hook.LastEntry()
	if entry == nil || entry.Level != logrus.WarnLevel {
		t.Fatal("空文字の環境変数で警告が表示されませんでした")
	}
	if entry.Message != "環境変数 DIGICORE_ENV_TEST_UNSET が設定されていません" {
		t.Fatalf("警告メッセージが想定と異なります: %s", entry.Message)
	}

	t.Setenv("DIGICORE_ENV_TEST_UNSET", "value")
	Get("DIGICORE_ENV_TEST_UNSET")
	if len(hook.AllEntries()) != 1 {
		t.Fatal("設定済みの環境変数で警告が表示されました")
	}
}
