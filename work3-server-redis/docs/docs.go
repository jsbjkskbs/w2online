// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/author/ping": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "token前面要添加Bearer",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "author"
                ],
                "summary": "token测试api",
                "responses": {}
            }
        },
        "/author/todolist/add": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "token前面要添加Bearer",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "author"
                ],
                "summary": "添加备忘录api",
                "parameters": [
                    {
                        "description": "标题,内容,截止日期[yyyy-mm-dd hh:mm:ss]",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/datastruct.TodolistBindJSONReceive"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/author/todolist/delete": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "token前面要添加Bearer;method(允许叠加,使用或运算):1[isdone],2[idlist],4[all]",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "author"
                ],
                "summary": "删除备忘录api",
                "parameters": [
                    {
                        "description": "是否完成,keyword不填,id数组,查找方法",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/datastruct.TodolistBindRedisCondition"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/author/todolist/modify": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "token前面要添加Bearer;method(允许叠加,使用或运算):2[idlist],4[all]",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "author"
                ],
                "summary": "更新备忘录api",
                "parameters": [
                    {
                        "description": "id数组,更新状态,查找方法",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/datastruct.TodolistBindRedisUpdate"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/author/todolist/search": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "token前面要添加Bearer;method(允许叠加,使用或运算):1[isdone],2[keyword],4[all]",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "author"
                ],
                "summary": "查找备忘录api",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "页码",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "description": "是否完成,关键字,idlist不填,查找方法",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/datastruct.TodolistBindRedisCondition"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/login": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "登录api",
                "parameters": [
                    {
                        "type": "string",
                        "description": "用户名",
                        "name": "username",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "密码",
                        "name": "password",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/register": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "注册api",
                "parameters": [
                    {
                        "type": "string",
                        "description": "用户名",
                        "name": "username",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "密码",
                        "name": "password",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "datastruct.TodolistBindJSONReceive": {
            "type": "object",
            "properties": {
                "deadline": {
                    "type": "string",
                    "example": "2077-01-01 01:01:01"
                },
                "text": {
                    "type": "string",
                    "example": "文本"
                },
                "title": {
                    "type": "string",
                    "example": "标题"
                }
            }
        },
        "datastruct.TodolistBindRedisCondition": {
            "type": "object",
            "properties": {
                "idlist": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "isdone": {
                    "type": "boolean",
                    "example": false
                },
                "keyword": {
                    "type": "string",
                    "example": "我超OP"
                },
                "method": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "datastruct.TodolistBindRedisUpdate": {
            "type": "object",
            "properties": {
                "idlist": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "isdone": {
                    "type": "boolean",
                    "example": false
                },
                "method": {
                    "type": "integer",
                    "example": 1
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.1-redis",
	Host:             "127.0.0.1:8080",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "test",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
