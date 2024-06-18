package net_utils

import (
	"github.com/sxydh/mgo-util/json_utils"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type SSHTunnel struct {
	Id         string `json:"id"`
	sshClient  *ssh.Client
	listener   *net.Listener
	SSHIp      string `json:"sshIp"`
	SSHPort    int    `json:"sshPort"`
	SSHUser    string `json:"sshUser"`
	ListenPort int    `json:"listenPort"`
	TargetIp   string `json:"targetIp"`
	TargetPort int    `json:"targetPort"`
	Status     int    `json:"status"`
}

func BuildSSHTunnel(tunnels *[]*SSHTunnel) {
	var todoTunnels = make(chan *SSHTunnel, 20)
	var doingTunnels = make(chan *SSHTunnel, 20)
	for _, tunnel := range *tunnels {
		todoTunnels <- tunnel
	}

	go func() {
		for {
			todoTunnel := <-todoTunnels
			err := dialSSHTunnel(todoTunnel)
			if err != nil {
				todoTunnels <- todoTunnel
				time.Sleep(2 * time.Second)
				continue
			}
			todoTunnel.Status = 1
			doingTunnels <- todoTunnel
			go acceptSSHTunnel(todoTunnel)
			time.Sleep(5 * time.Second)
		}
	}()
	go func() {
		keepaliveSSHTunnel(&doingTunnels, &todoTunnels)
	}()

	var done = make(chan bool)
	<-done
}

func dialSSHTunnel(tunnel *SSHTunnel) error {
	userHomeDir, _ := os.UserHomeDir()
	privateKeyPath := filepath.Join(userHomeDir, ".ssh", "id_rsa")
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Printf("Read ssh private key file error: privateKeyPath=%v, err=%v", privateKeyPath, err)
		return err
	}
	privateKey, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		log.Printf("Parse ssh private key error: privateKeyPath=%v, err=%v", privateKeyPath, err)
	}

	clientConfig := &ssh.ClientConfig{
		User: tunnel.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(privateKey),
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", tunnel.SSHIp+":"+strconv.Itoa(tunnel.SSHPort), clientConfig)
	if err != nil {
		log.Printf("Dial tcp to ssh host error: config=%v, err=%v", json_utils.ToJsonStr(tunnel), err)
		return err
	}
	listener, err := sshClient.Listen("tcp", ":"+strconv.Itoa(tunnel.ListenPort))
	if err != nil {
		_ = sshClient.Close()
		log.Printf("Listen tcp to ssh host error: config=%v, err=%v", json_utils.ToJsonStr(tunnel), err)
		return err
	}
	log.Printf("Listening tcp to ssh host: %v", json_utils.ToJsonStr(tunnel))

	tunnel.sshClient = sshClient
	tunnel.listener = &listener
	return nil
}

//goland:noinspection GoUnhandledErrorResult
func acceptSSHTunnel(tunnel *SSHTunnel) {
	listener := *tunnel.listener
	sshClient := tunnel.sshClient
	defer listener.Close()
	defer sshClient.Close()

	for {
		if tunnel.Status == 0 {
			return
		}
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept user connection error: config=%v, err=%v", json_utils.ToJsonStr(tunnel), err)
			return
		}
		go copySSHTunnelData(tunnel, &conn)
	}
}

//goland:noinspection GoUnhandledErrorResult
func copySSHTunnelData(tunnel *SSHTunnel, conn *net.Conn) {
	targetConn, err := net.Dial("tcp", tunnel.TargetIp+":"+strconv.Itoa(tunnel.TargetPort))
	if err != nil {
		log.Printf("Dial tcp to target host error: config=%v, err=%v", json_utils.ToJsonStr(tunnel), err)
		return
	}
	defer (*conn).Close()
	defer targetConn.Close()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		_, err := io.Copy(targetConn, *conn)
		if err != nil {
			log.Printf("Copy user to target error: config=%v, err=%v", json_utils.ToJsonStr(tunnel), err)
		}
		wg.Done()
	}()
	go func() {
		_, err := io.Copy(*conn, targetConn)
		if err != nil {
			log.Printf("Copy target to user error: config=%v, err=%v", json_utils.ToJsonStr(tunnel), err)
		}
		wg.Done()
	}()

	wg.Wait()
}

func keepaliveSSHTunnel(todoTunnels *chan *SSHTunnel, doingTunnels *chan *SSHTunnel) {
	for {
		checkTunnel := <-*doingTunnels
		go func() {
			session, err := checkTunnel.sshClient.NewSession()
			if err != nil {
				log.Printf("NewSession error: tunnel=%v, err=%v", json_utils.ToJsonStr(checkTunnel), err)
				checkTunnel.Status = 0
				*todoTunnels <- checkTunnel
				return
			}
			_, err = session.CombinedOutput("echo 1")
			if err != nil {
				log.Printf("CombinedOutput error: tunnel=%v, err=%v", json_utils.ToJsonStr(checkTunnel), err)
				checkTunnel.Status = 0
				*todoTunnels <- checkTunnel
				return
			}
			*doingTunnels <- checkTunnel
		}()
		time.Sleep(5 * time.Second)
	}
}