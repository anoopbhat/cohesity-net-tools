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

	log.Println(r)

	host := r.FormValue("host")
	user := r.FormValue("user")
	pass := r.FormValue("password")
	keyauth := r.FormValue("keyauth");

	config := &ssh.ClientConfig{}

	// if keyauth is set to true. always use keyauth.
	if keyauth == "true"  {
	    log.Println("Setting up SSH with Key Based Authenticaion")
	    config = setupSSHConfigWithKey(user)
	} else {
	    log.Println("Segtting up SSH with password based authentication.")
	    config = setupSSHConfigWithPass(user, pass)
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

// return an ssh config using password as the auth method
func setupSSHConfigWithPass(user string, pass string) *ssh.ClientConfig {

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

    return config
}

// setup an ssh config using private key as the auth method
func setupSSHConfigWithKey(user string) *ssh.ClientConfig {

    // the private key
    pkey_bytes := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA0BaySMin5yh1/ffUalO/AyMKX3OSJoKjGa+995UisvrpEDlu
Qysv3L+6uYYFmimBKG1SgrCPzfJLqw+D/4NL3uAI07DsqWjEhNZViPZTuPI5mTfx
ikj/CBehOURbMLOt2G0KVxMypIVtGnzuLhcF+Xlz2y8JGjaSronuFgsBLvoOt/BU
sAcHU8vtDNFdbsLOczSdFEQOQx/WuVyBjKXOyxOU9GOJaUfsdPlW0f4ZnXeuCTEn
Mh1UacMV80keGafCaHFv4m6dCH4iNFZYQEh0han5tu/VHB3YHFkDDyc8mETWlyn1
5laNrGde8wP8MHxunpitZZt2bXbh+HKlUlqM4wIDAQABAoIBAQCKhLFFdh0e6XYy
C4mhBgJ/KhI7nAlMDWZZMP26E9K3ZgNDQ5e8qsD/p7m6yhZsmvhZWvyz9qijpYjt
ZDSwIEyfHm+By6Ke2xkGfE8QDzmIQeZJsk3did4LGv+9yV0SvGkbSuq5MBRkJFWO
bl922uilO03+N/9NLcrS2QpeLhEpnSOQMc8eZkapuhN1TkNlTKtbsecZzPVAgM/h
9omqvmito3Ktj0V9T7/HPU/LdcjgEawN4R+ZxxGRwTf31CKoufXPDvQIclrV7ps0
tQFv5vcZ33doVwrVKsrMbYoRYdJdQ3NGf8eU8wwwZzEyFdKWtbt5CG+jyeeIbjJC
9cs5m4hhAoGBAPimXckNXFMukuErM+KfsUz35+iss86PTXQ3ag+8LLd6MMRAVE+g
tmcoC9O48NAAtrLpZM0mBLfQncya2J4opUEtHr2ibuXhVukkqmXKpHjMwZOvADVs
saMEOFXehMXs4JtBxb7T2XwjOarI/q7Rf9iLsH4Efg6MUPiRJIZVAObbAoGBANY9
YxavCkzGDOBeDgxo5LGD9Beim4hmLN42qElFMDOQhLyU7IENsttMB4loSPC2SdXL
U5QQulzdwxqPJBJAJ5B4kcKCDyoVRUi+pO2VbmzsixzDdwE5EOyO0CPuPnkptV9F
FXQbhYwKa5o4Rizw5mRj7wX2FYjqKbvQS9HXA/yZAoGAHt9BG7pd8TICKJTdn1Cm
ieDp2VjABnCCdGCA+a0qfClerq8yCKTyoMI3HbWDqL+9717NFi+XPF9ZiFLdfF2d
jwcUHwVw8XfV+6KCyZqsaxc5HaYHx5pUP+JBQGAdahmsFXrIG5ZgFWqmOU81V+1J
C1Dku/DA2fuP/hy/RTJ+pysCgYAfvEAtYAh6juvhYI1cMT2PPiiuR5wafGgxEo+j
KuiU+tdux/CwvUK9UWncZOJJJfeR/+iFimTQ1NjN2l5RhcdWk0WkNnfgl/4HZJYx
y2zsHa4NuLasK7PiFtWmPOhsMk13q1geNuV1dSWzVpqulZDLVjJWA7n06hr8g0J3
9w3UIQKBgH1PcYSGmRK8wPEPTtUuOO1nf3Xey+BYvPmYcJ2GRg0EpxRyKQrUIe8h
RQeS/3jMBksUEN+qawGsFkaesCu4axjDWKwOkH/Y/ExNiGgS6Wfo7WJueXVEZdF+
mgr8v3UU92cGLWY8AU3WHRaw6jaOaBOxOm7NHe320hhYggdX6Oha
-----END RSA PRIVATE KEY-----`)

    // parse the key
    signer, err := ssh.ParsePrivateKey(pkey_bytes)

	if err != nil {
	    panic(err)
        }

    config := &ssh.ClientConfig{
	User: user,
        Auth: []ssh.AuthMethod{
	    ssh.PublicKeys(signer),
	},

    // need this apparently. no idea why. see https://github.com/golang/go/issues/19767
        HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
            return nil
        },
    }

    return config
}
