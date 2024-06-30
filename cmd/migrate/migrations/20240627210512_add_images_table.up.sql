create table if not exists images (
    id serial primary key,
    user_id uuid references auth.users,
    status int not null default 1,
    image_location text,
    prompt text not null,
    batch_id uuid not null,
    deleted boolean not null default 'false',
    created_at timestamp not null default now(),
    deleted_at timestamp
)