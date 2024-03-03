package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    pb "github.com/atullal/ecommerce-backend-protobuf/order" // Replace with the correct import path
    "rest-service/services"      // Adjust the import path based on your project structure
    "strconv"
)

type OrderHandler struct {
    OrderService *services.OrderService
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
    var req pb.CreateOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    userID, exists := c.Get("userID")
    if !exists {
        // userID not found in the context, handle this scenario
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
        return
    }

    // Cast userID to the appropriate type, e.g., uint or string
    // (This depends on how userID is stored in the context)
    id, ok := userID.(uint) // or .(string), depending on the type
    if !ok {
        // userID is not of the expected type, handle this scenario
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
        return
    }
    req.CustomerId = int64(id)

    resp, err := h.OrderService.CreateOrder(c, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
    // Extract the product ID from the URL parameter
    orderID := c.Param("id") // Assuming the URL is something like /products/:id

    // Convert the productID to int64 (or your desired format)
    id, err := strconv.ParseInt(orderID, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
        return
    }

    // Create a GetProductRequest with the product ID
    req := pb.GetOrderRequest{OrderId: id}

    // Call the ProductService with the context and request
    resp, err := h.OrderService.GetOrder(c, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Respond with the product details
    c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	// Extract the product ID from the URL parameter
    orderID := c.Param("id") // Assuming the URL is something like /products/:id

    // Convert the productID to int64 (or your desired format)
    id, err := strconv.ParseInt(orderID, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
        return
    }

    var req pb.UpdateOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    req.OrderId = id
    resp, err := h.OrderService.UpdateOrder(c, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	userID, exists := c.Get("userID")
    if !exists {
        // userID not found in the context, handle this scenario
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
        return
    }

    // Cast userID to the appropriate type, e.g., uint or string
    // (This depends on how userID is stored in the context)
    id, ok := userID.(uint) // or .(string), depending on the type
    if !ok {
        // userID is not of the expected type, handle this scenario
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
        return
    }

    // Create a GetProductRequest with the product ID
    req := pb.ListOrdersRequest{CustomerId: int64(id)}

    // Call the ProductService with the context and request
    resp, err := h.OrderService.ListOrders(c, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Respond with the product details
    c.JSON(http.StatusOK, resp)
}
