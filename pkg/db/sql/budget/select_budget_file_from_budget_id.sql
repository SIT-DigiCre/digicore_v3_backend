SELECT BIN_TO_UUID(file_id) AS file_id FROM budget_files WHERE budget_id = UUID_TO_BIN(/*budgetId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
