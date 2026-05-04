DELETE FROM event_reservations 
WHERE id = UUID_TO_BIN(/*reservationId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee') 
  AND event_id = UUID_TO_BIN(/*eventId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
