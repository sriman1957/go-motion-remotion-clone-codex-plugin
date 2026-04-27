package mcp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Request struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      any            `json:"id,omitempty"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params,omitempty"`
}

type Response struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      any            `json:"id,omitempty"`
	Result  map[string]any `json:"result,omitempty"`
	Error   *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ServeStdio(r io.Reader, w io.Writer, srv *Server) error {
	reader := bufio.NewReader(r)
	for {
		payload, err := readFrame(reader)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		var req Request
		if err := json.Unmarshal(payload, &req); err != nil {
			if err := writeFrame(w, Response{
				JSONRPC: "2.0",
				Error: &ResponseError{
					Code:    -32700,
					Message: fmt.Sprintf("invalid json: %v", err),
				},
			}); err != nil {
				return err
			}
			continue
		}

		result, shouldReply, dispatchErr := srv.DispatchRequest(req.Method, req.Params)
		if !shouldReply {
			continue
		}

		resp := Response{
			JSONRPC: "2.0",
			ID:      req.ID,
		}
		if dispatchErr != nil {
			resp.Error = &ResponseError{
				Code:    -32000,
				Message: dispatchErr.Error(),
			}
		} else {
			resp.Result = result
		}
		if err := writeFrame(w, resp); err != nil {
			return err
		}
	}
}

func readFrame(reader *bufio.Reader) ([]byte, error) {
	contentLength := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}
		if strings.HasPrefix(strings.ToLower(line), "content-length:") {
			value := strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:"))
			value = strings.TrimSpace(strings.TrimPrefix(value, ":"))
			contentLength, err = strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid content length: %w", err)
			}
		}
	}
	if contentLength <= 0 {
		return nil, io.EOF
	}

	payload := make([]byte, contentLength)
	if _, err := io.ReadFull(reader, payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func writeFrame(w io.Writer, v any) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("Content-Length: %d\r\n\r\n", len(body)))
	buf.Write(body)
	_, err = w.Write(buf.Bytes())
	return err
}
