package main

import (
    "crypto/tls"
    "fmt"

    irc "github.com/fluffle/goirc/client"
)

func main() {
    server := "irc.wetfish.net"

    // Or, create a config and fiddle with it first:
    config := irc.NewConfig("gobot")
    config.SSL = true
    config.SSLConfig = &tls.Config{ServerName: server}
    config.Server = server + ":6697"
    config.NewNick = func(n string) string { return n + "^" }
    client := irc.Client(config)

    // Add handlers to do things here!
    // e.g. join a channel on connect.
    client.HandleFunc(irc.CONNECTED,
        func(connection *irc.Conn, line *irc.Line) {
            fmt.Printf("Connected to %s \n", server)
            connection.Join("#botspam")
        })

    client.HandleFunc(irc.PRIVMSG,
        func(connection *irc.Conn, line *irc.Line) {
            fmt.Printf("<%s> %s \n", line.Nick, line.Text())
        })

    // And a signal on disconnect
    quit := make(chan bool)
    client.HandleFunc(irc.DISCONNECTED,
        func(connection *irc.Conn, line *irc.Line) { quit <- true })

    // Tell client to connect.
    error := client.Connect();

    if error != nil {
        fmt.Printf("Connection error: %s\n", error.Error())
    }

    // Wait for disconnect
    <-quit
}
