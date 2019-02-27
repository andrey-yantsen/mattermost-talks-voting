CREATE TABLE registrations (
  channel_id varchar(32) PRIMARY KEY not null,
  timezone varchar(32) not null,
  owner_id varchar(32) not null,
  talks_to_show_on_vote INT not null,
  hide_talks_watched_by_others TINYINT not null,
  schedule_dow TINYINT not null,
  schedule_time_in_minutes SMALLINT not null,
  minimal_people_required SMALLINT not null,
  first_reminder_hours SMALLINT not null,
  second_reminder_minutes smallint not null,
  add_random_talk_on_second_reminder tinyint not null,
  final_call_minutes tinyint not null,
  show_vote_result tinyint not null,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE supported_languages (
  channel_id varchar(32) not null,
  language_code varchar(2) not null
);

CREATE UNIQUE INDEX supported_languages_uidx ON supported_languages (channel_id, language_code);

CREATE TABLE topics (
  id integer primary key autoincrement not null,
  alias varchar(32) not null
);

CREATE TABLE subscriptions (
  channel_id varchar(32) not null,
  topic_id int not null
);

CREATE UNIQUE INDEX subscriptions_uidx ON subscriptions (channel_id, topic_id);

CREATE TABLE talks (
  id integer primary key autoincrement not null,
  topic_id int not null,
  added_by varchar(32) not null,
  language varchar(2) not null,
  url varchar(500) not null,
  title varchar(500) not null,
  author varchar(500) not null,
  annotation text not null,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX talks_url_idx ON talks (url);
CREATE INDEX talks_topic_id_idx ON talks (topic_id);

CREATE TABLE bulletin_for_a_day (
  date date not null,
  channel_id varchar(32) not null,
  talk_id int not null,
  is_extra tinyint not null
);

CREATE UNIQUE INDEX bulletin_uidx ON bulletin_for_a_day (date, channel_id, talk_id);

CREATE TABLE votes (
  date date not null,
  channel_id varchar(32) not null,
  talk_id int not null,
  user_id varchar(32) not null,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX votes_uidx ON votes (date, channel_id, talk_id, user_id);

CREATE TABLE overwrites (
  date date not null,
  channel_id varchar(32) not null,
  speaker_id varchar(32),
  title varchar(500),
  annotation text
);

CREATE UNIQUE INDEX overwrites_uidx ON overwrites (date, channel_id);
