CREATE EXTENSION pg_trgm;
CREATE EXTENSION btree_gin;

CREATE INDEX IF NOT EXISTS posts_score_idx ON posts USING GIN (score);
CREATE INDEX IF NOT EXISTS posts_title_idx ON posts USING GIN (to_tsvector('simple', title));
