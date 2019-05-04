package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/phalaaxx/milter"
)

var (
	protocol = flag.String("protocol", "tcp", "unix or tcp")
	address  = flag.String("address", "127.0.0.1:47256", "host:port (for tcp) or path to socket")
	target   = flag.String("target", ".", "target directory to write e-mails to")
	spam     = flag.Bool("spam", false, "when given, separate mails based on the X-Spam header into spam and ham directories")
)

func main() {
	flag.Parse()

	if *protocol != "unix" && *protocol != "tcp" {
		log.Fatalln("Invalid protocol name selected. Must be tcp or unix.")
	}

	if len(*address) == 0 {
		log.Fatalln("No address given.")
	}

	info, err := os.Stat(*target)
	if err != nil {
		log.Fatalf("Target directory %s cannot be found: %v", *target, err)
	}

	if !info.IsDir() {
		log.Fatalf("Target %s is not a valid directory.", *target)
	}

	if *protocol == "unix" {
		os.Remove(*address)
	}

	log.Printf("Listening on %s:%s ...", *protocol, *address)
	socket, err := net.Listen(*protocol, *address)
	if err != nil {
		log.Fatal(err)
	}
	defer socket.Close()

	if *protocol == "unix" {
		if err := os.Chmod(*address, 0660); err != nil {
			log.Fatal(err)
		}
		defer os.Remove(*address)
	}

	mkMilter := func() (milter.Milter, milter.OptAction, milter.OptProtocol) {
		m := newArchiveMilter(*target, *spam)
		optAction := milter.OptAction(0)
		optProtocol := milter.OptNoConnect | milter.OptNoHelo | milter.OptNoMailFrom | milter.OptNoRcptTo

		return m, optAction, optProtocol
	}

	log.Println("Accepting connections now.")
	err = milter.RunServer(socket, mkMilter)
	if err != nil {
		log.Fatal(err)
	}
}
