package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "Argo Demo"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "msg",
			Usage:   "the message to include in responses",
			EnvVars: []string{"MSG"},
			Value:   "hello world",
		},
		&cli.StringFlag{
			Name:  "listen",
			Usage: "the host:port pair to listen on",
			Value: "0.0.0.0:8080",
		},
	}

	app.Action = run

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

type Response struct {
	Msg     string              `json:"msg"`
	Ip      string              `json:"ip"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

type Server struct {
	Msg string
}

func (s *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	outResp := &Response{
		Msg:     s.Msg,
		Ip:      req.RemoteAddr,
		Headers: req.Header,
		Body:    string(body),
	}

	outBytes, err := json.MarshalIndent(outResp, "", "  ")
	if err != nil {
		panic(err)
	}

	resp.Header().Add("Content-type", "application/json")
	resp.Write(outBytes)
}

func run(c *cli.Context) error {
	msg := c.String("msg")
	addr := c.String("listen")
	server := NewServer(msg)
	http.ListenAndServe(addr, server)
	return nil
}

func NewServer(msg string) *Server {
	return &Server{
		Msg: msg,
	}
}
