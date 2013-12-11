package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

//批量建立变量,batch define variable
var (
	host   = "192.168.0.90"    //默认IP地址,the default IP address
	port   = "8080"            //默认端口号,the default port number
	remote = host + ":" + port //远程地址,remote address
	fName  = "config.txt"
	/*
			IP地址配置文件，格式172.0.0.1：8080
		    IP address configuration file,format 172.0.0.1:8080
	*/
)

//建立结构体,define struct
type ClientInfo struct {
	NiceName string   //昵称,nickname
	Conn     net.Conn //连接,connect
}

//建立map,define map
var clientDB map[net.Conn]ClientInfo

//定义全局变量,define global variable
var clientnum = 0 //初始化用户数量,initialize the number of users

func main() {

	readFile(fName) //读取配置文件信息,read the configuration file

	fmt.Println("Server address : ", remote) //输出服务器地址,output server address
	ColorEdit(12)                            //红色,12 represents red color
	fmt.Println("Initiating server... (Ctrl-C to stop)")

	clientDB = make(map[net.Conn]ClientInfo) //初始化map,initialize map
	/*
		监听远程地址的TCP协议
		monitor the TCP protocol of remote address
	*/
	lis, err := net.Listen("tcp", remote)
	defer lis.Close() //延迟处理监听关闭,delay the monitor close

	/*
		如果错误不为空，则输出错误信息并退出
		if the error is not void,output the error message and exit
	*/
	if err != nil {
		ColorEdit(12)
		fmt.Println("Error when listen: ", remote)
		os.Exit(-1)
	}

	//循环等待接收客户端,Loop waiting to receive client
	for {
		/*
			等待接收连接，并返回这个连接conn
			Waiting to receive connection and return the connection(conn)
		*/
		conn, err := lis.Accept()

		/*
			如果错误不为空，则输出错误信息并退出
			if the error is not void,output the error message and exit
		*/
		if err != nil {
			ColorEdit(12)
			fmt.Println("Error accepting client: ", err.Error())
			os.Exit(0)
		}

		//接收用户输入的昵称,receive a nickname which is the user input
		var (
			ndata    = make([]byte, 1024) //接收到的数据,receive the message which has arrived
			nicename string
		)
		/*
			读取客户端信息，返回信息长度和错误信息
			read the user information, return the length of message and error message
		*/
		length, err := conn.Read(ndata)
		nicename = string(ndata[0:length]) //类型转换,type conversion

		clientDB[conn] = ClientInfo{nicename, conn} //将用户信息放入MAP中,put the user message into Map

		clientnum++

		welmsg := fmt.Sprintf("welcome %s(%s) join.\nnow %d clients online", nicename, conn.RemoteAddr(), clientnum) //定义欢迎信息
		sendMsgAll(welmsg, nil)                                                                                      //向所有用户发送欢迎信息

		go receiveMsg(conn) //建立针对此用户的接收goroutine,creat a goroutine aim to receive

	}
}

//接收信息函数,receive message function
func receiveMsg(con net.Conn) {
	/*
		获取所连接客户端的信息，返回客户端的实例
		obtain the whole message of user,return the client instance
	*/
	client, ok := clientDB[con]

	if !ok {
		fmt.Println("Did not find client.")
	}

	var (
		data = make([]byte, 1024) //接收到的数据,receive the message which has arrived
		res  string
	)
	ColorEdit(11) //青色,11 represents cyan color
	fmt.Println("New connection: ", client.NiceName, "(", con.RemoteAddr(), ")")
	fmt.Printf("%d clients online!\n", clientnum)
	ColorEdit(10) //绿色,10 represents green color
	for {
		/*
			读取客户端信息，返回信息长度和错误信息
			read the user information, return the length of message and error message
		*/
		length, err := con.Read(data)

		/*
			如果报错，则用户已经掉线或退出，向所有用户发送此通告
			if error,the user offline or exit. sent the message to all user.
		*/
		if err != nil {
			ColorEdit(11)                                     //青色,11 represents cyan color
			fmt.Printf("Client %v quit.\n", con.RemoteAddr()) //输出用户退出信息,output the exit message of user
			clientnum--                                       //减少一个连接用户,reduce a connected user
			fmt.Printf("%d clients online!\n", clientnum)     //输出目前在线用户数量,output the number of online user
			delete(clientDB, con)                             //从MAP中删除此退出的用户 delete the user from Map
			sysMsg := fmt.Sprintf("Client %s(%v) quit.\nnow %d clients online", client.NiceName, con.RemoteAddr(), clientnum)
			/*
				将用户退出及现有用户数量的信息发送给其他全部用户
				the exit and online message of user sent to the other user
			*/
			sendMsgAll(sysMsg, con)
			con.Close()   //关闭此用户的连接,close user's connect
			ColorEdit(10) //绿色,10 represents green color
			return
		}
		res = string(data[0:length]) //类型转换,type conversion

		//判断是否发送查询命令,whether to send a query command
		if res == "query" {
			queryMsg(con) //执行查询，并发送查询结果,carry out query and send the query result
			continue
		}
		/*
			将用户发来的信息显示到屏幕上
			display the message which is user send on screen
		*/
		fmt.Printf("%s said: %s\n", client.NiceName, res)
		//向所有用户发送此用户所写信息,send the write message of current user to all user
		saidMsg := fmt.Sprintf("%s said: %s", client.NiceName, res) //声明一个用户所发信息的变量,declare a variable which is the send message of user
		/*
			向所有用户发送此用户所发信息
			send the write message of current user to all user
		*/
		sendMsgAll(saidMsg, con)
	}
}

//群发信息函数,mass information function
func sendMsgAll(str string, con net.Conn) {
	for _, v := range clientDB { //遍历map,traversal map
		/*
			如果遍历所得用户端为当前发送信息用户端则不对自己发送信息
			if the message belongs to itself ,doesn't send the message to itself
		*/
		if v.Conn == con {
			continue
		}
		if v.Conn != nil {
			v.Conn.Write([]byte(str)) //发送信息,send the message
		}
	}
}

//发送查询信息,send query information
func queryMsg(con net.Conn) {
	str := fmt.Sprintf("%d online clients :\n you", clientnum)
	for _, v := range clientDB { //遍历map,traversal map
		/*
			如果遍历所得用户端为当前发送信息用户端则不做记录
			if traversal result is itself,doesn't record
		*/
		if v.Conn == con {
			continue
		}
		if v.Conn != nil {
			str = str + " , " + v.NiceName //记录到查询结果变量str中,save into the query result variables(str)
		}
	}
	con.Write([]byte(str)) //发送查询结果,send a query results
}

/*
   控制台彩色输出，1深蓝，2深绿，3深青，10绿色，11青色，12红色
   The console color output,1 deep blue,2 deep green,3 deep cyan,10 green,11 cyan,12 red
*/
func ColorEdit(i int) {
	kernel32, _ := syscall.LoadDLL("kernel32.dll") //调用dll,call dll
	defer kernel32.Release()                       //延迟释放,delay release
	/*
		找到dll中的方法，并返回此方法
		find the method in the dll and return the method
	*/
	proc, _ := kernel32.FindProc("SetConsoleTextAttribute")
	proc.Call(uintptr(syscall.Stdout), uintptr(i)) //执行此方法,carry out the method
}

//读取配置文件,read configuration files
func readFile(filename string) {
	fin, err := os.Open(filename) //打开文件,opened-file
	defer fin.Close()             //延迟关闭,delay close
	if err != nil {               //是否报错,whether or not an error
		ColorEdit(12)              //红色,12 represents red color
		fmt.Println(filename, err) //输出错误信息,output the error information
		ColorEdit(11)              //青色,11 represents cyan color
		return
	}
	/*
		声明1024位byte型数组切片
		declare a array slice which size is 1024 and belong to byte type
	*/
	buf := make([]byte, 1024)
	/*
		读取文件信息，并存入buf，返回信息长度
		read file information and save to buf,return the length of message
	*/
	n, _ := fin.Read(buf)
	if 0 == n { //如果信息长度为0,if the length of message is 0
		ColorEdit(12) //红色,12 represents red color
		fmt.Println("The IPCONFIG file is null.")
		ColorEdit(11) //青色,11 represents cyan color
		return
	}
	/*
		将读取的信息赋值给远程地址变量
		have been read message assign to the remote address variable
	*/
	remote = string(buf[:n])
}
