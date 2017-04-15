package main

import (
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

const Version = "0.0.1"

type Data struct {
	toa     int64
	payload []byte
}

func handlePackets(ps *gopacket.PacketSource, num uint, ch chan Data) {
	var count uint
	for packet := range ps.Packets() {
		var k Data
		count++
		if count > num {
			break
		}
		elements := packet.Data()
		if len(elements) == 0 {
			continue
		}
		k = Data{toa: packet.Metadata().CaptureInfo.Timestamp.UnixNano(), payload: packet.Data()}
		ch <- k
	}
	close(ch)
	return
}

func availableInterfaces() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	for _, device := range devices {
		if len(device.Addresses) == 0 {
			continue
		}
		fmt.Println("Interface: ", device.Name)
		for _, address := range device.Addresses {
			fmt.Println("   IP address:  ", address.IP)
			fmt.Println("   Subnet mask: ", address.Netmask)
		}
		fmt.Println("")
	}
}

func main() {
	var err error
	var handle *pcap.Handle
	var data []Data
	var xMax int
	ch := make(chan Data)

	dev := flag.String("interface", "", "Choose an interface for online processing")
	file := flag.String("file", "", "Choose a file for offline processing")
	filter := flag.String("filter", "", "Set a specific filter")
	lst := flag.Bool("list_interfaces", false, "List available interfaces")
	vers := flag.Bool("version", false, "Show version")
	help := flag.Bool("help", false, "Show help")
	num := flag.Uint("count", 10, "Number of packets to process")
	output := flag.String("output", "image.png", "Name of the resulting image")
	flag.Parse()

	if flag.NFlag() < 1 {
		flag.PrintDefaults()
		return
	}

	if *lst {
		availableInterfaces()
		return
	}

	if *vers {
		fmt.Println("Version:", Version)
		return
	}

	if *help {
		flag.PrintDefaults()
		return
	}

	if len(*dev) > 0 {
		handle, err = pcap.OpenLive(*dev, 4096, true, pcap.BlockForever)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	} else if len(*file) > 0 {
		handle, err = pcap.OpenOffline(*file)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Source is missing")
		return
	}
	defer handle.Close()

	if len(*filter) != 0 {
		err = handle.SetBPFFilter(*filter)
		if err != nil {
			log.Fatal(err, "\tInvalid filter: ", *filter)
			os.Exit(1)
		}
	}

	packetSource := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	packetSource.DecodeOptions = gopacket.Lazy

	go handlePackets(packetSource, *num, ch)

	for i := range ch {
		data = append(data, i)
		if xMax < len(i.payload) {
			xMax = len(i.payload)
		}
	}
	xMax++

	img := image.NewNRGBA(image.Rect(0, 0, xMax/3+1, len(data)))

	for i := range data {
		var j int
		for j = 0; j+3 <= len(data[i].payload); j += 3 {
			img.Set(j/3, i, color.NRGBA{
				R: uint8(data[i].payload[j]),
				G: uint8(data[i].payload[j+1]),
				B: uint8(data[i].payload[j+2]),
				A: 255})
		}
		switch len(data[i].payload) - j {
		case 2:
			img.Set(j/3, i, color.NRGBA{
				R: uint8(data[i].payload[j]),
				G: uint8(data[i].payload[j+1]),
				B: uint8(0),
				A: 255})
		case 1:
			img.Set(j/3, i, color.NRGBA{
				R: uint8(data[i].payload[j]),
				G: uint8(0),
				B: uint8(0),
				A: 255})
		default:
		}
	}

	f, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
