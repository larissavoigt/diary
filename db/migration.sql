create database diary;

use diary;

create table users(
  id int unsigned not null auto_increment primary key,
  token text not null,
  name varchar(255)
);
