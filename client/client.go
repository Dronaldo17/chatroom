package main

// #include <stdio.h>
// #include <stdlib.h>
/*
void myscan(void *p) {
	char *str = (char *)p;
	gets(str);
}
*/
import "C"

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"unsafe"
)

var str string

var msg = make([]byte, 1024)

//定义channel ,define channel
var c chan int

//批量建立变量,batch define variable
var (
	host   = "192.168.0.90"    //默认IP地址, the default IP address
	port   = "8080"            //默认端口号,the default port number
	remote = host + ":" + port //远程地址,remote address
	/*
			IP地址配置文件，格式172.0.0.1：8080
		    IP address configuration file,format 172.0.0.1:8080
	*/
	fName = "config.txt"
)

func main() {

	readFile(fName) //读取配置文件,read the configuration file
	/*
		提示即将连接服务器
		Tip "the server will be connect"
	*/
	fmt.Println("You will connect the server : ", remote)

	ColorEdit(12) //红色,12 represents red color
	for {
		//输入昵称,enter nickname

		fmt.Printf("Enter a nicename:")
		/*
			声明一个大小为1024的byte型变量strv
			declare a variable(strv) which size is 1024 and belong to byte type
		*/
		var strv [1024]byte
		/*
			调用C语言gets函数获取用户输入的信息，将strv首位指针作为参数传递给myscan()
			call gets() function which belongs to C programing language
			aim to get the information of the user input.
			transfer a parameter to myscan(),the parameter is the first pointer of strv
		*/
		C.myscan(unsafe.Pointer(&strv[0]))
		/*
			声明一个整数变量L，用于将要取出的strv的长度
			declare a variable(L) aim to get the length of strv
		*/
		l := 1024
		/*
			从strv的末尾开始判断ASCII值，如果为0，则长度L减1
			judge ASCII value from the end of strv,if it is 0,the length decrease 1
		*/
		for {
			if strv[l-1] == 0 {
				l = l - 1
			} else {
				break
			}
		}
		/*
			类型转换，取strv的首位到L的长度转换为string
			type conversion,get a length from begining to L aim to convert to the type of string
		*/
		str = string(strv[:l])

		if str == "" {
			fmt.Printf("Please input a right nicename!\n")
			continue
		}
		if str == "query" {
			fmt.Printf("query is a keyword , please input a right nicename!\n")
			continue
		}
		break
	}
	/*
		与远程服务器建立连接，返回连接con
		a connection with the remote server，return connection value con
	*/
	con, err := net.Dial("tcp", remote)
	defer con.Close() //延迟关闭连接,delay to close connection

	if err != nil {
		ColorEdit(12) //红色,12 represents red color
		fmt.Println("Server not found.")
		os.Exit(-1)
	}
	/*
		将昵称发送给服务器端
		send the nickname to server
	*/
	in, err := con.Write([]byte(str))
	if err != nil {
		ColorEdit(12) //红色,12 represents red color
		fmt.Printf("Error when send to server: %d\n", in)
		os.Exit(0)
	}

	ColorEdit(11) //青色,11 represents cyan color
	fmt.Println("Connection OK.")

	length, err := con.Read(msg) //接收欢迎信息,receive welcome message
	if err != nil {
		ColorEdit(12) //红色,12 represents red color
		fmt.Printf("Error when read from server.\n")
		os.Exit(0)
	}
	str = string(msg[0:length]) //格式转换,format conversion
	fmt.Println(str)            //将信息显示到屏幕,display the message on screen

	ColorEdit(10)      //绿色,10 represents green color
	go receiveMsg(con) //接收信息goroutine,receive the message
	go sendMsg(con)    //发送信息goroutine,send the message

	<-c
	<-c

}

//接收信息函数,receive message function
func receiveMsg(con net.Conn) {
	for {
		/*
			接收服务器的响应信息
			receive the response message of server
		*/
		length, err := con.Read(msg)
		if err != nil {
			ColorEdit(12) //红色，12 represents red color
			fmt.Printf("Error when read from server.\n")
			os.Exit(0)
		}
		str = string(msg[0:length]) //格式转换，format conversion
		ColorEdit(11)               //青色，11 represents cyan color
		fmt.Println(str)            //将信息显示到屏幕，display the message on screen
		ColorEdit(10)               //绿色，10 represents green color
	}
	c <- 1

}

//发送信息函数，send message function
func sendMsg(con net.Conn) {
	for {
		/*
			声明一个大小为1024的byte型变量strv
			declare a variable(strv) which size is 1024 and belong to byte type
		*/
		var strv [1024]byte
		/*
			调用C语言函数从屏幕输入，将strv首位指针作为参数传递给myscan()
			call gets() function which belongs to C programing language
			aim to get the information of the user input.
			transfer a parameter to myscan(),the parameter is the first pointer of strv
		*/
		C.myscan(unsafe.Pointer(&strv[0]))
		/*
			声明一个整数变量L，用于将要取出的strv的长度
			declare a variable(L) aim to get the length of strv
		*/
		l := 1024

		/*
			从strv的末尾开始判断ASCII值，如果为0，则长度L减1
			judge ASCII value from the end of strv,if it is 0,the length decrease 1
		*/
		for {
			if strv[l-1] == 0 {
				l = l - 1
			} else {
				break
			}
		}
		/*
			类型转换，取strv的首位到L的长度转换为string
			type conversion,get a length from begining to L aim to convert to the type of string
		*/
		str = string(strv[:l])

		if str == "quit" { //如果输入quit则结束连接
			ColorEdit(12) //红色,12 represents red color
			fmt.Println("Communication terminated.")
			os.Exit(1)
		}

		in, err := con.Write([]byte(str)) //将输入的内容发送给服务器端,send the input message to server
		if err != nil {
			ColorEdit(12) //红色,12 represents red color
			fmt.Printf("Error when send to server: %d\n", in)
			os.Exit(0)
		}
	}
	c <- 1
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
