package service

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/domain"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/dto"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/event"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/repository"
	"github.com/google/uuid"
)

type DiscoveryService struct {
	deviceRepo *repository.DiscoveredDeviceRepository
	publisher  event.Publisher
}

func NewDiscoveryService(
	deviceRepo *repository.DiscoveredDeviceRepository,
	publisher event.Publisher,
) *DiscoveryService {
	return &DiscoveryService{
		deviceRepo: deviceRepo,
		publisher:  publisher,
	}
}

func (s *DiscoveryService) DiscoverCIDR(
	ctx context.Context,
	req dto.DiscoverCamerasRequest,
) ([]domain.DiscoveredDevice, error) {
	if req.TimeoutMs <= 0 {
		req.TimeoutMs = 500
	}

	if req.MaxWorkers <= 0 {
		req.MaxWorkers = 64
	}

	if len(req.Ports) == 0 {
		req.Ports = []int{80, 443, 554, 8554, 8000, 8080, 8899}
	}

	ips, err := expandCIDR(req.NetworkCIDR)
	if err != nil {
		return nil, err
	}

	timeout := time.Duration(req.TimeoutMs) * time.Millisecond

	jobs := make(chan string)
	results := make(chan domain.DiscoveredDevice)

	var wg sync.WaitGroup

	for i := 0; i < req.MaxWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for ip := range jobs {
				device, ok := scanIP(ip, req.Ports, timeout)
				if !ok {
					continue
				}

				saved, err := s.deviceRepo.Upsert(ctx, device)
				if err != nil {
					fmt.Printf("failed to save discovered device %s: %v\n", device.IPAddress, err)
					continue
				}

				fmt.Printf("saved discovered device %s type=%s confidence=%.2f reason=%v\n",
					saved.IPAddress,
					saved.DeviceType,
					saved.Confidence,
					saved.DetectionReason,
				)

				if s.publisher != nil {
					if err := s.publisher.PublishCameraEvent(ctx, event.CameraEvent{
						EventID:            uuid.New(),
						EventType:          event.CameraDiscoveredEvent,
						DiscoveredDeviceID: &saved.ID,
						IPAddress:          saved.IPAddress,
						Timestamp:          time.Now().UTC(),
					}); err != nil {
						fmt.Printf("failed to publish camera.discovered event for %s: %v\n", saved.IPAddress, err)
					}
				}

				results <- *saved
			}
		}()
	}

	go func() {
		for _, ip := range ips {
			jobs <- ip
		}

		close(jobs)
		wg.Wait()
		close(results)
	}()

	devices := make([]domain.DiscoveredDevice, 0)

	for device := range results {
		devices = append(devices, device)
	}

	sort.Slice(devices, func(i, j int) bool {
		return devices[i].IPAddress < devices[j].IPAddress
	})

	return devices, nil
}

func (s *DiscoveryService) GetDiscoveredDevices(
	ctx context.Context,
) ([]domain.DiscoveredDevice, error) {
	return s.deviceRepo.FindAll(ctx)
}

func (s *DiscoveryService) GetDiscoveredDeviceByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.DiscoveredDevice, error) {
	return s.deviceRepo.FindByID(ctx, id)
}

func scanIP(ip string, ports []int, timeout time.Duration) (domain.DiscoveredDevice, bool) {
	var httpPort *int
	var rtspPort *int
	var onvifPort *int

	httpSupported := false
	rtspSupported := false
	onvifSupported := false

	deviceType := domain.DeviceTypeUnknown
	confidence := 0.0
	detectionReason := "no_camera_fingerprint_detected"

	for _, port := range ports {
		if !canConnect(ip, port, timeout) {
			continue
		}

		p := port

		switch port {
		case 554, 8554:
			rtspSupported = true
			rtspPort = &p

			deviceType = domain.DeviceTypePossibleCamera
			confidence = maxConfidence(confidence, 0.5)
			detectionReason = appendReason(detectionReason, "rtsp_port_open")

		case 80, 443, 8000, 8080, 8899:
			httpSupported = true
			httpPort = &p

			if port == 80 || port == 8899 || port == 8000 {
				onvifSupported = true
				onvifPort = &p
			}

			url := fmt.Sprintf("http://%s:%d/", ip, port)
			fp := fingerprintHTTP(url, timeout)

			if fp.Confidence > confidence {
				deviceType = fp.DeviceType
				confidence = fp.Confidence
				detectionReason = fp.DetectionReason
			}
		}
	}

	if !httpSupported && !rtspSupported && !onvifSupported {
		return domain.DiscoveredDevice{}, false
	}

	if deviceType != domain.DeviceTypeCamera &&
		deviceType != domain.DeviceTypePossibleCamera &&
		!rtspSupported {
		return domain.DiscoveredDevice{}, false
	}

	return domain.DiscoveredDevice{
		IPAddress:       ip,
		ONVIFSupported:  onvifSupported,
		RTSPSupported:   rtspSupported,
		HTTPSupported:   httpSupported,
		HTTPPort:        httpPort,
		RTSPPort:        rtspPort,
		ONVIFPort:       onvifPort,
		DiscoveryMethod: domain.DiscoveryMethodCIDRScan,
		Status:          domain.DiscoveredDeviceStatusDiscovered,
		DeviceType:      deviceType,
		Confidence:      confidence,
		DetectionReason: &detectionReason,
	}, true
}

type FingerprintResult struct {
	DeviceType      domain.DeviceType
	Confidence      float64
	DetectionReason string
}

func fingerprintHTTP(url string, timeout time.Duration) FingerprintResult {
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return FingerprintResult{
			DeviceType:      domain.DeviceTypeUnknown,
			Confidence:      0.1,
			DetectionReason: "http_port_open_but_no_response",
		}
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	body := strings.ToLower(string(bodyBytes))

	server := strings.ToLower(resp.Header.Get("Server"))
	auth := strings.ToLower(resp.Header.Get("WWW-Authenticate"))
	contentType := strings.ToLower(resp.Header.Get("Content-Type"))

	score := 0.0
	reasons := make([]string, 0)

	keywords := []string{
		"ip webcam",
		"ip camera",
		"onvif",
		"rtsp",
		"mjpeg",
		"video feed",
		"videofeed",
		"snapshot",
		"shot.jpg",
		"axis",
		"hikvision",
		"dahua",
		"reolink",
		"foscam",
		"amcrest",
		"vivotek",
	}

	for _, keyword := range keywords {
		if strings.Contains(body, keyword) ||
			strings.Contains(server, keyword) ||
			strings.Contains(auth, keyword) {
			score += 0.2
			reasons = append(reasons, keyword)
		}
	}

	if strings.Contains(body, "ip webcam") {
		score += 0.5
		reasons = append(reasons, "android_ip_webcam_detected")
	}

	if strings.Contains(body, "videofeed") ||
		strings.Contains(body, "video") ||
		strings.Contains(body, "shot.jpg") {
		score += 0.3
		reasons = append(reasons, "stream_endpoint_detected")
	}

	if strings.Contains(contentType, "multipart") ||
		strings.Contains(contentType, "image/jpeg") {
		score += 0.2
		reasons = append(reasons, "camera_like_content_type")
	}

	if score > 1 {
		score = 1
	}

	reason := strings.Join(reasons, ",")

	if score >= 0.7 {
		return FingerprintResult{
			DeviceType:      domain.DeviceTypeCamera,
			Confidence:      score,
			DetectionReason: reason,
		}
	}

	if score >= 0.3 {
		return FingerprintResult{
			DeviceType:      domain.DeviceTypePossibleCamera,
			Confidence:      score,
			DetectionReason: reason,
		}
	}

	return FingerprintResult{
		DeviceType:      domain.DeviceTypeNonCamera,
		Confidence:      score,
		DetectionReason: "no_camera_fingerprint_detected",
	}
}

func canConnect(ip string, port int, timeout time.Duration) bool {
	address := net.JoinHostPort(ip, fmt.Sprintf("%d", port))

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}

	_ = conn.Close()
	fmt.Printf("OPEN %s\n", address)
	return true
}

func expandCIDR(cidr string) ([]string, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	ips := make([]string, 0)

	for currentIP := ip.Mask(ipNet.Mask); ipNet.Contains(currentIP); incIP(currentIP) {
		ipCopy := make(net.IP, len(currentIP))
		copy(ipCopy, currentIP)

		ips = append(ips, ipCopy.String())
	}

	if len(ips) <= 2 {
		return ips, nil
	}

	return ips[1 : len(ips)-1], nil
}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++

		if ip[j] > 0 {
			break
		}
	}
}

func maxConfidence(a float64, b float64) float64 {
	if a > b {
		return a
	}

	return b
}

func appendReason(existing string, reason string) string {
	if existing == "" || existing == "no_camera_fingerprint_detected" {
		return reason
	}

	return existing + "," + reason
}
