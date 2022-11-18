SELECT BIN_TO_UUID(events.id) as event_id, events.name, events.description, calendar_view, IF(NOW() <= max(reservation_finish_date),true,false) as reservable, IF(1 <= count(user_id), true, false) as reservated FROM events LEFT JOIN event_reservations ON events.id = event_reservations.event_id LEFT JOIN event_reservation_users ON event_reservations.id = event_reservation_users.reservation_id AND user_id = UUID_TO_BIN(/*userId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee') GROUP BY events.id LIMIT 50 /* IF offset*/ OFFSET /*offset*/0 /* END */
