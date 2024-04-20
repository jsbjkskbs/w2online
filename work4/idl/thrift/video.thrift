namespace go video

include "base.thrift"

struct VideoFeedRequest {
    1: string latest_time (api.query="latest_time");
}

struct VideoFeedResponseData {
    list<base.Video> items;
}
struct VideoFeedResponse {
    1: base.Status base;
    2: VideoFeedResponseData data;
}

struct VideoPublishStartRequest {
    1: string access_token (api.header="Access-Token");
    2: string refresh_token (api.header="Refresh-Token");
    3: string title (api.go_tag='json:"title" vd:"len($)>0; msg:\"title can not be empty\""');
    4: string description (api.body="description");
    5: i64 chunk_total_number (api.go_tag='json:"chunk_total_number" vd:"$>0; msg:\"chunk number can not be lower than zero\""');
}

struct VideoPublishStartResponse {
    1: base.Status base;
    2: string uuid;
}

struct VideoPublishUploadingRequest {
    1: string access_token (api.header="Access-Token");
    2: string refresh_token (api.header="Refresh-Token");
    3: string uuid (api.go_tag='json:"uuid,required"');
    4: binary data (api.go_tag='json:"data"');
    5: string md5 (api.go_tag='json:"md5,required"');
    6: bool is_m3u8 (api.go_tag='json:"is_m3u8,required"');
    7: string filename (api.go_tag='json:"filename,required"');
    8: i64 chunk_number (api.go_tag='json:"chunk_number,required"');
}

struct VideoPublishUploadingResponse {
    1: base.Status base;
}

struct VideoPublishCompleteRequest {
    1: string access_token (api.header="Access-Token");
    2: string refresh_token (api.header="Refresh-Token");
    3: string uuid (api.go_tag='json:"uuid,required"');
}

struct VideoPublishCompleteResponse {
    1: base.Status base;
}

struct VideoPublishCancleRequest {
    1: string access_token (api.header="Access-Token");
    2: string refresh_token (api.header="Refresh-Token");
    3: string uuid (api.go_tag='json:"uuid,required"');
}

struct VideoPublishCancleResponse {
    1: base.Status base;
}

struct VideoListRequest {
    1: string user_id (api.query="user_id");
    2: i64 page_num (api.query="page_num");
    3: i64 page_size (api.query="page_size");
    4: string access_token (api.header="Access-Token");
    5: string refresh_token (api.header="Refresh-Token");
}

struct VideoListResponseData {
    1: list<base.Video> data;
    2: i64 total;
}
struct VideoListResponse {
    1: base.Status base;
    2: VideoListResponseData data;
}

struct VideoPopularRequest {
    1: i64 page_num (api.query="page_num");
    2: i64 page_size (api.query="page_size");
    3: string access_token (api.header="Access-Token");
    4: string refresh_token (api.header="Refresh-Token");
}

struct VideoPopularResponseData {
    list<base.Video> items;
}
struct VideoPopularResponse {
    1: base.Status base;
    2: VideoPopularResponseData data;
}

struct VideoSearchRequest {
    1: string access_token (api.header="Access-Token");
    2: string refresh_token (api.header="Refresh-Token");
    3: string keywords (api.go_tag='json:"keywords"');
    4: i64 page_num (api.go_tag='json:"page_num,required"');
    5: i64 page_size (api.go_tag='json:"page_size,required"');
    6: i64 from_date (api.body="from_date");
    7: i64 to_date (api.body="to_date");
    8: string username (api.body="username");
}

struct VideoSearchResponseData {
    1: list<base.Video> items;
    2: i64 total;
}
struct VideoSearchResponse {
    1: base.Status base;
    2: VideoSearchResponseData data;
}

struct VideoVisitRequest {
}

struct VideoVisitResponse {
    1: base.Status base;
    2: base.Video item;
}

service VideoFeed {
    VideoFeedResponse VideoFeedMethod(1: VideoFeedRequest request) (api.get="/video/feed/");
}

service VideoPublishStart {
    VideoPublishStartResponse VideoPublishStartMethod(1: VideoPublishStartRequest request) (api.post="/video/publish/start/");
}

service VideoPublishUploading {
    VideoPublishUploadingResponse VideoPublishUploadingMethod(1: VideoPublishUploadingRequest request) (api.post="/video/publish/uploading/");
}

service VideoPublishComplete {
    VideoPublishCompleteResponse VideoPublishCompleteMethod(1: VideoPublishCompleteRequest request) (api.post="/video/publish/complete/");
}

service VideoPublishCancle {
    VideoPublishCancleResponse VideoPublishCancleMethod(1: VideoPublishCancleRequest request) (api.post="/video/publish/cancle/");
}

service VideoList {
    VideoListResponse VideoListMethod(1: VideoListRequest request) (api.get="/video/list/");
}

service VideoPopular {
    VideoPopularResponse VideoPopularMethod(1: VideoPopularRequest request) (api.get="/video/popular/");
}

service VideoSearch {
    VideoSearchResponse VideoSearchMethod(1: VideoSearchRequest request) (api.get="/video/search/");
}

service VideoVisit {
    VideoVisitResponse VideoVisitMethod(1: VideoVisitRequest request) (api.get="/video/visit/:id");
}