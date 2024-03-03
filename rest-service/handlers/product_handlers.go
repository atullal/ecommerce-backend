package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    pb "github.com/atullal/ecommerce-backend-protobuf/product" // Replace with the correct import path
    "rest-service/services"      // Adjust the import path based on your project structure
    "strconv"
)

type ProductHandler struct {
    ProductService *services.ProductService
}

func (h *ProductHandler) AddProduct(c *gin.Context) {
    var req pb.AddProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := h.ProductService.AddProduct(c, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
    // Extract the product ID from the URL parameter
    productID := c.Param("id") // Assuming the URL is something like /products/:id

    // Convert the productID to int64 (or your desired format)
    id, err := strconv.ParseInt(productID, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    // Create a GetProductRequest with the product ID
    req := pb.GetProductRequest{Id: id}

    // Call the ProductService with the context and request
    resp, err := h.ProductService.GetProduct(c, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Respond with the product details
    c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	// Extract the product ID from the URL parameter
    productID := c.Param("id") // Assuming the URL is something like /products/:id

    // Convert the productID to int64 (or your desired format)
    id, err := strconv.ParseInt(productID, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    var req pb.UpdateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    req.Id = id
    resp, err := h.ProductService.UpdateProduct(c, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
    // Extract the product ID from the URL parameter
    productID := c.Param("id") // Assuming the URL is something like /products/:id

    // Convert the productID to int64 (or your desired format)
    id, err := strconv.ParseInt(productID, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    // Create a GetProductRequest with the product ID
    req := pb.DeleteProductRequest{Id: id}

    // Call the ProductService with the context and request
    resp, err := h.ProductService.DeleteProduct(c, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Respond with the product details
    c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
    // Parse query parameters
    page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
    pageSize, _ := strconv.ParseInt(c.DefaultQuery("pageSize", "10"), 10, 32)
    searchKeyword := c.Query("searchKeyword")
    categories := c.QueryArray("categories") // This is for parameters like ?categories=cat1&categories=cat2

    // Create ListProductsRequest
    req := pb.ListProductsRequest{
        Page:         int32(page),
        PageSize:     int32(pageSize),
        SearchKeyword: searchKeyword,
        Categories:   categories,
    }

    // Call the ProductService with the context and request
    resp, err := h.ProductService.ListProducts(c, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Respond with the product details
    c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) UpdateInventory(c *gin.Context) {
	// Extract the product ID from the URL parameter
    productID := c.Param("id") // Assuming the URL is something like /products/:id

    // Convert the productID to int64 (or your desired format)
    id, err := strconv.ParseInt(productID, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    var req pb.UpdateInventoryRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    req.ProductId = id
    resp, err := h.ProductService.UpdateInventory(c, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) GetInventory(c *gin.Context) {
    // Extract the product ID from the URL parameter
    productID := c.Param("id") // Assuming the URL is something like /products/:id

    // Convert the productID to int64 (or your desired format)
    id, err := strconv.ParseInt(productID, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    // Create a GetProductRequest with the product ID
    req := pb.GetInventoryRequest{ProductId: id}

    // Call the ProductService with the context and request
    resp, err := h.ProductService.GetInventory(c, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Respond with the product details
    c.JSON(http.StatusOK, resp)
}
