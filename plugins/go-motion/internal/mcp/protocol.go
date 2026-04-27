package mcp

import "fmt"

func (s *Server) DispatchRequest(method string, params map[string]any) (map[string]any, bool, error) {
	switch method {
	case "initialize":
		return map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]any{
				"tools": map[string]any{
					"listChanged": false,
				},
			},
			"serverInfo": map[string]any{
				"name":    "go-motion",
				"version": "0.1.0",
			},
		}, true, nil
	case "notifications/initialized":
		return nil, false, nil
	case "tools/list":
		return map[string]any{
			"tools": s.ListTools(),
		}, true, nil
	case "tools/call":
		if params == nil {
			return nil, true, fmt.Errorf("missing params")
		}
		name, _ := params["name"].(string)
		args, _ := params["arguments"].(map[string]any)
		return s.callTool(name, args)
	default:
		return nil, true, fmt.Errorf("unsupported method: %s", method)
	}
}

func (s *Server) callTool(name string, args map[string]any) (map[string]any, bool, error) {
	switch name {
	case "generate_video":
		prompt, _ := args["prompt"].(string)
		result, err := s.GenerateVideo(GenerateVideoRequest{Prompt: prompt})
		if err != nil {
			return nil, true, err
		}
		message := "Video job planned successfully."
		if result.Status == "rendered" {
			message = "Video rendered successfully."
		}
		if result.Status == "runtime-missing" {
			message = "Video composition created, but a required runtime is missing."
		}
		if result.Status == "render-failed" {
			message = "Video composition created, but rendering failed."
		}
		return map[string]any{
			"content": []map[string]any{
				{
					"type": "text",
					"text": message,
				},
			},
			"structuredContent": result,
		}, true, nil
	case "list_styles":
		return map[string]any{
			"content": []map[string]any{
				{
					"type": "text",
					"text": "Available styles: cinematic, kinetic-typography, product-promo, launch, minimal-tech.",
				},
			},
			"structuredContent": map[string]any{
				"styles": []string{"cinematic", "kinetic-typography", "product-promo", "launch", "minimal-tech"},
			},
		}, true, nil
	default:
		return nil, true, fmt.Errorf("unsupported tool: %s", name)
	}
}
