CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    mentor_email VARCHAR(255) NOT NULL,
    rating NUMERIC(3, 2) NOT NULL CHECK (rating >= 0 AND rating <= 5),
    comment TEXT,
    user_contact TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);