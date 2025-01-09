package pretty_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPrettyChatBot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PrettyChatBot Suite")
}
