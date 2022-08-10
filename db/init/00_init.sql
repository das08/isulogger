CREATE DATABASE isulogger;

\c isulogger

CREATE TABLE contest (
    contest_id int NOT NULL,
    contest_name varchar(255) NOT NULL default '',
    PRIMARY KEY (contest_id)
);

CREATE TABLE entry (
    id serial NOT NULL,
    contest_id integer NOT NULL,
    timestamp timestamp without time zone NOT NULL,
    score int not null default 0,
    message text NOT NULL default '',
    access_log_path text NOT NULL default '',
    slow_log_path text NOT NULL default '',
    image_path text NOT NULL default '',
    PRIMARY KEY (id),
    FOREIGN KEY (contest_id) REFERENCES contest (contest_id)
);

CREATE TABLE settings (
    selected_contest int NOT NULL default 0,
    FOREIGN KEY (selected_contest) REFERENCES contest (contest_id)
);