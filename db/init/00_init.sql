CREATE DATABASE isulogger;

\c isulogger

CREATE TABLE entry (
    id serial NOT NULL,
    isucon_id integer NOT NULL,
    timestamp timestamp without time zone NOT NULL,
    score int not null default 0,
    message text NOT NULL default '',
    access_log_path text default NULL,
    slow_log_path text default NULL,
    image_path text default NULL,
    PRIMARY KEY (id)
);