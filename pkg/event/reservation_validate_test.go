package event

import (
	"net/http"
	"testing"
	"time"
)

func TestValidateReservationPeriod(t *testing.T) {
	now := time.Now()
	newReservation := func(reservationStartDate, reservationFinishDate, finishDate time.Time) eventReservation {
		return eventReservation{
			ReservationStartDate:  reservationStartDate.Format(time.RFC3339),
			ReservationFinishDate: reservationFinishDate.Format(time.RFC3339),
			FinishDate:            finishDate.Format(time.RFC3339),
		}
	}

	cases := []struct {
		name        string
		reservation eventReservation
		wantCode    int
	}{
		{
			name:        "予約期間内かつイベント終了前は許容される",
			reservation: newReservation(now.Add(-1*time.Hour), now.Add(1*time.Hour), now.Add(2*time.Hour)),
			wantCode:    0,
		},
		{
			name:        "予約開始前は拒否される",
			reservation: newReservation(now.Add(1*time.Hour), now.Add(2*time.Hour), now.Add(3*time.Hour)),
			wantCode:    http.StatusBadRequest,
		},
		{
			name:        "予約受付期間終了後は拒否される",
			reservation: newReservation(now.Add(-2*time.Hour), now.Add(-1*time.Hour), now.Add(1*time.Hour)),
			wantCode:    http.StatusBadRequest,
		},
		{
			name:        "イベント終了後は拒否される",
			reservation: newReservation(now.Add(-2*time.Hour), now.Add(1*time.Hour), now.Add(-1*time.Hour)),
			wantCode:    http.StatusBadRequest,
		},
		{
			name:        "日付の解析に失敗した場合は500を返す",
			reservation: eventReservation{ReservationStartDate: "invalid"},
			wantCode:    http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateReservationPeriod(tc.reservation)
			if tc.wantCode == 0 {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("expected error but got nil")
			}
			if err.Code != tc.wantCode {
				t.Fatalf("unexpected status code: %d", err.Code)
			}
		})
	}
}
