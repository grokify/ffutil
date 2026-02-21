package ffutil

import (
	"testing"
)

func TestEncoderAvailable(t *testing.T) {
	// libx264 should be available on most systems with ffmpeg
	if !FFmpegAvailable() {
		t.Skip("ffmpeg not available")
	}

	// Test a common encoder that should exist
	// We don't assert true/false since it depends on the system
	_ = EncoderAvailable("libx264")
}

func TestListEncoders(t *testing.T) {
	if !FFmpegAvailable() {
		t.Skip("ffmpeg not available")
	}

	encoders, err := ListEncoders()
	if err != nil {
		t.Fatalf("ListEncoders() error: %v", err)
	}

	if len(encoders) == 0 {
		t.Error("ListEncoders() returned empty list")
	}

	// Should have at least one encoder
	found := false
	for _, enc := range encoders {
		if enc.Name != "" {
			found = true
			break
		}
	}
	if !found {
		t.Error("ListEncoders() returned no named encoders")
	}
}

func TestBestH264Encoder(t *testing.T) {
	if !FFmpegAvailable() {
		t.Skip("ffmpeg not available")
	}

	enc := BestH264Encoder()

	if enc.Name == "" {
		t.Error("BestH264Encoder() returned empty name")
	}

	if enc.Type != "hardware" && enc.Type != "software" {
		t.Errorf("BestH264Encoder() invalid type: %s", enc.Type)
	}
}

func TestBestHEVCEncoder(t *testing.T) {
	if !FFmpegAvailable() {
		t.Skip("ffmpeg not available")
	}

	enc := BestHEVCEncoder()

	if enc.Name == "" {
		t.Error("BestHEVCEncoder() returned empty name")
	}

	if enc.Type != "hardware" && enc.Type != "software" {
		t.Errorf("BestHEVCEncoder() invalid type: %s", enc.Type)
	}
}

func TestHardwareEncoderAvailable(t *testing.T) {
	if !FFmpegAvailable() {
		t.Skip("ffmpeg not available")
	}

	// Just verify it doesn't panic
	_ = HardwareEncoderAvailable()
}

func TestIsHardwareEncoder(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"h264_videotoolbox", true},
		{"h264_nvenc", true},
		{"h264_qsv", true},
		{"h264_amf", true},
		{"h264_vaapi", true},
		{"h264_v4l2m2m", true},
		{"libx264", false},
		{"libx265", false},
		{"mpeg4", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isHardwareEncoder(tt.name)
			if got != tt.expected {
				t.Errorf("isHardwareEncoder(%q) = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}

func TestCommonEncoders(t *testing.T) {
	// Verify all common encoders have proper names and types
	encoders := []Encoder{
		CommonEncoders.H264VideoToolbox,
		CommonEncoders.H264NVENC,
		CommonEncoders.H264QSV,
		CommonEncoders.H264AMF,
		CommonEncoders.H264VAAPI,
		CommonEncoders.Libx264,
		CommonEncoders.HEVCVideoToolbox,
		CommonEncoders.HEVCNVENC,
		CommonEncoders.HEVCQSV,
		CommonEncoders.HEVCAMF,
		CommonEncoders.HEVCVAAPI,
		CommonEncoders.Libx265,
	}

	for _, enc := range encoders {
		if enc.Name == "" {
			t.Error("CommonEncoder has empty name")
		}
		if enc.Description == "" {
			t.Errorf("CommonEncoder %q has empty description", enc.Name)
		}
		if enc.Type != "hardware" && enc.Type != "software" {
			t.Errorf("CommonEncoder %q has invalid type: %s", enc.Name, enc.Type)
		}
	}
}
