CREATE TABLE public.todos
(
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    active   BOOLEAN  NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);