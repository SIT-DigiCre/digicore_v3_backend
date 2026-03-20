package activity

import (
	"errors"
	"testing"
	"time"
)

func TestGetActivityRecordsReturnsDefaultPagination(t *testing.T) {
	now := time.Date(2026, 3, 19, 10, 0, 0, 0, time.Local)
	client := &fakeClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			switch queryPath {
			case "sql/activity/select_activity_records.sql":
				queryParams, ok := params.(*struct {
					Offset int `twowaysql:"offset"`
					Limit  int `twowaysql:"limit"`
				})
				if !ok {
					t.Fatal("failed to cast query params")
				}
				if queryParams.Offset != 0 {
					t.Fatalf("unexpected offset: %d", queryParams.Offset)
				}
				if queryParams.Limit != 50 {
					t.Fatalf("unexpected limit: %d", queryParams.Limit)
				}
				records := dest.(*[]activityRecordListItem)
				*records = []activityRecordListItem{{
					RecordId:           "record-id",
					UserId:             "user-id",
					Username:           "テストユーザー",
					Place:              "club_room",
					CheckedInAt:        now,
					InitialCheckedInAt: now,
				}}
			case "sql/activity/select_activity_records_count.sql":
				counts := dest.(*[]activityRecordsCount)
				*counts = []activityRecordsCount{{Total: 1}}
			default:
				t.Fatalf("unexpected query path: %s", queryPath)
			}
			return nil
		},
	}

	res, err := GetActivityRecords(nil, client, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Offset != 0 || res.Limit != 50 || res.Total != 1 {
		t.Fatalf("unexpected pagination: %+v", res)
	}
	if len(res.Records) != 1 {
		t.Fatalf("unexpected records length: %d", len(res.Records))
	}
	if res.Records[0].Username != "テストユーザー" {
		t.Fatalf("unexpected username: %s", res.Records[0].Username)
	}
}

func TestGetActivityRecordsReturnsEmptySliceWhenNoRows(t *testing.T) {
	client := &fakeClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			if queryPath == "sql/activity/select_activity_records_count.sql" {
				counts := dest.(*[]activityRecordsCount)
				*counts = []activityRecordsCount{{Total: 0}}
			}
			return nil
		},
	}

	res, err := GetActivityRecords(nil, client, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Records == nil {
		t.Fatal("expected records to be initialized")
	}
	if len(res.Records) != 0 {
		t.Fatalf("unexpected records length: %d", len(res.Records))
	}
}

func TestSelectActivityRecordsReturnsInternalServerErrorOnDBFailure(t *testing.T) {
	client := &fakeClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			return errors.New("db error")
		},
	}

	_, err := selectActivityRecords(client, 0, 50)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Code != 500 {
		t.Fatalf("unexpected status code: %d", err.Code)
	}
}

func TestSelectActivityRecordsCountReturnsInternalServerErrorOnDBFailure(t *testing.T) {
	client := &fakeClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			return errors.New("db error")
		},
	}

	_, err := selectActivityRecordsCount(client)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Code != 500 {
		t.Fatalf("unexpected status code: %d", err.Code)
	}
}

func TestSelectActivityRecordsCountUsesSameJoinAsList(t *testing.T) {
	client := &fakeClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			if queryPath != "sql/activity/select_activity_records_count.sql" {
				t.Fatalf("unexpected query path: %s", queryPath)
			}
			counts := dest.(*[]activityRecordsCount)
			*counts = []activityRecordsCount{{Total: 3}}
			return nil
		},
	}

	total, err := selectActivityRecordsCount(client)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 3 {
		t.Fatalf("unexpected total: %d", total)
	}
}
