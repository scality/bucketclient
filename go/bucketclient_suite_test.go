package bucketclient_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBucketclient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bucketclient Suite")
}
