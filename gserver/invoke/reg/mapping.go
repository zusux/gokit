package reg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/samber/lo"
	"github.com/zusux/gokit/gserver/invoke/api"
	"github.com/zusux/gokit/gserver/zrpc"
	"github.com/zusux/gokit/zlog"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

var mu sync.Mutex
var r = -1
var methodHandlers = make(map[uint32]func(context.Context, json.RawMessage) (interface{}, error))

type ServiceOptions struct {
	ServiceId           uint32
	GatewayEndpointType uint32 // 0: ip端口, 2:域名  3:服务发现
	GatewayServiceAuth  uint32 // 0未设置, 1 跳过鉴权  2需要鉴权
	GatewayServiceLimit uint32 // 限流数量
}

type MethodOptions struct {
	MethodOptionId         uint32 // 方法id
	GatewayMethodAuth      uint32 // 0未设置, 1 跳过鉴权  2需要鉴权
	GatewayMethodRateLimit uint32 //限流数量
	GatewayMethodSwitch    uint32 //代理开关 0 http  1 grpc
}

func GetHandler(fullId uint32) (func(context.Context, json.RawMessage) (interface{}, error), bool) {
	if r == -1 {
		return nil, false
	}
	h, ok := methodHandlers[fullId]
	return h, ok
}
func AutoRegisterGRPCServiceMethods(descProto []byte, app string, auth *api.Auth, httpTarget, grpcTarget []*api.Target, timeOut int64, impls ...any) (map[uint32]*api.Endpoint, error) {
	defer func() {
		r = 0
	}()
	var err error
	ret := make(map[uint32]*api.Endpoint, 0)
	for _, impl := range impls {
		// step1: 提取实际实现的方法名称与参数类型
		typeMap := ExtractMethodParamTypes(impl)

		// step2: 读取 .desc 并匹配方法名
		// 假设服务名为 user.v1.UserSrv
		serviceName := ExtractServiceFullNameFromImplByMatch(descProto, impl)     // 比如 "user.v1.UserSrv"
		serviceOptions := GetServiceOptionsFromDescriptor(descProto, serviceName) // 例如 0x1001
		endpoint, ok := ret[serviceOptions.ServiceId]
		if !ok {
			var tmpAuth *api.Auth
			if auth != nil {
				tmpAuth = auth
				switch serviceOptions.GatewayServiceAuth {
				case 1:
					tmpAuth.AuthSkip = true
				case 2:
					tmpAuth.AuthSkip = false
				}
			} else {
				switch serviceOptions.GatewayServiceAuth {
				case 1:
					tmpAuth = &api.Auth{}
					tmpAuth.AuthSkip = true
				case 2:
					tmpAuth = &api.Auth{}
					tmpAuth.AuthSkip = false
				default:
					tmpAuth = &api.Auth{}
				}
			}
			endpoint = &api.Endpoint{
				App:         app,
				ServiceId:   serviceOptions.ServiceId,
				Host:        "",
				Auth:        tmpAuth,
				AllowOrigin: "*",
				Rate: &api.Rate{
					Limit: lo.If(serviceOptions.GatewayServiceLimit > 0, float32(serviceOptions.GatewayServiceLimit)).Else(100000),
					Burst: 100000,
				},
				Timeout:             timeOut,
				HttpTarget:          httpTarget,
				GrpcTarget:          grpcTarget,
				Requests:            make(map[string]*api.Request),
				GatewayEndpointType: serviceOptions.GatewayEndpointType,
			}
		}
		// 构建 user.v1.UserSrv.Ping → 入参类型
		fullMap := make(map[string]any)
		for method, param := range typeMap {
			key := fmt.Sprintf("%s.%s", serviceName, method)
			fullMap[key] = param
		}

		e := RegisterFromDescriptor(descProto, impl, serviceOptions.ServiceId, auth, timeOut, fullMap, endpoint)
		if e != nil {
			err = errors.Join(err, e)
		} else {
			ret[serviceOptions.ServiceId] = endpoint
		}
	}
	return ret, err
}

func ExtractMethodParamTypes(impl any) map[string]any {
	methodMap := map[string]any{}
	typ := reflect.TypeOf(impl)
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		if m.Type.NumIn() != 3 || m.Type.NumOut() != 2 {
			continue
		}
		argType := m.Type.In(2)
		if argType.Kind() != reflect.Ptr {
			continue
		}
		methodMap[m.Name] = reflect.New(argType.Elem()).Interface()
	}
	return methodMap
}

func ExtractServiceFullNameFromImplByMatch(desc []byte, impl any) string {
	implType := reflect.TypeOf(impl).Elem().Name()           // e.g. UserSrvService
	expectService := strings.TrimSuffix(implType, "Service") // → UserSrv

	fds := &descriptorpb.FileDescriptorSet{}
	_ = proto.Unmarshal(desc, fds)

	for _, file := range fds.File {
		for _, svc := range file.Service {
			if svc.GetName() == expectService {
				return fmt.Sprintf("%s.%s", file.GetPackage(), svc.GetName())
			}
		}
	}

	panic("no matching service name found for: " + expectService)
}

func GetServiceOptionsFromDescriptor(desc []byte, serviceName string) ServiceOptions {
	fds := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(desc, fds); err != nil {
		panic(fmt.Sprintf("unmarshal desc error: %v", err))
	}

	rfd, err := protodesc.NewFiles(fds)
	if err != nil {
		panic(fmt.Sprintf("desc load error: %v", err))
	}

	var serviceOptions ServiceOptions
	rfd.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		for i := 0; i < fd.Services().Len(); i++ {
			svc := fd.Services().Get(i)
			if string(svc.FullName()) == serviceName {
				opts := svc.Options().(*descriptorpb.ServiceOptions)
				ext := proto.GetExtension(opts, zrpc.E_ServiceOptionId)
				if id, ok := ext.(uint32); ok {
					serviceOptions.ServiceId = id
				}

				extType := proto.GetExtension(opts, zrpc.E_GatewayEndpointType)
				if id, ok := extType.(uint32); ok {
					serviceOptions.GatewayEndpointType = id
				}
				extAuth := proto.GetExtension(opts, zrpc.E_GatewayServiceAuth)
				if id, ok := extAuth.(uint32); ok {
					serviceOptions.GatewayServiceAuth = id
				}
				extLimit := proto.GetExtension(opts, zrpc.E_GatewayServiceLimit)
				if id, ok := extLimit.(uint32); ok {
					serviceOptions.GatewayServiceLimit = id
				}
				return false
			}
		}
		return true
	})
	return serviceOptions
}

// RegisterFromDescriptor 读取 descriptor 文件，根据服务映射注册方法处理器
func RegisterFromDescriptor(desc []byte, serviceImpl any, serviceID uint32, auth *api.Auth, timeout int64, methodTypeMap map[string]any, endpoint *api.Endpoint) error {

	fds := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(desc, fds); err != nil {
		return fmt.Errorf("unmarshal desc error: %w", err)
	}

	rfd, err := protodesc.NewFiles(fds)
	if err != nil {
		return fmt.Errorf("convert to protoregistry.Files error: %w", err)
	}

	rfd.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		for i := 0; i < fd.Services().Len(); i++ {
			service := fd.Services().Get(i)
			serviceFullName := string(service.FullName()) // e.g. "user.v1.Aggregator"

			for j := 0; j < service.Methods().Len(); j++ {
				method := service.Methods().Get(j)
				methodName := string(method.Name())
				opts := method.Options().(*descriptorpb.MethodOptions)
				if opts == nil {
					continue
				}
				methodIDPtr := proto.GetExtension(opts, zrpc.E_MethodOptionId)
				methodID, ok := methodIDPtr.(uint32)
				if !ok {
					zlog.Printf("warn: no method_id for %s.%s\n", serviceFullName, methodName)
					continue
				}

				methodAuthPtr := proto.GetExtension(opts, zrpc.E_GatewayMethodAuth)
				methodAuth := methodAuthPtr.(uint32)
				var tmpAuth *api.Auth
				if auth != nil {
					switch methodAuth {
					case 1:
						tmpAuth = auth
						tmpAuth.AuthSkip = true
					case 2:
						tmpAuth = auth
						tmpAuth.AuthSkip = false
					}
				} else {
					switch methodAuth {
					case 1:
						tmpAuth = &api.Auth{}
						tmpAuth.AuthSkip = true
					case 2:
						tmpAuth = &api.Auth{}
						tmpAuth.AuthSkip = false
					}
				}

				methodRatePtr := proto.GetExtension(opts, zrpc.E_GatewayMethodRateLimit)
				methodRate := methodRatePtr.(uint32)

				methodSwitchPtr := proto.GetExtension(opts, zrpc.E_GatewayMethodSwitch)
				methodSwitch := methodSwitchPtr.(uint32)

				ext := proto.GetExtension(opts, annotations.E_Http)
				httpServerFullName := strings.Split(strings.Trim(serviceFullName, "/ "), ".")
				httpServer := httpServerFullName[0]
				httpMethod := ""
				httpRouter := ""
				if httpRule, ok := ext.(*annotations.HttpRule); ok {
					switch pattern := httpRule.Pattern.(type) {
					case *annotations.HttpRule_Get:
						httpMethod = "GET"
						httpRouter = pattern.Get
					case *annotations.HttpRule_Post:
						httpMethod = "POST"
						httpRouter = pattern.Post
					case *annotations.HttpRule_Put:
						httpMethod = "Put"
						httpRouter = pattern.Put
					case *annotations.HttpRule_Delete:
						httpMethod = "Delete"
						httpRouter = pattern.Delete
					case *annotations.HttpRule_Patch:
						httpMethod = "Patch"
						httpRouter = pattern.Patch
					}
				}
				var outerRouter string
				if strings.HasPrefix(httpRouter, "/") {
					outerRouter = "/" + strings.ToLower(httpServer) + httpRouter
				} else {
					outerRouter = "/" + strings.ToLower(httpServer) + "/" + httpRouter
				}

				methodKey := fmt.Sprintf("%s.%s", serviceFullName, methodName)
				reqPrototype, ok := methodTypeMap[methodKey]
				if !ok {
					zlog.Printf("warn: no request type for method %s\n", methodKey)
					continue
				}

				endpoint.Requests[outerRouter] = &api.Request{
					Router:      outerRouter,
					Methods:     []string{httpMethod},
					Auth:        tmpAuth,
					AllowOrigin: "*",
					Timeout:     timeout,
					MethodId:    methodID,
					Rate: &api.Rate{
						Limit: lo.If(methodRate > 0, float32(methodRate)).Else(100000),
						Burst: 100000,
					},
					Robin:      "simple",
					HttpTarget: make([]*api.Target, 0),
					GrpcTarget: make([]*api.Target, 0),
					Switch:     lo.If(methodSwitch == 0, "http").Else("grpc"),
					Http: &api.Http{
						Protocol: "http",
						Method:   httpMethod,
						Path:     httpRouter,
						Timeout:  timeout,
					},
					Grpc: &api.Grpc{
						Path:    methodKey,
						Timeout: timeout,
					},
				}
				fullID := serviceID<<16 | methodID
				RegisterAutoMethod(fullID, serviceImpl, methodName, reqPrototype)
				zlog.Infof("✅ Registered %s as 0x%08X\n", methodKey, fullID)
			}
		}
		return true
	})

	return nil
}

// RegisterAutoMethod 使用反射注册 gRPC 方法，自动将 JSON 解码为 proto message 并调用服务方法
func RegisterAutoMethod(methodID uint32, impl any, methodName string, reqPrototype any) {
	implValue := reflect.ValueOf(impl)
	method := implValue.MethodByName(methodName)
	if !method.IsValid() {
		panic(fmt.Errorf("method %s not found on %T", methodName, impl))
	}

	handler := func(ctx context.Context, raw json.RawMessage) (interface{}, error) {
		// 实例化参数类型
		reqVal := reflect.New(reflect.TypeOf(reqPrototype).Elem()).Interface()
		if err := json.Unmarshal(raw, reqVal); err != nil {
			return nil, fmt.Errorf("json decode error: %w", err)
		}
		// 调用服务方法
		results := method.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(reqVal),
		})

		if len(results) != 2 {
			return nil, fmt.Errorf("method %s must return (resp, error)", methodName)
		}

		if errVal := results[1].Interface(); errVal != nil {
			return nil, errVal.(error)
		}

		return results[0].Interface(), nil
	}

	// 注册
	mu.Lock()
	methodHandlers[methodID] = handler
	mu.Unlock()
}

// 使用方式：
// RegisterFromDescriptor("all.desc", aggregatorSrv, 0x1001, map[string]any{
//   "user.v1.Aggregator.Ping": &v1.PingRequest{},
//   "user.v1.Aggregator.Summary": &v1.SummaryRequest{},
// })
