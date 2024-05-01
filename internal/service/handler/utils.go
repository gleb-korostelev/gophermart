package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gleb-korostelev/gophermart.git/internal/config"
	"github.com/gleb-korostelev/gophermart.git/internal/models"
	"github.com/gleb-korostelev/gophermart.git/internal/workerpool"
	"github.com/gleb-korostelev/gophermart.git/tools/logger"
)

func (svc *APIService) CheckOrderStatus(login string, orders []models.OrdersData) []models.OrdersData {
	baseURL := config.AccuralSystemAddress + "/api/orders/%s"
	var ordersRes []models.OrdersData
	for _, order := range orders {
		if order.Status == "PROCESSED" {
			logger.Infof("Order already processed: %s", order.Number)
			ordersRes = append(ordersRes, order)
			continue
		}

		orderURL := fmt.Sprintf(baseURL, string(order.Number))
		resp, err := http.Get(orderURL)
		if err != nil {
			ordersRes = append(ordersRes, order)
			logger.Errorf("Failed to fetch order: %s, error: %v", order.Number, err)
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			ordersRes = append(ordersRes, order)
			logger.Errorf("Failed to read response body for order: %s, error: %v", order.Number, err)
			continue
		}
		if resp.StatusCode == http.StatusNoContent {
			ordersRes = append(ordersRes, order)
			logger.Infof("No data for order number: %s", order.Number)
			continue
		} else if resp.StatusCode != http.StatusOK {
			ordersRes = append(ordersRes, order)
			logger.Errorf("Bad response status: %s for order: %s", resp.Status, order.Number)
			continue
		}
		var orderInfo models.OrderResponse
		if err := json.Unmarshal(body, &orderInfo); err != nil {
			ordersRes = append(ordersRes, order)
			logger.Errorf("Failed to unmarshal order data for order: %s, error: %v", order.Number, err)
			continue
		}

		updatedOrder := models.OrdersData{
			Number:     orderInfo.Order,
			Status:     orderInfo.Status,
			Accrual:    orderInfo.Accrual,
			UploadedAt: order.UploadedAt,
		}

		ordersRes = append(ordersRes, updatedOrder)

		svc.worker.AddTask(workerpool.Task{
			Action: func(ctx context.Context) error {
				return svc.store.UpdateOrderInfo(ctx, login, orderInfo)
			},
		})
	}
	return ordersRes
}
