package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    pb "github.com/atullal/ecommerce-backend/user-service/gen" // Replace with the correct import path
    "rest-service/services"      // Adjust the import path based on your project structure
)

type UserHandler struct {
    UserService *services.UserService
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var req pb.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := h.UserService.CreateUser(c, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) AuthenticateUser(c *gin.Context) {
    var req pb.AuthenticateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := h.UserService.AuthenticateUser(c, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}
