package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strconv"
    "strings"
)


type ConnectInfo struct {
    ip_addr     net.IPAddr
    udp_port    uint16
    conn_type   string      // udp/tcp
    auto        bool        // autosending
    data_type   bool        // true - string; false - hex
    wait_recv   bool
    file_path   string
}


/**
-ip - ip address of the node
-pt - port
-s - string data type
-d - path to file with data
-a - auto sending
-w - wait answer ???
-f - file with data
*/
// ./net_client -ip 192.168.0.12 -pt 6000 -t udp -s -a
// ./net_client -ip 192.168.0.12 -pt 6000 -t udp -w -a y/n
// ./net_client -ip 11.200.1.1 -pt 12345 -t tcp ...
func main () {
    fmt.Println("\tUDP client v0.1")

    argsWithoutProg := os.Args[1:]
    var conn ConnectInfo
    if (len(argsWithoutProg) > 0) {
        conn, _ = GrepParams(argsWithoutProg)
    } else {
        FillParams(&conn)
    }

    fmt.Printf("IP %s:%d\r\n", conn.ip_addr.IP.String(), conn.udp_port)
    if conn.conn_type == "udp" {
        WorkUdp(conn)
    } else {
        WorkTcp(conn)
    }
}


/**
 @brief
 @param

 @return
 @return
 */
func GrepParams(args []string) (conn ConnectInfo, err error) {
    for i, arg := range args {
        switch arg {
        case "-ip":
            conn.ip_addr.IP = net.ParseIP(args[i+1])
            if conn.ip_addr.IP == nil {
                return conn, fmt.Errorf("GrepParams: illegal IP value - %s", args[i+1])
            }
        case "-pt":
            val, _ := strconv.ParseUint(args[i+1], 10, 16)
            conn.udp_port = uint16(val)
        case "-t":
            conn.conn_type = args[i+1]
        case "-s":
            if args[i+1] == "y" {
                conn.data_type = true
            }
        case "-w":
            if args[i+1] == "y" {
                conn.wait_recv = true
            }
        case "-f":
            // Read from file
        case "-h":
            fmt.Println("\tWelcome to NetClient v0.1")
            fmt.Println("NetClient - udp/tcp client, that allows to send specific data to specific node.")
            fmt.Println("Commands:")
            fmt.Println("\t-ip - IP address of destination node;")
            fmt.Println("\t-pt - Port of UDP/TCP connection with destination node;")
            fmt.Println("\t-t  - Type of connection (udp/tcp);")
            fmt.Println("\t-s  - Type of sending data (y - string data; n - raw (hex) data;")
            fmt.Println("\t-a  - Autosending mode, when enable - will send data without waiting answer;")
            fmt.Println("Example: ./net_client -ip 192.168.0.25 -pt 52344 -s y")
        }
    }

    if conn.conn_type != "tcp" {
        conn.conn_type = "udp"
    }

    return conn, err
}


/**
  @brief
  @param

  @return
 */
func FillParams(conn *ConnectInfo) {
    fmt.Print("Enter IP address: ")
    conn.ip_addr.IP = net.ParseIP(ReadKeybrdData())
    if conn.ip_addr.IP == nil {
        panic("illegal IP value")
    }

    fmt.Print("Enter port: ")
    val, _ := strconv.ParseUint(ReadKeybrdData(), 10, 16)
    conn.udp_port = uint16(val)

    fmt.Print("TCP or UDP? (y/n): ")
    if ReadKeybrdData() == "y" {
        conn.conn_type = "tcp"
    } else {
        conn.conn_type = "udp"
    }

    fmt.Print("Wait answer? (y/n): ")
    if ReadKeybrdData() == "y" {
        conn.wait_recv = true
    }


}


/**
  @brief
  @param[in]

  @return none
 */
func WorkUdp(conn ConnectInfo) {
    udp_conn, err := net.Dial(conn.conn_type, conn.ip_addr.String() + ":" + strconv.Itoa(int(conn.udp_port)))
    if err != nil {
        fmt.Printf("Error on Dial %v\r\n", err)
        return
    }
    defer udp_conn.Close()

    /*if err = udp_conn.SetReadDeadline(time.Now().Add(300)); err != nil { // 180 seconds for Receiving timeout
        panic(err)
    }*/
    for {
        _data := ReadKeybrdData()
        var data []byte
        if conn.data_type {
            data = []byte(_data)
        } else {
            _sep_data := strings.Fields(_data)
            data = StrDigitToBytes(_sep_data)
        }

        if err = Send(udp_conn, data); err != nil {
            panic(err)
        }

        if conn.wait_recv == true {
            recv_data, _, err := Recv(udp_conn); if err != nil {
                break
            }
            for i, recv := range recv_data {
                fmt.Printf("%d) 0x%02x\r\n", i, recv)
            }
        }

        fmt.Print("Reapeat? (y/n): ")
        if ReadKeybrdData() != "y" {
            break
        }
    }
}


func WorkTcp(conn ConnectInfo) {

}


/**
  @brief Send data from network
  @param[in] udp_conn - interface of net connection
  @param[in] data     - data for sending

  @return
*/
func Send(conn net.Conn, data []byte) (err error) {
    _, err = conn.Write(data)
    if err != nil {
        fmt.Printf("Error when send data: %v\r\n", err)
    }

    return
}


/**
  @brief Receive data from network
  @param[in] udp_conn - interface of net connection

  @return
*/
func Recv(conn net.Conn) (data []byte, count int, err error){
    count, err = conn.Read(data); if err != nil {
        fmt.Printf("Error when receive data: %v\r\n", err)
        return
    }

    return
}


/**
  @brief Convert from string representation of hex values to array of hex bytes
  @param[in] numbers - array (slice?) of strings

  @return data - array of bytes (with hex representation of digits)
*/
func StrDigitToBytes(numbers []string) (data []byte) {
    for _, str := range numbers {
        _prt_data, _ := strconv.ParseUint(str, 16, 8)
        data = append(data, byte(_prt_data))
    }

    return
}


/**
  @brief Get data from user input (via stdin)
  @param none

  @return data - string from user input
  @note If get error will call a panic
 */
func ReadKeybrdData() (data string) {
    reader := bufio.NewReader(os.Stdin)
    data, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error on ReadString -", err)
        panic(err)
    }

    data = data[:len(data)-1] // drop \n
    return
}