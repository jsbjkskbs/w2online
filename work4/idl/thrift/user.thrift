namespace go user

include "base.thrift"

struct UserRegisterRequest {
    1: string username (api.go_tag='json:"username,required"');
    2: string password (api.go_tag='json:"password,required" vd:"len($)>5 && len($)<18; msg:\"invalid password format\""');
}

struct UserRegisterResponse {
    1: base.Status base;
    2: string access_token;
    3: string refresh_token;
}

struct UserLoginRequest {
    1: string username (api.go_tag='json:"username,required"');
    2: string password (api.go_tag='json:"password,required" vd:"len($)>5 && len($)<18; msg:\"invalid password format\""');
    3: string code;
}

struct UserLoginResponse {
    1: base.Status base;
    2: base.User data;
    3: string access_token;
    4: string refresh_token;
}

struct UserInfoRequest {
    1: string user_id (api.query="user_id");
    2: string token (api.query="token");
    3: string access_token (api.header="Access-Token");
    4: string refresh_token (api.header="Refresh-Token");
}

struct UserInfoResponse {
    1: base.Status base;
    2: base.User data;
}

struct UserAvatarUploadRequest {
    1: string access_token (api.header="Access-Token");
    2: string refresh_token (api.header="Refresh-Token");
    3: binary data (api.body="data");
}


struct UserAvatarUploadResponseData {
    1: string id;
    2: string username;
    3: string password;
    4: string avatar_url;
    5: string created_at;
    6: string updated_at;
    7: string deleted_at;
}
struct UserAvatarUploadResponse {
    1: base.Status base;
    2: UserAvatarUploadResponseData data;
}

struct AuthMfaQrcodeRequest { 
    1: string access_token (api.header="Access-Token");
    2: string refresh_token (api.header="Refresh-Token");
}

struct Qrcode {
    1: string secret;
    2: string qrcode;
}
struct AuthMfaQrcodeResponse {
    1: base.Status base;
    2: Qrcode data;
}

struct AuthMfaBindRequest {
    1: string access_token (api.header="Access-Token");
    2: string refresh_token (api.header="Refresh-Token");
    3: string code (api.body="code");
    4: string secret (api.body="secret");
}

struct AuthMfaBindResponse {
    1: base.Status base;
}

service UserRegister {
    UserRegisterResponse UserRegisterMethod(1: UserRegisterRequest request) (api.post="/user/register/");
}

service UserLogin {
    UserLoginResponse UserLoginMethod(1: UserLoginRequest request) (api.post="/user/login/");
}

service UserInfo {
    UserInfoResponse UserInfoMethod(1: UserInfoRequest request) (api.get="/user/info/");
}

service UserAvatarUpload {
    UserAvatarUploadRequest UserAvatarUploadMethod(1: UserAvatarUploadRequest request) (api.put="/user/avatar/upload/");
}

service AuthMfaQrcode {
    AuthMfaQrcodeResponse AuthMfaQrcodeMethod(1: AuthMfaQrcodeRequest request) (api.get="/auth/mfa/qrcode/");
}

service AuthMfaBind {
    AuthMfaBindResponse AuthMfaBindMethod(1: AuthMfaBindRequest request) (api.post="/auth/mfa/bind/");
}