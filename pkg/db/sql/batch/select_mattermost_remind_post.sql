SELECT BIN_TO_UUID(id) as id, user_name, channel_name, body FROM mattermost_remind_post WHERE posted = false AND remind_date <=  /*remind_date*/'2023-01-01 00:50:00'
