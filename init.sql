CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    team_a VARCHAR(255) NOT NULL,
    team_b VARCHAR(255) NOT NULL,
    score_a INT NOT NULL,
    score_b INT NOT NULL
);
