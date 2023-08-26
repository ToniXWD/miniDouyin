namespace go miniDouyin.api

// /douyin/feed/ - 视频流接口
struct FeedRequest {
    1: optional i64 latest_time, // 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
    2: optional string token,    // 可选参数，登录用户设置
}

struct FeedResponse {
    1: required i32 status_code,        // 状态码，0-成功，其他值-失败
    2: optional string status_msg,      // 返回状态描述
    3: required list<Video> video_list, // 视频列表
    4: optional i64 next_time,          // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
}

struct Video {
    1: required i64 id,             // 视频唯一标识
    2: required User author,        // 视频作者信息
    3: required string play_url,    // 视频播放地址
    4: required string cover_url,   // 视频封面地址
    5: required i64 favorite_count, // 视频的点赞总数
    6: required i64 comment_count,  // 视频的评论总数
    7: required bool is_favorite,   // true-已点赞，false-未点赞
    8: required string title,       // 视频标题
}

// /douyin/user/register/ - 用户注册接口
struct UserRegisterRequest {
    1: required string username, // 注册用户名，最长32个字符
    2: required string password, // 密码，最长32个字符
}

struct UserRegisterResponse {
    1: required i32 status_code,   // 状态码，0-成功，其他值-失败
    2: optional string status_msg, // 返回状态描述
    3: required i64 user_id,       // 用户id
    4: required string token,      // 用户鉴权token
}

// /douyin/user/login/ - 用户登录接口
struct UserLoginRequest {
    1: required string username, // 登录用户名
    2: required string password, // 登录密码
}

struct UserLoginResponse {
    1: required i32 status_code,   // 状态码，0-成功，其他值-失败
    2: optional string status_msg, // 返回状态描述
    3: required i64 user_id,       // 用户id
    4: required string token,      // 用户鉴权token
}

struct User {
    1: required i64 id,                  // 用户id
    2: required string name,             // 用户名称
    3: optional i64 follow_count,        // 关注总数
    4: optional i64 follower_count,      // 粉丝总数
    5: required bool is_follow,          // true-已关注，false-未关注
    6: optional string avatar,           // 用户头像
    7: optional string background_image, // 用户个人页顶部大图
    8: optional string signature,        // 个人简介
    9: optional i64 total_favorited,     // 获赞数量
    10: optional i64 work_count,         // 作品数量
    11: optional i64 favorite_count,     // 点赞数量
}

///douyin/user/ - 用户信息
// GET
struct UserRequest {
    1: required i64 user_id,  // 用户id
    2: required string token, // 用户鉴权token
}

struct UserResponse {
    1: required i32 status_code,   // 状态码，0-成功，其他值-失败
    2: optional string status_msg, // 返回状态描述
    3: required User user,         // 用户信息
}

// /douyin/publish/action/ - 视频投稿
struct PublishActionRequest {
    1: required string token, // 用户鉴权token
    2: required binary data,  // 视频数据
    3: required string title, // 视频标题
}

struct PublishActionResponse {
    1: required i32 status_code,   // 状态码，0-成功，其他值-失败
    2: optional string status_msg, // 返回状态描述
}

// /douyin/publish/list/ - 发布列表
struct PublishListRequest {
    1: required i64 user_id,  // 用户id
    2: required string token, // 用户鉴权token
}

struct PublishListResponse {
    1: required i32 status_code,        // 状态码，0-成功，其他值-失败
    2: optional string status_msg,      // 返回状态描述
    3: required list<Video> video_list, // 用户发布的视频列表
}

// /douyin/favorite/action/ - 赞操作
struct FavoriteActionRequest {
    1: required string token,    // 用户鉴权token
    2: required i64 video_id,    // 视频id
    3: required i64 action_type, // 1-点赞，2-取消点赞
}

struct FavoriteActionResponse {
    1: required i32 status_code,   // 状态码，0-成功，其他值-失败
    2: optional string status_msg, // 返回状态描述
}

// /douyin/favorite/list/ - 喜欢列表
struct FavoriteListRequest {
    1: required i64 user_id,  // 用户id
    2: required string token, // 用户鉴权token
}

struct FavoriteListResponse {
    1: required i32 status_code,        // 状态码，0-成功，其他值-失败
    2: optional string status_msg,      // 返回状态描述
    3: required list<Video> video_list, // 用户点赞视频列表
}

// /douyin/comment/action/ - 评论操作
struct CommentActionRequest {
    1: required string token,        // 用户鉴权token
    2: required i64 video_id,        // 视频id
    3: required i64 action_type,     // 1-发布评论，2-删除评论
    4: optional string comment_text, // 用户填写的评论内容，在action_type=1的时候使用
    5: optional i64 comment_id,      // 要删除的评论id，在action_type=2的时候使用
}

struct CommentActionResponse {
    1: required i32 status_code,   // 状态码，0-成功，其他值-失败
    2: optional string status_msg, // 返回状态描述
    3: optional Comment comment,   // 评论成功返回评论内容，不需要重新拉取整个列表
}

struct Comment {
    1: required i64 id,             // 视频评论id
    2: required User user,          // 评论用户信息
    3: required string content,     // 评论内容
    4: required string create_date, // 评论发布日期，格式 mm-dd
}

// /douyin/comment/list/ - 视频评论列表
struct CommentListRequest {
    1: required string token, // 用户鉴权token
    2: required i64 video_id, // 视频id
}

struct CommentListResponse {
    1: required i32 status_code,            // 状态码，0-成功，其他值-失败
    2: optional string status_msg,          // 返回状态描述
    3: required list<Comment> comment_list, // 评论列表
}

// /douyin/relation/action/ - 关系操作
struct RelationActionRequest {
    1: required string token,    // 用户鉴权token
    2: required i64 to_user_id,  // 对方用户id
    3: required i64 action_type, // 1-关注，2-取消关注
}

struct RelationActionResponse {
    1: required i32 status_code,   // 状态码，0-成功，其他值-失败
    2: optional string status_msg, // 返回状态描述
}

// /douyin/relation/follow/list/ - 用户关注列表
struct RelationFollowListRequest {
    1: required i64 user_id,  // 用户id
    2: required string token, // 用户鉴权token
}

struct RelationFollowListResponse {
    1: required i32 status_code,      // 状态码，0-成功，其他值-失败
    2: optional string status_msg,    // 返回状态描述
    3: required list<User> user_list, // 用户信息列表
}

// /douyin/relation/follower/list/ - 用户粉丝列表
struct RelationFollowerListRequest {
    1: required i64 user_id,  // 用户id
    2: required string token, // 用户鉴权token
}

struct RelationFollowerListResponse {
    1: required i32 status_code,      // 状态码，0-成功，其他值-失败
    2: optional string status_msg,    // 返回状态描述
    3: required list<User> user_list, // 用户列表
}

// /douyin/relation/friend/list/ - 用户好友列表
struct RelationFriendListRequest {
    1: required i64 user_id,  // 用户id
    2: required string token, // 用户鉴权token
}

struct RelationFriendListResponse {
    1: required i32 status_code,      // 状态码，0-成功，其他值-失败
    2: optional string status_msg,    // 返回状态描述
    3: required list<User> user_list, // 用户列表
}

// /douyin/struct/chat/ - 聊天记录
struct ChatRecordRequest {
    1: required string token,     // 用户鉴权token
    2: required i64 to_user_id,   // 对方用户id
    3: required i64 pre_msg_time, // 上次最新消息的时间（新增字段-apk更新中）
}

struct ChatRecordResponse {
    1: required i32 status_code,           // 状态码，0-成功，其他值-失败
    2: optional string status_msg,         // 返回状态描述
    3: required list<Message> message_list, // 消息列表
}

struct Message {
    1: required i64 id,             // 消息id
    2: required i64 to_user_id,     // 该消息接收者的id
    3: required i64 from_user_id,   // 该消息发送者的id
    4: required string content,     // 消息内容
    5: optional i64 create_time, // 消息创建时间
}

struct SendMsgRequest {
    1: required string token,    // 用户鉴权token
    2: required i64 to_user_id,  // 对方用户id
    3: required i64 action_type, // 1-发送消息
    4: required string content,  // 消息内容
}

struct SendMsgResponse {
    1: required i32 status_code,
    2: optional string status_msg,  //返回状态描述
}

service miniDouyin {
    // 视频流接口
    FeedResponse Feed(1: FeedRequest request) (api.get = "/douyin/feed/"),
    // 用户注册接口
    UserRegisterResponse Register(1: UserRegisterRequest request) (api.post = "/douyin/user/register/"),
    // 用户登录接口
    UserLoginResponse Login(1: UserLoginRequest request) (api.post = "/douyin/user/login/"),
    // 用户信息
    UserResponse GetUserInfo(1: UserRequest request) (api.get = "/douyin/user/"),
    // 视频投稿
    PublishActionResponse VideoPublishAction(1: PublishActionRequest request) (api.post = "/douyin/publish/action/"),
    // 发布列表
    PublishListResponse PublishList(1: PublishListRequest request) (api.get = "/douyin/publish/list/"),
    // 赞操作
    FavoriteActionResponse FavoriteAction(1: FavoriteActionResponse request) (api.post = "/douyin/favorite/action/"),
    // 喜欢列表
    FavoriteListResponse FavoriteList(1: FavoriteListRequest request) (api.get = "/douyin/favorite/list/"),
    // 评论操作
    CommentActionResponse CommentAction(1: CommentActionResponse request) (api.post = "/douyin/comment/action/"),
    // 视频评论列表
    CommentListResponse CommentList(1: CommentListRequest request) (api.get = "/douyin/comment/list/"),
    // 关系操作
    RelationActionResponse RelationAction(1: RelationActionRequest request) (api.post = "/douyin/relation/action/"),
    // 用户关注列表
    RelationFollowListResponse FollowList(1: RelationFollowerListRequest request) (api.get = "/douyin/relation/follow/list/"),
    // 用户粉丝列表
    RelationFollowerListResponse FollowerList(1: RelationFollowerListRequest request) (api.get = "/douyin/relation/follower/list/"),
    // 用户好友列表
    RelationFriendListResponse FriendList(1: RelationFriendListRequest request) (api.get = "/douyin/relation/friend/list/"),
    // 聊天记录
    ChatRecordResponse ChatRec(1: ChatRecordResponse request) (api.get = "/douyin/message/chat/"),
    // 消息操作
    SendMsgRequest SendMsg(1: SendMsgRequest request) (api.post = "/douyin/message/action/"),
}