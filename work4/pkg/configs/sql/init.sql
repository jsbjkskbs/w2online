SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

drop table if exists `users`;
create table `users`(
    `uid` bigint not null auto_increment comment '用户ID',
    `username` varchar(255) not null comment '用户名称',
    `password` varchar(255) not null comment '用户密码',
    `avatar_url` varchar(255) comment '用户头像url',
    `created_at` bigint not null comment '创建账号时间',
    `updated_at` bigint not null comment '最近登录时间',
    `deleted_at` bigint not null comment '账号删除时间',
    `mfa_enable` bool not null comment '是否使用mfa',
    `mfa_secret` varchar(255) comment 'mfa秘钥',
    primary key (uid),
    key `username_password_index` (username,password) using btree comment '用户名与密码索引'
) engine =InnoDB auto_increment=10000 default charset =utf8mb4 comment ='用户表';

drop table if exists `videos`;
create table `videos`(
    `id` bigint not null auto_increment comment '视频ID',
    `user_id` bigint not null comment '作者ID',
    `video_url` varchar(255) not null comment '视频url',
    `cover_url` varchar(255) not null comment '封面url',
    `title` varchar(255) not null comment '标题',
    `description` varchar(255) not null comment '简介',
    `visit_count` bigint not null comment '浏览量',
    `created_at` bigint not null comment '创建时间',
    `updated_at` bigint not null comment '修改时间',
    `deleted_at` bigint not null comment '删除时间',
    primary key (id),
    foreign key (user_id) references users(uid) on delete cascade on update cascade,
    key `time` (created_at) using btree comment '时间查询索引',
    key `author` (user_id) using btree comment '作者查询索引'
) engine =InnoDB auto_increment=10000 default charset=utf8mb4 comment '视频表';

drop table if exists `comments`;
create table `comments`(
    `id` bigint not null auto_increment comment '评论ID',
    `user_id` bigint not null comment '作者ID',
    `video_id` bigint not null comment '视频ID',
    `parent_id` bigint not null comment '父评论ID', /* 没有则为-1 */
    `content` varchar(255) not null comment '评论内容',
    `created_at` bigint not null comment '创建时间',
    `updated_at` bigint not null comment '更新时间',
    `deleted_at` bigint not null comment '删除时间',
    primary key (id),
    foreign key (user_id) references users(uid) on delete cascade on update cascade,
    foreign key (video_id) references videos(id) on delete cascade on update cascade,
    foreign key (parent_id) references comments(id) on delete cascade on update cascade,
    key `video_index` (video_id) using btree comment '视频ID索引'
) engine =InnoDB auto_increment=10000 default charset =utf8mb4 comment '评论表';
insert ignore into comments values (-1,0,0,-1,'',0,0,0); /* 预留id,防止无法插入 */

drop table if exists `video_likes`;
create table `video_likes`(
    `id` bigint not null auto_increment comment '自增关系序号',
    `user_id` bigint not null comment '点赞者ID',
    `video_id` bigint not null comment '被点赞视频ID',
    `created_at` bigint not null comment '确认时间',
    `deleted_at` bigint not null comment '取消时间', /* 默认设为时间戳最小值 */
    primary key (id),
    unique key `user_id_video_id_no_duplicate` (user_id,video_id),
    foreign key (user_id) references users(uid) on delete cascade on update cascade,
    foreign key (video_id) references videos(id) on delete cascade on update cascade,
    key `user_id_video_id_index` (user_id,video_id) using btree comment '点赞者与视频索引',
    key `user_id_index` (user_id) using btree comment '点赞者ID索引',
    key `videos_id_index` (video_id) using btree comment '视频ID索引'
) engine =InnoDB auto_increment =10000 default charset =utf8mb4 comment '视频点赞表';

drop table if exists `comment_likes`;
create table `comment_likes`(
    `id`         bigint not null auto_increment comment '自增关系序号',
    `user_id`    bigint not null comment '点赞者ID',
    `comment_id` bigint not null comment '评论ID',
    `created_at` bigint not null comment '确认时间',
    `deleted_at` bigint not null comment '取消时间',
    primary key (id),
    unique key `user_id_comment_id_no_duplicate` (user_id,comment_id),
    foreign key (user_id) references users(uid) on delete cascade on update cascade,
    foreign key (comment_id) references comments(id) on delete cascade on update cascade,
    key `user_id_comment_id_index` (user_id, comment_id) using btree comment '点赞者与评论索引',
    key `user_id_index` (user_id) using btree comment '点赞者索引',
    key `comment_id_index` (comment_id) using btree comment '评论ID索引'
) engine =InnoDB auto_increment =10000 default charset =utf8mb4 comment '评论点赞表';

drop table if exists `follows`;
create table `follows`(
    `id` bigint not null auto_increment comment '自增关系序号',
    `followed_id` bigint not null comment '被关注者ID',
    `follower_id` bigint not null comment '粉丝ID',
    `created_at` bigint not null comment '关系创建时间',
    `deleted_at` bigint not null comment '关系取消时间',
    primary key (id),
    unique key `follower_followed_no_duplicate` (followed_id,follower_id),
    foreign key (follower_id) references users(uid) on delete cascade on update cascade,
    foreign key (followed_id) references users(uid) on delete cascade on update cascade,
    key `followed_id_follower_id_index` (followed_id,follower_id) using btree comment '被关注者与粉丝索引',
    key `followed_id_index` (followed_id) using btree comment '被关注者索引',
    key `follower_id_index` (follower_id) using btree comment '粉丝索引'
) engine =InnoDB auto_increment =10000 default charset =utf8mb4 comment '关注表';

drop table if exists `messages`;
create table `messages`(
    `id`           bigint       not null auto_increment comment '自增记录序号',
    `from_user_id` bigint       not null comment '发送者ID',
    `to_user_id`   bigint       not null comment '接受者ID',
    `content`      varchar(255) not null comment '内容',
    `created_at`   bigint    not null comment '创建时间',
    `deleted_at`   bigint    not null comment '删除时间',
    primary key (`id`),
    foreign key (from_user_id) references users(uid) on delete cascade on update cascade,
    foreign key (to_user_id) references users(uid) on delete cascade on update cascade,
    key `from_user_id_to_user_id_index` (`from_user_id`,`to_user_id`) using btree comment '发送者与接受者索引',
    key `from_user_id_to_user_id_created_at_index` (`from_user_id`,`to_user_id`,`created_at`) using btree comment '发送者与接受者的时间段索引',
    key `from_user_id_created_at_index` (`from_user_id`,`created_at`) using btree comment '发送者与发送时间索引', /* 一般不会用到 */
    key `created_at_index` (`created_at`) using btree comment '创建时间索引' /* 一般不会用到 */
) engine =InnoDB auto_increment =10000 default charset =utf8mb4 comment '消息表';
