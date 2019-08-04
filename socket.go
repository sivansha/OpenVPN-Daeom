package main

import "net"
import "fmt"
import "bufio"
import "strings"

type OpenVpnProcess struct {
	Ip   string
	Port int
}

func main() {
	processList := []OpenVpnProcess{
		{Ip: "127.0.0.1", Port: 11940},
		{Ip: "127.0.0.1", Port: 11941},
	}
	commonName := "07d1ccc455a21c2d5ac6068d4af727ca"

	killCount := 0
	finishedProcesses := make(chan int)
	successfulProcesses := make(chan int)
	for _, p := range processList {
		// XXX make go routine here with channel to get back whether or not it
		// succeeded instead of running them serially...
		// wasDisconnected := kill(&p, commonName)
		// if wasDisconnected {
		// killCount++	
		}
		go func() { 
			wasDisconnected := kill(&p, commonName)
			if wasDisconnected {
			++successfulProcesses
			}
		++finishedProcesses
		}()
	}

	fmt.Println(fmt.Sprintf("We disconnected %d clients with CN \"%s\"!", killCount, commonName))
}

func kill(p *OpenVpnProcess, commonName string) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", p.Ip, p.Port))
	if err != nil {
		// unable to connect, no matter, maybe the process is temporary away,
		// so no need to disconnect clients there ;-)
		return false
	}

	defer conn.Close()

	// turn off live OpenVPN log that can confuse our output parsing
	fmt.Fprintf(conn, fmt.Sprint("log off\n"))

	// read until we find ^END ^SUCCESS or ^ERROR
	// this is very ugly...
	reader := bufio.NewReader(conn)
	txt, _ := reader.ReadString('\n')
	for 0 != strings.Index(txt, "END") && 0 != strings.Index(txt, "SUCCESS") && 0 != strings.Index(txt, "END") {
		txt, _ = reader.ReadString('\n')
	}

	wasDisconnected := false
	fmt.Fprintf(conn, fmt.Sprintf("kill %s\n", commonName))
	// we keep reading until we read either SUCCESS or ERROR, a bit ugly, we
	// should already not have anything to read anymore...
	text, _ := reader.ReadString('\n')
	if 0 == strings.Index(text, "SUCCESS") {
		wasDisconnected = true
	}

	fmt.Fprintf(conn, "quit\n")

	return wasDisconnected
}

func status(p *OpenVpnProcess) {
}
