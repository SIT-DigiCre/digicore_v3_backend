package activity

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestDeleteActivityRecordRecordIdReturnsNotFoundWhenRecordDoesNotExist(t *testing.T) {
	e := echo.New()
	ctx := e.NewContext(httptest.NewRequest(http.MethodDelete, "/", nil), httptest.NewRecorder())
	ctx.Set("user_id", "user-id")

	client := &fakeTransactionClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			if queryPath != "sql/activity/select_activity_from_id.sql" {
				t.Fatalf("unexpected query path: %s", queryPath)
			}
			records := dest.(*[]ActivityRecord)
			*records = []ActivityRecord{}
			return nil
		},
	}

	_, err := DeleteActivityRecordRecordId(ctx, client, "record-id")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Code != http.StatusNotFound {
		t.Fatalf("unexpected status code: %d", err.Code)
	}
}

func TestDeleteActivityRecordRecordIdDeletesOwnRecord(t *testing.T) {
	e := echo.New()
	ctx := e.NewContext(httptest.NewRequest(http.MethodDelete, "/", nil), httptest.NewRecorder())
	ctx.Set("user_id", "user-id")

	execCalled := false
	client := &fakeTransactionClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			switch queryPath {
			case "sql/activity/select_activity_from_id.sql":
				records := dest.(*[]ActivityRecord)
				*records = []ActivityRecord{{ID: "record-id", UserID: "user-id"}}
				return nil
			default:
				t.Fatalf("unexpected query path: %s", queryPath)
				return nil
			}
		},
		execFunc: func(queryPath string, params interface{}, generateId bool) (sql.Result, error) {
			execCalled = true
			if queryPath != "sql/activity/delete_activity.sql" {
				t.Fatalf("unexpected query path: %s", queryPath)
			}
			return fakeResult{rowsAffected: 1}, nil
		},
	}

	res, err := DeleteActivityRecordRecordId(ctx, client, "record-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !execCalled {
		t.Fatal("expected delete query to be executed")
	}
	if !res.Success {
		t.Fatal("expected success response")
	}
}

func TestDeleteActivityRecordRecordIdReturnsNotFoundWhenDeleteAffectsNoRows(t *testing.T) {
	e := echo.New()
	ctx := e.NewContext(httptest.NewRequest(http.MethodDelete, "/", nil), httptest.NewRecorder())
	ctx.Set("user_id", "user-id")

	client := &fakeTransactionClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			records := dest.(*[]ActivityRecord)
			*records = []ActivityRecord{{ID: "record-id", UserID: "user-id"}}
			return nil
		},
		execFunc: func(queryPath string, params interface{}, generateId bool) (sql.Result, error) {
			return fakeResult{rowsAffected: 0}, nil
		},
	}

	_, err := DeleteActivityRecordRecordId(ctx, client, "record-id")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Code != http.StatusNotFound {
		t.Fatalf("unexpected status code: %d", err.Code)
	}
}

func TestDeleteActivityRecordRecordIdReturnsForbiddenForOtherUsersWithoutInfraClaim(t *testing.T) {
	e := echo.New()
	ctx := e.NewContext(httptest.NewRequest(http.MethodDelete, "/", nil), httptest.NewRecorder())
	ctx.Set("user_id", "request-user-id")

	client := &fakeTransactionClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			switch queryPath {
			case "sql/activity/select_activity_from_id.sql":
				records := dest.(*[]ActivityRecord)
				*records = []ActivityRecord{{ID: "record-id", UserID: "owner-user-id"}}
			case "sql/admin/select_user_has_claim.sql":
				result := dest.(*[]struct {
					HasClaim bool `db:"has_claim"`
				})
				*result = []struct {
					HasClaim bool `db:"has_claim"`
				}{}
			default:
				t.Fatalf("unexpected query path: %s", queryPath)
			}
			return nil
		},
	}

	_, err := DeleteActivityRecordRecordId(ctx, client, "record-id")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Code != http.StatusForbidden {
		t.Fatalf("unexpected status code: %d", err.Code)
	}
}

func TestDeleteActivityRecordRecordIdAllowsInfraUserToDeleteOtherUsersRecord(t *testing.T) {
	e := echo.New()
	ctx := e.NewContext(httptest.NewRequest(http.MethodDelete, "/", nil), httptest.NewRecorder())
	ctx.Set("user_id", "infra-user-id")

	execCalled := false
	client := &fakeTransactionClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			switch queryPath {
			case "sql/activity/select_activity_from_id.sql":
				records := dest.(*[]ActivityRecord)
				*records = []ActivityRecord{{ID: "record-id", UserID: "owner-user-id"}}
			case "sql/admin/select_user_has_claim.sql":
				result := dest.(*[]struct {
					HasClaim bool `db:"has_claim"`
				})
				*result = []struct {
					HasClaim bool `db:"has_claim"`
				}{
					{HasClaim: true},
				}
			default:
				t.Fatalf("unexpected query path: %s", queryPath)
			}
			return nil
		},
		execFunc: func(queryPath string, params interface{}, generateId bool) (sql.Result, error) {
			execCalled = true
			return fakeResult{rowsAffected: 1}, nil
		},
	}

	res, err := DeleteActivityRecordRecordId(ctx, client, "record-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !execCalled {
		t.Fatal("expected delete query to be executed")
	}
	if !res.Success {
		t.Fatal("expected success response")
	}
}

func TestDeleteActivityRecordRecordIdReturnsInternalServerErrorWhenDeleteFails(t *testing.T) {
	e := echo.New()
	ctx := e.NewContext(httptest.NewRequest(http.MethodDelete, "/", nil), httptest.NewRecorder())
	ctx.Set("user_id", "user-id")

	client := &fakeTransactionClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			records := dest.(*[]ActivityRecord)
			*records = []ActivityRecord{{ID: "record-id", UserID: "user-id"}}
			return nil
		},
		execFunc: func(queryPath string, params interface{}, generateId bool) (sql.Result, error) {
			return nil, errors.New("db error")
		},
	}

	_, err := DeleteActivityRecordRecordId(ctx, client, "record-id")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status code: %d", err.Code)
	}
}

func TestDeleteActivityRecordRecordIdReturnsInternalServerErrorWhenRowsAffectedFails(t *testing.T) {
	e := echo.New()
	ctx := e.NewContext(httptest.NewRequest(http.MethodDelete, "/", nil), httptest.NewRecorder())
	ctx.Set("user_id", "user-id")

	client := &fakeTransactionClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			records := dest.(*[]ActivityRecord)
			*records = []ActivityRecord{{ID: "record-id", UserID: "user-id"}}
			return nil
		},
		execFunc: func(queryPath string, params interface{}, generateId bool) (sql.Result, error) {
			return fakeResult{rowsErr: errors.New("rows affected error")}, nil
		},
	}

	_, err := DeleteActivityRecordRecordId(ctx, client, "record-id")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status code: %d", err.Code)
	}
}
