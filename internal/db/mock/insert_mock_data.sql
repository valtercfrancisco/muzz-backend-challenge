CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DELETE FROM likes;
DELETE FROM decisions;
DELETE FROM users;

-- Insert mock data into users table with fixed UUIDs
INSERT INTO users (user_id, username)
VALUES
    ('00000000-0000-0000-0000-000000000001', 'Alice'),
    ('00000000-0000-0000-0000-000000000002', 'Bob'),
    ('00000000-0000-0000-0000-000000000003', 'Charlie'),
    ('00000000-0000-0000-0000-000000000004', 'David'),
    ('00000000-0000-0000-0000-000000000005', 'Eva'),
    ('00000000-0000-0000-0000-000000000006', 'Frank'),
    ('00000000-0000-0000-0000-000000000007', 'Grace'),
    ('00000000-0000-0000-0000-000000000008', 'Hannah'),
    ('00000000-0000-0000-0000-000000000009', 'Ivy'),
    ('00000000-0000-0000-0000-000000000010', 'Jack'),
    ('00000000-0000-0000-0000-000000000011', 'Katherine'),
    ('00000000-0000-0000-0000-000000000012', 'Liam'),
    ('00000000-0000-0000-0000-000000000013', 'Mia'),
    ('00000000-0000-0000-0000-000000000014', 'Noah'),
    ('00000000-0000-0000-0000-000000000015', 'Olivia'),
    ('00000000-0000-0000-0000-000000000016', 'Paul'),
    ('00000000-0000-0000-0000-000000000017', 'Quincy'),
    ('00000000-0000-0000-0000-000000000018', 'Rachel'),
    ('00000000-0000-0000-0000-000000000019', 'Sam'),
    ('00000000-0000-0000-0000-000000000020', 'Tina');

-- Insert mock data into likes table with fixed relationships, including mutual likes
INSERT INTO likes (actor_user_id, recipient_user_id)
VALUES
    ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002'), -- Alice likes Bob
    ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001'), -- Bob likes Alice (mutual)
    ('00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000004'), -- Charlie likes David
    ('00000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000003'), -- David likes Charlie (mutual)
    ('00000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000006'), -- Eva likes Frank
    ('00000000-0000-0000-0000-000000000007', '00000000-0000-0000-0000-000000000008'), -- Grace likes Hannah
    ('00000000-0000-0000-0000-000000000008', '00000000-0000-0000-0000-000000000007'), -- Hannah likes Grace (mutual)
    ('00000000-0000-0000-0000-000000000009', '00000000-0000-0000-0000-000000000010'), -- Ivy likes Jack
    ('00000000-0000-0000-0000-000000000011', '00000000-0000-0000-0000-000000000012'), -- Katherine likes Liam
    ('00000000-0000-0000-0000-000000000012', '00000000-0000-0000-0000-000000000011'), -- Liam likes Katherine (mutual)
    ('00000000-0000-0000-0000-000000000013', '00000000-0000-0000-0000-000000000014'), -- Mia likes Noah
    ('00000000-0000-0000-0000-000000000015', '00000000-0000-0000-0000-000000000016'), -- Olivia likes Paul
    ('00000000-0000-0000-0000-000000000017', '00000000-0000-0000-0000-000000000018'), -- Quincy likes Rachel
    ('00000000-0000-0000-0000-000000000018', '00000000-0000-0000-0000-000000000017'), -- Rachel likes Quincy (mutual)
    ('00000000-0000-0000-0000-000000000019', '00000000-0000-0000-0000-000000000020'), -- Sam likes Tina
    ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000003'), -- Bob likes Charlie
    ('00000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000001'), -- David likes Alice
    ('00000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000009'), -- Frank likes Ivy
    ('00000000-0000-0000-0000-000000000010', '00000000-0000-0000-0000-000000000005'), -- Jack likes Eva
    ('00000000-0000-0000-0000-000000000013', '00000000-0000-0000-0000-000000000011'), -- Mia likes Katherine
    ('00000000-0000-0000-0000-000000000014', '00000000-0000-0000-0000-000000000015'); -- Noah likes Olivia

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
  AND random() > 0.7 -- Randomly users to simulate those who haven't seen the liker yet
LIMIT 100; -- Limit the number of inserts for testing
