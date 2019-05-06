package usb

import (
	"errors"
	"fmt"
	"io"
)

const (
	VendorT1            = 0x313a
	ProductT1Bootloader = 0x0000
	ProductT1Firmware   = 0x0001
	VendorT2            = 0x1209
	ProductT2Bootloader = 0x53C0
	ProductT2Firmware   = 0x53C1
)

var (
	ErrNotFound = fmt.Errorf("device not found")
)

type Info struct {
	Path      string
	VendorID  int
	ProductID int
}

type Device interface {
	io.ReadWriteCloser
}

type Bus interface {
	// Enumerate returns a list of all the devices accessible in the the system
	// - If the vendor id is set to 0 then any vendor matches.
	// - If the product id is set to 0 then any product matches.
	// - If the vendor and product id are both 0, all devices are returned.
	Enumerate(vendorID uint16, productID uint16) ([]Info, error)
	Connect(path string) (Device, error)
	Has(path string) bool
}

type USB struct {
	buses []Bus
}

func Init(buses ...Bus) *USB {
	return &USB{
		buses: buses,
	}
}

func (b *USB) Enumerate(vendorID uint16, productID uint16) ([]Info, error) {
	var infos []Info

	for _, b := range b.buses {
		l, err := b.Enumerate(vendorID, productID)
		if err != nil {
			return nil, err
		}
		infos = append(infos, l...)
	}
	return infos, nil
}

func (b *USB) Connect(path string) (Device, error) {
	for _, b := range b.buses {
		if b.Has(path) {
			return b.Connect(path)
		}
	}
	return nil, ErrNotFound
}

var errDisconnect = errors.New("Device disconnected during action")
var errClosedDevice = errors.New("Closed device")
