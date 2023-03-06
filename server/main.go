package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)

var clients = make(map[string]net.Conn)

func main() {
    // 创建监听端口
    l, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        return
    }
    defer l.Close()
    fmt.Println("Listening on 0.0.0.0:8080...")

    // 循环处理客户端请求
    for {
        // 接受客户端连接
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting:", err.Error())
            continue
        }

        // 读取客户端发送的信息，并将其加入到客户端列表中
        clientInfo, err := bufio.NewReader(conn).ReadString('\n')
        if err != nil {
            fmt.Println("Error reading info:", err.Error())
            conn.Close()
            continue
        }
        clients[strings.Trim(clientInfo, "\n")] = conn
        fmt.Println(clientInfo)

        // 启用goroutine处理客户端请求
        go handleClient(conn)
    }
}

func handleClient(conn net.Conn) {
    defer conn.Close()

    // 读取服务端命令并发送到客户端
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("Enter command to execute on client: ")
        cmdStr, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Error reading command:", err.Error())
            continue
        }

        // 将命令发送到所有客户端
        for _, clientConn := range clients {
            _, err = clientConn.Write([]byte(cmdStr))
            if err != nil {
                fmt.Println("Error sending command:", err.Error())
            }

            // 读取并输出客户端响应
            output := make([]byte, 1024)
            n, err := clientConn.Read(output)
            if err != nil {
                fmt.Println("Error receiving output:", err.Error())
                continue
            }
            fmt.Println("Client output:", string(output[:n]))
        }
    }
}
