package activity

import (
	"errors"
	"testing"
	"time"
)

type fakeClient struct {
	selectFunc func(dest interface{}, queryPath string, params interface{}) error
}

func (f *fakeClient) Select(dest interface{}, queryPath string, params interface{}) error {
	if f.selectFunc != nil {
		return f.selectFunc(dest, queryPath, params)
	}
	return nil
}

func TestGetActivityPlacePlaceHistoryReturnsBadRequestWhenStartAtAfterEndAt(t *testing.T) {
	startAt := time.Date(2026, 3, 20, 0, 0, 0, 0, time.Local)
	endAt := time.Date(2026, 3, 19, 23, 59, 59, 0, time.Local)

	_, err := GetActivityPlacePlaceHistory(nil, &fakeClient{}, "club_room", startAt, endAt)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Code != 400 {
		t.Fatalf("unexpected status code: %d", err.Code)
	}
}

func TestGetActivityPlacePlaceHistoryReturnsUsersForSingleDayRange(t *testing.T) {
	startAt := time.Date(2026, 3, 19, 0, 0, 0, 0, time.Local)
	endAt := time.Date(2026, 3, 19, 23, 59, 59, 0, time.Local)

	client := &fakeClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			if queryPath != "sql/activity/select_place_history.sql" {
				t.Fatalf("unexpected query path: %s", queryPath)
			}

			queryParams, ok := params.(*struct {
				SameDay   bool      `twowaysql:"sameDay"`
				Place     string    `twowaysql:"place"`
				StartDate time.Time `twowaysql:"startDate"`
				EndDate   time.Time `twowaysql:"endDate"`
			})
			if !ok {
				t.Fatal("failed to cast query params")
			}
			if !queryParams.SameDay {
				t.Fatal("expected sameDay to be true")
			}
			if queryParams.Place != "club_room" {
				t.Fatalf("unexpected place: %s", queryParams.Place)
			}

			users := dest.(*[]placeHistoryUser)
			*users = []placeHistoryUser{
				{
					UserId:            "user-id",
					Username:          "テストユーザー",
					ShortIntroduction: "紹介",
					IconUrl:           "https://example.com/icon.png",
					CheckInCount:      3,
				},
			}
			return nil
		},
	}

	res, err := GetActivityPlacePlaceHistory(nil, client, "club_room", startAt, endAt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Users) != 1 {
		t.Fatalf("unexpected users length: %d", len(res.Users))
	}
	if res.Users[0].CheckInCount != 3 {
		t.Fatalf("unexpected checkInCount: %d", res.Users[0].CheckInCount)
	}
}

func TestGetActivityPlacePlaceHistoryUsesDistinctDaysForMultiDayRange(t *testing.T) {
	startAt := time.Date(2026, 3, 17, 0, 0, 0, 0, time.Local)
	endAt := time.Date(2026, 3, 23, 23, 59, 59, 0, time.Local)

	client := &fakeClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			queryParams, ok := params.(*struct {
				SameDay   bool      `twowaysql:"sameDay"`
				Place     string    `twowaysql:"place"`
				StartDate time.Time `twowaysql:"startDate"`
				EndDate   time.Time `twowaysql:"endDate"`
			})
			if !ok {
				t.Fatal("failed to cast query params")
			}
			if queryParams.SameDay {
				t.Fatal("expected sameDay to be false")
			}

			users := dest.(*[]placeHistoryUser)
			*users = []placeHistoryUser{
				{
					UserId:            "user-id",
					Username:          "テストユーザー",
					ShortIntroduction: "紹介",
					IconUrl:           "https://example.com/icon.png",
					CheckInCount:      2,
				},
			}
			return nil
		},
	}

	res, err := GetActivityPlacePlaceHistory(nil, client, "club_room", startAt, endAt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := res.Users[0].CheckInCount, 2; got != want {
		t.Fatalf("unexpected checkInCount: got %d want %d", got, want)
	}
}

func TestGetActivityPlacePlaceHistoryReturnsEmptyUsersWhenNoRows(t *testing.T) {
	startAt := time.Date(2026, 3, 19, 0, 0, 0, 0, time.Local)
	endAt := time.Date(2026, 3, 19, 23, 59, 59, 0, time.Local)

	res, err := GetActivityPlacePlaceHistory(nil, &fakeClient{}, "club_room", startAt, endAt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Users == nil {
		t.Fatal("expected users to be initialized")
	}
	if len(res.Users) != 0 {
		t.Fatalf("unexpected users length: %d", len(res.Users))
	}
}

func TestSelectPlaceHistoryReturnsInternalServerErrorOnDBFailure(t *testing.T) {
	startAt := time.Date(2026, 3, 19, 0, 0, 0, 0, time.Local)
	endAt := time.Date(2026, 3, 19, 23, 59, 59, 0, time.Local)

	client := &fakeClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			return errors.New("db error")
		},
	}

	_, err := selectPlaceHistory(client, "club_room", startAt, endAt)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Code != 500 {
		t.Fatalf("unexpected status code: %d", err.Code)
	}
}
