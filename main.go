package main

import (
	"flag"
	"regexp"
	"time"
	"net"

	log "github.com/cihub/seelog"
	"github.com/stampzilla/stampzilla-go/nodes/basenode"
	"github.com/stampzilla/stampzilla-go/pkg/notifier"
	"github.com/stampzilla/stampzilla-go/protocol"
)

// MAIN - This is run when the init function is done

var notify *notifier.Notify

var ip = flag.String("ip", "10.21.10.158", "Ip to the device")

func main() { /*{{{*/
	log.Info("Starting stamp-presence node")

	// Parse all commandline arguments, host and port parameters are added in the basenode init function
	flag.Parse()

	//Get a config with the correct parameters
	config := basenode.NewConfig()

	//Activate the config
	basenode.SetConfig(config)

	node := protocol.NewNode("stamp-presence")

	//Start communication with the server
	connection := basenode.Connect()
	notify = notifier.New(connection)
	notify.SetSource(node)

	// Thit worker keeps track on our connection state, if we are connected or not
	go monitorState(node, connection)

	state := NewState()
	node.SetState(state)

	// This worker recives all incomming commands
	go serverRecv(node, connection)

	go socketConnection(state, node, connection)
	select {}
} /*}}}*/

// WORKER that monitors the current connection state
func monitorState(node *protocol.Node, connection basenode.Connection) {
	for s := range connection.State() {
		switch s {
		case basenode.ConnectionStateConnected:
			connection.Send(node.Node())
		case basenode.ConnectionStateDisconnected:
		}
	}
}

// WORKER that recives all incomming commands
func serverRecv(node *protocol.Node, connection basenode.Connection) {
	for d := range connection.Receive() {
		processCommand(node, connection, d)
	}
}

// THis is called on each incomming command
func processCommand(node *protocol.Node, connection basenode.Connection, cmd protocol.Command) {
}

func socketConnection(state *State, node *protocol.Node, connection basenode.Connection) {
	r, _ := regexp.Compile("<([01]+|DOOR)>")

	for {
		<-time.After(time.Second)

		log.Infof("Connecting to %s", *ip)

		s, err := net.Dial("tcp", *ip+":23")
		if err != nil {
			log.Error("Failed to open port: ", err)
			continue
		}

		log.Info("Connected to device")

		var buff string

	readLoop:
		for {

			// Read data
			buf := make([]byte, 128)
			n, err := s.Read(buf)
			if err != nil {
				log.Error(err)
				break readLoop
			}

			buff += string(buf[:n])

			res := r.FindAllStringSubmatchIndex(buff, -1)
			for _, match := range res {
				data := buff[match[0]+1 : match[1]-1]

				log.Infof("Data: %#v", data)

				if data == "DOOR" {
					state.Door = true
					connection.Send(node.Node())
					<-time.After(time.Second)
					state.Door = false
					connection.Send(node.Node())
				} else {

					for key, val := range data {
						switch key {
						case 0:
							state.Sensor1 = (val == 48)
						case 1:
							state.Sensor2 = (val == 48)
						case 2:
							state.Sensor3 = (val == 48)
						case 3:
							state.Sensor4 = (val == 48)
						}
					}
				}
				connection.Send(node.Node())
			}

			if len(res) > 0 {
				buff = buff[res[len(res)-1][1]:]
			}

		}
	}
}
