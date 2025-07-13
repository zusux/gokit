package isrv

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zusux/gokit/gserver/invoke/api"
	"github.com/zusux/gokit/gserver/invoke/reg"
	"github.com/zusux/gokit/zlog"
)

type InvokeService struct {
	api.UnimplementedInvokeRouterServer
}

func NewInvokeService() *InvokeService {
	return &InvokeService{}
}

func (s *InvokeService) Invoke(ctx context.Context, req *api.InvokeReq) (*api.InvokeResp, error) {
	zlog.Infof("Invoke: %v", req)
	fullID := req.ServiceId<<16 | req.MethodId

	handler, ok := reg.GetHandler(fullID)
	if !ok {
		return nil, fmt.Errorf("method not found: 0x%08X", fullID)
	}
	var tmp any
	if err := json.Unmarshal([]byte(req.PayloadJson), &tmp); err != nil {
		return nil, fmt.Errorf("PayloadJson must be a raw JSON string, got base64?")
	}
	// 转发给 handler，传入上下文与 JSON payload
	resp, err := handler(ctx, json.RawMessage(req.PayloadJson))
	if err != nil {
		return nil, err
	}

	// 将返回值序列化为 JSON
	respJson, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("response marshal error: %w", err)
	}
	return &api.InvokeResp{
		PayloadJson: string(respJson),
	}, nil
}
