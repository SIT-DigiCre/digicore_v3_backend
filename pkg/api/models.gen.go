// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
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

// ReqPostMattermostCreateuser defines model for ReqPostMattermostCreateuser.
type ReqPostMattermostCreateuser struct {
	Nickname string `ja:"ニックネーム" json:"nickname" validate:"required,min=3,max=22"`
	Password string `ja:"パスワード" json:"password" validate:"required,min=8,max=64"`
	Username string `ja:"ユーザー名" json:"username" validate:"required,min=3,max=22"`
}

// ReqPostSignupCallback defines model for ReqPostSignupCallback.
type ReqPostSignupCallback struct {
	Code string `ja:"認証コード" json:"code" validate:"required"`
}

// ReqPostStorageMyfile defines model for ReqPostStorageMyfile.
type ReqPostStorageMyfile struct {
	File     string `ja:"ファイル" json:"file" validate:"required,max=104857600"`
	IsPublic bool   `ja:"公開" json:"isPublic" validate:""`
	Name     string `ja:"ファイル名" json:"name" validate:"required,max=255"`
}

// ReqPostWorkTag defines model for ReqPostWorkTag.
type ReqPostWorkTag struct {
	Description string `ja:"説明" json:"description" validate:"required"`
	Name        string `ja:"タグ名" json:"name" validate:"required"`
}

// ReqPostWorkWork defines model for ReqPostWorkWork.
type ReqPostWorkWork struct {
	Authors     []string `ja:"作者" json:"authors" validate:"dive,uuid"`
	Description string   `ja:"説明" json:"description" validate:"required"`
	Files       []string `ja:"ファイル" json:"files" validate:"dive,uuid"`
	Name        string   `ja:"作品名" json:"name" validate:"required"`
	Tags        []string `ja:"タグ" json:"tags" validate:"dive,uuid"`
}

// ReqPutEventEventIdReservationIdMe defines model for ReqPutEventEventIdReservationIdMe.
type ReqPutEventEventIdReservationIdMe struct {
	Comment string `ja:"コメント" json:"comment" validate:"max=255"`
	Url     string `ja:"URL" json:"url" validate:"max=255"`
}

// ReqPutUserMe defines model for ReqPutUserMe.
type ReqPutUserMe struct {
	IconUrl           string `ja:"アイコンURL" json:"iconUrl" validate:"required,min=1,max=255"`
	SchoolGrade       int    `ja:"学年" json:"schoolGrade" validate:"required,min=1,max=9"`
	ShortIntroduction string `ja:"短い自己紹介" json:"shortIntroduction" validate:"required,min=1,max=255"`
	Username          string `ja:"ユーザー名" json:"username" validate:"required,min=1,max=255"`
}

// ReqPutUserMeDiscordCallback defines model for ReqPutUserMeDiscordCallback.
type ReqPutUserMeDiscordCallback struct {
	Code string `ja:"認証コード" json:"code" validate:"required"`
}

// ReqPutUserMeIntroduction defines model for ReqPutUserMeIntroduction.
type ReqPutUserMeIntroduction struct {
	Introduction string `ja:"自己紹介" json:"introduction" validate:""`
}

// ReqPutUserMePayment defines model for ReqPutUserMePayment.
type ReqPutUserMePayment struct {
	TransferName string `ja:"振込名義" json:"transferName" validate:"required,min=1,max=255"`
}

// ReqPutUserMePrivate defines model for ReqPutUserMePrivate.
type ReqPutUserMePrivate struct {
	Address               string  `ja:"住所" json:"address" validate:"required,min=1,max=255"`
	FirstName             string  `ja:"名前" json:"firstName" validate:"required,min=1,max=255"`
	FirstNameKana         string  `ja:"名前(カナ)" json:"firstNameKana" validate:"required,min=1,max=255"`
	IsMale                bool    `ja:"性別" json:"isMale" validate:""`
	LastName              string  `ja:"名字" json:"lastName" validate:"required,min=1,max=255"`
	LastNameKana          string  `ja:"名字(カナ)" json:"lastNameKana" validate:"required,min=1,max=255"`
	ParentAddress         string  `ja:"緊急連絡先住所" json:"parentAddress" validate:"required,min=1,max=255"`
	ParentCellphoneNumber string  `ja:"緊急連絡先携帯電話番号" json:"parentCellphoneNumber" validate:"required,phonenumber"`
	ParentHomephoneNumber *string `ja:"緊急連絡先固定電話番号" json:"parentHomephoneNumber,omitempty" validate:"phonenumber"`
	ParentName            string  `ja:"緊急連絡先氏名" json:"parentName" validate:"required,min=1,max=255"`
	PhoneNumber           string  `ja:"電話番号" json:"phoneNumber" validate:"required,phonenumber"`
}

// ReqPutWorkTagTagId defines model for ReqPutWorkTagTagId.
type ReqPutWorkTagTagId struct {
	Description string `ja:"説明" json:"description" validate:"required"`
	Name        string `ja:"タグ名" json:"name" validate:"required"`
}

// ReqPutWorkWorkWorkId defines model for ReqPutWorkWorkWorkId.
type ReqPutWorkWorkWorkId struct {
	Authors     []string `ja:"作者" json:"authors" validate:"dive,uuid"`
	Description string   `ja:"説明" json:"description" validate:"required"`
	Files       []string `ja:"ファイル" json:"files" validate:"dive,uuid"`
	Name        string   `ja:"作品名" json:"name" validate:"required"`
	Tags        []string `ja:"タグ" json:"tags" validate:"dive,uuid"`
}

// ResGetBudget defines model for ResGetBudget.
type ResGetBudget struct {
	Applicant  ResGetBudgetObjectApplicant `json:"applicant"`
	BudgetId   string                      `json:"budgetId"`
	Class      string                      `json:"class"`
	Settlement string                      `json:"settlement"`
	Status     string                      `json:"status"`
	Title      string                      `json:"title"`
	UpdatedAt  string                      `json:"updatedAt"`
}

// ResGetBudgetObjectApplicant defines model for ResGetBudgetObjectApplicant.
type ResGetBudgetObjectApplicant struct {
	IconUrl  string `json:"iconUrl"`
	UserId   string `json:"userId"`
	Username string `json:"username"`
}

// ResGetEvent defines model for ResGetEvent.
type ResGetEvent struct {
	Events []ResGetEventObjectEvent `json:"events"`
}

// ResGetEventEventId defines model for ResGetEventEventId.
type ResGetEventEventId struct {
	CalendarView bool                                  `json:"calendarView"`
	Description  string                                `json:"description"`
	EventId      string                                `json:"eventId"`
	Name         string                                `json:"name"`
	Reservable   bool                                  `json:"reservable"`
	Reservated   bool                                  `json:"reservated"`
	Reservations []ResGetEventEventIdObjectReservation `json:"reservations"`
}

// ResGetEventEventIdObjectReservation defines model for ResGetEventEventIdObjectReservation.
type ResGetEventEventIdObjectReservation struct {
	Capacity              int    `json:"capacity"`
	Description           string `json:"description"`
	FinishDate            string `json:"finishDate"`
	FreeCapacity          int    `json:"freeCapacity"`
	Name                  string `json:"name"`
	Reservable            bool   `json:"reservable"`
	Reservated            bool   `json:"reservated"`
	ReservationFinishDate string `json:"reservationFinishDate"`
	ReservationId         string `json:"reservationId"`
	ReservationStartDate  string `json:"reservationStartDate"`
	StartDate             string `json:"startDate"`
}

// ResGetEventEventIdReservationId defines model for ResGetEventEventIdReservationId.
type ResGetEventEventIdReservationId struct {
	Capacity              int                                         `json:"capacity"`
	Description           string                                      `json:"description"`
	EventId               string                                      `json:"eventId"`
	FinishDate            string                                      `json:"finishDate"`
	FreeCapacity          int                                         `json:"freeCapacity"`
	Name                  string                                      `json:"name"`
	Reservable            bool                                        `json:"reservable"`
	Reservated            bool                                        `json:"reservated"`
	ReservationFinishDate string                                      `json:"reservationFinishDate"`
	ReservationId         string                                      `json:"reservationId"`
	ReservationStartDate  string                                      `json:"reservationStartDate"`
	StartDate             string                                      `json:"startDate"`
	Users                 []ResGetEventEventIdReservationIdObjectUser `json:"users"`
}

// ResGetEventEventIdReservationIdObjectUser defines model for ResGetEventEventIdReservationIdObjectUser.
type ResGetEventEventIdReservationIdObjectUser struct {
	Comment  string `json:"comment"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	UserIcon string `json:"userIcon"`
	UserId   string `json:"userId"`
}

// ResGetEventObjectEvent defines model for ResGetEventObjectEvent.
type ResGetEventObjectEvent struct {
	CalendarView bool   `json:"calendarView"`
	Description  string `json:"description"`
	EventId      string `json:"eventId"`
	Name         string `json:"name"`
	Reservable   bool   `json:"reservable"`
	Reservated   bool   `json:"reservated"`
}

// ResGetGroup defines model for ResGetGroup.
type ResGetGroup struct {
	Groups []ResGetGroupObjectGroup `json:"groups"`
}

// ResGetGroupGroupId defines model for ResGetGroupGroupId.
type ResGetGroupGroupId struct {
	Description string                         `json:"description"`
	GroupId     string                         `json:"groupId"`
	Joinable    bool                           `json:"joinable"`
	Joined      bool                           `json:"joined"`
	Name        string                         `json:"name"`
	UserCount   int                            `json:"userCount"`
	Users       []ResGetGroupGroupIdObjectUser `json:"users"`
}

// ResGetGroupGroupIdObjectUser defines model for ResGetGroupGroupIdObjectUser.
type ResGetGroupGroupIdObjectUser struct {
	Name     string `json:"name"`
	UserIcon string `json:"userIcon"`
	UserId   string `json:"userId"`
}

// ResGetGroupObjectGroup defines model for ResGetGroupObjectGroup.
type ResGetGroupObjectGroup struct {
	GroupId   string `json:"groupId"`
	Joinable  bool   `json:"joinable"`
	Joined    bool   `json:"joined"`
	Name      string `json:"name"`
	UserCount int    `json:"userCount"`
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

// ResGetStorageFileId defines model for ResGetStorageFileId.
type ResGetStorageFileId struct {
	CreatedAt string `json:"createdAt"`
	Extension string `json:"extension"`
	FileId    string `json:"fileId"`
	IsPublic  bool   `json:"isPublic"`
	KSize     string `json:"kSize"`
	Name      string `json:"name"`
	UpdatedAt string `json:"updatedAt"`
	Url       string `json:"url"`
	UserId    string `json:"userId"`
}

// ResGetStorageMyfile defines model for ResGetStorageMyfile.
type ResGetStorageMyfile struct {
	Files []ResGetStorageMyfileObjectFile `json:"files"`
}

// ResGetStorageMyfileObjectFile defines model for ResGetStorageMyfileObjectFile.
type ResGetStorageMyfileObjectFile struct {
	CreatedAt string `json:"createdAt"`
	Extension string `json:"extension"`
	FileId    string `json:"fileId"`
	IsPublic  bool   `json:"isPublic"`
	KSize     string `json:"kSize"`
	Name      string `json:"name"`
	UpdatedAt string `json:"updatedAt"`
	UserId    string `json:"userId"`
}

// ResGetTool defines model for ResGetTool.
type ResGetTool struct {
	DiscordUrl string `json:"discordUrl"`
}

// ResGetUser defines model for ResGetUser.
type ResGetUser struct {
	Users []ResGetUserObjectUser `json:"users"`
}

// ResGetUserMe defines model for ResGetUserMe.
type ResGetUserMe struct {
	ActiveLimit       string `json:"activeLimit"`
	DiscordUserId     string `json:"discordUserId"`
	IconUrl           string `json:"iconUrl"`
	SchoolGrade       int    `json:"schoolGrade"`
	ShortIntroduction string `json:"shortIntroduction"`
	StudentNumber     string `json:"studentNumber"`
	UserId            string `json:"userId"`
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
	Histories []ResGetUserMePaymentObjectHistory `json:"histories"`
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
	UserId            string `json:"userId"`
	Username          string `json:"username"`
}

// ResGetUserUserId defines model for ResGetUserUserId.
type ResGetUserUserId struct {
	ActiveLimit       string `json:"activeLimit"`
	DiscordUserId     string `json:"discordUserId"`
	IconUrl           string `json:"iconUrl"`
	SchoolGrade       int    `json:"schoolGrade"`
	ShortIntroduction string `json:"shortIntroduction"`
	StudentNumber     string `json:"studentNumber"`
	UserId            string `json:"userId"`
	Username          string `json:"username"`
}

// ResGetUserUserIdIntroduction defines model for ResGetUserUserIdIntroduction.
type ResGetUserUserIdIntroduction struct {
	Introduction string `json:"introduction"`
}

// ResGetWorkTag defines model for ResGetWorkTag.
type ResGetWorkTag struct {
	Tags []ResGetWorkTagObjectTag `json:"tags"`
}

// ResGetWorkTagObjectTag defines model for ResGetWorkTagObjectTag.
type ResGetWorkTagObjectTag struct {
	Name  string `json:"name"`
	TagId string `json:"tagId"`
}

// ResGetWorkTagTagId defines model for ResGetWorkTagTagId.
type ResGetWorkTagTagId struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	TagId       string `json:"tagId"`
}

// ResGetWorkWork defines model for ResGetWorkWork.
type ResGetWorkWork struct {
	Works []ResGetWorkWorkObjectWork `json:"works"`
}

// ResGetWorkWorkObjectWork defines model for ResGetWorkWorkObjectWork.
type ResGetWorkWorkObjectWork struct {
	Authors []ResGetWorkWorkObjectWorkObjectAuthor `json:"authors"`
	Name    string                                 `json:"name"`
	Tags    []ResGetWorkWorkObjectWorkObjectTag    `json:"tags"`
	WorkId  string                                 `json:"workId"`
}

// ResGetWorkWorkObjectWorkObjectAuthor defines model for ResGetWorkWorkObjectWorkObjectAuthor.
type ResGetWorkWorkObjectWorkObjectAuthor struct {
	IconUrl  string `json:"iconUrl"`
	UserId   string `json:"userId"`
	Username string `json:"username"`
}

// ResGetWorkWorkObjectWorkObjectTag defines model for ResGetWorkWorkObjectWorkObjectTag.
type ResGetWorkWorkObjectWorkObjectTag struct {
	Name  string `json:"name"`
	TagId string `json:"tagId"`
}

// ResGetWorkWorkWorkId defines model for ResGetWorkWorkWorkId.
type ResGetWorkWorkWorkId struct {
	Authors     []ResGetWorkWorkWorkIdObjectAuthor `json:"authors"`
	Description string                             `json:"description"`
	Files       []ResGetWorkWorkWorkIdObjectTag    `json:"files"`
	Name        string                             `json:"name"`
	Tags        []ResGetWorkWorkWorkIdObjectFile   `json:"tags"`
	WorkId      string                             `json:"workId"`
}

// ResGetWorkWorkWorkIdObjectAuthor defines model for ResGetWorkWorkWorkIdObjectAuthor.
type ResGetWorkWorkWorkIdObjectAuthor struct {
	IconUrl  string `json:"iconUrl"`
	UserId   string `json:"userId"`
	Username string `json:"username"`
}

// ResGetWorkWorkWorkIdObjectFile defines model for ResGetWorkWorkWorkIdObjectFile.
type ResGetWorkWorkWorkIdObjectFile struct {
	Name  string `json:"name"`
	TagId string `json:"tagId"`
}

// ResGetWorkWorkWorkIdObjectTag defines model for ResGetWorkWorkWorkIdObjectTag.
type ResGetWorkWorkWorkIdObjectTag struct {
	FileId string `json:"fileId"`
	Name   string `json:"name"`
}

// ResPostLoginCallback defines model for ResPostLoginCallback.
type ResPostLoginCallback struct {
	Jwt string `json:"jwt"`
}

// ResPostMattermostCreateuser defines model for ResPostMattermostCreateuser.
type ResPostMattermostCreateuser struct {
	Username string `json:"username"`
}

// ResPostSignupCallback defines model for ResPostSignupCallback.
type ResPostSignupCallback struct {
	Jwt string `json:"jwt"`
}

// Success defines model for Success.
type Success struct {
	Success bool `json:"success"`
}

// BadRequest defines model for BadRequest.
type BadRequest = Error

// BlankSuccess defines model for BlankSuccess.
type BlankSuccess = Success

// InternalServer defines model for InternalServer.
type InternalServer = Error

// NotFound defines model for NotFound.
type NotFound = Error

// Unauthorized defines model for Unauthorized.
type Unauthorized = Error

// GetBudgetParams defines parameters for GetBudget.
type GetBudgetParams struct {
	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`
}

// GetEventParams defines parameters for GetEvent.
type GetEventParams struct {
	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`
}

// GetGroupParams defines parameters for GetGroup.
type GetGroupParams struct {
	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`
	Seed   *int `form:"seed,omitempty" json:"seed,omitempty"`
}

// GetUserParams defines parameters for GetUser.
type GetUserParams struct {
	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`
	Seed   *int `form:"seed,omitempty" json:"seed,omitempty"`
}

// GetWorkTagParams defines parameters for GetWorkTag.
type GetWorkTagParams struct {
	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`
}

// GetWorkWorkParams defines parameters for GetWorkWork.
type GetWorkWorkParams struct {
	Offset   *int    `form:"offset,omitempty" json:"offset,omitempty"`
	AuthorId *string `form:"authorId,omitempty" json:"authorId,omitempty"`
}

// PutEventEventIdReservationIdMeJSONRequestBody defines body for PutEventEventIdReservationIdMe for application/json ContentType.
type PutEventEventIdReservationIdMeJSONRequestBody = ReqPutEventEventIdReservationIdMe

// PostLoginCallbackJSONRequestBody defines body for PostLoginCallback for application/json ContentType.
type PostLoginCallbackJSONRequestBody = ReqPostLoginCallback

// PostMattermostCreateUserJSONRequestBody defines body for PostMattermostCreateUser for application/json ContentType.
type PostMattermostCreateUserJSONRequestBody = ReqPostMattermostCreateuser

// PostSignupCallbackJSONRequestBody defines body for PostSignupCallback for application/json ContentType.
type PostSignupCallbackJSONRequestBody = ReqPostSignupCallback

// PostStorageMyfileJSONRequestBody defines body for PostStorageMyfile for application/json ContentType.
type PostStorageMyfileJSONRequestBody = ReqPostStorageMyfile

// PutUserMeJSONRequestBody defines body for PutUserMe for application/json ContentType.
type PutUserMeJSONRequestBody = ReqPutUserMe

// PutUserMeDiscordCallbackJSONRequestBody defines body for PutUserMeDiscordCallback for application/json ContentType.
type PutUserMeDiscordCallbackJSONRequestBody = ReqPutUserMeDiscordCallback

// PutUserMeIntroductionJSONRequestBody defines body for PutUserMeIntroduction for application/json ContentType.
type PutUserMeIntroductionJSONRequestBody = ReqPutUserMeIntroduction

// PutUserMePaymentJSONRequestBody defines body for PutUserMePayment for application/json ContentType.
type PutUserMePaymentJSONRequestBody = ReqPutUserMePayment

// PutUserMePrivateJSONRequestBody defines body for PutUserMePrivate for application/json ContentType.
type PutUserMePrivateJSONRequestBody = ReqPutUserMePrivate

// PostWorkTagJSONRequestBody defines body for PostWorkTag for application/json ContentType.
type PostWorkTagJSONRequestBody = ReqPostWorkTag

// PutWorkTagTagIdJSONRequestBody defines body for PutWorkTagTagId for application/json ContentType.
type PutWorkTagTagIdJSONRequestBody = ReqPutWorkTagTagId

// PostWorkWorkJSONRequestBody defines body for PostWorkWork for application/json ContentType.
type PostWorkWorkJSONRequestBody = ReqPostWorkWork

// PutWorkWorkWorkIdJSONRequestBody defines body for PutWorkWorkWorkId for application/json ContentType.
type PutWorkWorkWorkIdJSONRequestBody = ReqPutWorkWorkWorkId
