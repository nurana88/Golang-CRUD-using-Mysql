
CREATE TABLE warehouse (
  id    INTEGER AUTO_INCREMENT,
  items   VARCHAR(255) NOT NULL,
  unit_price  DECIMAL(7,2) NOT NULL,
  item_category VARCHAR(255) NOT NULL,
  quantity  INTEGER  NOT NULL,
  sold_quantity INTEGER NOT NULL,
  PRIMARY KEY(id));


INSERT INTO warehouse (items, unit_price, item_category, quantity, sold_quantity) VALUES
('Samsung Galaxy Watch', '279.99', 'smartwatch', '45','18'),
('LIEBHERR T 1404-21', '255.99', 'refrigerator', '95','36'),
('DELONGHI ECAM', '349', 'coffee machine', '36','9'),
('BOSCH WAU 28', '629.99', 'washing machine', '50','22'),
('APPLE Watch SE', '278.36', 'smartwatch', '63','39'),
('BEKO RSNE415T34XPN', '399.99', 'refrigerator', '52','26'),
('PHILIPS Sonicare', '245.99', 'electronic tootbrush', '77','45'),
('FITBIT FB507BKBK Versa', '129', 'smartwatch', '27','12'),
('JURA E8', '849', 'coffee machine', '17','2');


INSERT INTO warehouse (items, unit_price, item_category, quantity, sold_quantity) VALUES ('BOSCH WAU 28', '629.99', 'washing machine', '50','22');


UPDATE warehouse SET items = 'PHILIPS Sonicare', unit_price='245.99', item_category='electronic tootbrush', quantity=77, sold_quantity=45 WHERE id=4;


DELETE FROM  warehouse WHERE id=5;
