CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    homeTeam VARCHAR(255) NOT NULL,
    awayTeam VARCHAR(255) NOT NULL,
    matchDate DATE NOT NULL,
    goals INT DEFAULT 0,            
    yellowCards INT DEFAULT 0,     
    redCards INT DEFAULT 0,        
    extraTime INT DEFAULT 0        
);
