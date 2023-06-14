SELECT BIN_TO_UUID(id) as id, body, remind_date FROM mattermost_remind_post WHERE posted = false AND user_name = /*userName*/'user_name' AND channel_name =  /*channelName*/'channel_name'
