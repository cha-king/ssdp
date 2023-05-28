package ssdp

import "net"

var ssdpUdpAddr = &net.UDPAddr{IP: net.IPv4(239, 255, 255, 250), Port: 1900}

const All = "ssdp:all"
