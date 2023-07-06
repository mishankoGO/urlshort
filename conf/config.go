package conf

import (
	"flag"
	"log"
	"os"
)

type ShortenerConfig struct {
	Path string
}

func InitFlags(conf *ShortenerConfig) error {
	flags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// parse file name
	FileName := flags.String("f", conf.Path, "file name")

	// parse the flags
	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Printf("error parsing flags: %v", err)
		return err
	}

	conf.Path = *FileName
	return nil
}
