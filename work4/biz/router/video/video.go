// Code generated by hertz generator. DO NOT EDIT.

package video

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	video "work/biz/handler/video"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *server.Hertz) {

	root := r.Group("/", rootMw()...)
	{
		_video := root.Group("/video", _videoMw()...)
		{
			_feed := _video.Group("/feed", _feedMw()...)
			_feed.GET("/", append(_videofeedMw(), video.VideoFeed)...)
		}
		{
			_list := _video.Group("/list", _listMw()...)
			_list.GET("/", append(_videolistMw(), video.VideoList)...)
		}
		{
			_popular := _video.Group("/popular", _popularMw()...)
			_popular.GET("/", append(_videopopularMw(), video.VideoPopular)...)
		}
		{
			_publish := _video.Group("/publish", _publishMw()...)
			_publish.POST("/cancle", append(_videopublishcancleMw(), video.VideoPublishCancle)...)
			_publish.POST("/complete", append(_videopublishcompleteMw(), video.VideoPublishComplete)...)
			_publish.POST("/start", append(_videopublishstartMw(), video.VideoPublishStart)...)
			_publish.POST("/uploading", append(_videopublishuploadingMw(), video.VideoPublishUploading)...)
		}
		{
			_search := _video.Group("/search", _searchMw()...)
			_search.POST("/", append(_videosearchMw(), video.VideoSearch)...)
		}
		{
			_visit := _video.Group("/visit", _visitMw()...)
			_visit.GET("/:id", append(_videovisitMw(), video.VideoVisit)...)
		}
	}
}
