create table if not exists "user"
(
    id            int       not null GENERATED BY DEFAULT AS IDENTITY primary key,
    email         varchar   not null,
    nickname      varchar   not null,
    fullname      varchar   not null,
    about         text      not null default ''
);

create table if not exists "forum"
(
    id         int         not null GENERATED BY DEFAULT AS IDENTITY primary key,
    posts      int         not null default 0,
    threads    int         not null default 0,
    title      varchar     not null,
    slug       varchar     not null,
    user_id    int         not null,
    author     varchar     not null,

    FOREIGN KEY (user_id) REFERENCES "user" (id)
);

create table if not exists "thread"
(
    id          int         not null GENERATED BY DEFAULT AS IDENTITY primary key,
    forum       varchar     not null,
    forum_id    int         not null,
    user_id     int         not null,
    author      varchar     not null,
    created     varchar     not null,
    message     varchar     not null,
    slug        varchar,
    title       varchar     not null,
    votes       int         not null default 0,

    foreign key (forum_id) references "forum" (id),
    foreign key (user_id) references "user" (id)
);

create table if not exists "post"
(
    id          int         not null GENERATED BY DEFAULT AS IDENTITY (START WITH 1) primary key,
    forum_id    int         not null, -- use
    forum       varchar     not null,
    thread_id   int         not null,
    user_id     int         not null,
    author      varchar     not null,
    created     varchar     not null,
    message     varchar     not null,
    isEdited    bool        not null default false,
    parent      int         not null default 0,

    foreign key (forum_id) references "forum" (id),
    foreign key (thread_id) references "thread" (id),
    foreign key (user_id) references "user" (id)
);

create table if not exists "vote"
(
    id          int     not null GENERATED BY DEFAULT AS IDENTITY primary key,
    thread_id   int     not null,
    nickname    varchar not null,
    voice       int     not null
);

create table if not exists "forum_user"
(
    forum   text not null,
    "user"  text not null
);
create unique index on "forum_user" ("user", "forum");

-- create index post_forum_id_idx on post using btree (forum_id);
-- create index post_thread_id_idx on post using btree (thread_id);
-- create index post_user_id_idx on post using btree (user_id);
-- create index post_forum_idx on post (forum);

create index vote_thread_user_idx on vote (thread_id, nickname);

create index thread_id_idx on thread (id);
create index thread_forum_slug_idx on thread (forum);
create index thread_user_id_idx on thread (user_id);
create index thread_slug_idx on thread (lower(slug));

create index forum_id_idx on forum (id);
create index forum_slug_idx on forum (lower(slug));

create index user_id_idx on "user" (id);
create index user_nickname_idx on "user" (lower(nickname));
create index user_email_idx on "user" (lower(email));

CREATE FUNCTION fn_insert_thread_votes()
    RETURNS TRIGGER AS '
    BEGIN
        UPDATE thread
        SET
            votes = votes + NEW.voice
        WHERE id = NEW.thread_id;
        RETURN NULL;
    END;
' LANGUAGE plpgsql;

CREATE TRIGGER on_vote_insert
    AFTER INSERT ON vote
    FOR EACH ROW EXECUTE PROCEDURE fn_insert_thread_votes();

CREATE FUNCTION fn_update_thread_votes()
    RETURNS TRIGGER AS '
    BEGIN
        IF OLD.voice = NEW.voice
            THEN RETURN NULL;
        END IF;
        UPDATE thread
        SET votes = votes +
            CASE
                WHEN NEW.voice = -1 THEN -2
                ELSE 2
            END
        WHERE id = NEW.thread_id;
        RETURN NULL;
    END;
' LANGUAGE plpgsql;

CREATE TRIGGER on_vote_update
    AFTER UPDATE ON vote
    FOR EACH ROW EXECUTE PROCEDURE fn_update_thread_votes();