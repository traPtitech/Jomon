-- Modify "requests" table
ALTER TABLE `requests` DROP FOREIGN KEY `requests_groups_request`;
-- Modify "transactions" table
ALTER TABLE `transactions` DROP FOREIGN KEY `transactions_group_budgets_transaction`;
