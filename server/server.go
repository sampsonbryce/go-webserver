package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	host     string
	port     string
	connType string
	handlers map[string]func(Request) Response
}

type ServerOptions struct {
	Host     string
	Port     string
	ConnType string
}

type Request struct {
	Method   string
	Uri      string
	Protocol string
	Headers  map[string]string
	Body     string
}

type Response struct {
	StatusCode int
}

func (s *Server) Listen() error {
	l, err := net.Listen(s.connType, s.host+":"+s.port)

	if err != nil {
		return err
	}

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			continue
		}

		fmt.Println("Client connected.")

		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

		go s.handleConnection(c)
	}
}

func (s *Server) HandleFunc(path string, handler func(Request) Response) {
	s.handlers[path] = handler
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := s.parseRequest(conn)

	if err != nil {
		log.Println(err)
		return
	}

	handler := s.findHandler(req.Uri)

	response := handler(req)

	s.writeResponse(response, conn)
}

func (s *Server) findHandler(path string) func(Request) Response {
	if val, ok := s.handlers[path]; ok {
		return val
	}

	// TODO allow for custom handler
	return func(req Request) Response {
		return Response{StatusCode: 404}
	}
}

func (s *Server) parseRequest(conn net.Conn) (Request, error) {
	request := Request{}
	reader := bufio.NewReader(conn)

	// Request Line
	err := s.parseRequestLine(reader, &request)
	if err != nil {
		return request, err
	}

	// Headers
	err = s.parseHeaders(reader, &request)
	if err != nil {
		return request, err
	}

	if contentLengthStr, ok := request.Headers["Content-Length"]; ok {
		contentLength, err := strconv.Atoi(contentLengthStr)
		if err != nil {
			return request, err
		}

		if contentLength > 0 {
			err = s.parseBody(reader, contentLength, &request)
			if err != nil {
				return request, err
			}
		}
	}

	return request, nil
}

func (s *Server) parseRequestLine(reader *bufio.Reader, request *Request) error {
	requestLineBuffer, err := reader.ReadString('\n')

	if err != nil {
		return err
	}

	requestLine := string(requestLineBuffer[:len(requestLineBuffer)-1])

	splitRequestLine := strings.Split(requestLine, " ")

	if len(splitRequestLine) != 3 {
		return fmt.Errorf("invalid request line")
	}

	request.Method = splitRequestLine[0]
	request.Uri = splitRequestLine[1]
	request.Protocol = splitRequestLine[2]

	return nil
}

func (s *Server) parseHeaders(reader *bufio.Reader, request *Request) error {
	request.Headers = make(map[string]string)

	for {
		buffer, err := reader.ReadString('\n')

		if err != nil {
			return err
		}

		line := string(buffer[:len(buffer)-1])
		lineSplit := strings.Split(line, ":")

		headerName := strings.TrimSpace(lineSplit[0])
		headerValue := strings.TrimSpace(strings.Join(lineSplit[1:], ":"))

		request.Headers[headerName] = headerValue

		if len(buffer) == 2 {
			break
		}
	}

	return nil
}

func (s *Server) parseBody(reader *bufio.Reader, contentLength int, request *Request) error {
	readLength := 0
	for {
		line, err := reader.ReadBytes('\n')

		if err != nil {
			return err
		}

		readLength += len(line)

		request.Body += string(line)

		if readLength >= contentLength {
			break
		}
	}

	return nil
}

func (s *Server) writeResponse(response Response, conn net.Conn) {
	conn.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", response.StatusCode, "OK")))

	layout := "Mon, 02 Jan 2006 15:04:05"
	loc, _ := time.LoadLocation("UTC")

	conn.Write([]byte(fmt.Sprintf("Date: %s GMT\r\n", time.Now().In(loc).Format(layout))))
	conn.Write([]byte("Server: Super Cool Awesome Server\r\n"))

	conn.Write([]byte("\r\n"))
}

func CreateServer(options *ServerOptions) Server {
	connType := "tcp"

	if options.ConnType != "" {
		connType = options.ConnType
	}

	handlers := make(map[string]func(Request) Response)
	return Server{
		host:     options.Host,
		port:     options.Port,
		connType: connType,
		handlers: handlers,
	}
}
