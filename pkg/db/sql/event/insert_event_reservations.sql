INSERT INTO event_reservations
      (id, event_id, name, description,
       start_date, finish_date,
       reservation_start_date, reservation_finish_date, capacity)
VALUES (UUID_TO_BIN(/*reservationId*/''), UUID_TO_BIN(/*eventId*/''), /*name*/'name', /*description*/'description', /*startDate*/'1970-01-01 00:00:00', /*finishDate*/'1970-01-01 00:00:00', /*reservationStart*/'1970-01-01 00:00:00', /*reservationFinish*/'1970-01-01 00:00:00', /*capacity*/0);
