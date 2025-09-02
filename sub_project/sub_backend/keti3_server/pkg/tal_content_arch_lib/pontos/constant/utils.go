package constant

const (
	TraceIDCtxKey       = "x_trace_id" //大鱼网关trace
	RpcIdCtxKey         = "x_rpc_id"
	OTELTraceIdKey      = "otel_trace_id"
	OTELSpanIdKey       = "otel_span_id"
	TalIDCtxKey         = "p_tal_id"
	DeviceSNCtxKey      = "p_device_sn"
	DeviceVersionCtxKey = "p_device_version"
	GWTraceKey          = "Traceid" //大鱼网关
)

const (
	SKIPLOGCtxKey = "SKIPLOG"
)

const (
	SnHeaderKey      = "X-Tal-Sn"
	VersionHeaderKey = "X-Tal-Version"
)

// 压测标志
const (
	PressureTestKey   = "tal-request-type"
	PressureTestValue = "performance-testing"
	IsPressureTestKey = "UC_IS_Pressure"
)

// 登录Token
const (
	LoginConfigKey       = "client_login"
	LoginTokenKey        = "Authorization"
	ClientUserDataCtxKey = "client_user_data"
	TalId                = "tal_id"
)

const (
	FORMContentType = "application/x-www-form-urlencoded"
	JSONContentType = "application/json"
)

// 数据类型
const (
	RAW = iota
	Encrypt
)

// K8S相关ENV
const (
	K8SNamespaceEnvKey      = "CLOUD_NAMESPACE_NAME"
	K8SPodIPEnvKey          = "CLOUD_POD_IP"
	K8SNodeNameEnvKey       = "CLOUD_NODE_NAME"
	K8SNodeIPEnvKey         = "CLOUD_NODE_IP"
	K8SEnvKey               = "KUBERNETES_ENV"
	K8SClusterEnvKey        = "KUBERNETES_CLUSTER"
	K8SDeploymentNameEnvKey = "CLOUD_DEPLOYMENT_NAME"
)

const (
	PingPath = "/pontos/ping"
)

const (
	OtelDeviceSNAttr      = "device.sn"
	OtelDeviceVersionAttr = "device.version"
	OtelDeviceTalIDAttr   = "device.tal_id"
)

// Version for pontos
func Version() string {
	return "v1.1.19-bugfix2"
}
