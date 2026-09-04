package cvm

import (
	"testing"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cvm/v20170312"
)

func TestCvmInstanceHasPublicIPAddress(t *testing.T) {
	tests := []struct {
		name string
		inst *cvm.Instance
		want bool
	}{
		{name: "nil instance", inst: nil, want: false},
		{name: "no public IP", inst: &cvm.Instance{}, want: false},
		{name: "empty public IP", inst: &cvm.Instance{PublicIpAddresses: []*string{nil, common.StringPtr("")}}, want: false},
		{name: "has public IP", inst: &cvm.Instance{PublicIpAddresses: common.StringPtrs([]string{"203.0.113.10"})}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cvmInstanceHasPublicIPAddress(tt.inst); got != tt.want {
				t.Fatalf("cvmInstanceHasPublicIPAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
