package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sealway-strava/domain/strava"
	"strconv"
	"strings"
	"time"
)

type StravaConfig struct {
	ClientId string
	SecretId string
}

type StravaQuota struct {
	Limit15min int       `json:"limit_15_min"`
	LimitDay   int       `json:"limit_day"`
	Usage15min int       `json:"usage_15_min"`
	UsageDay   int       `json:"usage_day"`
	LastUpdate time.Time `json:"last_update"`
}

type StravaClient struct {
	*strava.APIClient
	*StravaConfig

	Quota            *StravaQuota
	QueueUpdateQuota chan StravaQuota
}

func MakeStravaClient(config *StravaConfig) *StravaClient {
	apiClient := strava.NewAPIClient(strava.NewConfiguration())

	client := &StravaClient{
		StravaConfig:     config,
		APIClient:        apiClient,
		Quota:            &StravaQuota{},
		QueueUpdateQuota: make(chan StravaQuota),
	}

	go client.RunQuotaUpdate()

	return client
}

func (client *StravaClient) Close() {
	close(client.QueueUpdateQuota)
}

func (client *StravaClient) RefreshToken(ctx context.Context, refreshToken string) (*string, error) {
	url := "https://www.strava.com/api/v3/oauth/token"

	values := map[string]string{
		"client_id":     client.ClientId,
		"client_secret": client.SecretId,
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}
	data, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.WithContext(ctx)
	if err != nil {
		return nil, err
	}

	httpClient := http.DefaultClient
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	accessToken := res["access_token"].(string)

	return &accessToken, nil
}

func (client *StravaClient) CheckQuota() error {
	quota := client.Quota
	diff := time.Since(quota.LastUpdate)

	if quota.Limit15min-quota.Usage15min <= 0 && diff.Minutes() < 15 {
		return fmt.Errorf("CheckQuota - 15min quota [%d/%d]", quota.Usage15min, quota.Limit15min)
	}

	if quota.LimitDay-quota.UsageDay <= 0 && diff.Hours() < 24 {
		return fmt.Errorf("CheckQuota - daily quota [%d/%d]", quota.UsageDay, quota.LimitDay)
	}

	return nil
}

func (client *StravaClient) UpdateQuota(response *http.Response) {
	if response == nil {
		return
	}

	limitHeader := "X-Ratelimit-Limit"
	usageHeader := "X-Ratelimit-Usage"

	limits := response.Header[limitHeader]
	usages := response.Header[usageHeader]

	if len(limits) == 0 || len(usages) == 0 {
		return
	}

	limitValues := strings.Split(limits[0], ",")
	usageValues := strings.Split(usages[0], ",")

	limit15min, _ := strconv.Atoi(limitValues[0])
	limitDay, _ := strconv.Atoi(limitValues[1])
	usage15min, _ := strconv.Atoi(usageValues[0])
	usageDay, _ := strconv.Atoi(usageValues[1])

	client.QueueUpdateQuota <- StravaQuota{
		Limit15min: limit15min,
		LimitDay:   limitDay,
		Usage15min: usage15min,
		UsageDay:   usageDay,
		LastUpdate: time.Now(),
	}
}

func (client *StravaClient) RunQuotaUpdate() {

	for {
		newQuota, ok := <-client.QueueUpdateQuota
		if !ok {
			break
		}

		diffNew := time.Since(newQuota.LastUpdate)
		diffOld := time.Since(client.Quota.LastUpdate)
		if diffNew.Milliseconds() < diffOld.Milliseconds() {
			client.Quota = &newQuota
		}
	}
}
