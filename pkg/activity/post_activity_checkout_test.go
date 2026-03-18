package activity

import (
	"database/sql"
	"errors"
	"testing"
	"time"
)

type fakeTransactionClient struct {
	selectFunc func(dest interface{}, queryPath string, params interface{}) error
	execFunc   func(queryPath string, params interface{}, generateId bool) (sql.Result, error)
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
	return nil, false, nil
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

func TestExecuteCheckoutReturnsFalseWhenRowsAreNotUpdated(t *testing.T) {
	now := time.Date(2026, 3, 19, 10, 0, 0, 0, time.Local)
	client := &fakeTransactionClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			if queryPath != "sql/activity/select_current_activity.sql" {
				t.Fatalf("unexpected query path: %s", queryPath)
			}
			records := dest.(*[]ActivityRecord)
			*records = []ActivityRecord{
				{
					ID:                 "record-id",
					InitialCheckedInAt: now.Add(-time.Hour),
					CheckedInAt:        now.Add(-time.Hour),
				},
			}
			return nil
		},
		execFunc: func(queryPath string, params interface{}, generateId bool) (sql.Result, error) {
			if queryPath != "sql/activity/update_activity_checkout.sql" {
				t.Fatalf("unexpected query path: %s", queryPath)
			}
			return fakeResult{rowsAffected: 0}, nil
		},
	}

	executed, err := executeCheckout(client, "user-id", "room", &now, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if executed {
		t.Fatal("expected checkout to be treated as not executed")
	}
}

func TestExecuteCheckoutReturnsErrorWhenRowsAffectedFails(t *testing.T) {
	now := time.Date(2026, 3, 19, 10, 0, 0, 0, time.Local)
	client := &fakeTransactionClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			if queryPath != "sql/activity/select_current_activity.sql" {
				t.Fatalf("unexpected query path: %s", queryPath)
			}
			records := dest.(*[]ActivityRecord)
			*records = []ActivityRecord{
				{
					ID:                 "record-id",
					InitialCheckedInAt: now.Add(-time.Hour),
					CheckedInAt:        now.Add(-time.Hour),
				},
			}
			return nil
		},
		execFunc: func(queryPath string, params interface{}, generateId bool) (sql.Result, error) {
			return fakeResult{rowsErr: errors.New("rows affected error")}, nil
		},
	}

	_, err := executeCheckout(client, "user-id", "room", &now, nil)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestExecuteCheckoutUsesCurrentOpenRecordEvenIfLatestEditedRecordWouldDiffer(t *testing.T) {
	now := time.Date(2026, 3, 19, 10, 0, 0, 0, time.Local)
	client := &fakeTransactionClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			if queryPath != "sql/activity/select_current_activity.sql" {
				t.Fatalf("unexpected query path: %s", queryPath)
			}
			records := dest.(*[]ActivityRecord)
			*records = []ActivityRecord{
				{
					ID:                 "open-record-id",
					InitialCheckedInAt: now.Add(-2 * time.Hour),
					CheckedInAt:        now.Add(-2 * time.Hour),
				},
			}
			return nil
		},
		execFunc: func(queryPath string, params interface{}, generateId bool) (sql.Result, error) {
			updateParams, ok := params.(*struct {
				Id           string    `twowaysql:"id"`
				CheckedOutAt time.Time `twowaysql:"checkedOutAt"`
				Note         *string   `twowaysql:"note"`
			})
			if !ok {
				t.Fatal("failed to cast update params")
			}
			if updateParams.Id != "open-record-id" {
				t.Fatalf("unexpected record id: %s", updateParams.Id)
			}
			return fakeResult{rowsAffected: 1}, nil
		},
	}

	executed, err := executeCheckout(client, "user-id", "room", &now, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !executed {
		t.Fatal("expected checkout to succeed for current open record")
	}
}
