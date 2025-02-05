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
    category VARCHAR(50)
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

CREATE TABLE IF NOT EXISTS leaderboard (
                                           id SERIAL PRIMARY KEY,
                                           user_id INT NOT NULL,
                                           total_win FLOAT NOT NULL,
                                           updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Пример данных
INSERT INTO leaderboard (user_id, total_win) VALUES (1, 1000.50);
INSERT INTO leaderboard (user_id, total_win) VALUES (2, 1500.75);
INSERT INTO leaderboard (user_id, total_win) VALUES (3, 2000.00);


INSERT INTO users (username, email, password, role) VALUES
('admin', 'admin@example.com', 'hashedpassword1', 'admin'),
('user1', 'user1@example.com', 'hashedpassword2', 'user'),
('user2', 'user2@example.com', 'hashedpassword3', 'user'),
('Aday','adaydhx@gmail.com','Aday2004','user');

INSERT INTO categories (name, description) VALUES
('Sports', 'Sports events like football or basketball'),
('Esports', 'Esports events like Dota 2 or League of Legends'),
('Politics', 'Political events and elections');

INSERT INTO events (name, description, start_time, category) VALUES
('Football World Cup Final', 'The final match of the Football World Cup', '2025-07-10 18:00:00', 'Sports'),
('Dota 2 International', 'The annual Dota 2 championship', '2025-08-15 12:00:00', 'Esports'),
('US Presidential Election', '2028 US Presidential Election', '2028-11-05 00:00:00', 'Politics');
