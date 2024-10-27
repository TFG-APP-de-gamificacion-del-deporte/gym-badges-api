package session_service

import (
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
		service = NewSessionService()
	})

	Context("Generate Session", func() {

		var (
			userID string
		)

		BeforeEach(func() {
			userID = "ironman"
		})

		It("CASE: Successful session generation", func() {
			response, err := service.GenerateSession(userID)
			Expect(err).To(BeNil())
			Expect(response).ToNot(BeNil())
		})

	})

	Context("Validate Session", func() {

		var (
			userID    string
			sessionID string
		)

		BeforeEach(func() {
			userID = "ironman"
			sessionID = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiaXJvbm1hbiIsImV4cCI6MTc2MTU5MTU1NH0.xpHM5ykrlnWEEfJDQ5jKuNOmzCz0wIFBhHxMoWVK_7w"
		})

		It("CASE: Successful session validation", func() {
			err := service.ValidateSession(userID, sessionID)
			Expect(err).To(BeNil())
		})

		It("CASE: Fail session validation because the token is invalid", func() {
			sessionID = "invalid"
			err := service.ValidateSession(userID, sessionID)
			Expect(err).To(BeAssignableToTypeOf(customErrors.UnauthorizedError{}))
		})

		It("CASE: Fail session validation because the token does not belong to the user", func() {
			userID = "thanos"
			err := service.ValidateSession(userID, sessionID)
			Expect(err).To(BeAssignableToTypeOf(customErrors.UnauthorizedError{}))
		})

	})

})
