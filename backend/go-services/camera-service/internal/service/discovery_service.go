package service

import (
	"context"
	"fmt"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/domain"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/dto"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/repository"
)

type DiscoveryService struct {
	deviceRepo *repository.DiscoveredDeviceRepository
}

func NewDiscoveryService(deviceRepo *repository.DiscoveredDeviceRepository) *DiscoveryService {
	return &DiscoveryService{deviceRepo: deviceRepo}
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

				fmt.Printf("saved discovered device %s\n", saved.IPAddress)
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

func scanIP(ip string, ports []int, timeout time.Duration) (domain.DiscoveredDevice, bool) {
	var httpPort *int
	var rtspPort *int
	var onvifPort *int

	httpSupported := false
	rtspSupported := false
	onvifSupported := false

	for _, port := range ports {
		if canConnect(ip, port, timeout) {
			p := port

			switch port {
			case 554, 8554:
				rtspSupported = true
				rtspPort = &p
			case 80, 443, 8000, 8080, 8899:
				httpSupported = true
				httpPort = &p
			}

			if port == 80 || port == 8899 || port == 8000 {
				onvifSupported = true
				onvifPort = &p
			}
		}
	}

	if !httpSupported && !rtspSupported && !onvifSupported {
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
	}, true
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
