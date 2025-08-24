package lcd

import (
	"fmt"
	"log"
	"strings"

	"periph.io/x/conn/v3/display"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/hd44780"
	"periph.io/x/devices/v3/mcp23xxx"
	"periph.io/x/host/v3"
)

const (
	d4           = 3
	d5           = 4
	d6           = 5
	d7           = 6
	rsPin        = 1
	enablePin    = 2
	backlightPin = 7
)

type AdafruitLCDDevice struct {
	mcp          *mcp23xxx.Dev
	bus          i2c.BusCloser
	dev          *hd44780.Dev
	backlightPin gpio.PinOut
	columns      uint8
	rows         uint8
	format       string
}

func NewAdafruitI2CBackpack(bus i2c.BusCloser, address uint16, columns uint8, rows uint8) (*AdafruitLCDDevice, error) {
	mcp, err := mcp23xxx.NewI2C(bus, mcp23xxx.MCP23008, address)
	if err != nil {
		return nil, err
	}

	gr := *mcp.Group(0, []int{d4, d5, d6, d7, rsPin, enablePin, backlightPin})
	reset, _ := gr.ByOffset(4).(gpio.PinOut)
	enable, _ := gr.ByOffset(5).(gpio.PinOut)
	bl := gr.ByOffset(6).(gpio.PinOut)

	dataPins := []gpio.PinOut{
		gr.ByOffset(0).(gpio.PinOut),
		gr.ByOffset(1).(gpio.PinOut),
		gr.ByOffset(2).(gpio.PinOut),
		gr.ByOffset(3).(gpio.PinOut),
	}

	dev, err := hd44780.New(dataPins, reset, enable)

	return &AdafruitLCDDevice{
		mcp:          mcp,
		bus:          bus,
		backlightPin: bl,
		dev:          dev,
		columns:      columns,
		rows:         rows,
		format:       fmt.Sprintf("%%%ds", columns),
	}, nil
}

func NewDevice(address uint16, columns uint8, rows uint8) (*AdafruitLCDDevice, error) {
	if _, err := host.Init(); err != nil {
		return nil, err
	}

	bus, err := i2creg.Open("")
	if err != nil {
		return nil, err
	}

	dev, err := NewAdafruitI2CBackpack(bus, address, columns, rows)
	if err != nil {
		return nil, err
	}

	return dev, nil
}

func (d *AdafruitLCDDevice) Close() error {
	if err := d.mcp.Close(); err != nil {
		return err
	}

	if err := d.bus.Close(); err != nil {
		return err
	}

	return nil
}

func (d *AdafruitLCDDevice) MustClose() {
	err := d.Close()
	if err != nil {
		panic(err)
	}
}

func (d *AdafruitLCDDevice) Backlight(intensity display.Intensity) (err error) {
	if intensity == 0 {
		err = d.backlightPin.Out(gpio.Low)
	} else {
		err = d.backlightPin.Out(gpio.High)
	}
	return
}

func (d *AdafruitLCDDevice) UpdateText(lines ...string) error {
	log.Println("Updating text...")
	if len(lines) > int(d.rows) {
		lines = lines[:d.rows]
	}

	log.Println("Enabling backlight...")
	if err := d.Backlight(1); err != nil {
		return err
	}

	log.Println(">---------------")
	for i, line := range lines {
		if err := d.dev.SetCursor(uint8(i), 0); err != nil {
			return err
		}

		line = strings.Trim(line, " \t")

		if len(line) > int(d.columns) {
			line = fmt.Sprintf(d.format, line[:d.columns])
		} else {
			line = fmt.Sprintf(d.format, line)
		}

		log.Println(line)
		if err := d.dev.Print(line); err != nil {
			return err
		}
	}
	log.Println(">---------------")

	return nil
}
