package api

type ReqPostMattermostCmd struct {
	ChannelName string `form:"channel_name" ja:"チャンネル名" validate:"required"`
	Command     string `form:"command" ja:"コマンド" validate:"required"`
	Text        string `form:"text" ja:"テキスト"`
	Token       string `form:"token" ja:"トークン" validate:"required"`
	UserName    string `form:"user_name" ja:"ユーザー名" validate:"required"`
}

type ResPostMattermostCmd struct {
	IconEmoji string `json:"icon_emoji"`
	Text      string `json:"text"`
	Username  string `json:"username"`
}
