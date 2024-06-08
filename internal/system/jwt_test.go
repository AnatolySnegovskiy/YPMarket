package system

import (
	"fmt"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	token, err := CreateToken(100)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

// Helper function to create a token for testing
func createTestJWT(userID int, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(userID),                       // jwt-go expects claims to be of certain types
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // set expiration to one day ahead
	})

	return token.SignedString([]byte(secretKey))
}

func TestGetUserID(t *testing.T) {
	secretKey := "your_secret_key"
	validUserID := 123
	wrongMethodToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.e30.invalidsignature"
	noUserIDToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{}).SignedString([]byte("your_secret_key"))
	assert.NoError(t, err, "creating test JWT should not produce an error")
	wrongTypeUserIDToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": "not-a-float64"}).SignedString([]byte("your_secret_key"))
	assert.NoError(t, err, "creating test JWT should not produce an error")
	invalidToken := "invalid.token.parts"

	validToken, err := createTestJWT(validUserID, secretKey)
	assert.NoError(t, err, "creating test JWT should not produce an error")

	testCases := []struct {
		name          string
		signedToken   string
		expectedID    int
		expectedError error
	}{
		{
			name:          "Valid Token",
			signedToken:   validToken,
			expectedID:    validUserID,
			expectedError: nil,
		},
		{
			name:          "Invalid Token - Wrong Secret",
			signedToken:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.HW1GlqIzEVUg5gdbDRM3p1h0YVGsRrCtvC6Ol6R8X0Y",
			expectedID:    0,
			expectedError: fmt.Errorf("signature is invalid"),
		},
		{
			name:          "Invalid Token - No UserID Claim",
			signedToken:   "", // Insert an invalid token without user_id claim here
			expectedID:    0,
			expectedError: fmt.Errorf("token contains an invalid number of segments"),
		},
		{
			name:          "Unexpected Signing Method",
			signedToken:   wrongMethodToken,
			expectedID:    0,
			expectedError: fmt.Errorf("unexpected signing method: RS256"),
		},
		{
			name:          "Missing UserID Claim",
			signedToken:   noUserIDToken,
			expectedID:    0,
			expectedError: fmt.Errorf("user_id claim is missing or not of type float64"),
		},
		{
			name:          "Wrong Type UserID Claim",
			signedToken:   wrongTypeUserIDToken,
			expectedID:    0,
			expectedError: fmt.Errorf("user_id claim is missing or not of type float64"),
		},
		{
			name:          "Invalid Token",
			signedToken:   invalidToken,
			expectedError: fmt.Errorf("invalid character '\\u008a' looking for beginning of value"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userID, err := GetUserID(tc.signedToken)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedID, userID)
		})
	}
}
