package  main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func CreateConn(addr string,debug bool)  {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", addr)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Connected server failed :",err.Error())
		return
	}

	fmt.Println("Connected server:" + conn.RemoteAddr().String())

	defer func() {
		_ = conn.Close()
	}()

	buffer := make([]byte, 1024)
	for{
		_,err = conn.Write([]byte(fmt.Sprintf("ping %d",time.Now().Unix())))
		if err != nil {
			fmt.Println("Disconnected from server:" + conn.RemoteAddr().String())
			return
		}

		n,err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Read error:",err.Error(), "Disconnected from server:" + conn.RemoteAddr().String())
			return
		}

		msg := string(buffer[0:n])
		if debug{
			fmt.Println("Received msg: ",msg)
		}

		time.Sleep(time.Second * 1)
	}
}

func Server(addr string,debug bool){
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", addr)
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("tcp监听端口异常，程序退出！"+addr)
		return
	}

	defer func(){
		_ = tcpListener.Close()
	}()

	fmt.Println("Start server ... ")

	//循环接收客户端的连接，创建一个协程具体去处理连接
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("A client connected["+ addr +"] :" + tcpConn.RemoteAddr().String())

		go func(){
			defer func(){
				_ = tcpConn.Close()
			}()

			buffer := make([]byte,1024)
			for{
				n,err := tcpConn.Read(buffer)
				if err != nil {
					fmt.Println("Read error:",err.Error(), "Disconnected client:" + tcpConn.RemoteAddr().String())
					return
				}

				msg := string(buffer[0:n])
				if debug{
					fmt.Println("Received msg: ",msg)
				}

				_,err = tcpConn.Write([]byte(fmt.Sprintf("pong %d",time.Now().Unix())))
				if err != nil {
					fmt.Println("Disconnected client:" + tcpConn.RemoteAddr().String())
					return
				}
				//time.Sleep(time.Second * 1)
			}
		}()
	}

}


func main() {


	var server bool
	flag.BoolVar(&server, "s",false,"Run as server mode")

	var client bool
	flag.BoolVar(&client, "c",false,"Run as client mode")

	var ip string
	flag.StringVar(&ip, "i", "0.0.0.0", "Specify ip to use.  defaults to 0.0.0.0")

	var port string
	flag.StringVar(&port, "p", "8000", "Specify port to use.  defaults to 8000")

	var num int
	flag.IntVar(&num, "n", 10, "Total connection count.  defaults to 10")

	var debug bool
	flag.BoolVar(&debug, "d",false,"Debug with log info")


	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("    example: %s -m server -i 192.168.56.200 -p 8000 -c 2000\n", os.Args[0])

		flag.PrintDefaults()
	}

	flag.Parse()

	if server {
		Server(fmt.Sprintf("%s:%s", ip, port),debug)
	} else if client {
		for i:= 0;i< num;i++ {
			go CreateConn(fmt.Sprintf("%s:%s", ip, port),debug)
		}
	} else {
		fmt.Println("Please specify mode!")
		flag.Usage()
		os.Exit(1)
	}
	select {}
}
