CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    team_a VARCHAR(255) NOT NULL,
    team_b VARCHAR(255) NOT NULL,
    match_date DATE NOT NULL
);
