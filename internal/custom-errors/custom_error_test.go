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

})
