package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/tr3ee/librjsocks"
)

var (
	username, password, device string
	listdev, verbose           bool
)

func init() {
	flag.StringVar(&username, "u", "", "username")
	flag.StringVar(&password, "p", "", "password")
	flag.StringVar(&device, "d", "", "determine the network device")
	flag.BoolVar(&listdev, "l", false, "list all network devices")
	flag.BoolVar(&verbose, "v", false, "verbose mode")
	flag.Parse()
}

func listalldevs() error {
	devs, err := librjsocks.FindAllAdapters()
	if err != nil {
		return err
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"adapter", "device", "desc", "mac"})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	for _, dev := range devs {
		table.Append([]string{dev.AdapterName, dev.DeviceName, dev.DeviceDesc, dev.Mac.String()})
	}
	table.Render()
	return nil
}

func validate() error {
	if len(username) == 0 {
		return errors.New("empty username")
	}
	if len(password) == 0 {
		return errors.New("empty password")
	}
	if len(device) == 0 {
		return errors.New("network device is required")
	}
	return nil
}

func main() {
	if listdev {
		if err := listalldevs(); err != nil {
			log.Println(err)
		}
		return
	}
	if err := validate(); err != nil {
		flag.Usage()
		log.Fatal(err)
	}
	var dev *librjsocks.NwAdapterInfo
	devs, err := librjsocks.FindAllAdapters()
	if err != nil {
		log.Fatal(err)
	}
	for _, d := range devs {
		if d.AdapterName == device {
			dev = &d
			break
		}
	}
	if dev == nil {
		log.Fatalf("no such device: %s", device)
	}
	srv, err := librjsocks.NewService(username, password, dev)
	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan librjsocks.Event, 16)
	srv.Notify(ch)
	go func() {
		for e := range ch {
			if verbose {
				fmt.Println(e.String())
			} else {
				if e == librjsocks.EventSuccess {
					fmt.Println("SUCCESS!")
				}
			}
		}
	}()
	if err := srv.Run(); err != nil {
		log.Println(err)
		return
	}
}
