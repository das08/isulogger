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
    FOREIGN KEY (contest_id) REFERENCES contest (contest_id)
);

CREATE TABLE attached_file(
    id INTEGER PRIMARY KEY,
    entry_id INTEGER NOT NULL,
    file_type VARCHAR(32) NOT NULL,
    source VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    FOREIGN KEY (entry_id) REFERENCES entry (id) ON DELETE CASCADE ON UPDATE CASCADE,
    UNIQUE (entry_id, file_type, source)
);