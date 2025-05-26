package session_service

import (
	configs "gym-badges-api/config/gym-badges-server"
	customErrors "gym-badges-api/internal/custom-errors"
	toolsTesting "gym-badges-api/tools/testing"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestServiceSessionSuite(t *testing.T) {
	toolsTesting.ConfigureTestSuite(t, "SERVICE: Session Test Suite")
}

var _ = Describe("SERVICE: Session Test Suite", func() {

	var (
		service ISessionService
	)

	BeforeEach(func() {
		configs.Basic.SessionDuration = 10
		service = NewSessionService()
	})

	Context("GenerateSession", func() {

		It("CASE: Successfully generates a valid JWT token", func() {
			userID := "test-user"
			token, err := service.GenerateSession(userID)
			Expect(err).To(BeNil())
			Expect(token).ToNot(BeEmpty())

			// Verify the token is valid
			err = service.ValidateSession(userID, token)
			Expect(err).To(BeNil())
		})

		It("CASE: Generated token expires after configured duration", func() {

			// Set expiration time to 1 second ago
			configs.Basic.SessionDuration = -1

			userID := "test-user"
			token, err := service.GenerateSession(userID)
			Expect(err).To(BeNil())

			err = service.ValidateSession(userID, token)
			Expect(err).To(BeAssignableToTypeOf(customErrors.UnauthorizedError{}))
		})
	})

	Context("ValidateSession", func() {
		It("CASE: Successfully validates a valid token", func() {
			userID := "test-user"
			token, err := service.GenerateSession(userID)
			Expect(err).To(BeNil())

			err = service.ValidateSession(userID, token)
			Expect(err).To(BeNil())
		})

		It("CASE: Fails to validate token with wrong user ID", func() {
			userID := "test-user"
			token, err := service.GenerateSession(userID)
			Expect(err).To(BeNil())

			err = service.ValidateSession("wrong-user", token)
			Expect(err).To(BeAssignableToTypeOf(customErrors.UnauthorizedError{}))
		})

		It("CASE: Fails to validate invalid token", func() {
			err := service.ValidateSession("test-user", "invalid-token")
			Expect(err).To(BeAssignableToTypeOf(customErrors.UnauthorizedError{}))
		})
	})
})
