CREATE TABLE web_user (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255),
    password VARCHAR(255)
);

CREATE TABLE channel (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  title VARCHAR(255),
  photourl TEXT NOT NULL
);

CREATE TABLE tg_user (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) UNIQUE NOT NULL,
  fullname VARCHAR(255) UNIQUE NOT NULL,
  photourl TEXT NOT NULL 
);

CREATE TABLE message (
  id SERIAL PRIMARY KEY,
  channel_id INT NOT NULL,
  user_id INT NOT NULL,
  title TEXT,
  CONSTRAINT fk_channel FOREIGN KEY(channel_id) REFERENCES channel(id),
  CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES tg_user(id)
);

CREATE TABLE replie (
  id SERIAL PRIMARY KEY,
  message_id INT NOT NULL,
  user_id INT NOT NULL,
  title TEXT,
  CONSTRAINT fk_message FOREIGN KEY(message_id) REFERENCES message(id),
  CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES tg_user(id)
);
