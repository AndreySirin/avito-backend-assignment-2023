CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id uuid PRIMARY KEY default uuid_generate_v4(),

    full_name VARCHAR NOT NULL,
    gender VARCHAR NOT NULL,
    date_of_birth DATE NOT NULL,

    create_at TIMESTAMP NOT NULL DEFAULT now(),
    update_at TIMESTAMP NOT NULL DEFAULT now(),
    delete_at TIMESTAMP
);

CREATE TABLE segments (
    id uuid PRIMARY KEY default uuid_generate_v4(),

    title VARCHAR NOT NULL,
    description VARCHAR NOT NULL ,
    auto_user_prc SMALLINT NOT NULL,

    create_at TIMESTAMP NOT NULL DEFAULT now(),
    update_at TIMESTAMP NOT NULL DEFAULT now(),
    delete_at TIMESTAMP
);

CREATE TABLE subscriptions (
    PRIMARY KEY (user_id, segment_id),
    user_id uuid NOT NULL REFERENCES users(id),
    segment_id uuid NOT NULL REFERENCES segments(id),

    ttl TIMESTAMP,
    is_auto_add BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    update_at TIMESTAMP NOT NULL DEFAULT now(),
    delete_at TIMESTAMP
);

-- 🔹 Вставка юзеров
INSERT INTO users (id, full_name, gender, date_of_birth)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'Alice Ivanova', 'female', '1990-05-01'),
    ('22222222-2222-2222-2222-222222222222', 'Bob Petrov', 'male', '1985-09-15'),
    ('33333333-3333-3333-3333-333333333333', 'Charlie Smirnov', 'male', '1992-12-25'),
    ('44444444-4444-4444-4444-444444444444', 'Sirin ANDREY', 'male', '1995-12-25'),
    ('55555555-5555-5555-5555-555555555555', 'Sirin IGOR', 'male', '1997-12-25'),
    ('66666666-6666-6666-6666-666666666666', 'Sirin oleg', 'male', '1968-12-25'),
    ('77777777-7777-7777-7777-777777777777', 'pussan', 'female', '2000-12-25'),
    ('88888888-8888-8888-8888-888888888888', 'ronaldo', 'male', '1985-12-25');
-- 🔹 Вставка сегментов
INSERT INTO segments (id, title, description, auto_user_prc)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'Segment A', 'First segment for testing', 0),
    ('22222222-2222-2222-2222-222222222222', 'Segment B', 'Second segment for testing', 0),
    ('33333333-3333-3333-3333-333333333333', 'Segment C', 'Third segment for testing', 0),
    ('44444444-4444-4444-4444-444444444444', 'Segment d', 'Third segment for testing', 0),
    ('55555555-5555-5555-5555-555555555555', 'Segment e', 'Third segment for testing', 0),
    ('66666666-6666-6666-6666-666666666666', 'Segment f', 'Third segment for testing', 0),
    ('77777777-7777-7777-7777-777777777777', 'Segment g', 'Third segment for testing', 0);

INSERT INTO subscriptions (user_id, segment_id, ttl, is_auto_add)
VALUES
    (
        (SELECT id FROM users WHERE full_name = 'Alice Ivanova'),
        (SELECT id FROM segments WHERE title = 'Segment A'),
        '2025-05-01 12:00:00',
        true
    ),
    (
        (SELECT id FROM users WHERE full_name = 'Bob Petrov'),
        (SELECT id FROM segments WHERE title = 'Segment B'),
        now() + interval '15 days',
        false
    ),
    (
        (SELECT id FROM users WHERE full_name = 'Charlie Smirnov'),
        (SELECT id FROM segments WHERE title = 'Segment C'),
        '2025-07-01 12:00:00',
        true
    ),
    (
        (SELECT id FROM users WHERE full_name = 'Sirin ANDREY'),
        (SELECT id FROM segments WHERE title = 'Segment A'),
        '2022-07-01 12:00:00',
        false
    );
