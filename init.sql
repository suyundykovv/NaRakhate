CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user'
);

CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    start_time TIMESTAMP NOT NULL,
    category VARCHAR(50),
    referee VARCHAR(100),
    -- Venue details
    venue_name VARCHAR(100),
    venue_city VARCHAR(100),
    -- Teams details
    home_team_name VARCHAR(100),
    home_team_logo VARCHAR(255),
    away_team_name VARCHAR(100),
    away_team_logo VARCHAR(255),
    -- Match result
    home_goals INT DEFAULT 0,
    away_goals INT DEFAULT 0,
    -- Match status
    match_status VARCHAR(50),
    -- Periods and score details (if available)
    halftime_home_goals INT DEFAULT 0,
    halftime_away_goals INT DEFAULT 0,
    fulltime_home_goals INT DEFAULT 0,
    fulltime_away_goals INT DEFAULT 0
);

CREATE TABLE data (
    key VARCHAR(255) PRIMARY KEY,
    value TEXT
);
CREATE TABLE bets (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    event_id INT REFERENCES events(id) ON DELETE CASCADE,
    amount NUMERIC(10, 2) NOT NULL,
    outcome VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE IF NOT EXISTS wheel_rewards (
    id SERIAL PRIMARY KEY,
    reward_name VARCHAR(255) NOT NULL,
    reward_type VARCHAR(50) NOT NULL,
    reward_value INT NOT NULL,
    probability FLOAT NOT NULL
    );

ALTER TABLE users ADD COLUMN IF NOT EXISTS last_spin_time TIMESTAMP DEFAULT '1970-01-01 00:00:00';
ALTER TABLE users ADD COLUMN IF NOT EXISTS spin_count INT DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS wincash FLOAT DEFAULT 0;


INSERT INTO wheel_rewards (reward_name, reward_type, reward_value, probability)
VALUES
    ('Фрибет 100', 'bonus_money', 100, 50.0),
    ('Кэшбек 50%', 'cashback', 50, 30.0),
    ('Джекпот 1000', 'bonus_money', 1000, 20.0);

INSERT INTO users (username, email, password, role) VALUES
('admin', 'admin@example.com', 'hashedpassword1', 'admin'),
('user1', 'user1@example.com', 'hashedpassword2', 'user'),
('user2', 'user2@example.com', 'hashedpassword3', 'user'),
('Aday','adaydhx@gmail.com','Aday2004','user');

INSERT INTO categories (name, description) VALUES
('Sports', 'Sports events like football or basketball'),
('Esports', 'Esports events like Dota 2 or League of Legends'),
('Politics', 'Political events and elections');

INSERT INTO events (
    name, description, start_time, category, referee, venue_name, venue_city, 
    home_team_name, home_team_logo, away_team_name, away_team_logo, 
    home_goals, away_goals, match_status, halftime_home_goals, halftime_away_goals, 
    fulltime_home_goals, fulltime_away_goals
) VALUES (
    'Burnley vs Manchester City', 
    'Football match description', 
    '2023-08-11 19:00:00', 
    'Sports', 
    'C. Pawson', 
    'Turf Moor', 
    'Burnley', 
    'Burnley', 
    'https://media.api-sports.io/football/teams/44.png',
    'Manchester City', 
    'https://media.api-sports.io/football/teams/50.png',
    0, 
    3, 
    'Match Finished',
    0,
    2,
    0,
    3
);