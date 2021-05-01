# 主调api使用
> 简单的做了三api供使用

## api的docker环境变量配置

+ WebApiStatus 是否开启api功能。这里填 "true"。默认为"false"
+ WebApiHost api的监听地址。这里留空即可监听 ipv4和ipv6
+ WebApiPort api监听端口。eg "9001"
+ WebApiToken api鉴权的token。就是下方url中的 access_token

## api接口
### 通过username 拉取用户id和chatid
> POST  [application/json]  http://ip:port/get_chat_info?access_token=你的token
>
> body内容： 
```json
{
    "name":"@LvanLamCommitCodeBot"
}
```

> response
```json
{
    "data": {
        "@type": "chat",
        "@extra": "QutcnQfoDmVkQWONMBgeFtrCHWHysEWp",
        "id": 1473260455,
        "type": {
            "@type": "chatTypePrivate",
            "@extra": "",
            "user_id": 1473260455
        },
        "title": "Commit Code Bot",
        "photo": {
            "@type": "chatPhotoInfo",
            "@extra": "",
            "small": {
                "@type": "file",
                "@extra": "",
                "id": 78,
                "size": 0,
                "expected_size": 0,
                "local": {
                    "@type": "localFile",
                    "@extra": "",
                    "path": "",
                    "can_be_downloaded": true,
                    "can_be_deleted": false,
                    "is_downloading_active": false,
                    "is_downloading_completed": false,
                    "download_offset": 0,
                    "downloaded_prefix_size": 0,
                    "downloaded_size": 0
                },
                "remote": {
                    "@type": "remoteFile",
                    "@extra": "",
                    "id": "AQADBQADEKsxG9pt6FcACMDFy2x0AAMCAAOnK9BXAAROGCIFw8huqDr_BQABHgQ",
                    "unique_id": "AQADwMXLbHQAAzr_BQAB",
                    "is_uploading_active": false,
                    "is_uploading_completed": true,
                    "uploaded_size": 0
                }
            },
            "big": {
                "@type": "file",
                "@extra": "",
                "id": 79,
                "size": 0,
                "expected_size": 0,
                "local": {
                    "@type": "localFile",
                    "@extra": "",
                    "path": "",
                    "can_be_downloaded": true,
                    "can_be_deleted": false,
                    "is_downloading_active": false,
                    "is_downloading_completed": false,
                    "download_offset": 0,
                    "downloaded_prefix_size": 0,
                    "downloaded_size": 0
                },
                "remote": {
                    "@type": "remoteFile",
                    "@extra": "",
                    "id": "AQADBQADEKsxG9pt6FcACMDFy2x0AAMDAAOnK9BXAAROGCIFw8huqDz_BQABHgQ",
                    "unique_id": "AQADwMXLbHQAAzz_BQAB",
                    "is_uploading_active": false,
                    "is_uploading_completed": true,
                    "uploaded_size": 0
                }
            },
            "has_animation": false
        },
        "permissions": {
            "@type": "chatPermissions",
            "@extra": "",
            "can_send_messages": true,
            "can_send_media_messages": true,
            "can_send_polls": true,
            "can_send_other_messages": true,
            "can_add_web_page_previews": true,
            "can_change_info": false,
            "can_invite_users": false,
            "can_pin_messages": true
        },
        "last_message": {
            "@type": "message",
            "@extra": "",
            "id": 1163919360,
            "sender": {
                "@type": "messageSenderUser",
                "@extra": "",
                "user_id": 1473260455
            },
            "chat_id": 1473260455,
            "sending_state": null,
            "scheduling_state": null,
            "is_outgoing": false,
            "is_pinned": false,
            "can_be_edited": false,
            "can_be_forwarded": true,
            "can_be_deleted_only_for_self": true,
            "can_be_deleted_for_all_users": true,
            "can_get_statistics": false,
            "can_get_message_thread": false,
            "is_channel_post": false,
            "contains_unread_mention": false,
            "date": 1619798701,
            "edit_date": 0,
            "forward_info": null,
            "interaction_info": null,
            "reply_in_chat_id": 0,
            "reply_to_message_id": 0,
            "message_thread_id": 0,
            "ttl": 0,
            "ttl_expires_in": 0,
            "via_bot_user_id": 0,
            "author_signature": "",
            "media_album_id": 0,
            "restriction_reason": "",
            "content": {
                "@type": "messageText",
                "@extra": "",
                "text": {
                    "@type": "formattedText",
                    "@extra": "",
                    "text": "你已经达到本周期内最大提交值了，下次再来吧",
                    "entities": []
                },
                "web_page": null
            },
            "reply_markup": null
        },
        "positions": [],
        "is_marked_as_unread": false,
        "is_blocked": false,
        "has_scheduled_messages": false,
        "can_be_deleted_only_for_self": true,
        "can_be_deleted_for_all_users": false,
        "can_be_reported": true,
        "default_disable_notification": false,
        "unread_count": 0,
        "last_read_inbox_message_id": 1163919360,
        "last_read_outbox_message_id": 1157627904,
        "unread_mention_count": 0,
        "notification_settings": {
            "@type": "chatNotificationSettings",
            "@extra": "",
            "use_default_mute_for": true,
            "mute_for": 0,
            "use_default_sound": true,
            "sound": "default",
            "use_default_show_preview": true,
            "show_preview": false,
            "use_default_disable_pinned_message_notifications": true,
            "disable_pinned_message_notifications": false,
            "use_default_disable_mention_notifications": true,
            "disable_mention_notifications": false
        },
        "action_bar": null,
        "voice_chat_group_call_id": 0,
        "is_voice_chat_empty": false,
        "reply_markup_message_id": 0,
        "draft_message": null,
        "client_data": ""
    },
    "retcode": 0,
    "status": "ok"
}
```

### 通过chatID发送文本消息（当前仅做了文本消息支持）
> ps: chat_id 是聊天id，一般不会变，可自行存储映射。
> 
> POST [application/json] http://ip:port/send_msg?access_token=你的token
>
> body json内容：
```json
{
    "chat_id":"1473260455",
    "message":{
        "msgtype":"messageText",
        "content":"123"
    }
}
```

### 获取当前登录用户的信息
> GET http://ip:port/getme?access_token=你的token
