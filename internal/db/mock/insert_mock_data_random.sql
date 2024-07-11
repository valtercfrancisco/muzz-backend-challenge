CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DELETE FROM likes;
DELETE FROM decisions;
DELETE FROM users;

-- Insert mock data into users table with UUIDs
INSERT INTO users (user_id, username)
VALUES
    (uuid_generate_v4(), 'Alice'),
    (uuid_generate_v4(), 'Bob'),
    (uuid_generate_v4(), 'Charlie'),
    (uuid_generate_v4(), 'David'),
    (uuid_generate_v4(), 'Eva'),
    (uuid_generate_v4(), 'Frank'),
    (uuid_generate_v4(), 'Grace'),
    (uuid_generate_v4(), 'Hannah'),
    (uuid_generate_v4(), 'Ivy'),
    (uuid_generate_v4(), 'Jack'),
    (uuid_generate_v4(), 'Katherine'),
    (uuid_generate_v4(), 'Liam'),
    (uuid_generate_v4(), 'Mia'),
    (uuid_generate_v4(), 'Noah'),
    (uuid_generate_v4(), 'Olivia'),
    (uuid_generate_v4(), 'Paul'),
    (uuid_generate_v4(), 'Quincy'),
    (uuid_generate_v4(), 'Rachel'),
    (uuid_generate_v4(), 'Sam'),
    (uuid_generate_v4(), 'Tina');

-- Insert mock data into likes table with UUIDs, randomly skipping some users and allowing multiple likes per user
WITH user_uuids AS (
    SELECT user_id, username FROM users
)
INSERT INTO likes (actor_user_id, recipient_user_id)
SELECT
    u1.user_id AS actor_user_id,
    u2.user_id AS recipient_user_id
FROM
    user_uuids u1
        CROSS JOIN user_uuids u2
WHERE
    u1.user_id != u2.user_id -- Ensure actor and recipient are different
  AND random() > 0.2 -- Randomly skip some users
  AND (
          SELECT COUNT(*) FROM likes
          WHERE actor_user_id = u1.user_id
      ) < (random() * 10) + 1 -- Randomly determine the number of likes per user, between 1 and 10
  AND NOT EXISTS (
    SELECT 1 FROM likes
    WHERE actor_user_id = u1.user_id
      AND recipient_user_id = u2.user_id
)
LIMIT 100; -- Limit the number of inserts for testing

-- Insert mock data into decisions table based on likes table
INSERT INTO decisions (actor_user_id, recipient_user_id, liked_recipient)
SELECT
    actor_user_id,
    recipient_user_id,
    TRUE AS liked
FROM
    likes;

-- Insert mock data into decisions table for users not in the likes table
WITH user_uuids AS (
    SELECT user_id, username FROM users
)
INSERT INTO decisions (actor_user_id, recipient_user_id, liked_recipient)
SELECT
    u1.user_id AS actor_user_id,
    u2.user_id AS recipient_user_id,
    FALSE AS liked
FROM
    user_uuids u1
        CROSS JOIN user_uuids u2
WHERE
    u1.user_id != u2.user_id
  AND NOT EXISTS (
    SELECT 1 FROM decisions
    WHERE actor_user_id = u1.user_id
      AND recipient_user_id = u2.user_id
)
  AND random() > 0.7 -- Randomly include some users with liked = FALSE
LIMIT 100; -- Limit the number of inserts for testing
