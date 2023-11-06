package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"

	"github.com/akamensky/argparse"
)

var versionCode = "v0.3"
var proxyBufferSize = 4096
var httpListenPort = 80
var httpsConnectingPort = 443
var allowInsecure = false
var debugLog = false

func handler(responseToRequest http.ResponseWriter, incomingRequest *http.Request) {

	host := incomingRequest.Host
	url := incomingRequest.URL
	remote := incomingRequest.RemoteAddr

	log.Printf("Request from %s to host %s and url %s", remote, host, url)

	// Get the raw request bytes
	requestDump, err := httputil.DumpRequest(incomingRequest, true)
	if err != nil {
		log.Printf("cannot dump %s", err)
		http.Error(responseToRequest, "Cannot dump request", http.StatusBadRequest)
	}

	//ioutil.WriteFile("input.txt", requestDump, 0644)

	if debugLog {
		log.Printf("Dump:\n%s\n", string(requestDump))
	}

	conf := &tls.Config{}

	if allowInsecure {
		conf = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	conn, err := tls.Dial("tcp", host+":"+strconv.Itoa(httpsConnectingPort), conf)
	if err != nil {
		log.Printf("Cannot dial host %s", err)
		http.Error(responseToRequest, "Cannot dial host", http.StatusGatewayTimeout)
		return
	}
	defer conn.Close()

	n, err := conn.Write(requestDump)
	if err != nil {
		log.Printf("Cannot write request %d %s\n", n, err)
		http.Error(responseToRequest, "Cannot write request"+err.Error(), http.StatusBadGateway)
		return
	}

	// Prepare the requesting socket for writing. Access raw socket by hijacking
	// Reference: https://stackoverflow.com/questions/29531993/accessing-the-underlying-socket-of-a-net-http-response

	hj, ok := responseToRequest.(http.Hijacker)
	if !ok {
		http.Error(responseToRequest, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	returnConn, _, err := hj.Hijack()

	if err != nil {
		http.Error(responseToRequest, err.Error(), http.StatusInternalServerError)
		return
	}

	defer returnConn.Close()

	readBuf := make([]byte, proxyBufferSize)

	for {
		//Read from response socket from external server and pass data back
		bytesRead, err := conn.Read(readBuf)

		if err != nil {

			if err == io.EOF {
				log.Printf("EOF reached")
			} else {
				log.Printf("Error getting bytes from server %d %s", bytesRead, err)
			}

			break
		}

		bytesWritten, err := returnConn.Write(readBuf[:bytesRead])

		//ioutil.WriteFile("output.txt", readBuf[:bytesRead], 0644)

		if err != nil {
			log.Printf("Error writing bytes to requester %d %s", bytesWritten, err)
			break
		}

	}

	log.Println("End of handler")

}

func main() {
	parser := argparse.NewParser("http-to-https-proxy", "A proxy that upgrades HTTP connections to HTTPS for systems which cannot make HTTPS requests.")

	var parsedHTTPPort *int = parser.Int("l", "listen", &argparse.Options{Help: "HTTP port to listen on", Default: httpListenPort})
	var parsedHTTPSPort *int = parser.Int("c", "connect", &argparse.Options{Help: "HTTPS port to connect to", Default: httpsConnectingPort})
	var parsedProxyBuffer *int = parser.Int("b", "buffer", &argparse.Options{Help: "Buffer size", Default: proxyBufferSize})
	var parsedAllowInsecure *bool = parser.Flag("i", "insecure", &argparse.Options{Help: "Allow insecure TLS certificates", Default: allowInsecure})
	var parsedDebugLog *bool = parser.Flag("d", "debug", &argparse.Options{Help: "Enable debug console logging", Default: debugLog})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Print(parser.Usage(err))
	}

	httpListenPort = int(*parsedHTTPPort)
	httpsConnectingPort = int(*parsedHTTPSPort)
	proxyBufferSize = int(*parsedProxyBuffer)
	allowInsecure = bool(*parsedAllowInsecure)
	debugLog = bool(*parsedDebugLog)

	log.Printf("HTTP to HTTPS proxy %s listening to %d, forward to %d with listening buffer %d", versionCode, httpListenPort, httpsConnectingPort, proxyBufferSize)

	if allowInsecure {
		log.Printf("Allow insecure TLS certificates")
	}

	http.HandleFunc("/", handler)

	if err := http.ListenAndServe(":"+strconv.Itoa(httpListenPort), nil); err != nil {
		log.Fatal(err)
	}
}
