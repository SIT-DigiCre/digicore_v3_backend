// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package api

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Error defines model for Error.
type Error struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

// ReqPostLoginCallback defines model for ReqPostLoginCallback.
type ReqPostLoginCallback struct {
	Code string `ja:"認証コード" json:"code" validate:"required"`
}

// ReqPostSignupCallback defines model for ReqPostSignupCallback.
type ReqPostSignupCallback struct {
	Code string `ja:"認証コード" json:"code" validate:"required"`
}

// ReqPostUserMeDiscordCallback defines model for ReqPostUserMeDiscordCallback.
type ReqPostUserMeDiscordCallback struct {
	Code string `ja:"認証コード" json:"code" validate:"required"`
}

// ReqPutUserMe defines model for ReqPutUserMe.
type ReqPutUserMe struct {
	IconUrl           string `ja:"アイコンURL" json:"iconUrl" validate:"required,min=1,max=255"`
	SchoolGrade       int    `ja:"学年" json:"schoolGrade" validate:"required,min=1,max=9"`
	ShortIntroduction string `ja:"短い自己紹介" json:"shortIntroduction" validate:"required,min=1,max=255"`
	Username          string `ja:"ユーザー名" json:"username" validate:"required,min=1,max=255"`
}

// ReqPutUserMeIntroduction defines model for ReqPutUserMeIntroduction.
type ReqPutUserMeIntroduction struct {
	Introduction string `ja:"自己紹介" json:"introduction" validate:"required,min=1"`
}

// ReqPutUserMePayment defines model for ReqPutUserMePayment.
type ReqPutUserMePayment struct {
	TransferName string `ja:"振込名義" json:"transferName" validate:"required,min=1,max=255"`
}

// ReqPutUserMePrivate defines model for ReqPutUserMePrivate.
type ReqPutUserMePrivate struct {
	Address               string `ja:"住所" json:"address" validate:"required,min=1,max=255"`
	FirstName             string `ja:"名前" json:"firstName" validate:"required,min=1,max=255"`
	FirstNameKana         string `ja:"名前(カナ)" json:"firstNameKana" validate:"required,min=1,max=255"`
	IsMale                bool   `ja:"性別" json:"isMale" validate:""`
	LastName              string `ja:"名字" json:"lastName" validate:"required,min=1,max=255"`
	LastNameKana          string `ja:"名字(カナ)" json:"lastNameKana" validate:"required,min=1,max=255"`
	ParentAddress         string `ja:"緊急連絡先住所" json:"parentAddress" validate:"required,min=1,max=255"`
	ParentCellphoneNumber string `ja:"緊急連絡先携帯電話番号" json:"parentCellphoneNumber" validate:"required,numeric,min=1,max=15"`
	ParentHomephoneNumber string `ja:"緊急連絡先固定電話番号" json:"parentHomephoneNumber" validate:"required,numeric,min=1,max=15"`
	ParentName            string `ja:"緊急連絡先氏名" json:"parentName" validate:"required,min=1,max=255"`
	PhoneNumber           string `ja:"電話番号" json:"phoneNumber" validate:"required,numeric,min=1,max=15"`
}

// ResGetLogin defines model for ResGetLogin.
type ResGetLogin struct {
	Url string `json:"url"`
}

// ResGetSignup defines model for ResGetSignup.
type ResGetSignup struct {
	Url string `json:"url"`
}

// ResGetStatus defines model for ResGetStatus.
type ResGetStatus struct {
	Status bool `json:"status"`
}

// ResGetUser defines model for ResGetUser.
type ResGetUser struct {
	User []ResGetUserObjectUser `json:"user"`
}

// ResGetUserMe defines model for ResGetUserMe.
type ResGetUserMe struct {
	ActiveLimit       string `json:"activeLimit"`
	DiscordUserID     string `json:"discordUserID"`
	IconUrl           string `json:"iconUrl"`
	SchoolGrade       int    `json:"schoolGrade"`
	ShortIntroduction string `json:"shortIntroduction"`
	StudentNumber     string `json:"studentNumber"`
	UserID            string `json:"userID"`
	Username          string `json:"username"`
}

// ResGetUserMeDiscord defines model for ResGetUserMeDiscord.
type ResGetUserMeDiscord struct {
	Url string `json:"url"`
}

// ResGetUserMeIntroduction defines model for ResGetUserMeIntroduction.
type ResGetUserMeIntroduction struct {
	Introduction string `json:"introduction"`
}

// ResGetUserMePayment defines model for ResGetUserMePayment.
type ResGetUserMePayment struct {
	History []ResGetUserMePaymentObjectHistory `json:"history"`
}

// ResGetUserMePaymentObjectHistory defines model for ResGetUserMePaymentObjectHistory.
type ResGetUserMePaymentObjectHistory struct {
	Checked      bool   `json:"checked"`
	TransferName string `json:"transferName"`
	UpdatedAt    string `json:"updatedAt"`
	Year         int    `json:"year"`
}

// ResGetUserMePrivate defines model for ResGetUserMePrivate.
type ResGetUserMePrivate struct {
	Address               string `json:"address"`
	FirstName             string `json:"firstName"`
	FirstNameKana         string `json:"firstNameKana"`
	IsMale                bool   `json:"isMale"`
	LastName              string `json:"lastName"`
	LastNameKana          string `json:"lastNameKana"`
	ParentAddress         string `json:"parentAddress"`
	ParentCellphoneNumber string `json:"parentCellphoneNumber"`
	ParentHomephoneNumber string `json:"parentHomephoneNumber"`
	ParentName            string `json:"parentName"`
	PhoneNumber           string `json:"phoneNumber"`
}

// ResGetUserObjectUser defines model for ResGetUserObjectUser.
type ResGetUserObjectUser struct {
	IconUrl           string `json:"iconUrl"`
	ShortIntroduction string `json:"shortIntroduction"`
	UserID            string `json:"userID"`
	Username          string `json:"username"`
}

// ResGetUserUserID defines model for ResGetUserUserID.
type ResGetUserUserID struct {
	ActiveLimit       string `json:"activeLimit"`
	DiscordUserID     string `json:"discordUserID"`
	IconUrl           string `json:"iconUrl"`
	SchoolGrade       int    `json:"schoolGrade"`
	ShortIntroduction string `json:"shortIntroduction"`
	StudentNumber     string `json:"studentNumber"`
	UserID            string `json:"userID"`
	Username          string `json:"username"`
}

// ResPostLoginCallback defines model for ResPostLoginCallback.
type ResPostLoginCallback struct {
	Jwt string `json:"jwt"`
}

// ResPostSignupCallback defines model for ResPostSignupCallback.
type ResPostSignupCallback struct {
	Jwt string `json:"jwt"`
}

// BadRequest defines model for BadRequest.
type BadRequest = Error

// InternalServer defines model for InternalServer.
type InternalServer = Error

// NotFound defines model for NotFound.
type NotFound = Error

// Unauthorized defines model for Unauthorized.
type Unauthorized = Error

// PostLoginCallbackJSONBody defines parameters for PostLoginCallback.
type PostLoginCallbackJSONBody = ReqPostLoginCallback

// PostSignupCallbackJSONBody defines parameters for PostSignupCallback.
type PostSignupCallbackJSONBody = ReqPostSignupCallback

// GetUserParams defines parameters for GetUser.
type GetUserParams struct {
	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`
	Seed   *int `form:"seed,omitempty" json:"seed,omitempty"`
}

// PutUserMeJSONBody defines parameters for PutUserMe.
type PutUserMeJSONBody = ReqPutUserMe

// PutUserMeDiscordCallbackJSONBody defines parameters for PutUserMeDiscordCallback.
type PutUserMeDiscordCallbackJSONBody = ReqPostUserMeDiscordCallback

// PutUserMeIntroductionJSONBody defines parameters for PutUserMeIntroduction.
type PutUserMeIntroductionJSONBody = ReqPutUserMeIntroduction

// PutUserMePaymentJSONBody defines parameters for PutUserMePayment.
type PutUserMePaymentJSONBody = ReqPutUserMePayment

// PutUserMePrivateJSONBody defines parameters for PutUserMePrivate.
type PutUserMePrivateJSONBody = ReqPutUserMePrivate

// PostLoginCallbackJSONRequestBody defines body for PostLoginCallback for application/json ContentType.
type PostLoginCallbackJSONRequestBody = PostLoginCallbackJSONBody

// PostSignupCallbackJSONRequestBody defines body for PostSignupCallback for application/json ContentType.
type PostSignupCallbackJSONRequestBody = PostSignupCallbackJSONBody

// PutUserMeJSONRequestBody defines body for PutUserMe for application/json ContentType.
type PutUserMeJSONRequestBody = PutUserMeJSONBody

// PutUserMeDiscordCallbackJSONRequestBody defines body for PutUserMeDiscordCallback for application/json ContentType.
type PutUserMeDiscordCallbackJSONRequestBody = PutUserMeDiscordCallbackJSONBody

// PutUserMeIntroductionJSONRequestBody defines body for PutUserMeIntroduction for application/json ContentType.
type PutUserMeIntroductionJSONRequestBody = PutUserMeIntroductionJSONBody

// PutUserMePaymentJSONRequestBody defines body for PutUserMePayment for application/json ContentType.
type PutUserMePaymentJSONRequestBody = PutUserMePaymentJSONBody

// PutUserMePrivateJSONRequestBody defines body for PutUserMePrivate for application/json ContentType.
type PutUserMePrivateJSONRequestBody = PutUserMePrivateJSONBody
