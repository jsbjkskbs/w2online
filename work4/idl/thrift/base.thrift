namespace go base

struct Status {
    1: i64 code (api.go_tag='json:"code,required"');
    2: string msg (api.go_tag='json:"msg,required"');
}

struct User {
    1: string uid;
    2: string username;
    3: string avatar_url(api.go_tag='json:"avatar_url,required"');
    4: string created_at;
    5: string updated_at;
    6: string deleted_at (api.go_tag='json:"deleted_at,required"');
}

struct UserLite {
    1: string uid;
    2: string username;
    3: string avatar_url (api.go_tag='json:"avatar_url,required"');
}

struct Video {
    1:  string id (api.go_tag='json:"id,required"');
    2:  string user_id (api.go_tag='json:"user_id,required"');
    3:  string video_url;
    4:  string cover_url;
    5:  string title;
    6:  string description;
    7:  i64 visit_count (api.go_tag='json:"visit_count,required"');
    8:  i64 like_count (api.go_tag='json:"like_count,required"');
    9:  i64 comment_count (api.go_tag='json:"comment_count,required"');
    10: string created_at;
    11: string updated_at;
    12: string deleted_at (api.go_tag='json:"deleted_at,required"');
}

struct Comment {
    1:  string id;
    2:  string user_id;
    3:  string video_id;
    4:  string parent_id;
    5:  i64 like_count (api.go_tag='json:"like_count,required"');
    6:  i64 child_count (api.go_tag='json:"child_count,required"');
    7:  string content;
    8:  string created_at;
    9:  string updated_at;
    10: string deleted_at (api.go_tag='json:"deleted_at,required"');
}