create table if not exists "user"
(
    id            int       not null primary key,
    email         varchar   not null,
    nickname      varchar   not null,
    fullname      varchar   not null,
    about         varchar   not null
);

create table if not exists "forum"
(
    id         int         not null primary key,
    posts      int         not null,
    threads    int         not null,
    title      varchar     not null,
    slug       varchar     not null,
    user_id    int     not null,

    FOREIGN KEY (user_id) REFERENCES "user" (id)
);

create table if not exists "thread"
(
    id          int         not null primary key,
    forum_id    int         not null,
    user_id     int         not null,
    created     date        not null,
    message     varchar     not null,
    slug        varchar     not null,
    title       varchar     not null,
    votes       int         not null,

    foreign key (forum_id) references "forum" (id),
    foreign key (user_id) references "user" (id)
);

create table if not exists "post"
(
    id          int         not null primary key,
    forum_id    int         not null,
    thread_id   int         not null,
    user_id     int         not null,
    created     date        not null,
    message     varchar     not null,
    title       varchar     not null,
    isEdited    bool        not null,
    parent      int         not null,

    foreign key (forum_id) references "forum" (id),
    foreign key (thread_id) references "thread" (id),
    foreign key (user_id) references "user" (id)
);

create table if not exists "vote"
(
    thread_id   int     not null,
    nickname    varchar not null,
    vote        int,

    foreign key (thread_id) references "thread" (id)
);
