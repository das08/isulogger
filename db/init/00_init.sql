CREATE TABLE contest (
    contest_id INTEGER PRIMARY KEY AUTOINCREMENT,
    contest_name varchar(255) NOT NULL default ''
);

CREATE TABLE entry (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    contest_id integer NOT NULL,
    timestamp TEXT NOT NULL,
    branch_name varchar(255) NOT NULL default '',
    score int not null default 0,
    message text NOT NULL default '',
    access_log_path text NOT NULL default '',
    slow_log_path text NOT NULL default '',
    image_path text NOT NULL default '',
    FOREIGN KEY (contest_id) REFERENCES contest (contest_id)
);