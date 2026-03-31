UPDATE event_reservations 
SET name = /*name*/'name', 
    description = /*description*/'description', 
    start_date = /*startDate*/'2024-01-01 00:00:00',
    finish_date = /*finishDate*/'2024-01-01 00:00:00',
    reservation_start_date = /*reservationStartDate*/'2024-01-01 00:00:00',
    reservation_finish_date = /*reservationFinishDate*/'2024-01-01 00:00:00',
    capacity = /*capacity*/1
WHERE id = UUID_TO_BIN(/*reservationId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee') 
  AND event_id = UUID_TO_BIN(/*eventId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')
  AND /*capacity*/1 >= (
    SELECT COALESCE(COUNT(event_reservation_users.id), 0) 
    FROM event_reservation_users 
    WHERE reservation_id = UUID_TO_BIN(/*reservationId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')
  );
