create database diary;

use diary;

create table users(
  id int unsigned not null auto_increment primary key,
  facebook_id varchar(30) not null unique,
  token text not null,
  name varchar(255)
);
