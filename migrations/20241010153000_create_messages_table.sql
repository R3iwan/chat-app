CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    sender_id INT REFERENCES users(id),
    receiver_id INT REFERENCES users(id),
    content TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT NOW()
);
