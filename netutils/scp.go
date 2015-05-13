package netsq // import "github.com/leejansq/pkg/netutils"

import (
	//"bytes"
	"code.google.com/p/go.crypto/ssh"
	//"fmt"
	"fmt"
	scp "github.com/leejansq/pkg/netutils/ssh"
	//scp "github.com/leejansq/pkg/ssh"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

type ScpClient interface {
	UploadDir(string, string, []string) error
	Upload(string, io.Reader, *os.FileInfo) error
}

type scpcomm struct {
	robot ScpClient
}

func (s *scpcomm) Upload(dst, src string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	fs, err := f.Stat()
	if err != nil {
		return err
	}
	if fs.Mode().IsDir() {
		return s.robot.UploadDir(dst, src, nil)
	}
	return s.robot.Upload(dst, f, &fs)
}

type ScpUpload interface {
	Upload(string, string) error
}

func NewScpClient(address, password string) (ScpUpload, error) {
	// Create client config
	var user string
	if !strings.Contains(address, "@") {
		return nil, fmt.Errorf("address formate is error!")
	} else {
		spits := strings.SplitN(address, "@", 2)
		user = spits[0]
		address = spits[1]
		if strings.Contains(address, ":") {
			reg := regexp.MustCompile(`:[0-9]{2,5}`)
			if !reg.MatchString(address) {
				return nil, fmt.Errorf("address formate is error!")
			}
		} else {
			address += ":22"
		}

	}
	log.Printf("[Address=\"%s\" | user=\"%s\" | password=\"%s\"]", address, user, password)

	clientConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	conn := func() (net.Conn, error) {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			log.Fatalf("unable to dial to remote side: %s", err)
		}
		return conn, err
	}

	config := &scp.Config{
		Connection: conn,
		SSHConfig:  clientConfig,
	}

	cli, err := scp.New(address, config)
	if err != nil {
		return nil, err
	}

	return &scpcomm{robot: cli}, nil

	//var cmd packer.RemoteCmd
	//stdout := new(bytes.Buffer)
	//cmd.Command = "echo foo"
	//cmd.Stdout = stdout

	//client.Start(&cmd)

}
