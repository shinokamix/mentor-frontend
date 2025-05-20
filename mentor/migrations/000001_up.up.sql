CREATE TABLE IF NOT EXISTS mentors (
    id SERIAL PRIMARY KEY,
    mentor_email TEXT UNIQUE NOT NULL,
    contact TEXT,
    count_reviews INTEGER NOT NULL DEFAULT 0,
    sum_rating FLOAT NOT NULL DEFAULT 0,
    average_rating FLOAT GENERATED ALWAYS AS (
        CASE 
            WHEN count_reviews = 0 THEN 0 
            ELSE ROUND( (sum_rating::NUMERIC / count_reviews)::NUMERIC, 1 ) 
        END
    ) STORED
);