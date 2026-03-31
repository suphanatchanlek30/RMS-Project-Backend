package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/suphanatchanlek30/rms-project-backend/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	EmployeeID int    `json:"employeeId"`
	RoleID     int    `json:"roleId"`
	RoleName   string `json:"roleName"`
	Email      string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(employeeID, roleID int, roleName, email string) (string, int, error) {
	secret := config.GetEnv("JWT_SECRET", "super-secret-rms-key")
	expiresInStr := config.GetEnv("JWT_EXPIRES_IN_SECONDS", "3600")

	expiresIn, err := strconv.Atoi(expiresInStr)
	if err != nil {
		expiresIn = 3600
	}

	now := time.Now()
	expireAt := now.Add(time.Duration(expiresIn) * time.Second)

	claims := JWTClaims{
		EmployeeID: employeeID,
		RoleID:     roleID,
		RoleName:   roleName,
		Email:      email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(employeeID),
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}

	return signedToken, expiresIn, nil
}

func ParseJWT(tokenString string) (*JWTClaims, error) {
	secret := config.GetEnv("JWT_SECRET", "super-secret-rms-key")

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
