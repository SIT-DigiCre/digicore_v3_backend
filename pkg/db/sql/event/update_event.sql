UPDATE events 
SET name = /*name*/'name', 
    description = /*description*/'description', 
    calendar_view = /*calendarView*/0 
WHERE id = UUID_TO_BIN(/*eventId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
