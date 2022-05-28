package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"go.bug.st/serial"
)

var commands = []string{
	"POF",     // Power off
	"PON",     // Power on,
	"QPW",     // Query power status
	"QMI",     // Query input
	"IMS:HM1", // HDMI
	"IMS:HM2", // HDMI2
	"IMS:DP1", // Display Port
	"IMS:DV1", // DVI
	"IMS:PC1", // PC/VGA
	"IMS:YP1", // Component
	"IMS:VD1", // Composite
	"IMS:UD1", // USB
	"QAV",     // Query audio volume
	"AUU",     // Audio volume up
	"AUD",     // Audio volume down
	"QAM",     // Query mute
	"AMT:0",   // Audio mute - on
	"AMT:1",   // Audio mute - off
	"QVM",
	"VMT:0",
}

type Queryable string

const (
	QueryPowerStatus Queryable = "QPW"
	QueryInput       Queryable = "QMI"
	QueryAudioVolume Queryable = "QAV"
	QueryAudioMute   Queryable = "QAM"
	QueryVideoMute   Queryable = "QVM"
)

type remote struct {
	port serial.Port
}

func (r *remote) work() {
	for {
		time.Sleep(time.Millisecond * 50)
		byteBuffer := make([]byte, 2028)
		n, err := r.port.Read(byteBuffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				continue
			}
			panic(err)
		}
		if n > 0 {
			log.Println(string(byteBuffer[0:n]))
		}
	}
}

var (
	stx   byte = '\x02'
	colon byte = ':'
	etx   byte = '\x03'
)

func (r *remote) SendCommand(cmdBytes []byte) error {
	cmd := []byte{}
	cmd = append(cmd, stx)

	cmd = append(cmd, cmdBytes...)

	cmd = append(cmd, etx)

	log.Println("CMD >>> ", hex.EncodeToString(cmd))
	_, err := r.port.Write(cmd)
	if err != nil {
		return err
	}
	return nil
}

func (r *remote) derpAround() {
	//shouldTurnOn := true

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd := scanner.Text()
		err := r.SendCommand([]byte(cmd))
		if err != nil {
			panic(err)
		}
	}

	//for {

	//cmdText := "PON"
	//if !shouldTurnOn {
	//	cmdText = "POF"
	//}
	//log.Println("Setting TV power", cmdText)
	//if err := r.SendCommand([]byte(cmdText)); err != nil {
	//	panic(err)
	//}
	//
	//shouldTurnOn = !shouldTurnOn
	//	select {
	//	case <-time.After(time.Millisecond):
	//		continue
	//	}
	//}
}

func main() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		fmt.Println("No serial ports found!")
	} else {
		for _, port := range ports {
			fmt.Printf("Found port: %v\n", port)
			if strings.Contains(port, "ttyUSB0") {
				sp, err := serial.Open(port, &serial.Mode{
					BaudRate: 9600,
					DataBits: 8,
					Parity:   serial.NoParity,
					StopBits: serial.OneStopBit,
				})
				if err != nil {
					panic(err)
				}
				r := remote{
					port: sp,
				}
				go r.work()
				r.derpAround()
			}
		}
	}
}
