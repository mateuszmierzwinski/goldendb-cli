package main

import (
	"flag"
	"fmt"
	"goldencli/internal/protocol"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultBuffSize = 65535
)

func main() {

	var addr string
	flag.StringVar(&addr, "a", "localhost:8091", "Address of the database")
	flag.Parse()

	if addr == "" {
		flag.Usage()
		os.Exit(-1)
	}

	if len(os.Args) < 2 {
		fmt.Printf("%s -a <address> command operators\n", filepath.Base(os.Args[0]))
		flag.Usage()
		os.Exit(-1)
	}

	resultFun := commandProcessor(os.Args)

	if resultFun == nil {
		log.Println("Unknown command given")
		os.Exit(-1)
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()
	ts := time.Now()
	resultFun(conn)
	tSince := time.Since(ts)
	fmt.Println(fmt.Sprintf("(Done in %d ms / %d ns)", tSince.Milliseconds(), tSince.Nanoseconds()))

}

func commandProcessor(args []string) func(conn net.Conn) {
	switch strings.ToLower(args[1]) {
	case "ping":
		return pingFunc
	case "get":
		return getFunc
	case "put":
		return putFunc
	}
	return nil
}

func putFunc(conn net.Conn) {
	if len(os.Args) < 4 {
		fmt.Printf("%s put <input file> <object name>\n", filepath.Base(os.Args[0]))
		return
	}

	f, err := os.OpenFile(os.Args[2], os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(">> starting transfer...")
	// we're creating object
	conn.Write([]byte{protocol.Create})

	fmt.Println(">> writing key size")
	conn.Write(protocol.IntToBytes(len(os.Args[3])))

	fmt.Println(">> writing key")
	conn.Write([]byte(os.Args[3]))

	fmt.Println(">> Writing object size")
	// we're sending object size in bytes
	conn.Write(protocol.Int64toBytes(stat.Size()))

	fmt.Println(">> Sending data")
	// copy the data
	if _, err := io.Copy(conn, f); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(">> Sent! Waiting for response")
	// get response
	op := readStatus(conn)
	if op == protocol.StatusOK {
		fmt.Println(fmt.Sprintf("OK - sent %d bytes", stat.Size()))
		return
	} else {
		data, err := ioutil.ReadAll(conn)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(data))
	}
}

func getFunc(conn net.Conn) {
	if len(os.Args) < 4 {
		fmt.Printf("%s get <object name> <output file>\n", filepath.Base(os.Args[0]))
		return
	}

	conn.Write([]byte{protocol.Read})
	conn.Write(protocol.IntToBytes(len(os.Args[2])))
	conn.Write([]byte(os.Args[2]))

	op := readStatus(conn)
	if op == protocol.StatusOK {
		sizeBuff := make([]byte, 8)
		if _, err := conn.Read(sizeBuff); err != nil {
			fmt.Println(err.Error())
			return
		}
		size := protocol.BytesArrayToUint64(sizeBuff)

		log.Println(fmt.Sprintf("Downloading %d bytes...\n", size))

		f, err := os.OpenFile(os.Args[3], os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer f.Close()

		valBuff := make([]byte, defaultBuffSize)
		copied := 0
		for {
			sz, e := conn.Read(valBuff)
			copied = copied + sz
			f.Write(valBuff[:sz])
			if int64(copied) >= size || e != nil {
				break
			}
		}
	} else {
		data, err := ioutil.ReadAll(conn)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(data))
	}
}

func pingFunc(conn net.Conn) {
	conn.Write([]byte{protocol.Ping})

	op := readStatus(conn)
	if op == protocol.StatusOK {
		data, err := ioutil.ReadAll(conn)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(data))
	}
}

func readStatus(conn net.Conn) byte {
	bufr := make([]byte, 1)
	if _, err := conn.Read(bufr); err != nil {
		fmt.Println(err.Error())
		return 255
	}
	return bufr[0]
}
