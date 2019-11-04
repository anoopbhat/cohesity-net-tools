package main

import (
   "fmt"
    "io"
    "log"
    "net/http"
    "github.com/sparrc/go-ping"
    "time"
    "net"
    "strings"
    "bytes"

    "golang.org/x/crypto/ssh"
)

func main() {

    // how we serve up static pages like the forms
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/", fs)
    // todo help page
    // http.Handle("/help", fs)

    // function to handle a ping test
    pingFunc := func(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
            log.Println(w, "ParseForm() err: %v", err)
            return
        }

	pingHost := r.FormValue("pingHost")

	pinger, err := ping.NewPinger(pingHost)
	if err != nil {
		panic(err)
	}

	pinger.Count = 3
	pinger.Timeout =  5 * time.Second
	pinger.Run() // blocks until finished
	stats := pinger.Statistics() // get send/receive/rtt stats

	statsStr := fmt.Sprintf(" %v - %d packets transmitted, %d packets received, %v%% packet loss. round-trip min/avg/max/stddev = %v/%v/%v/%v",
                stats.IPAddr, stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss, stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt);
	io.WriteString(w, statsStr);
    }

    // function to handle the testing of an open port from the node
    // todo handle a timeout appropriately and send back to the client.
    portTestFunc := func(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
            log.Println(w, "ParseForm() err: %v", err)
            return
        }

	host := r.FormValue("Host")
	port := r.FormValue("Port")

	connStr := net.JoinHostPort(host, port);

	dialer := net.Dialer{Timeout: 3 * time.Second}

	conn, err := dialer.Dial("tcp", connStr)
		if err != nil {
		    io.WriteString(w,err.Error())
		} else {
		    io.WriteString(w, "Connected to " +  connStr);
		    defer conn.Close()
		}
    }


    // function to handle the looking up of a name. 
    // todo lookup an IP and get a name back
    DNSLookupFunc := func(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
            log.Println(w, "ParseForm() err: %v", err)
            return
        }

	fqdn := r.FormValue("fqdn")

	IPs, err := net.LookupIP(fqdn)

	    if err != nil {
		panic(err)
	    }

	// convert each IP, a slice, in IPs which is also a slice. it's a slice of slices.
	sliceIPs := []string{}

	for _, IP := range IPs {
	    sliceIPs = append(sliceIPs, IP.String())
	}
	// make a string out of it
	strIPs := strings.Join(sliceIPs, ", ")

	io.WriteString(w,  "IPs: " +  strIPs);
    }

    // function to handle ssh'ing into a host
    sshFunc := func(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
	    log.Println(w, "ParseForm() err: %v", err)
	    return
	}

	host := r.FormValue("host")
	user := r.FormValue("user")
	pass := r.FormValue("password")
	/*pemkey := r.FormValue("sshPrivKey")

	key, err := ssh.ParsePrivateKey([]byte(pemkey))

	if err != nil {
	    panic(err)
	}*/

	// setup the client configuration
	config := &ssh.ClientConfig{
	    User: user,
	    Auth: []ssh.AuthMethod{
	    ssh.Password(pass),
	    },

	    // need this apparently. no idea why. see https://github.com/golang/go/issues/19767
	    HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	    },
	}

	// connect to the host
	client, err := ssh.Dial("tcp", host, config) 

	if err != nil {
	    panic(err)
	}

	// establish a session
	session, err := client.NewSession()

	if err != nil {
	    panic(err)
	}

	defer session.Close()

	var b bytes.Buffer

	session.Stdout = &b
	if err := session.Run("/usr/bin/id"); err != nil {
	    panic("Failed to run: " + err.Error())
	}

	io.WriteString(w, "/usr/bin/id output: " + b.String())
    }

    http.HandleFunc("/ping", pingFunc)
    http.HandleFunc("/port", portTestFunc)
    http.HandleFunc("/dns", DNSLookupFunc)
    http.HandleFunc("/ssh", sshFunc)
//    http.HandleFunc("/trace", traceFunc)


    log.Println("Listening...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
