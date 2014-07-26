package main

import (
	"flag"
	"fmt"
	"github.com/bioothod/elliptics-go/elliptics"
	"github.com/vmihailenco/msgpack"
	"log"
	"strconv"
	"strings"
	"time"
)

func main() {
	remote := flag.String("remote", "", "remote address in the standard elliptics format: addr:port:family")
	logfile := flag.String("log", "/tmp/backrunner-delete-go.log", "elliptics log file")
	loglevel := flag.Int("log-level", int(elliptics.INFO), "elliptics log level")
	group_str := flag.String("metadata-groups", "11:12", "elliptics *metadata* groups")
	flag.Parse()

	if *remote == "" {
		log.Fatal("You must provide remote address")
	}

	elog, err := elliptics.NewFileLogger(*logfile, *loglevel)
	if err != nil {
		log.Fatalln("NewFileLogger: ", err)
	}
	defer elog.Free()
	elog.Log(elliptics.INFO, "started: %v, logfile: %s, level: %d, remote: %s", time.Now(), *logfile, *loglevel, *remote)

	// Create elliptics node
	node, err := elliptics.NewNode(elog)
	if err != nil {
		log.Fatal(err)
	}
	defer node.Free()

	node.SetTimeouts(100, 1000)
	if err = node.AddRemote(*remote); err != nil {
		log.Printf("AddRemote: %v", err)
	}

	metadata_session, err := elliptics.NewSession(node)
	if err != nil {
		log.Fatal("Failed to create metadata session", err)
	}

	var metadata_groups []int32
	for _, group := range strings.Split(*group_str, ":") {
		g, err := strconv.Atoi(group)
		if err != nil {
			log.Fatal("Invalid element in group string:", group)
		}

		metadata_groups = append(metadata_groups, int32(g))
	}

	metadata_session.SetGroups(metadata_groups)
	metadata_session.SetNamespace(BucketNamespace)

	for _, bucket := range Buckets {
		fmt.Printf("bucket: %s: reading metadata\n", bucket)
		var meta BucketMeta
		var out []interface{}
		for rd := range metadata_session.ReadData(bucket) {
			if rd.Error() == nil {
				err = msgpack.Unmarshal([]byte(rd.Data()), &out)
				if err != nil {
					log.Fatal("Could not parse bucket metadata: ", err)
				}
			}

			err = meta.ExtractMsgpack(out)
			if err != nil {
				log.Fatal("Unsupported msgpack data:", err)
			}

			fmt.Printf("bucket: %s: groups: %v, flags: %v, acl: %v\n", meta.bucket, meta.groups, meta.flags, meta.acl)
		}
	}
}
