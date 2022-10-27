-- +goose Up
create table chat_user
(
    id integer not null
        constraint chat_user_pk
            primary key autoincrement,
    chat_id integer not null,
    user_id integer not null,
    value integer default 1 not null,
    created_at datetime default current_timestamp not null,
    updated_at datetime default current_timestamp not null
);

create unique index chat_user_chat_id_user_id_uindex
    on chat_user (chat_id, user_id);

-- +goose Down
DROP TABLE chat_user;
