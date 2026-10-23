package event

import (
	"testing"
	"time"
)

func TestGetReservationFromReservationIdReservable(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name   string
		mutate func(*eventReservation)
		want   bool
	}{
		{
			name:   "予約期間内かつ予約枠終了前は予約可能",
			mutate: func(r *eventReservation) {},
			want:   true,
		},
		{
			name: "予約枠の終了後は予約不可",
			mutate: func(r *eventReservation) {
				r.FinishDate = now.Add(-1 * time.Hour).Format(time.RFC3339)
			},
			want: false,
		},
		{
			name: "予約受付期間の終了後は予約不可",
			mutate: func(r *eventReservation) {
				r.ReservationFinishDate = now.Add(-1 * time.Hour).Format(time.RFC3339)
			},
			want: false,
		},
		{
			name: "満員の場合は予約不可",
			mutate: func(r *eventReservation) {
				r.FreeCapacity = 0
			},
			want: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reservation := validEventReservation()
			tc.mutate(&reservation)
			client := &fakeTransactionClient{
				selectFunc: reservationSelectFunc(t, []eventReservation{reservation}),
			}

			res, err := getReservationFromReservationId(client, "event-id", "reservation-id", "user-id")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res.Reservable != tc.want {
				t.Fatalf("unexpected reservable: %v", res.Reservable)
			}
		})
	}
}
