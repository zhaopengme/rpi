package epd

import (
	"fmt"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/host"
)

// Epd is basic struc for Waveshare eps2in13bc
type Epd struct {
	width   int
	height  int
	port    spi.PortCloser
	spiConn spi.Conn
	rstPin  gpio.PinIO
	dcPin   gpio.PinIO
	csPin   gpio.PinIO
	busyPin gpio.PinIO
}

// CreateEpd is constructor for Epd
func CreateEpd() Epd {
	e := Epd{
		width:  800,
		height: 480,
	}

	var err error

	host.Init()

	// SPI
	e.port, err = spireg.Open("")
	if err != nil {
		fmt.Println(err)
	}

	e.spiConn, err = e.port.Connect(4*physic.MegaHertz, spi.Mode0, 8)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(e.spiConn)

	// GPIO - read
	e.rstPin = gpioreg.ByName("GPIO17")  // out
	e.dcPin = gpioreg.ByName("GPIO25")   // out
	e.csPin = gpioreg.ByName("GPIO8")    // out
	e.busyPin = gpioreg.ByName("GPIO24") // in

	return e
}

// Close is closing pariph.io port
func (e *Epd) Close() {
	e.port.Close()
}

// reset epd
func (e *Epd) reset() {
	e.rstPin.Out(true)
	time.Sleep(200 * time.Millisecond)
	e.rstPin.Out(false)
	time.Sleep(2 * time.Millisecond)
	e.rstPin.Out(true)
	time.Sleep(200 * time.Millisecond)
}

// sendCommand sets DC ping low and sends byte over SPI
func (e *Epd) sendCommand(command byte) {
	e.dcPin.Out(false)
	e.csPin.Out(false)
	c := []byte{command}
	r := make([]byte, len(c))
	e.spiConn.Tx(c, r)
	e.csPin.Out(true)
}

// sendData sets DC ping high and sends byte over SPI
func (e *Epd) sendData(data byte) {
	e.dcPin.Out(true)
	e.csPin.Out(false)
	c := []byte{data}
	r := make([]byte, len(c))
	e.spiConn.Tx(c, r)
	e.csPin.Out(true)
}

// ReadBusy waits for epd
func (e *Epd) readBusy() {
	fmt.Println("e-Paper busy")
	// Low for busy
	for e.busyPin.Read() == gpio.Low {
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("e-Paper busy release")
}

// Sleep powers off the epd
func (e *Epd) Sleep() {
	// POWER_OFF
	e.sendCommand(0x02)
	e.readBusy()
	// DEEP_SLEEP
	e.sendCommand(0x07)
	// check code
	e.sendData(0xA5)
}

// Display sends []byte data to the epd
func (e *Epd) Display(imagered []byte) {
	e.sendCommand(0x13)
	for i := 0; i < (e.width*e.height)/8; i++ {
		e.sendData(imagered[i])
	}

	// REFRESH
	e.sendCommand(0x12)
	time.Sleep(100 * time.Millisecond)
	e.readBusy()
}

// Init starts the epd
func (e *Epd) Init() {
	fmt.Println("reset")
	e.reset()

	fmt.Println("BOOSTER_SOFT_START")
	e.sendCommand(0x01) // POWER SETTING
	e.sendData(0x07)
	e.sendData(0x07) // VGH=20V,VGL=-20V
	e.sendData(0x3f) // VDH=15V
	e.sendData(0x3f) // VDL=-15V

	e.sendData(0x17)

	//POWER_ON
	fmt.Println("POWER_ON")
	e.sendCommand(0x04)
	time.Sleep(100 * time.Millisecond)
	e.readBusy()

	//PANEL_SETTING
	fmt.Println("PANEL_SETTING")
	e.sendCommand(0x00)
	e.sendData(0x1F) // KW-3f   KWR-2F	BWROTP 0f	BWOTP 1f

	// VCOM_AND_DATA_INTERVAL_SETTING
	fmt.Println("VCOM_AND_DATA_INTERVAL_SETTING")
	e.sendCommand(0x61) // tres
	e.sendData(0x03)    // source 800
	e.sendData(0x20)
	e.sendData(0x01) // gate 480
	e.sendData(0xE0)

	e.sendCommand(0x15)
	e.sendData(0x00)

	e.sendCommand(0X50) // VCOM AND DATA INTERVAL SETTING
	e.sendData(0x10)
	e.sendData(0x07)

	e.sendCommand(0x60) // TCON SETTING
	e.sendData(0x22)
}

// Clear sets epd display to white
func (e *Epd) Clear() {

	fmt.Println("CLEAR 1")
	e.sendCommand(0x10)
	for i := 0; i < (e.width*e.height)/8; i++ {
		e.sendData(0x00)
	}

	fmt.Println("CLEAR 2")
	e.sendCommand(0x13)
	for i := 0; i < (e.width*e.height)/8; i++ {
		e.sendData(0x00)
	}

	// REFRESH
	fmt.Println("REFRESH")
	e.sendCommand(0x12)
	time.Sleep(100 * time.Millisecond)
	e.readBusy()

}
