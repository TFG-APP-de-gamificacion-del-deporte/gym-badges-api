package custom_errors

import (
	toolsTesting "gym-badges-api/tools/testing"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCustomErrorsSuite(t *testing.T) {
	toolsTesting.ConfigureTestSuite(t, "ERRORS: Custom Errors Test Suite")
}

var _ = Describe("ERRORS: Custom Errors Test Suite", func() {

	Context("Unauthorized Error", func() {

		It("BuildUnauthorizedError", func() {
			err := BuildUnauthorizedError("unauthorized")
			Expect(err.Error()).To(Equal("unauthorized"))
		})

		It("BuildUnauthorizedError with parameters", func() {
			err := BuildUnauthorizedError("unauthorized %d", http.StatusUnauthorized)
			Expect(err.Error()).To(Equal("unauthorized 401"))
		})

	})

	Context("Conflict Error", func() {

		It("BuildConflictError", func() {
			err := BuildConflictError("conflict")
			Expect(err.Error()).To(Equal("conflict"))
		})

		It("BuildConflictError with parameters", func() {
			err := BuildConflictError("conflict %d", http.StatusConflict)
			Expect(err.Error()).To(Equal("conflict 409"))
		})

	})

	Context("Not Found Error", func() {

		It("BuildNotFoundError", func() {
			err := BuildNotFoundError("not found")
			Expect(err.Error()).To(Equal("not found"))
		})

		It("BuildNotFoundError with parameters", func() {
			err := BuildNotFoundError("not found %d", http.StatusNotFound)
			Expect(err.Error()).To(Equal("not found 404"))
		})

	})

})
