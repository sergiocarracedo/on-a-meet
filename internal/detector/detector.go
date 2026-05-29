package detector

import "fmt"

func New(method string) (Detector, error) {
	switch method {
	case "v4l2":
		return NewV4L2Detector(), nil
	case "lsof":
		return NewLsofDetector(), nil
	default:
		return nil, fmt.Errorf("unknown detection method: %q (supported: v4l2, lsof)", method)
	}
}
