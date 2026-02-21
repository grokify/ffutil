package ffutil

import (
	"bytes"
	"os/exec"
	"runtime"
	"strings"
)

// Encoder represents a video encoder.
type Encoder struct {
	Name        string // Codec name (e.g., "h264_videotoolbox")
	Description string // Human-readable description
	Type        string // "software" or "hardware"
}

// CommonEncoders lists well-known encoders by preference order.
var CommonEncoders = struct {
	// H264 encoders
	H264VideoToolbox Encoder // macOS hardware
	H264NVENC        Encoder // NVIDIA hardware
	H264QSV          Encoder // Intel QuickSync
	H264AMF          Encoder // AMD hardware
	H264VAAPI        Encoder // Linux VA-API
	Libx264          Encoder // Software (universal)

	// HEVC/H265 encoders
	HEVCVideoToolbox Encoder // macOS hardware
	HEVCNVENC        Encoder // NVIDIA hardware
	HEVCQSV          Encoder // Intel QuickSync
	HEVCAMF          Encoder // AMD hardware
	HEVCVAAPI        Encoder // Linux VA-API
	Libx265          Encoder // Software (universal)
}{
	H264VideoToolbox: Encoder{"h264_videotoolbox", "Apple VideoToolbox H.264", "hardware"},
	H264NVENC:        Encoder{"h264_nvenc", "NVIDIA NVENC H.264", "hardware"},
	H264QSV:          Encoder{"h264_qsv", "Intel QuickSync H.264", "hardware"},
	H264AMF:          Encoder{"h264_amf", "AMD AMF H.264", "hardware"},
	H264VAAPI:        Encoder{"h264_vaapi", "VA-API H.264", "hardware"},
	Libx264:          Encoder{"libx264", "x264 H.264 (software)", "software"},

	HEVCVideoToolbox: Encoder{"hevc_videotoolbox", "Apple VideoToolbox HEVC", "hardware"},
	HEVCNVENC:        Encoder{"hevc_nvenc", "NVIDIA NVENC HEVC", "hardware"},
	HEVCQSV:          Encoder{"hevc_qsv", "Intel QuickSync HEVC", "hardware"},
	HEVCAMF:          Encoder{"hevc_amf", "AMD AMF HEVC", "hardware"},
	HEVCVAAPI:        Encoder{"hevc_vaapi", "VA-API HEVC", "hardware"},
	Libx265:          Encoder{"libx265", "x265 HEVC (software)", "software"},
}

// EncoderAvailable checks if a specific encoder is available.
func EncoderAvailable(name string) bool {
	cmd := exec.Command("ffmpeg", "-hide_banner", "-encoders")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), name)
}

// ListEncoders returns all available video encoders.
func ListEncoders() ([]Encoder, error) {
	cmd := exec.Command("ffmpeg", "-hide_banner", "-encoders")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var encoders []Encoder
	lines := strings.Split(stdout.String(), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Encoder lines start with capability flags like "V..... "
		if len(line) < 8 || line[0] != 'V' {
			continue
		}

		// Parse: "V..... codec_name    Description here"
		parts := strings.Fields(line[7:])
		if len(parts) < 2 {
			continue
		}

		name := parts[0]
		desc := strings.Join(parts[1:], " ")

		encType := "software"
		if isHardwareEncoder(name) {
			encType = "hardware"
		}

		encoders = append(encoders, Encoder{
			Name:        name,
			Description: desc,
			Type:        encType,
		})
	}

	return encoders, nil
}

// BestH264Encoder returns the best available H.264 encoder.
// Prefers hardware encoders based on platform, falls back to libx264.
func BestH264Encoder() Encoder {
	// Platform-specific hardware encoder preference
	switch runtime.GOOS {
	case "darwin":
		if EncoderAvailable("h264_videotoolbox") {
			return CommonEncoders.H264VideoToolbox
		}
	case "linux":
		// Check NVIDIA first (most common discrete GPU)
		if EncoderAvailable("h264_nvenc") {
			return CommonEncoders.H264NVENC
		}
		// Intel QuickSync
		if EncoderAvailable("h264_qsv") {
			return CommonEncoders.H264QSV
		}
		// VA-API (generic Linux hardware)
		if EncoderAvailable("h264_vaapi") {
			return CommonEncoders.H264VAAPI
		}
		// AMD
		if EncoderAvailable("h264_amf") {
			return CommonEncoders.H264AMF
		}
	case "windows":
		if EncoderAvailable("h264_nvenc") {
			return CommonEncoders.H264NVENC
		}
		if EncoderAvailable("h264_qsv") {
			return CommonEncoders.H264QSV
		}
		if EncoderAvailable("h264_amf") {
			return CommonEncoders.H264AMF
		}
	}

	// Fallback to software encoder
	return CommonEncoders.Libx264
}

// BestHEVCEncoder returns the best available HEVC/H.265 encoder.
// Prefers hardware encoders based on platform, falls back to libx265.
func BestHEVCEncoder() Encoder {
	switch runtime.GOOS {
	case "darwin":
		if EncoderAvailable("hevc_videotoolbox") {
			return CommonEncoders.HEVCVideoToolbox
		}
	case "linux":
		if EncoderAvailable("hevc_nvenc") {
			return CommonEncoders.HEVCNVENC
		}
		if EncoderAvailable("hevc_qsv") {
			return CommonEncoders.HEVCQSV
		}
		if EncoderAvailable("hevc_vaapi") {
			return CommonEncoders.HEVCVAAPI
		}
		if EncoderAvailable("hevc_amf") {
			return CommonEncoders.HEVCAMF
		}
	case "windows":
		if EncoderAvailable("hevc_nvenc") {
			return CommonEncoders.HEVCNVENC
		}
		if EncoderAvailable("hevc_qsv") {
			return CommonEncoders.HEVCQSV
		}
		if EncoderAvailable("hevc_amf") {
			return CommonEncoders.HEVCAMF
		}
	}

	return CommonEncoders.Libx265
}

// HardwareEncoderAvailable returns true if any hardware encoder is available.
func HardwareEncoderAvailable() bool {
	best := BestH264Encoder()
	return best.Type == "hardware"
}

// isHardwareEncoder checks if an encoder name indicates hardware acceleration.
func isHardwareEncoder(name string) bool {
	hwSuffixes := []string{
		"_videotoolbox",
		"_nvenc",
		"_qsv",
		"_amf",
		"_vaapi",
		"_v4l2m2m",
		"_omx",
		"_cuvid",
		"_mf",
	}
	for _, suffix := range hwSuffixes {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}
