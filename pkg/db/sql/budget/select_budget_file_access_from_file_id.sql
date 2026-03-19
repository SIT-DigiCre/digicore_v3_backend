SELECT
  BIN_TO_UUID(budgets.id) AS budget_id,
  BIN_TO_UUID(budgets.proposer_user_id) AS proposer_user_id
FROM budget_files
INNER JOIN budgets ON budget_files.budget_id = budgets.id
WHERE budget_files.file_id = UUID_TO_BIN(/*fileId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')
LIMIT 1;
