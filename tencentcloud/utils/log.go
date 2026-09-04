package utils

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// LogRequest Format tencentcloud-sdk request into JSON and log
func LogRequest(ctx context.Context, action string, req any) {
	if !plugin.Logger(ctx).IsDebug() {
		return
	}
	bs, err := json.Marshal(req)
	if err != nil {
		plugin.Logger(ctx).Debug(action, "request", fmt.Sprintf("%+v", req))
		return
	}
	plugin.Logger(ctx).Debug(action, "request", string(bs))
}
