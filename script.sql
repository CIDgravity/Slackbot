create sequence users_seq;

alter sequence users_seq owner to twinquasar_sat;

create sequence miners_seq;

alter sequence miners_seq owner to twinquasar_sat;

create sequence registration_seq;

alter sequence registration_seq owner to twinquasar_sat;

create table if not exists users
(
    id            integer default nextval('users_seq'::regclass) not null,
    email         varchar,
    github_user   varchar,
    slack_user_id varchar,
    slack_team_id varchar,
    username      varchar,
    last_update   timestamp,
    constraint users_pkey
        primary key (id)
);

alter table users
    owner to twinquasar_sat;

create table if not exists miners
(
    id               integer default nextval('miners_seq'::regclass) not null,
    miner_address    varchar,
    owner_key        varchar,
    user_id          integer,
    is_active        boolean,
    sign_process_key varchar,
    constraint miners_pkey
        primary key (id),
    constraint miners_users_id_fk
        foreign key (user_id) references users
);

alter table miners
    owner to twinquasar_sat;

create table if not exists usage_limits
(
    id            serial not null,
    name          varchar,
    default_value integer default 2
);

alter table usage_limits
    owner to twinquasar_sat;

create unique index if not exists usage_limits_id_uindex
    on usage_limits (id);

create table if not exists users_has_usage_limits
(
    user_id        integer,
    usage_limit_id integer,
    value          integer,
    constraint users_has_usage_limits_users_id_fk
        foreign key (user_id) references users,
    constraint users_has_usage_limits_usage_limits_id_fk
        foreign key (usage_limit_id) references usage_limits (id)
);

alter table users_has_usage_limits
    owner to twinquasar_sat;

create table if not exists blocks
(
    cid             text    not null,
    height          varchar,
    block_timestamp varchar,
    miner           varchar,
    is_validated    boolean,
    win_count       bigint,
    reward          numeric not null,
    constraint table_name_pk
        primary key (cid)
);

alter table blocks
    owner to twinquasar_sat;

create unique index if not exists table_name_cid_uindex
    on blocks (cid);

create table if not exists block_messages
(
    id              serial not null,
    cid             text,
    miner_to        varchar,
    miner_from      varchar,
    actor           varchar,
    decoded_message text,
    method_name     varchar,
    gas_fee_cap     numeric,
    gas_limit       numeric,
    gas_premium     numeric,
    message_value   numeric,
    status          boolean,
    block_cid       text,
    nonce           integer,
    duplicate       boolean,
    constraint block_messages_pk
        primary key (id),
    constraint block_messages_blocks_cid_fk
        foreign key (block_cid) references blocks
);

alter table block_messages
    owner to twinquasar_sat;

create unique index if not exists block_messages_id_uindex
    on block_messages (id);

create or replace function notify_event() returns trigger
    language plpgsql
as
$$
DECLARE
    data         json;
    notification json;

BEGIN

    -- Convert the old or new row to JSON, based on the kind of action.
    -- Action = DELETE?             -> OLD row
    -- Action = INSERT or UPDATE?   -> NEW row
    IF (TG_OP = 'DELETE') THEN
        data = row_to_json(OLD);
    ELSE
        data = row_to_json(NEW);
    END IF;

    -- Contruct the notification as a JSON string.
    notification = json_build_object(
            'table', TG_TABLE_NAME,
            'action', TG_OP,
            'data', data);


    -- Execute pg_notify(channel, notification)
    PERFORM pg_notify('events', notification::text);

    -- Result is ignored since this is an AFTER trigger
    RETURN NULL;
END;

$$;

alter function notify_event() owner to twinquasar_sat;

create trigger blocks_notify_event
    after insert or update or delete
    on blocks
    for each row
execute procedure notify_event();
