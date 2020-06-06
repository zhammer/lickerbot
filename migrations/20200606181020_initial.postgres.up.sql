CREATE TABLE bootlickers (
    id UUID NOT NULL,
    PRIMARY KEY(id),
    twitter_user_id BIGINT NOT NULL UNIQUE,
    twitter_handle VARCHAR (255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
CREATE TABLE licks (
    id UUID NOT NULL,
    PRIMARY KEY(id),
    tweet_id BIGINT NOT NULL UNIQUE,
    tweet_text TEXT NOT NULL,
    bootlicker_id uuid NOT NULL REFERENCES bootlickers(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
CREATE TABLE pledged_donations (
    id UUID NOT NULL,
    PRIMARY KEY(id),
    bootlicker_id UUID NOT NULL REFERENCES bootlickers(id) ON DELETE CASCADE,
    amount integer NOT NULL CHECK (amount > 0),
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);
