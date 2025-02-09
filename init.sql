CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    cash NUMERIC(10, 2) DEFAULT 0.00
);
CREATE TABLE matches (
    id SERIAL PRIMARY KEY,               
    home_team_id INT NOT NULL,          
    away_team_id INT NOT NULL,           
    home_goals INT NOT NULL,             
    away_goals INT NOT NULL,             
    match_date TIMESTAMP NOT NULL,      
    league_id INT,                       
    referee VARCHAR(100),               
    venue_name VARCHAR(100),            
    venue_city VARCHAR(100)             
);
CREATE INDEX idx_home_team ON matches (home_team_id);
CREATE INDEX idx_away_team ON matches (away_team_id);
CREATE INDEX idx_match_date ON matches (match_date);

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
    fulltime_away_goals INT DEFAULT 0,
    -- Odds
    home_win_odds FLOAT,
    away_win_odds FLOAT,
    draw_odds FLOAT
);
CREATE TABLE data (
    key VARCHAR(255) PRIMARY KEY,
    value TEXT
);
CREATE TABLE bets (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    event_id INT REFERENCES events(id) ON DELETE CASCADE,
    odd_selection VARCHAR(50),
    odd_value NUMERIC(10, 2),
    amount NUMERIC(10, 2) NOT NULL,
    income NUMERIC(10, 2),
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT
);

INSERT INTO users (username, email, password, role, cash) VALUES
('admin', 'admin@example.com', 'hashedpassword1', 'admin', '1000'),
('user1', 'user1@example.com', 'hashedpassword2', 'user','100'),
('user2', 'user2@example.com', 'hashedpassword3', 'user','10000'),
('Aday','adaydhx@gmail.com','Aday2004','user','1000');

INSERT INTO categories (name, description) VALUES
('Sports', 'Sports events like football or basketball'),
('Esports', 'Esports events like Dota 2 or League of Legends'),
('Politics', 'Political events and elections');

INSERT INTO events (
    name, description, start_time, category, referee, venue_name, venue_city, 
    home_team_name, home_team_logo, away_team_name, away_team_logo, 
    home_goals, away_goals, match_status, halftime_home_goals, halftime_away_goals, 
    fulltime_home_goals, fulltime_away_goals,
    home_win_odds, away_win_odds, draw_odds
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
    3,
    5.0,  
    1.5,  
    3.8   
);

INSERT INTO matches (home_team_id, away_team_id, home_goals, away_goals, match_date, league_id, referee, venue_name, venue_city)
VALUES 
  (1, 2, 2, 1, '2023-08-11 19:00:00', 140, 'John Doe', 'Stadium A', 'City A'),
  (1, 4, 1, 1, '2023-08-11 20:00:00', 140, 'Jane Smith', 'Stadium B', 'City B'),
  (1, 6, 3, 2, '2023-08-11 21:00:00', 140, 'Jim Beam', 'Stadium C', 'City C'),
  (1, 8, 0, 0, '2023-08-11 22:00:00', 140, 'Alice Johnson', 'Stadium D', 'City D'),
  (1, 10, 2, 2, '2023-08-11 23:00:00', 140, 'Bob Brown', 'Stadium E', 'City E'),
  (1, 12, 1, 0, '2023-08-12 19:00:00', 140, 'Carol White', 'Stadium F', 'City F'),
  (1, 14, 0, 1, '2023-08-12 20:00:00', 140, 'David Green', 'Stadium G', 'City G'),
  (1, 16, 4, 2, '2023-08-12 21:00:00', 140, 'Eva Black', 'Stadium H', 'City H'),
  (1, 18, 1, 1, '2023-08-12 22:00:00', 140, 'Frank Blue', 'Stadium I', 'City I'),
  (1, 20, 3, 3, '2023-08-12 23:00:00', 140, 'Grace Yellow', 'Stadium J', 'City J');
INSERT INTO matches (home_team_id, away_team_id, home_goals, away_goals, match_date, league_id, referee, venue_name, venue_city)
VALUES 
  (2, 22, 2, 0, '2023-08-13 19:00:00', 140, 'Henry Orange', 'Stadium K', 'City K'),
  (2, 24, 1, 2, '2023-08-13 20:00:00', 140, 'Ivy Violet', 'Stadium L', 'City L'),
  (2, 26, 3, 1, '2023-08-13 21:00:00', 140, 'Jack Black', 'Stadium M', 'City M'),
  (2, 28, 0, 0, '2023-08-13 22:00:00', 140, 'Kelly White', 'Stadium N', 'City N'),
  (2, 30, 2, 1, '2023-08-13 23:00:00', 140, 'Leo Brown', 'Stadium O', 'City O'),
  (3, 32, 1, 1, '2023-08-14 19:00:00', 140, 'Mia Gray', 'Stadium P', 'City P'),
  (3, 34, 0, 2, '2023-08-14 20:00:00', 140, 'Nina Green', 'Stadium Q', 'City Q'),
  (3, 36, 4, 3, '2023-08-14 21:00:00', 140, 'Oscar Blue', 'Stadium R', 'City R'),
  (3, 38, 1, 0, '2023-08-14 22:00:00', 140, 'Paul Red', 'Stadium S', 'City S'),
  (3, 40, 3, 2, '2023-08-14 23:00:00', 140, 'Quinn Yellow', 'Stadium T', 'City T');
INSERT INTO bets (user_id, event_id, odd_selection, odd_value, amount, income, status) VALUES
(1, 1, 'home', 2.50, 100.00, 250.00, 'open'),
(2, 1, 'draw', 3.20, 50.00, 160.00, 'closed');