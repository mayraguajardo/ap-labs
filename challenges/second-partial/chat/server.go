// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"flag"
	"bytes"
	"strings"
	"time"
)

//!+broadcaster
type client chan<- string // an outgoing message channel

type userStruct struct{
	userName string
	msg string
}

type clientStruct struct {
	clientt client
	connection net.Conn
	lastConnection string
}

var (
	serverPrefix = "irc-server > "
	admin client
	globalUser string
	entering = make(chan clientStruct)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
	direct = make(chan userStruct)
	kick = make(chan string)
	users  map[string]client
	connections map[string]net.Conn
	lastConnections map[string]string
)

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	users = make(map[string]client)
	connections = make(map[string]net.Conn)
	lastConnections = make(map[string]string)


	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg
			}

		case usrS := <- direct:
			users[usrS.userName] <- usrS.msg

		case cliS := <-entering:
			if len(clients) == 0{
				cliS.clientt <- serverPrefix + "You are the first user"
				cliS.clientt <- serverPrefix + "You are tje new IRC Server Admin"
				fmt.Printf("[%s] was promoted as the channel admin\n", globalUser)
				admin = cliS.clientt
			}
			clients[cliS.clientt] = true
			users[globalUser] = cliS.clientt
			connections[globalUser] = cliS.connection
			lastConnections[globalUser] = cliS.lastConnection

		case cli := <-leaving:
			if admin == cli{
				for newAdmin := range clients{
					admin = newAdmin
					newAdmin <- serverPrefix + "You are the admin now"
					continue
				}
			}
			delete(clients, cli)
			close(cli)

		case user := <- kick:
			connection := connections[user]
			client := users[user]
			delete(clients,client)
			delete(connections,user)
			delete(users,user)
			delete(lastConnections,user)
			connection.Close()
		}
	}
}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn) {

	var buf = make([]byte, 1024)
	conn.Read(buf)
	localUser := string(bytes.Trim(buf, "\x00"))
	globalUser = string(bytes.Trim(buf, "\x00"))

	ch := make(chan string) // outgoing client messages

	if users[localUser] != nil{
		//Existing user
		fmt.Fprintln(conn, " " + localUser + "User name already exists")
		close(ch)
		conn.Close()
		return
	}

	fmt.Printf("%s New connected user [%s]\n", serverPrefix, localUser)
	go clientWriter(conn, ch)

	// welcome message
	ch <- serverPrefix + "Welcome to the IRC Server "
	ch <- serverPrefix + "User ["+ localUser + "] is succesfully logged"

	// entering messages
	messages <- serverPrefix + "New connected user ["+ localUser + "]"
	entering <- clientStruct{ch,conn,time.Now().Format("15:04:05\n")}

	input := bufio.NewScanner(conn)
	for input.Scan(){
		if len(input.Text()) > 0 && string(input.Text()[0]) == "/"{
			slice := strings.Split(input.Text(), " ")
			command := slice[0]
			switch command{
			case "/users":
				str := ""
				for usr := range users{
					userLastConnection := lastConnections[usr]
					str += serverPrefix + usr + " - last connection : " + userLastConnection
				}
				ch <- str
			case "/msg":
				if len(slice) < 2{
					ch <- "User no specified"
					continue
				}
				if len(slice) < 3{
					ch <- "Please enter a message"
					continue
				}
				addr := slice[1]
				if _, ok := users[addr]; ok{
					directMsg := input.Text()[strings.Index(input.Text(),addr) + len(addr)+1:]
					direct <- userStruct{addr, localUser + " > " + directMsg}
				} else{
					ch <- "User: " + addr + " does not exist"
				}

			case "/time":
				timeZone := "America/Mexico_City"
				loc, _ := time.LoadLocation(timeZone)
				timme := time.Now().In(loc).Format("15:04\n")
				ch <- serverPrefix + "Local time: " + timeZone + " " + strings.Trim(timme, "\n")

			case "/user":
				if len(slice) < 2{
					ch <- "User not specified"
					continue
				}
				user := slice[1]
				if _, ok := users[user];ok{
					ip := connections[user].RemoteAddr().String()
					userLastConnection := lastConnections[user]
					ch <- serverPrefix + "username: " + user + ", IP: " + ip + " - last connection: " + userLastConnection  
				} else {
					ch <- "User: " + user + " does not exist"
				}
			case "/kick":
				if len(slice) < 2{
					ch <- "Enter the user to kick"
					continue
				}
				if ch == admin{
					user := slice[1]
					if _, ok := users[user]; ok{
						messages <- "[" + user + "] was kicked from channel"
						fmt.Printf("%s[%s] was kicked from channel\n", serverPrefix, user)
						direct <- userStruct{user, serverPrefix + "You were kicked from this channel"}
						kick <- user
					}else{
						ch <- "User: " + user + " does not exist"
					}
				}else{
					ch <- "Only the admin can kick people out of the server"
				}
			default:
				ch <- "invalid command"

			}
		}else{
			messages <- localUser + " > " + input.Text()
		}
	}


	
	// NOTE: ignoring potential errors from input.Err()

	leaving <- ch
	messages <-"[" + localUser + "]  has left"
	fmt.Printf("%s[%s] left\n", serverPrefix, localUser )
	delete(users, localUser)
	delete(connections, localUser)
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

//!+main
func main() {

	if len(os.Args) < 5{
		fmt.Println("Usage: go run server.go -host localhost -port [port]")
		os.Exit(1)
	}
	host :=flag.String("host", "localhost", "localhost")
	port := flag.String("port", "9000", "port")
	flag.Parse()
	listener, err := net.Listen("tcp", *host+ ":" + *port)
	fmt.Println(serverPrefix + "Simple IRC Server started at " + *host + ":" + *port)
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	fmt.Println(serverPrefix + "Ready for recieving clients")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

//!-main
