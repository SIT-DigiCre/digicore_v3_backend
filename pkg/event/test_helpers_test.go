package event

import (
	"database/sql"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

var _ db.TransactionClient = (*fakeTransactionClient)(nil)

type fakeTransactionClient struct {
	selectFunc          func(dest interface{}, queryPath string, params interface{}) error
	execFunc            func(queryPath string, params interface{}, generateId bool) (sql.Result, error)
	duplicateUpdateFunc func(insertQueryPath string, updateQueryPath string, params interface{}) (sql.Result, bool, error)
}

func (f *fakeTransactionClient) Select(dest interface{}, queryPath string, params interface{}) error {
	if f.selectFunc != nil {
		return f.selectFunc(dest, queryPath, params)
	}
	return nil
}

func (f *fakeTransactionClient) Exec(queryPath string, params interface{}, generateId bool) (sql.Result, error) {
	if f.execFunc != nil {
		return f.execFunc(queryPath, params, generateId)
	}
	return fakeResult{}, nil
}

func (f *fakeTransactionClient) GetId() (string, error) {
	return "", nil
}

func (f *fakeTransactionClient) DuplicateUpdate(insertQueryPath string, updateQueryPath string, params interface{}) (sql.Result, bool, error) {
	if f.duplicateUpdateFunc != nil {
		return f.duplicateUpdateFunc(insertQueryPath, updateQueryPath, params)
	}
	return fakeResult{}, false, nil
}

type fakeResult struct {
	rowsAffected int64
	rowsErr      error
}

func (f fakeResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (f fakeResult) RowsAffected() (int64, error) {
	return f.rowsAffected, f.rowsErr
}

// 予約期間内かつイベント終了前の予約枠を返す。
func validEventReservation() eventReservation {
	now := time.Now()
	return eventReservation{
		EventId:               "event-id",
		ReservationId:         "reservation-id",
		Name:                  "テスト予約枠",
		Description:           "説明",
		Capacity:              10,
		FreeCapacity:          5,
		StartDate:             now.Add(1 * time.Hour).Format(time.RFC3339),
		FinishDate:            now.Add(2 * time.Hour).Format(time.RFC3339),
		ReservationStartDate:  now.Add(-1 * time.Hour).Format(time.RFC3339),
		ReservationFinishDate: now.Add(1 * time.Hour).Format(time.RFC3339),
	}
}

func newReservationContext(method string) echo.Context {
	e := echo.New()
	ctx := e.NewContext(httptest.NewRequest(method, "/", nil), httptest.NewRecorder())
	ctx.Set("user_id", "user-id")
	return ctx
}

// 予約枠取得→予約者一覧取得のSelectに応答するselectFuncを返す。
func reservationSelectFunc(t *testing.T, reservations []eventReservation) func(dest interface{}, queryPath string, params interface{}) error {
	t.Helper()
	return func(dest interface{}, queryPath string, params interface{}) error {
		switch queryPath {
		case "sql/event/select_event_reservation_from_event_id_reservation_id.sql":
			rows := dest.(*[]eventReservation)
			*rows = reservations
		case "sql/event/select_event_reservation_user_from_event_id_reservation_id.sql":
			users := dest.(*[]eventReservationObjectUser)
			*users = []eventReservationObjectUser{}
		default:
			t.Fatalf("unexpected query path: %s", queryPath)
		}
		return nil
	}
}
