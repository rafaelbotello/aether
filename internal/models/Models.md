CREATE TABLE manga (
    manga_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    canonical_id TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    cover_url TEXT NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE chapters (
    chapter_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    manga_id UUID REFERENCES manga(manga_id) ON DELETE CASCADE,
    chapter_number NUMERIC(7,2) NOT NULL,
    source_name TEXT NOT NULL,
    source_url TEXT NOT NULL,
    released_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (manga_id, chapter_number, source_name)
);