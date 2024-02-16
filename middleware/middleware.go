package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	// UserId   int    `json:"user_id"`
	UserUUID           string `json:"user_uuid"`
	AppRoleId          int    `json:"application_role_id"`
	DivisionTitle      string `json:"division_title"`
	RoleCode           string `json:"role_code"`
	jwt.StandardClaims        // Embed the StandardClaims struct

}

func DecryptJWE(jweToken string, secretKey string) (string, error) {
	// Dekripsi token JWE
	decrypted, _, err := jose.Decode(jweToken, secretKey)
	if err != nil {
		return "", err
	}
	return decrypted, nil
}

// func DecryptJWE(jweToken string, secretKey string) (string, error) {
// 	// Dekripsi token JWE
// 	decrypted, _, err := jose.Decode(jweToken, secretKey)
// 	if err != nil {
// 		return "", err
// 	}
// 	return decrypted, nil
// }

func ExtractClaims(jwtToken string) (JwtCustomClaims, error) {
	claims := &JwtCustomClaims{}
	secretKey := "secretJwToken" // Ganti dengan kunci yang benar

	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return JwtCustomClaims{}, err
	}

	return *claims, nil
}

// func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		tokenString := c.Request().Header.Get("Authorization")
// 		//secretKey := "secretJwToken" // Ganti dengan kunci yang benar

// 		// Periksa apakah tokenString tidak kosong
// 		if tokenString == "" {
// 			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 				"code":    401,
// 				"message": "Token tidak ditemukan!",
// 				"status":  false,
// 			})
// 		}

// 		// Periksa apakah tokenString mengandung "Bearer "
// 		if !strings.HasPrefix(tokenString, "Bearer ") {
// 			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 				"code":    401,
// 				"message": "Token tidak valid!",
// 				"status":  false,
// 			})
// 		}

// 		// Hapus "Bearer " dari tokenString
// 		tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

// 		// Tetapkan token yang telah didekripsi ke header konteks
// 		c.Request().Header.Set("DecryptedToken", tokenOnly)

//			// Token JWE valid, Anda dapat melanjutkan dengan pengolahan berikutnya
//			return next(c)
//		}
// //	}
// func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		tokenString := c.Request().Header.Get("Authorization")
// 		secretKey := "secretJwToken" // Ganti dengan kunci yang benar

// 		// Periksa apakah tokenString tidak kosong
// 		if tokenString == "" {
// 			log.Print("tidak ada token:", tokenString)
// 			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 				"code":    401,
// 				"message": "Token tidak ditemukan!",
// 				"status":  false,
// 			})
// 		}

// 		// Periksa apakah tokenString mengandung "Bearer "
// 		if !strings.HasPrefix(tokenString, "Bearer ") {
// 			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 				"code":    401,
// 				"message": "Token tidak valid!",
// 				"status":  false,
// 			})
// 		}

// 		// Hapus "Bearer " dari tokenString
// 		tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

// 		// Langkah 1: Mendekripsi token JWE
// 		decrypted, err := DecryptJWE(tokenOnly, secretKey)
// 		if err != nil {
// 			fmt.Println("Gagal mendekripsi token:", err)
// 			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 				"code":    401,
// 				"message": "Token tidak valid!",
// 				"status":  false,
// 			})
// 		}

// 		// Parse token JWT yang telah didekripsi
// 		var claims JwtCustomClaims
// 		err = json.Unmarshal([]byte(decrypted), &claims)
// 		if err != nil {
// 			fmt.Println("Gagal mengurai klaim:", err)
// 			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 				"code":    401,
// 				"message": "Token tidak valid!",
// 				"status":  false,
// 			})
// 		}

// 		// Set user UUID ke dalam header konteks
// 		c.Set("user_uuid", claims.UserUUID)
// 		c.Set("DecryptedToken", decrypted) // tambahkan ini untuk token yang sudah didekripsi

// 		// Token JWE valid, Anda dapat melanjutkan dengan pengolahan berikutnya
// 		return next(c)
// 	}
// }

// func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		tokenString := c.Request().Header.Get("Authorization")
// 		secretKey := "secretJwToken" // Ganti dengan kunci yang benar

// 		// Periksa apakah tokenString tidak kosong
// 		if tokenString == "" {
// 			log.Print("tidak ada token:", tokenString)
// 			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 				"code":    401,
// 				"message": "Token tidak ditemukan!",
// 				"status":  false,
// 			})
// 		}

// 		// Periksa apakah tokenString mengandung "Bearer "
// 		if !strings.HasPrefix(tokenString, "Bearer ") {
// 			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 				"code":    401,
// 				"message": "Token tidak valid!",
// 				"status":  false,
// 			})
// 		}

// 		// Hapus "Bearer " dari tokenString
// 		tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

// 		// Langkah 1: Mendekripsi token JWE
// 		decrypted, err := DecryptJWE(tokenOnly, secretKey)
// 		if err != nil {
// 			fmt.Println("Gagal mendekripsi token:", err)
// 			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 				"code":    401,
// 				"message": "Token tidak valid!",
// 				"status":  false,
// 			})
// 		}

// 		// Parse token JWT yang telah didekripsi
// 		var claims JwtCustomClaims
// 		err = json.Unmarshal([]byte(decrypted), &claims)
// 		if err != nil {
// 			fmt.Println("Gagal mengurai klaim:", err)
// 			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 				"code":    401,
// 				"message": "Token tidak valid!",
// 				"status":  false,
// 			})
// 		}

// 		// Set user UUID ke dalam header konteks
// 		c.Set("user_uuid", claims.UserUUID)
// 		c.Set("DecryptedToken", decrypted) // tambahkan ini untuk token yang sudah didekripsi

// 		// Token JWE valid, Anda dapat melanjutkan dengan pengolahan berikutnya
// 		return next(c)
// 	}
// }

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		secretKey := "secretJwToken" // Ganti dengan kunci yang benar

		// Periksa apakah tokenString tidak kosong
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak ditemukan!",
				"status":  false,
			})
		}

		// Periksa apakah tokenString mengandung "Bearer "
		if !strings.HasPrefix(tokenString, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		// Hapus "Bearer " dari tokenString
		tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

		// Langkah 1: Mendekripsi token JWE
		decrypted, err := DecryptJWE(tokenOnly, secretKey)
		if err != nil {
			fmt.Println("Gagal mendekripsi token:", err)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		fmt.Println("Token yang sudah dideskripsi:", decrypted)

		var claims JwtCustomClaims
		errJ := json.Unmarshal([]byte(decrypted), &claims)
		if errJ != nil {
			fmt.Println("Gagal mengurai klaim:", errJ)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		// Sekarang Anda memiliki data dalam struct JwtCustomClaims
		// Anda bisa mengakses UserId atau klaim lain sesuai kebutuhan
		// fmt.Println("UserID:", claims.UserId)

		userUUID := claims.UserUUID // Mengakses UserID langsung
		// roleID := claims.AppRoleId
		// divisionTitle := claims.DivisionTitle
		// roleCode := claims.RoleCode
		// if roleCode != "" {
		// 	log.Print(roleCode)
		// }

		fmt.Println("User UUID:", userUUID)
		// fmt.Println("Role Code:", roleCode)

		c.Set("user_uuid", userUUID)
		// c.Set("application_role_id", roleID)
		// c.Set("division_title", divisionTitle)
		// c.Set("role_code", roleCode)

		// Token JWE valid, Anda dapat melanjutkan dengan pengolahan berikutnya
		return next(c)
	}
}
