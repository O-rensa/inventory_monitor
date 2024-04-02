package ap_auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	ap_shared "github.com/o-rensa/iv/internal/adapters/primary/shared"
	ps_admin "github.com/o-rensa/iv/internal/ports/secondary/admin"
	"github.com/o-rensa/iv/pkg/configs"
)

type contextKey string

const UserKey contextKey = "adminID"

func WithJWTAuth(handlerFunc http.HandlerFunc, store ps_admin.AdminStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := ap_shared.GetTokenFromRequest(r)

		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Printf("invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["adminID"].(string)

		adminID, err := uuid.Parse(str)
		if err != nil {
			log.Printf("failed to convert adminID into uuid: %v", err)
			permissionDenied(w)
			return
		}

		a, err := store.GetAdminByID(adminID)
		if err != nil {
			log.Printf("failed to get admin by id: %v", err)
			permissionDenied(w)
			return
		}

		// add admin to the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, a.ID)
		r = r.WithContext(ctx)

		// call the function if the token is valid
		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(configs.Configs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	ap_shared.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func CreateJWT(secret []byte, adminID uuid.UUID) (string, error) {
	expiration := time.Second * time.Duration(configs.Configs.JWTExpirationinSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"adminID":   adminID.String(),
		"expiresAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, err
}
