package event

import (
	"database/sql"
	"net/http"
	"testing"
	"time"
)

func TestDeleteEventEventIdReservationIdMeDeletesReservationInPeriod(t *testing.T) {
	ctx := newReservationContext(http.MethodDelete)
	reservation := validEventReservation()
	reservation.Reservated = true
	execCalled := false
	client := &fakeTransactionClient{
		selectFunc: reservationSelectFunc(t, []eventReservation{reservation}),
		execFunc: func(queryPath string, params interface{}, generateId bool) (sql.Result, error) {
			if queryPath != "sql/event/delete_reservation_user.sql" {
				t.Fatalf("unexpected query path: %s", queryPath)
			}
			execCalled = true
			return fakeResult{rowsAffected: 1}, nil
		},
	}

	_, err := DeleteEventEventIdReservationIdMe(ctx, client, "event-id", "reservation-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !execCalled {
		t.Fatal("expected delete query to be executed")
	}
}

func TestDeleteEventEventIdReservationIdMeRejectsAfterEventFinish(t *testing.T) {
	ctx := newReservationContext(http.MethodDelete)
	reservation := validEventReservation()
	reservation.FinishDate = time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	execCalled := false
	client := &fakeTransactionClient{
		selectFunc: reservationSelectFunc(t, []eventReservation{reservation}),
		execFunc: func(queryPath string, params interface{}, generateId bool) (sql.Result, error) {
			execCalled = true
			return fakeResult{rowsAffected: 1}, nil
		},
	}

	_, err := DeleteEventEventIdReservationIdMe(ctx, client, "event-id", "reservation-id")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status code: %d", err.Code)
	}
	if execCalled {
		t.Fatal("delete query should not be executed")
	}
}

func TestDeleteEventEventIdReservationIdMeRejectsAfterReservationFinish(t *testing.T) {
	ctx := newReservationContext(http.MethodDelete)
	reservation := validEventReservation()
	reservation.ReservationFinishDate = time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	client := &fakeTransactionClient{
		selectFunc: reservationSelectFunc(t, []eventReservation{reservation}),
	}

	_, err := DeleteEventEventIdReservationIdMe(ctx, client, "event-id", "reservation-id")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status code: %d", err.Code)
	}
}

func TestDeleteEventEventIdReservationIdMeReturnsNotFound(t *testing.T) {
	ctx := newReservationContext(http.MethodDelete)
	execCalled := false
	client := &fakeTransactionClient{
		selectFunc: reservationSelectFunc(t, []eventReservation{}),
		execFunc: func(queryPath string, params interface{}, generateId bool) (sql.Result, error) {
			execCalled = true
			return fakeResult{rowsAffected: 1}, nil
		},
	}

	_, err := DeleteEventEventIdReservationIdMe(ctx, client, "event-id", "reservation-id")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Code != http.StatusNotFound {
		t.Fatalf("unexpected status code: %d", err.Code)
	}
	if execCalled {
		t.Fatal("delete query should not be executed")
	}
}
