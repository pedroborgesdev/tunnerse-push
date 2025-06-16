package utils

import "tunnerse/models"

type RegisterRequest struct {
	Name            string               `json:"name" binding:"required"`
	ClientPublicKey models.PublicKeyJSON `json:"clientPublicKey"`
}
