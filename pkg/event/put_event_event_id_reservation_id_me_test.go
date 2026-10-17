package event

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
)

func TestPutEventEventIdReservationIdMeCreatesReservationInPeriod(t *testing.T) {
	ctx := newReservationContext(http.MethodPut)
	reservation := validEventReservation()
	duplicateUpdateCalled := false
	client := &fakeTransactionClient{
		selectFunc: reservationSelectFunc(t, []eventReservation{reservation}),
		duplicateUpdateFunc: func(insertQueryPath string, updateQueryPath string, params interface{}) (sql.Result, bool, error) {
			if insertQueryPath != "sql/event/insert_reservation_user.sql" {
				t.Fatalf("unexpected insert query path: %s", insertQueryPath)
			}
			duplicateUpdateCalled = true
			return fakeResult{rowsAffected: 1}, false, nil
		},
	}

	_, err := PutEventEventIdReservationIdMe(ctx, client, "event-id", "reservation-id", api.ReqPutEventEventIdReservationIdMe{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !duplicateUpdateCalled {
		t.Fatal("expected insert/update query to be executed")
	}
}

func TestPutEventEventIdReservationIdMeRejectsFullReservation(t *testing.T) {
	ctx := newReservationContext(http.MethodPut)
	reservation := validEventReservation()
	reservation.FreeCapacity = 0
	duplicateUpdateCalled := false
	client := &fakeTransactionClient{
		selectFunc: reservationSelectFunc(t, []eventReservation{reservation}),
		duplicateUpdateFunc: func(insertQueryPath string, updateQueryPath string, params interface{}) (sql.Result, bool, error) {
			duplicateUpdateCalled = true
			return fakeResult{rowsAffected: 1}, false, nil
		},
	}

	_, err := PutEventEventIdReservationIdMe(ctx, client, "event-id", "reservation-id", api.ReqPutEventEventIdReservationIdMe{})
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status code: %d", err.Code)
	}
	if duplicateUpdateCalled {
		t.Fatal("insert/update query should not be executed")
	}
}

func TestPutEventEventIdReservationIdMeAllowsUpdateWhenReservatedEvenIfFull(t *testing.T) {
	ctx := newReservationContext(http.MethodPut)
	reservation := validEventReservation()
	reservation.Reservated = true
	reservation.FreeCapacity = 0
	duplicateUpdateCalled := false
	client := &fakeTransactionClient{
		selectFunc: reservationSelectFunc(t, []eventReservation{reservation}),
		duplicateUpdateFunc: func(insertQueryPath string, updateQueryPath string, params interface{}) (sql.Result, bool, error) {
			duplicateUpdateCalled = true
			return fakeResult{rowsAffected: 1}, false, nil
		},
	}

	_, err := PutEventEventIdReservationIdMe(ctx, client, "event-id", "reservation-id", api.ReqPutEventEventIdReservationIdMe{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !duplicateUpdateCalled {
		t.Fatal("expected insert/update query to be executed")
	}
}

func TestPutEventEventIdReservationIdMeRejectsAfterReservationFinish(t *testing.T) {
	ctx := newReservationContext(http.MethodPut)
	reservation := validEventReservation()
	reservation.ReservationFinishDate = time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	duplicateUpdateCalled := false
	client := &fakeTransactionClient{
		selectFunc: reservationSelectFunc(t, []eventReservation{reservation}),
		duplicateUpdateFunc: func(insertQueryPath string, updateQueryPath string, params interface{}) (sql.Result, bool, error) {
			duplicateUpdateCalled = true
			return fakeResult{rowsAffected: 1}, false, nil
		},
	}

	_, err := PutEventEventIdReservationIdMe(ctx, client, "event-id", "reservation-id", api.ReqPutEventEventIdReservationIdMe{})
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status code: %d", err.Code)
	}
	if duplicateUpdateCalled {
		t.Fatal("insert/update query should not be executed")
	}
}
