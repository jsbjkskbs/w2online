namespace go relation

include "base.thrift"

struct RelationActionRequest {
    1: string access_token (api.header="Access-Token");
    2: string refresh_token (api.header="Refresh-Token");
    3: string to_user_id (api.go_tag='json:"to_user_id,required"');
    4: i64 action_type (api.go_tag='json:"action_type,required"');
}

struct RelationActionResponse {
    1: base.Status base;
}

struct FollowingListRequest {
    1: string user_id (api.query="user_id");
    2: i64 page_num (api.query="page_num");
    3: i64 page_size (api.query="page_size");
    4: string access_token (api.header="Access-Token");
    5: string refresh_token (api.header="Refresh-Token");
}

struct FollowingListResponseData {
    1: list<base.UserLite> items (api.go_tag='json:"items,required"');
    2: i64 total (api.go_tag='json:"total,required"');
}
struct FollowingListResponse {
    1: base.Status base;
    2: FollowingListResponseData data;
}

struct FollowerListRequest {
    1: string user_id (api.query="user_id");
    2: i64 page_num (api.query="page_num");
    3: i64 page_size (api.query="page_size");
    4: string access_token (api.header="Access-Token");
    5: string refresh_token (api.header="Refresh-Token");
}

struct FollowerListResponseData {
    1: list<base.UserLite> items (api.go_tag='json:"items,required"');
    2: i64 total (api.go_tag='json:"total,required"');
}
struct FollowerListResponse {
    1: base.Status base;
    2: FollowerListResponseData data;
}

struct FriendListRequest {
    1: i64 page_num (api.query="page_num");
    2: i64 page_size (api.query="page_size");
    3: string access_token (api.header="Access-Token");
    4: string refresh_token (api.header="Refresh-Token");
}

struct FriendListResponseData {
    1: list<base.UserLite> items (api.go_tag='json:"items,required"');
    2: i64 total (api.go_tag='json:"total,required"');
}
struct FriendListResponse {
    1: base.Status base;
    2: FriendListResponseData data;
}

service RelationAction {
    RelationActionResponse RelationActionMethod(1: RelationActionRequest request) (api.post="/relation/action/");
}

service FollowingList {
    FollowingListResponse FollowingListMethod(1: FollowingListRequest request) (api.get="/following/list/");
}

service FollowerList {
    FollowerListResponse FollowerListMethod(1: FollowerListRequest request) (api.get="/follower/list/");
}

service FriendList {
    FriendListResponse FriendListMethod(1: FriendListRequest request) (api.get="/friend/list/");
}