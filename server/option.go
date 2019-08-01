package server

type Option interface {
	Apply(*Server)
}

//                   /---> TcpEnabled
// ServerModule --->
//					 \---> UdpEnabled

// Port ---> Protocol


