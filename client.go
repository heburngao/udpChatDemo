package main
 
import (
    "net"
    "fmt"
    "log"
    "os"
    "time"
    "strings"
)
/*
//发送信息
func sender(conn net.Conn) {
    words := "Hello Server!"
    conn.Write([]byte(words))
    fmt.Println("send over")
 
    //接收服务端反馈
    buffer := make([]byte, 2048)
 
    n, err := conn.Read(buffer)
    if err != nil {
        Log(conn.RemoteAddr().String(), "waiting server back msg error: ", err)
        return
    }
    Log(conn.RemoteAddr().String(), "receive server back msg: ", string(buffer[:n]))
 
}

*/



//日志
func Log(v ...interface{}) {
    log.Println(v...)
}
 
//==========================================

var PT = fmt.Printf //Println
//===============client============
type Client struct {
	conn * net.UDPConn
	gkey bool // 判断用户是否退出
	userID int
	userName string
	chanWrite chan []byte
	chanRead chan []byte
}

//==========================================
func (c *Client) sndMsg2 (sid int,msg string){
    str := fmt.Sprintf("###%d##%d##%s##%s###", sid, c.userID, c.userName, msg)
    _,err := c.conn.Write([]byte(str))
    checkError(err,"sndMsg2")
}

//==========================================
func (c * Client ) sndMsg() {
	for c.gkey {
		msg := <- c.chanWrite
		str := fmt.Sprintf("###2##%d##%s##%s###", c.userID, c.userName , msg)
		PT(" sndMsg len: %d  msg: %s  \r\n" , len(str), str)
		_,err := c.conn.Write([]byte(str))
		checkError(err,"sndMsg")
	}

}

func (c * Client) inputMsg(){
	var msg string
	for c.gkey {
		fmt.Println(" waitting for inputMsg ..." )
		_,err := fmt.Scanln(&msg)
		PT("inputMsg:: enter Msg: %s ", msg)
		checkError(err, "inputMsg")
		if msg == ":quit" {
			c.gkey = false
		}else{
			c.chanWrite <- []byte (encodeMsg(msg))
		}
	}
}
//==========================================
func (c * Client) readOut(){
	for c.gkey{
		msg := <- c.chanRead
		fmt.Println("<- readOut after rcvMsg : ", string(msg), " \r\n" )
	}
}

func (c * Client) rcvMsg () {
	var buf [1024]byte
	for c.gkey {
	   n,err := c.conn.Read(buf[0:])
	   checkError(err,"rcvMsg")
	   c.chanRead <- buf[0:n]
	}

}

//==========================================


//===============client============


func encodeMsg(msg string) (string){
	var str string =  strings.Join(strings.Split(strings.Join(strings.Split(msg,"\\"),"\\\\"),"#"),"\\#")
	fmt.Println("encodeMsg : " , str ) 
	return str
}

func nowtime()string{
	return time.Now().String()
}


func checkError(err error, funcName string) {
		if err != nil {
			fmt.Fprintf(os.Stderr,"Fatal  error: %s ---- in func: %s " , err.Error() , funcName)
			os.Exit(1)
		}

}


//==========================================

func main() {
//    server := "127.0.0.1:1024"
//===udp========

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s | host:port" , os.Args[0]," \r\n")
		os.Exit(1)
	}
	server := os.Args[1]
	udpAddr , err := net.ResolveUDPAddr("udp4",server)
	checkError(err, " 1 main")
	PT(" 1 udpServerAddr:%s\r\n " , udpAddr)	
	
	var c Client 
	c.gkey = true
	c.chanWrite = make(chan []byte)
	c.chanRead = make(chan []byte)
	
	PT(" 2 enter userID , udpAddr: %s\r\n " , udpAddr)
	_,err = fmt.Scanln(&c.userID)
	PT(" ok ,userID is : %d\r\n " , c.userID)
	checkError(err, " 2 main")
	
	PT("3 enter userName , udpAddr:%s\r\n " , udpAddr)
	_,err = fmt.Scanln(&c.userName)
	PT(" ok , userName is : %s\r\n " , c.userName)
	checkError(err, " 3 main")

	PT("4 udpAddr:%s\r\n " , udpAddr)
	c.conn, err = net.DialUDP("udp",nil, udpAddr)
	checkError(err, " 4 main")
	defer c.conn.Close()
	PT("5 udpAddr:%s\r\n " , udpAddr)
    //发送进入聊天室消息,类型1，###1##uid##uName##进入聊天室###
    //messagestr := fmt.Sprintf("###1##%d##%s###", c.userID, c.userName)
    //_,err = c.conn.Write([]byte(messagestr))
    //checkError(err)
	PT("发送进入聊天室消息,类型1，###1##uid##uName##进入聊天室###")

   	c.sndMsg2(1,c.userName + "进入聊天室")

	go c.readOut()
	go c.rcvMsg()

	go c.sndMsg()
	
	c.inputMsg()

	c.sndMsg2(3,c.userName + "离开聊天室")
	fmt.Println("退出成功")
	os.Exit(0)
	

//==========================================
//====tcp======
   // server := "183.192.30.60:1024"
/*    tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
 
    fmt.Println("connection success")
    sender(conn)
*/

}

