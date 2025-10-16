package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ProjectsTask/EasySwapBackend/src/service/svc"
	"github.com/ProjectsTask/EasySwapBackend/src/types/v1"
)

func GetAccountBalance(ctx context.Context, svcCtx *svc.ServerCtx, address string, chainId int) (*types.EtherScanBalanceResponse, error) {
	// 1. 构建基础 URL
	u, err := url.Parse(svcCtx.C.EtherScan.BaseUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid etherscan base url: %v", err)
	}

	// 2. 构建 query 参数
	params := url.Values{}
	params.Set("apiKey", svcCtx.C.EtherScan.ApiKey)
	params.Set("chainid", fmt.Sprintf("%d", chainId))
	params.Set("module", "account")
	params.Set("action", "balance")
	params.Set("address", address)
	params.Set("tag", "latest")

	u.RawQuery = params.Encode()

	// 3. 构建请求
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// 可选：设置 header（有些 API 需要）
	req.Header.Set("Accept", "application/json")

	// 4. 发送请求
	resp, err := svcCtx.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("etherscan request error: %v", err)
	}
	defer resp.Body.Close()

	// 5. 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// 6. 解析 JSON
	var result types.EtherScanBalanceResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// 7. 校验状态
	if result.Status != "1" {
		return nil, fmt.Errorf("etherscan API error: %s (%s)", result.Message, result.Result)
	}

	return &result, nil
}
