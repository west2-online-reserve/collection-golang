| 层级 | 字段名  | 数据类型 | 含义说明              |
| :--- | :------ | :------- | :-------------------- |
| 根级 | code    | number   | 接口状态码，0表示成功 |
| 根级 | message | string   | 状态消息              |
| 根级 | ttl     | number   | 缓存生存时间（秒）    |
| 根级 | data    | object   | 主要数据容器          |

## data

### cursor

| 字段名           | 数据类型 | 含义说明           |
| :--------------- | :------- | :----------------- |
| is_begin         | boolean  | 是否在第一页       |
| prev             | number   | 上一页标识         |
| next             | number   | 下一页标识         |
| is_end           | boolean  | 是否在最后一页     |
| pagination_reply | object   | 分页回复信息       |
| session_id       | string   | 会话ID             |
| mode             | number   | 显示模式           |
| mode_text        | string   | 模式文本           |
| all_count        | number   | 总评论数           |
| support_mode     | array    | 支持的显示模式列表 |
| name             | string   | 分页名称           |

### replies

| 字段名    | 数据类型 | 含义说明         |
| :-------- | :------- | :--------------- |
| rpid      | number   | 评论ID           |
| oid       | number   | 对象ID（视频ID） |
| type      | number   | 评论类型         |
| mid       | number   | 用户ID           |
| root      | number   | 根评论ID         |
| parent    | number   | 父评论ID         |
| dialog    | number   | 对话ID           |
| count     | number   | 子评论数量       |
| rcount    | number   | 回复数量         |
| state     | number   | 评论状态         |
| fansgrade | number   | 粉丝等级         |
| attr      | number   | 属性标记         |
| ctime     | number   | 创建时间戳       |
| *_str     | string   | 各ID的字符串形式 |
| like      | number   | 点赞数           |
| action    | number   | 用户操作状态     |

#### member

| 字段名           | 数据类型 | 含义说明     |
| :--------------- | :------- | :----------- |
| mid              | string   | 用户ID       |
| uname            | string   | 用户名       |
| sex              | string   | 性别         |
| sign             | string   | 个性签名     |
| avatar           | string   | 头像URL      |
| rank             | string   | 用户等级     |
| level_info       | object   | 等级详细信息 |
| pendant          | object   | 头像挂件     |
| nameplate        | object   | 成就勋章     |
| official_verify  | object   | 官方认证     |
| vip              | object   | 会员信息     |
| is_senior_member | number   | 是否正式会员 |

#### level_info

| 字段名        | 数据类型 | 含义说明         |
| :------------ | :------- | :--------------- |
| current_level | number   | 当前等级         |
| current_min   | number   | 当前等级最低经验 |
| current_exp   | number   | 当前经验值       |
| next_exp      | number   | 下一级所需经验   |

#### vip

| 字段名           | 数据类型 | 含义说明      |
| :--------------- | :------- | :------------ |
| vipType          | number   | VIP类型       |
| vipDueDate       | number   | VIP到期时间戳 |
| vipStatus        | number   | VIP状态       |
| label            | object   | VIP标签信息   |
| avatar_subscript | number   | 头像角标      |
| nickname_color   | string   | 昵称颜色      |

#### content

| 字段名   | 数据类型 | 含义说明     |
| :------- | :------- | :----------- |
| message  | string   | 评论内容     |
| members  | array    | @的用户列表  |
| emote    | object   | 表情信息     |
| jump_url | object   | 跳转链接     |
| max_line | number   | 最大显示行数 |

#### reply_control

| 字段名             | 数据类型 | 含义说明     |
| :----------------- | :------- | :----------- |
| max_line           | number   | 最大显示行数 |
| time_desc          | string   | 时间描述     |
| translation_switch | number   | 翻译开关     |
| support_share      | boolean  | 是否支持分享 |

#### folder

| 字段名     | 数据类型 | 含义说明   |
| :--------- | :------- | :--------- |
| has_folded | boolean  | 是否被折叠 |
| is_folded  | boolean  | 是否折叠   |
| rule       | string   | 折叠规则   |

### top

| 字段名 | 数据类型    | 含义说明   |
| :----- | :---------- | :--------- |
| admin  | null/object | 管理员置顶 |
| upper  | null/object | UP主置顶   |
| vote   | null/object | 投票置顶   |

### config

| 字段名       | 数据类型 | 含义说明     |
| :----------- | :------- | :----------- |
| showtopic    | number   | 显示话题     |
| show_up_flag | boolean  | 显示UP主标志 |
| read_only    | boolean  | 是否只读     |

### Control

| 字段名                    | 数据类型 | 含义说明         |
| :------------------------ | :------- | :--------------- |
| input_disable             | boolean  | 输入框是否禁用   |
| root_input_text           | string   | 根评论输入框提示 |
| child_input_text          | string   | 子评论输入框提示 |
| giveup_input_text         | string   | 放弃输入提示     |
| screenshot_icon_state     | number   | 截图图标状态     |
| upload_picture_icon_state | number   | 上传图片图标状态 |