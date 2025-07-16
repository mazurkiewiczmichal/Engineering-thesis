package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type HLSServer struct {
	rtspURL   string
	outputDir string
	ffmpegCmd *exec.Cmd
}

func NewHLSServer(rtspURL, outputDir string) *HLSServer {
	return &HLSServer{
		rtspURL:   rtspURL,
		outputDir: outputDir,
	}
}

func (h *HLSServer) Start() error {
	// Create output directory
	if err := os.MkdirAll(h.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Start FFmpeg process
	playlistPath := filepath.Join(h.outputDir, "stream.m3u8")
	segmentPath := filepath.Join(h.outputDir, "segment_%03d.ts")

	h.ffmpegCmd = exec.Command("ffmpeg",
		"-i", h.rtspURL,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-preset", "ultrafast",
		"-tune", "zerolatency",
		"-f", "hls",
		"-hls_time", "2",
		"-hls_list_size", "5",
		"-hls_flags", "delete_segments+split_by_time",
		"-hls_segment_filename", segmentPath,
		"-y", // Overwrite output files
		playlistPath,
	)

	h.ffmpegCmd.Stdout = os.Stdout
	h.ffmpegCmd.Stderr = os.Stderr

	go func() {
		if err := h.ffmpegCmd.Run(); err != nil {
			log.Printf("FFmpeg error: %v", err)
		}
	}()

	// Wait for the playlist file to be created with timeout
	playlistCreated := false
	for i := 0; i < 30; i++ { // Wait up to 30 seconds
		if _, err := os.Stat(playlistPath); err == nil {
			playlistCreated = true
			log.Printf("Playlist created at: %s", playlistPath)
			break
		}
		time.Sleep(1 * time.Second)
		log.Printf("Waiting for playlist creation... (%d/30)", i+1)
	}

	if !playlistCreated {
		return fmt.Errorf("playlist file not created after 30 seconds")
	}

	return nil
}

func (h *HLSServer) Stop() error {
	if h.ffmpegCmd != nil && h.ffmpegCmd.Process != nil {
		return h.ffmpegCmd.Process.Kill()
	}
	return nil
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {
	rtspURL := "rtsp://admin:admin123@192.168.5.132:554/live/ch00_0"
	outputDir := "./hls_output"

	server := NewHLSServer(rtspURL, outputDir)

	// Start the HLS conversion
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start HLS server: %v", err)
	}

	// Add a health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		playlistPath := filepath.Join(outputDir, "stream.m3u8")
		if _, err := os.Stat(playlistPath); err == nil {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Stream is ready")
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "Stream not ready yet")
		}
	})

	// Add a debug endpoint to list files
	http.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		files, err := os.ReadDir(outputDir)
		if err != nil {
			http.Error(w, "Cannot read output directory", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Files in %s:\n", outputDir)
		for _, file := range files {
			fmt.Fprintf(w, "- %s\n", file.Name())
		}
	})

	// Serve HLS files
	http.HandleFunc("/hls/", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)

		// Remove /hls/ prefix and serve from output directory
		path := strings.TrimPrefix(r.URL.Path, "/hls/")
		filePath := filepath.Join(outputDir, path)

		log.Printf("Serving HLS file: %s -> %s", r.URL.Path, filePath)

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			log.Printf("File not found: %s", filePath)
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		// Set appropriate content type
		if filepath.Ext(r.URL.Path) == ".m3u8" {
			w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		} else if filepath.Ext(r.URL.Path) == ".ts" {
			w.Header().Set("Content-Type", "video/MP2T")
		}

		http.ServeFile(w, r, filePath)
	})

	// Serve the HTML page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>RTSP Stream</title>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/hls.js/1.4.10/hls.min.js"></script>
</head>
<body>
    <h1>RTSP Stream</h1>
    <video id="video" controls width="800" height="600"></video>
    
    <script>
        const video = document.getElementById('video');
        const videoSrc = '/hls/stream.m3u8';
        
        if (Hls.isSupported()) {
            const hls = new Hls();
            hls.loadSource(videoSrc);
            hls.attachMedia(video);
            hls.on(Hls.Events.MANIFEST_PARSED, function() {
                video.play();
            });
        } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
            video.src = videoSrc;
            video.addEventListener('loadedmetadata', function() {
                video.play();
            });
        }
    </script>
</body>
</html>`
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})

	// Cleanup on exit
	defer server.Stop()

	fmt.Println("HLS server starting on :8080")
	fmt.Println("Stream will be available at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
