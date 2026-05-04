package authenticator

import "sync"

var (
	nonMemberAllowedPathMutex sync.RWMutex
	nonMemberAllowedPaths     = map[string]bool{}
)

// RegisterNonMemberAllowedPath は is_member=false のユーザーにも許可するパスを登録します。
// パスは Echo のルートパターン（例: /user/me/reentry）を指定してください。
func RegisterNonMemberAllowedPath(path string) {
	nonMemberAllowedPathMutex.Lock()
	defer nonMemberAllowedPathMutex.Unlock()
	nonMemberAllowedPaths[path] = true
}

func isNonMemberAllowedPath(path string) bool {
	nonMemberAllowedPathMutex.RLock()
	defer nonMemberAllowedPathMutex.RUnlock()
	return nonMemberAllowedPaths[path]
}
