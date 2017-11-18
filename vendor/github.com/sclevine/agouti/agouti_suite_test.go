package agouti_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAgouti(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Agouti Suite")
}
