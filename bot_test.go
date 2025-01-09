package pretty_test

import (
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Bot", func() {
	Context("Config validation", configValidation)
	Context("Update processing", updateProcessing)
})
