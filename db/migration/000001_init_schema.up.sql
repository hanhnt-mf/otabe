CREATE TABLE `restaurant` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(100) NOT NULL,
    `website` varchar(100) DEFAULT '',
    `phone` varchar(25) DEFAULT '',
    `description` varchar(255) NOT NULL,
    `postal_code` varchar(25) NOT NULL,
    `address` varchar(30) NOT NULL,
    `long` decimal(12,9) NOT NULL,
    `lat` decimal(12,9) NOT NULL,
    `created_at` date DEFAULT NULL,
    `updated_at` date DEFAULT NULL,
    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `user` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_name` varchar(100) NOT NULL,
    `nation` varchar(25) NOT NULL,
    `password` varchar(100) NOT NULL,
    `role` varchar(100) DEFAULT "user",
    `created_at` date DEFAULT NULL,
    `updated_at` date DEFAULT NULL,
    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `menu` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `restaurant_id` bigint(20) unsigned NOT NULL,
    `name` varchar(100) NOT NULL,
    `created_at` date DEFAULT NULL,
    `updated_at` date DEFAULT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`restaurant_id`) REFERENCES `restaurant`(`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `item` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `menu_id` bigint(20) unsigned NOT NULL,
    `name` varchar(100) NOT NULL,
    `description` varchar(255) NOT NULL,
    `price` varchar(20) NOT NULL,
    `created_at` date DEFAULT NULL,
    `updated_at` date DEFAULT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`menu_id`) REFERENCES `menu`(`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `item_feedback` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_id` bigint(20) unsigned NOT NULL,
    `item_id` bigint(20) unsigned NOT NULL,
    `comment` varchar(255) DEFAULT '',
    `rate` int NOT NULL,
    `created_at` date DEFAULT NULL,
    `updated_at` date DEFAULT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`),
    FOREIGN KEY (`item_id`) REFERENCES `item`(`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

insert into `restaurant`(`name`,`website`,`phone`,`description`,`postal_code`, `address`, `long`, `lat`)
values
    ("HaNoi & Hanoi","hanoi.com","09642540626","oishii","0600004", "1-3-28 Shaiba Tokyo", 23.088, 978.23423);

insert into `menu`(`restaurant_id`,`name`)
values
    (1, "Main");

insert into `item`(`menu_id`,`name`,`description`, `price`)
values
    (1, "Banh mi nhan thit", "With pork inside", "5300");

insert into `item`(`menu_id`,`name`,`description`, `price`)
values
    (1, "Banh mi rau", "Vetgetable", "3400");

insert into `user`(`user_name`, `nation`, `password`, `role`)
values
    ("HaNa", "Vietnamese", "$2a$10$d/xK0aG7NEo5BGPgjmbmUObD1.EucgeSRHKKzi9.UoqePZOwJxEQS", "admin");

insert into `item_feedback`(`user_id`,`item_id`,`comment`,`rate`)
values
    (1, 1, "Oishii ne", 5);


