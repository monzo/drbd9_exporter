package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type drbdConnection struct {
	// Infomation about the connection
	Resource   string
	RemoteHost string
	ResourceID string
	// Properties / Status of the connection
	ConnectionStatus string
	MyRole           string
	RemoteRole       string
	MyDiskStatus     string
	RemoteDiskStatus string
	Suspended        bool
	KVs              []drbdConnectionKV
}

type drbdConnectionKV struct {
	Name  string
	Value float64
}

func getAllDRDBstatues() []drbdConnection {
	connections := make([]drbdConnection, 0)

	filepath.Walk("/sys/kernel/debug/drbd/resources/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() != "proc_drbd" {
			return nil
		}

		// /sys/kernel/debug/drbd/resources/$Resource/connections/$RemoteHost/$ResourceID/proc_drbd
		//
		// []string{"", "sys", "kernel", "debug", "drbd", "resources", "$Resource", "connections", "$RemoteHost", "$ResourceID", "proc_drbd"}
		pathSegments := strings.Split(path, string(os.PathSeparator))

		dC := drbdConnection{
			Resource:   pathSegments[6],
			RemoteHost: pathSegments[8],
			ResourceID: pathSegments[9],
		}

		procDrbd, err := ioutil.ReadFile(path)
		if err != nil {
			log.Printf("failed to read %s: %s", path, err.Error())
			return nil
		}

		err = parseProcDRBD(procDrbd, &dC)
		if err != nil {
			log.Printf("failed to parse %s: %s", path, err.Error())
			return nil
		}

		connections = append(connections, dC)
		return nil
	})

	return connections
}

var errInvalidOutput = fmt.Errorf("Failed to parse proc_drbd data")

/*
 1: cs:Established ro:Secondary/Primary ds:UpToDate/UpToDate C r-----
    ns:0 nr:119598720 dw:161539164 dr:0 al:2125 bm:0 lo:0 pe:[0;0] ua:0 ap:[0;0] ep:1 wo:2 oos:0
    resync: used:0/61 hits:25476 misses:922 starving:0 locked:0 changed:592
    act_log: used:0/1237 hits:13283342 misses:15109 starving:0 locked:0 changed:12057
    blocked on activity log: 0/0/0

 0: cs:Established ro:Secondary/Secondary ds:UpToDate/UpToDate C r-----
    ns:0 nr:0 dw:158780236 dr:0 al:1950 bm:0 lo:0 pe:[0;0] ua:0 ap:[0;0] ep:1 wo:2 oos:0
    resync: used:0/61 hits:2 misses:285 starving:0 locked:0 changed:285
    act_log: used:0/1237 hits:13371654 misses:15913 starving:0 locked:0 changed:11885
    blocked on activity log: 0/0/0

 1: cs:Established ro:Primary/Secondary ds:UpToDate/UpToDate C r-----
    ns:119876764 nr:0 dw:154278296 dr:14701504 al:2070 bm:0 lo:0 pe:[0;0] ua:0 ap:[0;0] ep:1 wo:2 oos:0
    resync: used:0/61 hits:25566 misses:719 starving:0 locked:0 changed:485
    act_log: used:0/1237 hits:13318630 misses:41690 starving:0 locked:0 changed:12063
    blocked on activity log: 0/0/0
*/

var bannerRegexp = regexp.MustCompilePOSIX(` [0-9]: cs:([a-zA-Z]+) ro:([^/]+)/([^/]+) ds:([^/]+)/([^/]+) [a-zA-Z] ([\-rs])`)
var kvExtractionRegexp = regexp.MustCompilePOSIX(`(([a-z]+):([0-9]+))+`)

func parseProcDRBD(input []byte, dC *drbdConnection) error {

	if !bannerRegexp.Match(input) {
		return errInvalidOutput
	}

	regexMatches := bannerRegexp.FindAllStringSubmatch(string(input), -1)
	//  1: cs:Established ro:Primary/Secondary ds:UpToDate/UpToDate C r-----
	//[][]string{[]string{" 1: cs:Established ro:Secondary/Primary ds:UpToDate/UpToDate C r", "Established", "Secondary", "Primary", "UpToDate", "UpToDate", "r"}}

	dC.ConnectionStatus = regexMatches[0][1]
	dC.MyRole = regexMatches[0][2]
	dC.RemoteRole = regexMatches[0][3]
	dC.MyDiskStatus = regexMatches[0][4]
	dC.RemoteDiskStatus = regexMatches[0][5]
	if regexMatches[0][6] != "r" {
		dC.Suspended = true
	}

	kvMatches := kvExtractionRegexp.FindAllStringSubmatch(string(input), -1)

	kvs := make([]drbdConnectionKV, 0)
	for _, blob := range kvMatches {
		// blob will look something like:
		// 	[]string{"nr:119598720", "nr:119598720", "nr", "119598720"},
		i, err := strconv.ParseFloat(blob[3], 64)
		if err != nil {
			log.Printf("failed to parse KV number! %#v", blob)
			continue
		}

		kv := drbdConnectionKV{
			Name:  blob[2],
			Value: i,
		}

		kvs = append(kvs, kv)
	}

	dC.KVs = kvs
	return nil
}
