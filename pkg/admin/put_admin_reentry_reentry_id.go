package admin

import (
	"fmt"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/mail"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func PutAdminReentryReentryId(ctx echo.Context, dbClient db.TransactionClient, reentryId string, requestBody api.ReqPutAdminReentryReentryId) (api.BlankSuccess, string, *response.Error) {
	checkerId := ctx.Get("user_id").(string)

	detailParams := struct {
		ReentryId string `twowaysql:"reentryId"`
	}{ReentryId: reentryId}
	details := []reentryDetail{}
	err := dbClient.Select(&details, "sql/reentry/select_reentry_by_id.sql", &detailParams)
	if err != nil {
		return api.BlankSuccess{}, "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "再入部申請の取得に失敗しました", Log: err.Error()}
	}
	if len(details) == 0 {
		return api.BlankSuccess{}, "", &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "指定された再入部申請が見つかりません", Log: "reentry not found"}
	}

	note := ""
	if requestBody.Note != nil {
		note = *requestBody.Note
	}
	if requestBody.Status == "approved" {
		checked, cerr := isCurrentSchoolYearPaymentChecked(dbClient, details[0].UserId)
		if cerr != nil {
			return api.BlankSuccess{}, "", cerr
		}
		if !checked {
			return api.BlankSuccess{}, "", &response.Error{
				Code:    http.StatusBadRequest,
				Level:   "Info",
				Message: "今年度の部費振込報告が会計承認されていないため、再入部を承認できません",
				Log:     "payment is not approved",
			}
		}
	}

	updateParams := struct {
		ReentryId string `twowaysql:"reentryId"`
		Status    string `twowaysql:"status"`
		Note      string `twowaysql:"note"`
		CheckedBy string `twowaysql:"checkedBy"`
	}{
		ReentryId: reentryId,
		Status:    requestBody.Status,
		Note:      note,
		CheckedBy: checkerId,
	}
	result, rerr := dbClient.Exec("sql/reentry/update_reentry_status.sql", &updateParams, false)
	if rerr != nil {
		return api.BlankSuccess{}, "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "再入部申請の更新に失敗しました", Log: rerr.Error()}
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return api.BlankSuccess{}, "", &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "指定された申請が見つからないか、既に処理済みです", Log: "reentry not found or already processed"}
	}

	reentry := details[0]
	if requestBody.Status == "approved" {
		memberParams := struct {
			UserId   string `twowaysql:"userId"`
			IsMember bool   `twowaysql:"isMember"`
		}{
			UserId:   reentry.UserId,
			IsMember: true,
		}
		_, rerr = dbClient.Exec("sql/reentry/update_user_member_status.sql", &memberParams, false)
		if rerr != nil {
			return api.BlankSuccess{}, "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "ユーザー状態の更新に失敗しました", Log: rerr.Error()}
		}
	}

	studentNumber, serr := getStudentNumberByUserId(dbClient, reentry.UserId)
	if serr != nil {
		logrus.Errorf("再入部通知メール送信用の学籍番号取得に失敗しました(%s): %v", reentry.UserId, serr)
		return api.BlankSuccess{}, "", &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "再入部通知メール送信用の学籍番号取得に失敗しました",
			Log:     serr.Log,
		}
	}

	return api.BlankSuccess{Success: true}, studentNumber, nil
}

func NotifyReentryDecision(studentNumber string, status string, note string) {
	address := fmt.Sprintf("%s@shibaura-it.ac.jp", studentNumber)
	subject := "再入部申請結果のお知らせ"
	body := "再入部申請の審査結果をお知らせします。\n"
	if status == "approved" {
		body += "結果: 承認\n再入部手続きが完了しました。"
	} else {
		body += "結果: 却下"
		if note != "" {
			body += "\n備考: " + note
		}
	}
	if err := mail.SendEmail(address, subject, body); err != nil {
		logrus.Warnf("再入部判定メール送信に失敗しました(%s): %v", address, err)
	}
}

func getStudentNumberByUserId(dbClient db.Client, userId string) (string, *response.Error) {
	params := struct {
		UserId string `twowaysql:"userId"`
	}{
		UserId: userId,
	}
	rows := []struct {
		StudentNumber string `db:"student_number"`
	}{}
	err := dbClient.Select(&rows, "sql/reentry/select_student_number_by_user_id.sql", &params)
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "通知先の取得に失敗しました", Log: err.Error()}
	}
	if len(rows) == 0 {
		return "", &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "通知先の取得に失敗しました", Log: "student number not found"}
	}
	return rows[0].StudentNumber, nil
}

func isCurrentSchoolYearPaymentChecked(dbClient db.Client, userId string) (bool, *response.Error) {
	params := struct {
		UserId string `twowaysql:"userId"`
		Year   int    `twowaysql:"year"`
	}{
		UserId: userId,
		Year:   utils.GetSchoolYear(),
	}
	rows := []struct {
		Checked bool `db:"checked"`
	}{}
	err := dbClient.Select(&rows, "sql/reentry/select_payment_checked_by_user_id_year.sql", &params)
	if err != nil {
		return false, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "振込報告の確認に失敗しました", Log: err.Error()}
	}
	if len(rows) == 0 {
		return false, nil
	}
	return rows[0].Checked, nil
}
