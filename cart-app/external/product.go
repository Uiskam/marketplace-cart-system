package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ProductClient struct {
	baseURL    string
	httpClient *http.Client
}

type LockProductRequest struct {
	ProductUUID   string `json:"product_uuid"`
	LockingEntity string `json:"locking_entity"`
}

type LockProductResponse struct {
	Message string `json:"message"`
	Price   int    `json:"price"`
}

type UnlockProductRequest struct {
	ProductUUID   string `json:"product_uuid"`
	LockingEntity string `json:"locking_entity"`
}

type SellProductRequest struct {
	ProductUUID   string `json:"product_uuid"`
	LockingEntity string `json:"locking_entity"`
}

type ProductResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func NewProductClient(baseURL string) *ProductClient {
	return &ProductClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *ProductClient) LockProduct(productUUID, lockingEntity string) (int, error) {
	req := LockProductRequest{
		ProductUUID:   productUUID,
		LockingEntity: lockingEntity,
	}

	var resp LockProductResponse
	err := c.makeRequest("POST", "/products/lock", req, &resp)
	if err != nil {
		return 0, err
	}

	return resp.Price, nil
}

func (c *ProductClient) UnlockProduct(productUUID, lockingEntity string) error {
	req := UnlockProductRequest{
		ProductUUID:   productUUID,
		LockingEntity: lockingEntity,
	}

	var resp ProductResponse
	err := c.makeRequest("POST", "/products/unlock", req, &resp)
	if err != nil {
		return err
	}

	return nil
}

func (c *ProductClient) SellProduct(productUUID, lockingEntity string) error {
	req := SellProductRequest{
		ProductUUID:   productUUID,
		LockingEntity: lockingEntity,
	}

	var resp ProductResponse
	err := c.makeRequest("POST", "/products/sell", req, &resp)
	if err != nil {
		return err
	}

	return nil
}

func (c *ProductClient) makeRequest(method, endpoint string, reqBody interface{}, respBody interface{}) error {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp ProductResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return fmt.Errorf("request failed with status %d", resp.StatusCode)
		}
		return fmt.Errorf("product API error: %s", errorResp.Error)
	}

	if err := json.NewDecoder(resp.Body).Decode(respBody); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}