SELECT BIN_TO_UUID(events.id) as event_id FROM events WHERE events.id = UUID_TO_BIN(/*eventId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')
