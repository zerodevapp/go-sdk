package useropbuilder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zerodevapp/sdk-go/cmd/types"
)

// UseropBuilderClient represents a UserOp Builder API client
type UseropBuilderClient struct {
	projectID  string
	baseURL    string
	httpClient *http.Client
}

// NewUserOpBuilder creates a new UserOp Builder API client
func NewUserOpBuilder(projectID string, baseURL string) *UseropBuilderClient {
	return &UseropBuilderClient{
		projectID: projectID,
		baseURL:   baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewUserOpBuilderWithHTTPClient creates a new client with a custom HTTP client
func NewUserOpBuilderWithHTTPClient(projectID string, baseURL string, httpClient *http.Client) *UseropBuilderClient {
	return &UseropBuilderClient{
		projectID:  projectID,
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

// BuildUserOp builds a user operation
func (c *UseropBuilderClient) InitialiseKernelClient(chainID uint64, ctx context.Context) (bool, error) {
	url := fmt.Sprintf("%s/%s/%d/init-kernel-client", c.baseURL, c.projectID, chainID)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return false, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result types.BuildUserOpResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	return true, nil
}

// BuildUserOp builds a user operation
func (c *UseropBuilderClient) BuildUserOp(ctx context.Context, chainID uint64, req *types.BuildUserOpRequest) (*types.BuildUserOpResponse, error) {
	url := fmt.Sprintf("%s/%s/%d/build-userop", c.baseURL, c.projectID, chainID)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result types.BuildUserOpResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// SendUserOp sends a user operation
func (c *UseropBuilderClient) SendUserOp(ctx context.Context, chainID uint64, req *types.SendUserOpRequest) (*types.SendUserOpResponse, error) {
	url := fmt.Sprintf("%s/%s/%d/send-userop", c.baseURL, c.projectID, chainID)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result types.SendUserOpResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetUserOpReceipt gets the receipt for a user operation
func (c *UseropBuilderClient) GetUserOpReceipt(ctx context.Context, chainID uint64, req *types.GetUserOpReceiptRequest) (*types.UserOpReceipt, error) {
	url := fmt.Sprintf("%s/%s/%d/get-userop-receipt", c.baseURL, c.projectID, chainID)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Try to decode as receipt first
	var result types.UserOpReceipt
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check if response contains an error
	var errorCheck map[string]any
	if err := json.Unmarshal(bodyBytes, &errorCheck); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if _, hasError := errorCheck["error"]; hasError {
		return nil, fmt.Errorf("receipt not found yet")
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to decode receipt: %w", err)
	}

	return &result, nil
}

// WaitForUserOpReceipt polls for the user operation receipt until it's available or timeout
func (c *UseropBuilderClient) WaitForUserOpReceipt(ctx context.Context, chainID uint64, req *types.GetUserOpReceiptRequest, pollInterval time.Duration, timeout time.Duration) (*types.UserOpReceipt, error) {
	if pollInterval == 0 {
		pollInterval = 2 * time.Second
	}
	if timeout == 0 {
		timeout = 60 * time.Second
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	// Try immediately first

	receipt, err := c.GetUserOpReceipt(timeoutCtx, chainID, req)
	if err == nil {

		return receipt, nil
	}

	// Then poll
	attemptNum := 2
	for {
		select {
		case <-timeoutCtx.Done():
			return nil, fmt.Errorf("timed out waiting for user operation receipt after %d attempts", attemptNum-1)

		case <-ticker.C:
			receipt, err := c.GetUserOpReceipt(timeoutCtx, chainID, req)
			if err == nil {

				return receipt, nil
			}

			attemptNum++
		}
	}
}
