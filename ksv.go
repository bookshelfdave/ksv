package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
)

type v1secret struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string
	Metadata   map[string]string
	Type       string
	Data       map[string]string
	StringData map[string]string
}

func readInputOrFail(r io.Reader) []byte {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal("Can't read from stdin")
	}
	return bytes
}

func secretFromYaml(r io.Reader) (*v1secret, error) {
	bytesIn := readInputOrFail(r)
	s := v1secret{}
	s.StringData = make(map[string]string)

	if err := yaml.Unmarshal(bytesIn, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func secretToYamlString(s *v1secret) (string, error) {
	if d, err := yaml.Marshal(&s); err != nil {
		return "", err
	} else {
		return string(d), nil
	}
}

func decodeFromBase64(r io.Reader, decodeToStringData bool) (*v1secret, error) {
	s, err := secretFromYaml(r)
	if err != nil {
		return nil, err
	}
	for k, v := range s.Data {
		decoded, _ := base64.StdEncoding.DecodeString(v)
		s.Data[k] = string(decoded)
		if decodeToStringData {
			s.StringData[k] = string(decoded)
			delete(s.Data, k)
		}
	}
	return s, nil
}

func encodeToBase64(r io.Reader) (*v1secret, error) {
	s, err := secretFromYaml(r)
	if err != nil {
		return nil, err
	}
	for k, v := range s.Data {
		str := base64.StdEncoding.EncodeToString([]byte(v))
		s.Data[k] = str
	}
	return s, nil
}

func encodeCmd(c *cli.Context) {
	s, err := encodeToBase64(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	if ss, err := secretToYamlString(s); err != nil {
		log.Fatal("Can't convert back to yaml")
	} else {
		fmt.Println(ss)
	}
}

func decodeCmd(c *cli.Context) {
	s, err := decodeFromBase64(os.Stdin, c.Bool("s"))
	if err != nil {
		log.Fatal(err)
	}
	if ss, err := secretToYamlString(s); err != nil {
		log.Fatal("Can't convert back to yaml")
	} else {
		fmt.Println(ss)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "ksv"
	app.Usage = "decode base64-encoded K8s YAML secrets from stdin"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		{
			Name:    "encode",
			Aliases: []string{"e"},
			Usage:   "encode a secrets yaml file on stdin",
			Action:  encodeCmd,
		},
		{
			Name:    "decode",
			Aliases: []string{"d"},
			Usage:   "decode a secrets yaml file from stdin (default command)",
			Action:  decodeCmd,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "stringData, s",
					Usage: "convert Data to StringData",
				},
			},
		},
	}

	// default subcommand is decode without converting stringData
	app.Action = func(c *cli.Context) {
		decodeCmd(c)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
