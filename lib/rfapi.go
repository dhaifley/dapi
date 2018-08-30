// Package lib provides data structures and functions for interacting with
// the Royal Farms Application Programming Interface service library.
package lib

import (
	"github.com/dhaifley/dlib"
)

// ServiceInfo provides information about this service.
var ServiceInfo dlib.ServiceInfo

func init() {
	ServiceInfo = dlib.ServiceInfo{
		Name:    "dapi",
		Short:   "Royal Farms Application Programming Interface",
		Long:    "Provides a unified interface for interacting will all Royal Farms services.",
		Version: "1.0.1",
	}
}
