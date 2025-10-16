package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ProjectsTask/EasySwapBackend/src/service/svc"
	"github.com/ProjectsTask/EasySwapBackend/src/types/v1"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func SearchCoins(ctx context.Context, svcCtx *svc.ServerCtx, keyword string) ([]types.CoinsSearchApiResponse, error) {
	// generate url and query param
	baseURL, err := url.Parse(svcCtx.C.Coingecko.BaseUrl + "search")
	if err != nil {
		return nil, fmt.Errorf("coingecko base URL parsing error: %v", err)
	}

	params := url.Values{}
	params.Add("query", keyword)

	baseURL.RawQuery = params.Encode()

	// build request
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("coingecko NewRequestWithContext error: %v", err)
	}

	req.Header.Set(svcCtx.C.Coingecko.ApiHeaderKey, svcCtx.C.Coingecko.ApiKey)

	// 发送请求
	resp, err := svcCtx.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("coingecko request error: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var result struct {
		Coins []types.CoinsSearchApiResponse `json:"coins"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return result.Coins, nil
}

func GetCoinsPrice(ctx context.Context, svcCtx *svc.ServerCtx, coinIds []string, currency string) (map[string]types.CoinSimplePriceResponse, error) {
	// generate url and query param
	baseURL, err := url.Parse(svcCtx.C.Coingecko.BaseUrl + "simple/price")
	if err != nil {
		return nil, fmt.Errorf("coingecko base URL parsing error: %v", err)
	}

	params := url.Values{}
	params.Add("ids", strings.Join(coinIds, ","))
	params.Add("vs_currencies", currency)

	baseURL.RawQuery = params.Encode()

	// build request
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("coingecko NewRequestWithContext error: %v", err)
	}

	req.Header.Set(svcCtx.C.Coingecko.ApiHeaderKey, svcCtx.C.Coingecko.ApiKey)

	// 发送请求
	resp, err := svcCtx.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("coingecko request error: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var result map[string]types.CoinSimplePriceResponse

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return result, nil
}
