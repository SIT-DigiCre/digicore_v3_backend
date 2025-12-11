SET @event_uuid := ?;
INSERT INTO events (id, name, calendar_view, description)
VALUES (UUID_TO_BIN(@event_uuid), ?, 0, ?);

INSERT INTO event_reservations
      (id, event_id, name, description,
       start_date, finish_date,
       reservation_start_date, reservation_finish_date, capacity)
VALUES (?, UUID_TO_BIN(@event_uuid), ?, ?, ?, ?, ?, ?, ?);