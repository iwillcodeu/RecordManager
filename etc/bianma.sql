drop database bianma;

create database bianma;
use bianma;

create table bianmaindex
(
    id int auto_increment,
    category varchar(32) not null,
    origin  int not null,
    origintxt varchar(64) not null,
    details varchar(64) not null,
    code varchar(8) not null,
    primary key(id)
);