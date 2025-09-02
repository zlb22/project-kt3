package constant

// 日志中x_tag字段枚举值
const (
	LogRequestTag  = "log_request"      //request
	LogResponseTag = "log_response"     //response
	LogJsonTag     = "log_json_tag"     //json marshal/unmarshal tag
	LogRedisTag    = "log_redis_tag"    //redis operation
	LogBizTag      = "log_biz_log"      //business logic
	LogHttpRequest = "log_http_service" //http service
	LogSQLTag      = "log_sql_tag"      //sql
)
