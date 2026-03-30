SELECT COUNT(*) AS count FROM event_reservation_users 
WHERE reservation_id = UUID_TO_BIN(/*reservationId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
