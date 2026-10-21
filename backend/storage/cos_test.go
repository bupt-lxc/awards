package storage

import (
	"bytes"
	"io"
	"path/filepath"
	"testing"
)

func TestBuildPublicPath(t *testing.T) {
	tests := []struct {
		name       string
		publicURL  string
		objectName string
		want       string
	}{
		{
			name:       "default proxy path does not duplicate certificates segment",
			publicURL:  "/static/certificates",
			objectName: "certificates/example.png",
			want:       "/static/certificates/example.png",
		},
		{
			name:       "empty URL uses local proxy path",
			publicURL:  "",
			objectName: "certificates/example.png",
			want:       "/static/certificates/example.png",
		},
		{
			name:       "CDN URL keeps object prefix",
			publicURL:  "https://cdn.example.com/awards/",
			objectName: "certificates/example.png",
			want:       "https://cdn.example.com/awards/certificates/example.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildPublicPath(tt.publicURL, tt.objectName); got != tt.want {
				t.Fatalf("buildPublicPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewCOSServiceFallsBackWhenConfigurationIsIncomplete(t *testing.T) {
	service, err := NewCOSService(COSConfig{})
	if err != nil {
		t.Fatalf("NewCOSService() returned error: %v", err)
	}
	fallback, ok := service.(*localFallback)
	if !ok {
		t.Fatalf("NewCOSService() = %T, want *localFallback", service)
	}
	if fallback.publicURL != "/static/certificates" {
		t.Fatalf("fallback publicURL = %q, want /static/certificates", fallback.publicURL)
	}
}

func TestLocalFallbackUploadGetDelete(t *testing.T) {
	dir := t.TempDir()
	service := &localFallback{publicURL: "/static/certificates", dir: dir}
	want := []byte("certificate-image")

	publicPath, err := service.Upload("award.png", bytes.NewReader(want), int64(len(want)), "image/png")
	if err != nil {
		t.Fatalf("Upload() returned error: %v", err)
	}
	if filepath.Ext(publicPath) != ".png" {
		t.Fatalf("Upload() path = %q, want .png extension", publicPath)
	}

	reader, contentType, err := service.Get(publicPath)
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() returned error: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("Close() returned error: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("Get() content = %q, want %q", got, want)
	}
	if contentType != "image/png" {
		t.Fatalf("Get() content type = %q, want image/png", contentType)
	}

	if err := service.Delete(publicPath); err != nil {
		t.Fatalf("Delete() returned error: %v", err)
	}
	if _, _, err := service.Get(publicPath); err == nil {
		t.Fatal("Get() after Delete() returned nil error")
	}
}
