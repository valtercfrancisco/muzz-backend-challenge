-- Ensure the uuid-ossp extension is enabled to generate UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Delete existing data from tables
DELETE FROM likes;
DELETE FROM users;

-- Insert mock data into users table with UUIDs
INSERT INTO users (user_id, username) VALUES (uuid_generate_v4(), 'Alice');
INSERT INTO users (user_id, username) VALUES (uuid_generate_v4(), 'Bob');
INSERT INTO users (user_id, username) VALUES (uuid_generate_v4(), 'Charlie');
INSERT INTO users (user_id, username) VALUES (uuid_generate_v4(), 'David');
INSERT INTO users (user_id, username) VALUES (uuid_generate_v4(), 'Eva');
INSERT INTO users (user_id, username) VALUES (uuid_generate_v4(), 'Frank');
INSERT INTO users (user_id, username) VALUES (uuid_generate_v4(), 'Grace');
INSERT INTO users (user_id, username) VALUES (uuid_generate_v4(), 'Hannah');
INSERT INTO users (user_id, username) VALUES (uuid_generate_v4(), 'Ivy');
INSERT INTO users (user_id, username) VALUES (uuid_generate_v4(), 'Jack');
INSERT INTO users (user_id, username) VALUES (uuid_generate_v4(), 'Katherine');
INSERT INTO users (user_id, username) VALUES (uuid_generate_v4(), 'Liam');
INSERT INTO users (user_id, username) VALUES (uuid_generate_v4(), 'Mia');
INSERT INTO users (user_id, username) VALUES (uuid_generate_v4(), 'Noah');
INSERT INTO users (user_id, username) VALUES (uuid_generate_v4(), 'Olivia');
INSERT INTO users (user_id, username) VALUES (uuid_generate_v4(), 'Paul');
INSERT INTO users (user_id, username) VALUES (uuid_generate_v4(), 'Quincy');
INSERT INTO users (user_id, username) VALUES (uuid_generate_v4(), 'Rachel');
INSERT INTO users (user_id, username) VALUES (uuid_generate_v4(), 'Sam');
INSERT INTO users (user_id, username) VALUES (uuid_generate_v4(), 'Tina');

-- Manually get the UUIDs of the inserted users to use them in the likes table
WITH user_uuids AS (
    SELECT user_id, username FROM users
)
-- Insert mock data into likes table with UUIDs
INSERT INTO likes (actor_user_id, recipient_user_id, liked_recipient)
SELECT
    u1.user_id AS actor_user_id,
    u2.user_id AS recipient_user_id,
    (random() > 0.5) AS liked_recipient
FROM
    user_uuids u1
        CROSS JOIN user_uuids u2
WHERE
    u1.user_id != u2.user_id -- Ensure actor and recipient are different
  AND NOT EXISTS (
    SELECT 1 FROM likes
    WHERE actor_user_id = u1.user_id
      AND recipient_user_id = u2.user_id
)
LIMIT 100; -- Limit the number of inserts for testing
