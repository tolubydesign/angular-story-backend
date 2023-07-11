-- Remove preexisting tables;
DROP TABLE IF EXISTS city, contacts, country, order_status, product, sale, status_name, store, story, users;
-- TODO: create user tables
-- TODO: create dummy users
-- TODO: JWT tokens
-- TODO: JWT tokens - Configured to expire 
-- TODO: User login system

-- Creation of product table
CREATE TABLE IF NOT EXISTS product (
  product_id INT NOT NULL,
  name varchar(250) NOT NULL,
  PRIMARY KEY (product_id)
);

-- Creation of country table
CREATE TABLE IF NOT EXISTS country (
  country_id INT NOT NULL,
  country_name varchar(450) NOT NULL,
  PRIMARY KEY (country_id)
);

-- Creation of city table
CREATE TABLE IF NOT EXISTS city (
  city_id INT NOT NULL,
  city_name varchar(450) NOT NULL,
  country_id INT NOT NULL,
  PRIMARY KEY (city_id),
  CONSTRAINT fk_country
      FOREIGN KEY(country_id) 
	  REFERENCES country(country_id)
);

-- Creation of store table
CREATE TABLE IF NOT EXISTS store (
  store_id INT NOT NULL,
  name varchar(250) NOT NULL,
  city_id INT NOT NULL,

  PRIMARY KEY (store_id),
  CONSTRAINT fk_city
      FOREIGN KEY(city_id) 
	  REFERENCES city(city_id)
);

-- Creation of user table
CREATE TABLE IF NOT EXISTS users (
  user_id INT NOT NULL,
  name varchar(250) NOT NULL,
  PRIMARY KEY (user_id)
);

-- Creation of status_name table
CREATE TABLE IF NOT EXISTS status_name (
  status_name_id INT NOT NULL,
  status_name varchar(450) NOT NULL,
  PRIMARY KEY (status_name_id)
);

-- Creation of sale table
CREATE TABLE IF NOT EXISTS sale (
  sale_id varchar(200) NOT NULL,
  amount DECIMAL(20,3) NOT NULL,
  date_sale TIMESTAMP,
  product_id INT NOT NULL,
  user_id INT NOT NULL,
  store_id INT NOT NULL,

  PRIMARY KEY (sale_id),
  CONSTRAINT fk_product
      FOREIGN KEY(product_id) 
	  REFERENCES product(product_id),
  CONSTRAINT fk_user
      FOREIGN KEY(user_id) 
	  REFERENCES users(user_id),
  CONSTRAINT fk_store
      FOREIGN KEY(store_id) 
	  REFERENCES store(store_id)	  
);

-- Creation of order_status table
CREATE TABLE IF NOT EXISTS order_status (
  order_status_id varchar(200) NOT NULL,
  update_at TIMESTAMP,
  sale_id varchar(200) NOT NULL,
  status_name_id INT NOT NULL,

  PRIMARY KEY (order_status_id),
  CONSTRAINT fk_sale
      FOREIGN KEY(sale_id) 
	  REFERENCES sale(sale_id),
  CONSTRAINT fk_status_name
      FOREIGN KEY(status_name_id) 
	  REFERENCES status_name(status_name_id)  
);

-- Self created content
-- Create database
-- CREATE DATABASE testdb;
-- Resource: https://www.tutorialspoint.com/postgresql/postgresql_create_database.htm
-- $ createdb -h localhost -p 5432 -U postgres testdb
-- $ password ****** 

-- list database
-- $ \l
-- $ \d

-- connect to database 
-- $ \c dbname

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS contacts (
  contact_id uuid DEFAULT uuid_generate_v4 (),
  first_name VARCHAR NOT NULL,
  last_name VARCHAR NOT NULL,
  email VARCHAR NOT NULL,
  phone VARCHAR,
  PRIMARY KEY (contact_id)
);

-- Create table
CREATE TABLE IF NOT EXISTS Story (
  story_id uuid DEFAULT uuid_generate_v4(),
  title VARCHAR NOT NULL,
  description VARCHAR NOT NULL,
  content json,
  
  PRIMARY KEY (story_id)
);
