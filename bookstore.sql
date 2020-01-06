CREATE TABLE IF NOT EXISTS books (
  isbn    char(14) NOT NULL PRIMARY KEY,
  title   varchar(255) NOT NULL,
  author  varchar(255) NOT NULL,
  price   decimal(5,2) NOT NULL DEFAULT 0.00,
  created timestamp with time zone DEFAULT current_timestamp,
  updated timestamp with time zone
);
	
INSERT INTO books (isbn, title, author, price, created, updated) VALUES
('978-1503261969', 'Emma', 'Jayne Austen', 9.44, current_timestamp, NOW()),
('978-1505255607', 'The Time Machine', 'H. G. Wells', 5.99, current_timestamp, NOW()),
('978-1503379640', 'The Prince', 'Niccolò Machiavelli', 6.99, current_timestamp, NOW());