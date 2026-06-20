-- Medium example: Multiple tables with relationships
CREATE TABLE authors (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    bio TEXT
);

CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    isbn TEXT UNIQUE NOT NULL,
    published_date DATE,
    author_id INTEGER REFERENCES authors(id) ON DELETE CASCADE
);

CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL,
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Queries
-- name: GetAuthor :one
SELECT * FROM authors WHERE id = $1;

-- name: ListAuthors :many
SELECT * FROM authors ORDER BY name;

-- name: GetAuthorBooks :many
SELECT * FROM books WHERE author_id = $1 ORDER BY published_date;

-- name: GetBook :one
SELECT * FROM books WHERE id = $1;

-- name: ListBooks :many
SELECT b.*, a.name as author_name 
FROM books b
JOIN authors a ON b.author_id = a.id
ORDER BY b.published_date DESC;

-- name: GetBookReviews :many
SELECT * FROM reviews WHERE book_id = $1 ORDER BY created_at DESC;

-- name: CreateAuthor :one
INSERT INTO authors (name, email, bio)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateBook :one
INSERT INTO books (title, isbn, published_date, author_id)
VALUES ($1, $2, $3, $4)
RETURNING *;
