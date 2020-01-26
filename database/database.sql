drop table if exists "user";
drop table if exists "forum";
drop table if exists "thread";
drop table if exists "post";
drop table if exists "vote";

create table if not exists "user"
(
    id            int       not null GENERATED BY DEFAULT AS IDENTITY primary key,
    email         varchar   not null,
    nickname      varchar   not null unique,
    fullname      varchar   not null,
    about         text      not null default ''
);

create table if not exists "forum"
(
    id         int         not null GENERATED BY DEFAULT AS IDENTITY primary key,
    posts      int         not null default 0,
    threads    int         not null default 0,
    title      varchar     not null,
    slug       varchar     not null UNIQUE,
    author     varchar     NOT NULL REFERENCES "user" (nickname)
);

create table if not exists "thread"
(
    id          int         not null GENERATED BY DEFAULT AS IDENTITY primary key,
    forum       varchar     NOT NULL REFERENCES "forum" (slug),
    author      varchar     NOT NULL REFERENCES "user" (nickname),
    created     varchar     not null,
    message     varchar     not null,
    slug        varchar,
    title       varchar     not null,
    votes       int         not null default 0
);

-- CREATE OR REPLACE FUNCTION fn_check_post_before_insert(parent_id int, post_thread_id int) RETURNS BOOL AS $$
-- BEGIN
--     IF parent_id = 0 OR (SELECT COUNT(*) FROM post WHERE id=parent_id) > 0 AND (SELECT thread_id FROM post WHERE id=parent_id) = post_thread_id THEN
--         RETURN TRUE;
--     ELSE
--         RETURN FALSE;
--     end if;
-- END;
-- $$ LANGUAGE plpgsql;

create table if not exists "post"
(
    id          int         not null GENERATED BY DEFAULT AS IDENTITY (START WITH 1) primary key,
    forum       varchar     not null,
    thread_id   int         not null,
    author      varchar     not null,
    created     varchar     not null,
    message     varchar     not null,
    isEdited    bool        not null default false,
    parent      int         not null default 0, -- references "post" (id) on delete cascade on update restrict
        --CONSTRAINT post_parent_constraint CHECK (fn_check_parent_post_same_thread(parent)=thread_id),
    path        text,

    foreign key (thread_id) references "thread" (id),
    foreign key (author) references "user" (nickname)

    --CONSTRAINT post_parent_constraint CHECK (fn_check_post_before_insert(parent, thread_id) = true)
);

create unlogged table if not exists "vote"
(
    thread_id   int     not null,
    nickname    varchar not null,
    voice       int     not null
    --UNIQUE (thread_id, nickname)
);

create table if not exists "forum_user"
(
    forum           varchar not null,
    user_nickname   varchar not null references "user" (nickname),
    UNIQUE (forum, user_nickname)
);

-- Indexes

create index post_thread_id_idx on post (thread_id);
-- create index post_forum_idx on post (forum);

create index vote_thread_user_idx on vote (thread_id, nickname);

create index thread_id_idx on thread (id);
create index thread_forum_slug_idx on thread (forum);
create index thread_slug_idx on thread (lower(slug));

--create index forum_id_idx on forum (id);
create index forum_slug_idx on forum (lower(slug));
create index forum_author_idx on forum (author);

--create index user_id_idx on "user" (id);
create index user_nickname_idx on "user" (lower(nickname));
create index user_email_idx on "user" (lower(email));

create unique index forum_users_idx ON forum_user(forum, user_nickname);

-- Functions

CREATE FUNCTION fn_inc_forum_thread()
    RETURNS TRIGGER AS '
        BEGIN
            UPDATE forum
            SET threads = threads + 1
            WHERE slug = NEW.forum;
            RETURN NEW;
        END;
' LANGUAGE plpgsql;

CREATE TRIGGER thread_insert
    AFTER INSERT ON thread
    FOR EACH ROW EXECUTE PROCEDURE fn_inc_forum_thread();

CREATE FUNCTION fn_inc_forum_post()
    RETURNS TRIGGER AS '
        BEGIN
            UPDATE forum
            SET posts = posts + 1
            WHERE slug = NEW.forum;
            RETURN NEW;
        END;
' LANGUAGE plpgsql;

CREATE TRIGGER post_insert
    AFTER INSERT ON post
    FOR EACH ROW EXECUTE PROCEDURE fn_inc_forum_post();

CREATE OR REPLACE FUNCTION fn_update_thread_votes() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP='INSERT' THEN
        UPDATE thread SET votes=votes+NEW.voice WHERE id=NEW.thread_id;
        RETURN NEW;
    ELSIF TG_OP='UPDATE' THEN
        UPDATE thread SET votes=votes+(NEW.voice-OLD.voice) WHERE id=NEW.thread_id;
        RETURN NEW;
    ELSE
        RAISE EXCEPTION 'Invalid call fn_update_thread_votes()';
    end if;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_thread_vote ON vote;
CREATE TRIGGER update_thread_vote
    AFTER INSERT OR UPDATE ON vote
    FOR EACH ROW EXECUTE PROCEDURE fn_update_thread_votes();

CREATE OR REPLACE FUNCTION fn_users_forum_thread() RETURNS TRIGGER AS $$
begin
    IF NEW.forum IS NOT NULL then
        INSERT INTO forum_user(forum, user_nickname) VALUES (NEW.forum, NEW.author)
        ON CONFLICT DO NOTHING;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_forum
    AFTER INSERT ON thread
    FOR EACH ROW EXECUTE PROCEDURE fn_users_forum_thread();

CREATE OR REPLACE FUNCTION fn_users_forum_post() RETURNS TRIGGER AS $$
begin
    IF NEW.forum IS NOT NULL THEN
        INSERT INTO forum_user(forum, user_nickname) VALUES (NEW.forum, NEW.author)
        ON CONFLICT DO NOTHING;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_forum
    AFTER INSERT ON post
    FOR EACH ROW EXECUTE PROCEDURE fn_users_forum_post();

CREATE OR REPLACE FUNCTION fn_create_post_path() RETURNS TRIGGER AS
$$
BEGIN
    IF NEW.parent = 0 THEN
        NEW.path = CAST( NEW.id+10000000 AS varchar(100));
    ELSE
        NEW.path = CAST((SELECT path FROM post WHERE id = NEW.parent) || '->' || NEW.id+10000000 AS varchar(100));
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER post_path
    BEFORE INSERT ON post
    FOR EACH ROW EXECUTE PROCEDURE fn_create_post_path();