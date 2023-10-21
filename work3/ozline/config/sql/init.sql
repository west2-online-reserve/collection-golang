create table west2online.`user`
(
    `id`          bigint        not null auto_increment,
    `username`    varchar(255)  not null,
    `password`    varchar(255)  not null,
    `created_at`  timestamp     not null default current_timestamp,
    `updated_at`  timestamp     not null on update current_timestamp default current_timestamp,
    `deleted_at`  timestamp     null default null,
    constraint `id`
        primary key (`id`)
) engine=InnoDB auto_increment=10000 default charset=utf8mb4;

create table west2online.`task`
(
    `id`            bigint        not null auto_increment,
    `user_id`        bigint        not null,
    `title`         varchar(255)  not null,
    `content`       varchar(255)  not null,
    `status`        bigint        not null,
    `start_at`      timestamp     not null,
    `end_at`        timestamp     not null,
    `created_at`    timestamp     not null default current_timestamp,
    `updated_at`    timestamp     not null on update current_timestamp default current_timestamp,
    `deleted_at`    timestamp     null default null,
    constraint `id`
        primary key (`id`)
) engine=InnoDB auto_increment=10000 default charset=utf8mb4;