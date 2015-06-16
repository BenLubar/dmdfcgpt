//go:generate protoc --go_out=. network.proto

package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"io"
	"net"
	"strconv"
	"time"

	"gopkg.in/tomb.v2"

	"github.com/golang/protobuf/proto"
	"github.com/huin/goupnp"
	"github.com/huin/goupnp/dcps/internetgateway1"
)

var flagNoUPnP = flag.Bool("no-upnp", false, "Disable Universal Plug and Play. UPnP is used to automatically forward a port on your router and get your external IP address.")
var flagPort = flag.Int("host-port", 0, "Port number for hosting a server. 0 chooses a random unused port number.")

var ExternalAddr = make(chan string)

func network() error {
	type networkServiceClient interface {
		GetExternalIPAddress() (ip string, err error)
		AddPortMapping(externalHost string, externalPort uint16, protocol string, internalPort uint16, internalClient string, enabled bool, portMappingDescription string, leaseDuration uint32) error
		DeletePortMapping(externalHost string, externalPort uint16, protocol string) error
		GetNATRSIPStatus() (sip bool, nat bool, err error)
	}

	var serviceClient networkServiceClient
	serviceClientCh := make(chan networkServiceClient, 1)
	var externalAddr string
	var externalAddrCh chan string

	// find a UPnP device, but don't hang the game if we exit before the
	// device is found.
	go func() {
		if *flagNoUPnP {
			serviceClientCh <- nil
			return
		}

		var serviceClient networkServiceClient

		if devices, err := goupnp.DiscoverDevices(internetgateway1.URN_WANConnectionDevice_1); err == nil {
			for _, device := range devices {
				if serviceClient != nil {
					break
				}

				if device.Root == nil {
					continue
				}

				device.Root.Device.VisitServices(func(service *goupnp.Service) {
					if serviceClient != nil {
						return
					}

					sc := goupnp.ServiceClient{service.NewSOAPClient(), device.Root, device.Location, service}
					sc.SOAPClient.HTTPClient.Timeout = 10 * time.Second

					var found networkServiceClient
					switch sc.Service.ServiceType {
					case internetgateway1.URN_WANIPConnection_1:
						found = &internetgateway1.WANIPConnection1{sc}
					case internetgateway1.URN_WANPPPConnection_1:
						found = &internetgateway1.WANPPPConnection1{sc}
					default:
						return
					}

					if _, nat, err := found.GetNATRSIPStatus(); err == nil && nat {
						serviceClient = found
					}
				})
			}
		}
		serviceClientCh <- serviceClient
	}()

	accept := make(chan net.Conn)

	for {
		select {
		case <-Tomb.Dying():
			return nil

		case serviceClient = <-serviceClientCh:
			serviceClientCh = nil

			var ip string
			ifaces, err := net.Interfaces()
			if err != nil {
				return err
			}
		ifaceloop:
			for _, iface := range ifaces {
				if iface.Flags&net.FlagUp != net.FlagUp {
					continue
				}

				if iface.Flags&net.FlagLoopback == net.FlagLoopback {
					continue
				}

				addrs, _ := iface.Addrs()
				for _, addr := range addrs {
					if ipnet, ok := addr.(*net.IPNet); ok {
						ip = ipnet.IP.String()
						if ipnet.IP.To4() != nil {
							break ifaceloop
						}
					} else if ipaddr, ok := addr.(*net.IPAddr); ok {
						ip = ipaddr.IP.String()
						if ipaddr.IP.To4() != nil {
							break ifaceloop
						}
					}
				}
			}

			// init
			listener, err := net.Listen("tcp", net.JoinHostPort(ip, strconv.Itoa(*flagPort)))
			if err != nil {
				return err
			}
			defer listener.Close()

			Tomb.Go(func() error {
				for {
					conn, err := listener.Accept()
					if err != nil {
						if Tomb.Alive() {
							return err
						} else {
							return nil
						}
					}

					select {
					case <-Tomb.Dying():
						conn.Close()
						return nil

					case accept <- conn:
					}
				}
			})

			ip, portStr, err := net.SplitHostPort(listener.Addr().String())
			if err != nil {
				return err
			}
			port, err := net.LookupPort("tcp", portStr)
			if err != nil {
				return err
			}

			if serviceClient == nil {
				// no router. keep the local host/port.
			} else {
				externalIP, err := serviceClient.GetExternalIPAddress()
				if err != nil {
					return err
				}
				err = serviceClient.AddPortMapping(externalIP, uint16(port), "TCP", uint16(port), ip, true, "dmdfcgpt", 0)
				if err != nil {
					return err
				}
				defer serviceClient.DeletePortMapping(externalIP, uint16(port), "TCP")
				ip = externalIP
			}

			externalAddr = net.JoinHostPort(ip, portStr)
			externalAddrCh = ExternalAddr

		case externalAddrCh <- externalAddr:

		case conn := <-accept:
			Tomb.Go(func() error {
				serverHandler(conn)
				return nil
			})
		}
	}
}

func serverHandler(client net.Conn) {
	var connTomb tomb.Tomb

	defer client.Close()

	in, out := make(chan Packet), make(chan Packet)

	connTomb.Go(func() error {
		connTomb.Go(func() error {
			w := bufio.NewWriter(client)
			var length [binary.MaxVarintLen64]byte
			for {
				select {
				case packet := <-out:
					b, err := proto.Marshal(&packet)
					if err != nil {
						return err
					}
					l := length[:binary.PutUvarint(length[:], uint64(len(b)))]
					n, err := w.Write(l)
					if err == nil && n != len(l) {
						err = io.ErrShortWrite
					}
					if err != nil {
						return err
					}

					n, err = w.Write(b)
					if err == nil && n != len(b) {
						err = io.ErrShortWrite
					}
					if err != nil {
						return err
					}

					err = w.Flush()
					if err != nil {
						return err
					}

				case <-connTomb.Dying():
					return nil
				}
			}
		})
		connTomb.Go(func() error {
			r := bufio.NewReader(client)

			for {
				var packet Packet

				l, err := binary.ReadUvarint(r)
				if err != nil {
					return err
				}

				b := make([]byte, l)
				_, err = io.ReadFull(r, b)
				if err != nil {
					return err
				}

				err = proto.Unmarshal(b, &packet)
				if err != nil {
					return err
				}

				select {
				case in <- packet:

				case <-connTomb.Dying():
					return nil
				}
			}
		})

		return nil
	})

	for {
		select {
		case <-connTomb.Dying():
			Log.Println(client.RemoteAddr(), "fatal error:", connTomb.Err())
			return

		case <-Tomb.Dying():
			connTomb.Kill(nil)
			return

		case packet := <-in:
			_ = packet // TODO: handle packet
		}
	}
}
