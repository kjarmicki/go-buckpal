CREATE TABLE activities (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  timestamp DATETIME NOT NULL,
  owner_account_id INT NOT NULL,
  source_account_id INT NOT NULL, 
  target_account_id INT NOT NULL,
  amount BIGINT NOT NULL
);
