module github.com/yuyayang02/BlogLite-api/internal/categories

go 1.22.0



replace (
	github.com/yuyayang02/BlogLite-api/internal => ../common
	github.com/yuyayang02/BlogLite-api/internal/server/httpresponse => ../common/server/httpresponse
)