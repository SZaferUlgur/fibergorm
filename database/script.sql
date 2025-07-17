--||-- 
create table if not exists `products` (
    `id` int(11) not null auto_increment,
    `name` varchar(100) not null,
    `price` decimal(9,2) not null,
    primary key (`id`)
) ENGINE=InnoDb default CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
--||--
create table if not exists `users` (
    `id` int(11) not null auto_increment,
    `name` varchar(100) not null,
    `email` varchar(100) not null,
    primary key (`id`)
) ENGINE=InnoDb default CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
--||--
create procedure if not exists `get_all_products`()
begin
    select * from products;
end;
--||--
create table if not exists `product_log` (
    `id` int(11) not null auto_increment,
    `product_id` int(11) not null,
    `action` varchar(100) not null,
    `created_at` datetime not null,
    primary key (`id`)
) ENGINE=InnoDb default CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
--||--
create trigger if not exists `product_log_trigger` after insert on products
for each row
begin
    insert into product_log (product_id, action, created_at) values (new.id, 'created', now());
end;
--||--