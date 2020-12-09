create database EbookShop;
\c ebookshop;

create or replace function delete_books() returns trigger
    language plpgsql
as $$
BEGIN
    DELETE FROM comments
    WHERE comments.book_id = old.id;
    DELETE FROM carts
    WHERE carts.book_id = old.id;
    DELETE FROM purchased_books
    WHERE purchased_books.book_id = old.id;
    RETURN old;
end;
$$;

create or replace function delete_categories() returns trigger
    language plpgsql
as $$
BEGIN
    UPDATE books
    SET category_id = 1
    WHERE category_id = old.id;
    RETURN old;
END;
$$;


create table if not exists categories
(
    id bigserial not null
        constraint categories_pkey
            primary key,
    name text
);

alter table categories owner to postgres;

create table if not exists books
(
    id bigserial not null
        constraint books_pkey
            primary key,
    name text,
    price numeric,
    content text,
    author text,
    star_sum bigint,
    pay_number bigint,
    category_id bigint
        constraint fk_categories_books
            references categories
);

alter table books owner to postgres;

create trigger delete_book
    before delete
    on books
    for each row
execute procedure delete_books();

CREATE TRIGGER delete_category
    BEFORE DELETE
    ON categories
    FOR EACH ROW
EXECUTE PROCEDURE delete_categories();

create table if not exists users
(
    id bigserial not null
        constraint users_pkey
            primary key,
    mail text,
    username text not null
        constraint users_username_key
            unique,
    password_hash text not null,
    name text not null
);

alter table users owner to postgres;

create table if not exists carts
(
    user_id bigint not null
        constraint fk_users_carts
            references users,
    book_id bigint not null
        constraint fk_books_carts
            references books,
    constraint carts_pkey
        primary key (user_id, book_id)
);

alter table carts owner to postgres;

create table if not exists purchased_books
(
    time timestamp with time zone,
    user_id bigint not null
        constraint fk_users_purchased_books
            references users,
    book_id bigint not null
        constraint fk_books_purchased_books
            references books,
    constraint purchased_books_pkey
        primary key (user_id, book_id)
);

alter table purchased_books owner to postgres;

create table if not exists comments
(
    id bigserial not null
        constraint comments_pkey
            primary key,
    content text not null,
    star bigint,
    user_id bigint not null
        constraint fk_users_comments
            references users,
    book_id bigint not null
        constraint fk_books_comments
            references books
);

alter table comments owner to postgres;

alter function delete_categories() owner to postgres;

alter function delete_books() owner to postgres;

