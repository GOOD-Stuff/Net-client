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
}


/**
-ip - ip address of the node
-pt - port
-s - string data type
-d - path to file with data
-a - auto sending
-w - wait answer ???
*/
// ./net_client -ip 192.168.0.12 -pt 6000 -t udp -s -a
// ./net_client -ip 192.168.0.12 -pt 6000 -t udp -w -a y/n
// ./net_client -ip 11.200.1.1 -udp 12345 -t t ...
func main () {
    fmt.Println("\tUDP client v0.1")

    argsWithoutProg := os.Args[1:]
    var conn ConnectInfo
    if (len(argsWithoutProg) > 0) {
        conn, _ = GrepParams(argsWithoutProg)
    } else {
        fmt.Println("Enter IP address: ")
    }

    fmt.Printf("IP %s\r\n", conn.ip_addr.IP.String())
    fmt.Println("Please, enter your data:")

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
        //fmt.Printf("%d) %s\r\n", i, arg)
        switch arg {
        case "-ip":
//            fmt.Printf("Node IP: %s\r\n", args[i+1])
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
            } else {
                conn.data_type = false
            }
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

    if conn.conn_type != "udp" || conn.conn_type != "tcp" {
        conn.conn_type = "udp"
    }


    return conn, err
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

    reader := bufio.NewReader(os.Stdin)
    _data, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error on ReadString -", err)
        return
    }

    var data []byte
    if conn.data_type {
        data = []byte(_data)
    } else {
        _sep_data := strings.Fields(_data)
        data = StrDigToBytes(_sep_data)
    }

    if err = SendUdp(udp_conn, data); err != nil {
        panic(err)
    }

    var re []byte
    udp_conn.Read(re)
}


func WorkTcp(conn ConnectInfo) {

}


/**
  @brief Send data in UDP mode
  @param[in] udp_conn - interface of UDP connection
  @param[in] data     -
*/
func SendUdp(udp_conn net.Conn, data []byte) (err error) {
    _, err = udp_conn.Write(data)
    if err != nil {
        fmt.Printf("Error when send data: %v\r\n", err)
    }

    return
}


func RecvUdp() {

}


/**
  @brief Convert from string representation of hex values to array of hex bytes
  @param[in] numbers - array (slice?) of strings

  @return data - array of bytes (with hex representation of digits)
*/
func StrDigToBytes(numbers []string) (data []byte) {
    for _, str := range numbers {
        _prt_data, _ := strconv.ParseUint(str, 16, 8)
        data = append(data, byte(_prt_data))
    }

    return
}