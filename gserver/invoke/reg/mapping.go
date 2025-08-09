package reg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zusux/gokit/gserver/invoke/api"
	"github.com/zusux/gokit/gserver/zrpc"
	"github.com/zusux/gokit/zlog"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"reflect"
	"strings"
	"sync"
)

var mu sync.Mutex
var r = -1
var methodHandlers = make(map[uint32]func(context.Context, json.RawMessage) (interface{}, error))

func GetHandler(fullId uint32) (func(context.Context, json.RawMessage) (interface{}, error), bool) {
	if r == -1 {
		return nil, false
	}
	h, ok := methodHandlers[fullId]
	return h, ok
}
func AutoRegisterGRPCServiceMethods(descProto []byte, impls ...any) (map[uint32]*api.Endpoint, error) {
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
		serviceName := ExtractServiceFullNameFromImplByMatch(descProto, impl) // 比如 "user.v1.UserSrv"
		serviceID := GetServiceIDFromDescriptor(descProto, serviceName)       // 例如 0x1001
		endpoint, ok := ret[serviceID]
		if !ok {
			endpoint = &api.Endpoint{
				App:       "",
				ServiceId: serviceID,
				Host:      "",
				Auth: &api.Auth{
					AuthSkip: false,
					Secret:   "",
					TokenKey: "",
				},
				AllowOrigin: "*",
				Rate: &api.Rate{
					Limit: 0,
					Burst: 0,
				},
				Timeout:    3000,
				HttpTarget: make([]*api.Target, 0),
				GrpcTarget: make([]*api.Target, 0),
				Requests:   make(map[string]*api.Request),
			}
		}
		// 构建 user.v1.UserSrv.Ping → 入参类型
		fullMap := make(map[string]any)
		for method, param := range typeMap {
			key := fmt.Sprintf("%s.%s", serviceName, method)
			fullMap[key] = param
		}

		e := RegisterFromDescriptor(descProto, impl, serviceID, fullMap, endpoint)
		if e != nil {
			err = errors.Join(err, e)
		} else {
			ret[serviceID] = endpoint
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
func GetServiceIDFromDescriptor(desc []byte, serviceName string) uint32 {
	fds := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(desc, fds); err != nil {
		panic(fmt.Sprintf("unmarshal desc error: %v", err))
	}

	rfd, err := protodesc.NewFiles(fds)
	if err != nil {
		panic(fmt.Sprintf("desc load error: %v", err))
	}

	var foundID uint32
	rfd.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		for i := 0; i < fd.Services().Len(); i++ {
			svc := fd.Services().Get(i)
			if string(svc.FullName()) == serviceName {
				opts := svc.Options().(*descriptorpb.ServiceOptions)
				ext := proto.GetExtension(opts, zrpc.E_ServiceOptionId)
				id, ok := ext.(uint32)
				if ok {
					foundID = id
				}
				return false
			}
		}
		return true
	})

	if foundID == 0 {
		panic("service ID not found for: " + serviceName)
	}
	return foundID
}

// RegisterFromDescriptor 读取 descriptor 文件，根据服务映射注册方法处理器
func RegisterFromDescriptor(desc []byte, serviceImpl any, serviceID uint32, methodTypeMap map[string]any, endpoint *api.Endpoint) error {

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
				ext := proto.GetExtension(opts, zrpc.E_MethodOptionHttpApi)
				httpServerFullName := strings.Split(strings.Trim(serviceFullName, "/ "), ".")
				httpServer := httpServerFullName[0]
				httpMethod := ""
				httpRouter := ""
				outerRouter := strings.ToLower(httpServer) + "-" + strings.ToLower(methodName)
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

				methodKey := fmt.Sprintf("%s.%s", serviceFullName, methodName)
				reqPrototype, ok := methodTypeMap[methodKey]
				if !ok {
					zlog.Printf("warn: no request type for method %s\n", methodKey)
					continue
				}

				endpoint.Requests[outerRouter] = &api.Request{
					Router:  outerRouter,
					Methods: []string{httpMethod},
					Auth: &api.Auth{
						AuthSkip: false,
						Secret:   "",
						TokenKey: "",
					},
					AllowOrigin: "*",
					Timeout:     3000,
					MethodId:    methodID,
					Rate: &api.Rate{
						Limit: 0,
						Burst: 0,
					},
					Robin:      "simple",
					HttpTarget: make([]*api.Target, 0),
					GrpcTarget: make([]*api.Target, 0),
					Switch:     "http",
					Http: &api.Http{
						Protocol: "http",
						Method:   httpMethod,
						Path:     httpRouter,
						Timeout:  3000,
					},
					Grpc: &api.Grpc{
						Path:    methodKey,
						Timeout: 3000,
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
