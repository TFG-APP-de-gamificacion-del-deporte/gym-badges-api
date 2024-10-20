package utils

import (
	utils "gym-badges-api/tools/testing"
	"testing"

	. "github.com/onsi/ginkgo/v2"
)

func TestSuite(t *testing.T) {
	utils.ConfigureTestSuite(t, "")
}

var _ = Describe("Test", func() {

	BeforeEach(func() {
	})

	AfterEach(func() {
	})

	Context("Test", func() {

		It("Test 1", func() {

			type TestConfig struct {
				Port int `default:"8080" envconfig:"APP_PORT"`
			}

			var testConfig TestConfig

			LoadGenericConfig(&testConfig)
		})

	})

})
