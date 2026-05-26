package container

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
)

var (
	DockerCli      *client.Client
	containerMutex sync.Mutex
	upgrader       = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for dev
		},
	}

	// Connection tracking for self-healing idle pruning
	activeConnections = make(map[string]int)
	connectionsMutex  sync.Mutex
)

const GameImageName = "learn-linux-game-env:latest"
const FallbackImageName = "ubuntu:22.04"

// InitDocker initializes the connection to the Docker daemon.
func InitDocker() error {
	var err error
	DockerCli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to init Docker client: %w", err)
	}
	log.Println("Docker client initialized successfully.")
	return nil
}

// GetOrCreateContainer returns the container ID for a user. If not running/created, spins it up.
func GetOrCreateContainer(userID string, username string) (string, error) {
	containerMutex.Lock()
	defer containerMutex.Unlock()

	ctx := context.Background()
	// Use username instead of userID for readable container names
	containerName := fmt.Sprintf("learn-linux-player-%s", username)

	// Check if container already exists
	inspect, err := DockerCli.ContainerInspect(ctx, containerName)
	if err == nil {
		// Container exists
		if !inspect.State.Running {
			// Start it if stopped
			log.Printf("Starting existing container: %s", containerName)
			err = DockerCli.ContainerStart(ctx, inspect.ID, types.ContainerStartOptions{})
			if err != nil {
				return "", fmt.Errorf("failed to start container: %w", err)
			}
		}
		return inspect.ID, nil
	}

	// Container does not exist, create it.
	// First check if game environment image is available, otherwise use standard Ubuntu fallback
	useImage := GameImageName
	images, err := DockerCli.ImageList(ctx, types.ImageListOptions{})
	found := false
	if err == nil {
		for _, img := range images {
			for _, repoTag := range img.RepoTags {
				if repoTag == GameImageName {
					found = true
					break
				}
			}
		}
	}

	if !found {
		log.Printf("Game image '%s' not found locally. Checking fallback '%s'...", GameImageName, FallbackImageName)
		useImage = FallbackImageName
		// Pull fallback image if not present
		reader, err := DockerCli.ImagePull(ctx, FallbackImageName, types.ImagePullOptions{})
		if err == nil {
			// Must read the pull response to actually trigger the pull completely
			_, _ = io.Copy(io.Discard, reader)
			reader.Close()
		} else {
			log.Printf("Warning: Failed to pull fallback image %s: %v. Attempting to run anyway.", FallbackImageName, err)
		}
	}

	log.Printf("Creating new container '%s' using image '%s'", containerName, useImage)

	// Setup resource limits to prevent server overloading (512MB RAM limit to support apt update/install)
	resources := dockercontainer.Resources{
		Memory:   512 * 1024 * 1024, // 512MB RAM limit
		NanoCPUs: 200000000,         // 0.2 CPU limit
	}

	cc := &dockercontainer.Config{
		Image:        useImage,
		Tty:          true,
		OpenStdin:    true,
		StdinOnce:    false,
		WorkingDir:   "/home/player",
		User:         "player",
		Env:          []string{"TERM=xterm-256color"},
		Cmd:          []string{"/bin/bash"},
	}

	// If using fallback, we run as root initially, then configure player user inside
	if useImage == FallbackImageName {
		cc.User = "root"
		cc.WorkingDir = "/root"
	}

	hc := &dockercontainer.HostConfig{
		Resources: resources,
	}

	resp, err := DockerCli.ContainerCreate(ctx, cc, hc, nil, nil, containerName)
	if err != nil {
		// Fallback check in case container was created concurrently right before this create call
		if strings.Contains(err.Error(), "Conflict") || strings.Contains(err.Error(), "already in use") {
			inspectAgain, errInspect := DockerCli.ContainerInspect(ctx, containerName)
			if errInspect == nil {
				if !inspectAgain.State.Running {
					errStart := DockerCli.ContainerStart(ctx, inspectAgain.ID, types.ContainerStartOptions{})
					if errStart != nil {
						return "", fmt.Errorf("failed to start container on fallback: %w", errStart)
					}
				}
				return inspectAgain.ID, nil
			}
		}
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	err = DockerCli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to start container after creation: %w", err)
	}

	// If fallback image was used, dynamically configure it now
	if useImage == FallbackImageName {
		log.Println("Setting up game state inside fallback Ubuntu container...")
		setupCmds := []string{
			"apt-get update && apt-get install -y sudo curl nano vim grep findutils coreutils net-tools iputils-ping",
			"useradd -m -s /bin/bash player",
			"echo 'player:player' | chpasswd",
			"usermod -aG sudo player",
			"echo 'player ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers",
			"mkdir -p /var/mainframe/archives /var/mainframe/cores /var/mainframe/terminal_history",
			"echo 'Archive systems online.' > /var/mainframe/archives/system_backup.txt",
			"echo 'Core offline.' > /var/mainframe/cores/core_state.txt",
			"echo '[2026-05-26 00:00:01] Mainframe system starting up...\n[2026-05-26 00:00:15] Security modules initialized.\n[2026-05-26 00:00:30] WARNING: Breach detected in Sector 7.' > /var/log/mainframe.log",
			"echo '[2026-05-26 00:01:00] SYSTEM BOOTING: SUCCESS\n[2026-05-26 00:02:15] SECTOR 1: ACTIVE (TEMP 45C)\n[2026-05-26 00:03:30] SECTOR 2: ACTIVE (TEMP 48C)\n[2026-05-26 00:04:12] WARNING: FUEL LEVEL LOW\n[2026-05-26 00:05:00] SECTOR 3: ACTIVE (TEMP 52C)\n[2026-05-26 00:05:10] CRITICAL: REACTOR OVERHEAT PASSWORD: CORE_TEMP_CRITICAL_9982\n[2026-05-26 00:06:00] SECTOR 4: ACTIVE (TEMP 55C)\n[2026-05-26 00:07:05] ALERT: COOLING SYSTEM FAILURE\n[2026-05-26 00:08:00] SECTOR 5: INACTIVE (MAINTENANCE)' > /var/log/reactor.log",
			"echo 'GRID_PASS_99' > /home/player/.cyber_key",
			"chown -R player:player /home/player",
			"chmod 600 /home/player/.cyber_key",
		}
		for _, cmd := range setupCmds {
			_, _, _ = RunExecCommand(resp.ID, "root", strings.Split(cmd, " "))
		}
	}

	return resp.ID, nil
}

// RunExecCommand runs a quick execution command in a container and returns stdout/stderr and exit code.
func RunExecCommand(containerID string, user string, cmd []string) (string, int, error) {
	ctx := context.Background()

	execConfig := types.ExecConfig{
		User:         user,
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	execCreateResp, err := DockerCli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return "", -1, err
	}

	resp, err := DockerCli.ContainerExecAttach(ctx, execCreateResp.ID, types.ExecStartCheck{})
	if err != nil {
		return "", -1, err
	}
	defer resp.Close()

	outputBytes, err := io.ReadAll(resp.Reader)
	if err != nil {
		return "", -1, err
	}

	inspectResp, err := DockerCli.ContainerExecInspect(ctx, execCreateResp.ID)
	if err != nil {
		return string(outputBytes), -1, err
	}

	return string(outputBytes), inspectResp.ExitCode, nil
}

// HandleTerminalWS upgrades the HTTP request to WebSocket and proxies input/output to the container's shell.
func HandleTerminalWS(w http.ResponseWriter, r *http.Request, userID string, username string, containerID string) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer wsConn.Close()

	// Track connection increment
	connectionsMutex.Lock()
	activeConnections[userID]++
	log.Printf("[Network Hub] User %s connected. Active connections: %d", username, activeConnections[userID])
	connectionsMutex.Unlock()

	// Clean up connection decrement on termination
	defer func() {
		connectionsMutex.Lock()
		activeConnections[userID]--
		count := activeConnections[userID]
		log.Printf("[Network Hub] User %s disconnected. Active connections remaining: %d", username, count)
		connectionsMutex.Unlock()

		// If no active connections remain, trigger self-healing stop timer
		if count <= 0 {
			go func(uid string, uname string) {
				// 2-minute grace period
				time.Sleep(2 * time.Minute)

				connectionsMutex.Lock()
				currentCount := activeConnections[uid]
				connectionsMutex.Unlock()

				if currentCount <= 0 {
					log.Printf("[Network Watchdog] User %s (%s) offline for 2m. Stopping container...", uname, uid)
					ctxStop := context.Background()
					containerName := fmt.Sprintf("learn-linux-player-%s", uname)
					_ = DockerCli.ContainerStop(ctxStop, containerName, dockercontainer.StopOptions{})
				}
			}(userID, username)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create Exec session running /bin/bash as user 'player'
	execConfig := types.ExecConfig{
		User:         "player",
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   "/home/player",
		Cmd:          []string{"/bin/bash"},
	}

	execIDResp, err := DockerCli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		log.Printf("Docker Exec Create failed: %v", err)
		_ = wsConn.WriteMessage(websocket.TextMessage, []byte("\r\n[System Error: Failed to start bash shell inside sandbox]\r\n"))
		return
	}

	// Attach to Exec session
	resp, err := DockerCli.ContainerExecAttach(ctx, execIDResp.ID, types.ExecStartCheck{
		Tty: true,
	})
	if err != nil {
		log.Printf("Docker Exec Attach failed: %v", err)
		return
	}
	defer resp.Close()

	// Bidirectional piping
	errChan := make(chan error, 2)

	// Goroutine: Read from Container Exec output -> Write to WebSocket
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := resp.Reader.Read(buf)
			if n > 0 {
				err = wsConn.WriteMessage(websocket.BinaryMessage, buf[:n])
				if err != nil {
					errChan <- err
					return
				}
			}
			if err != nil {
				errChan <- err
				return
			}
		}
	}()

	// Goroutine: Read from WebSocket -> Write to Container Exec input
	go func() {
		for {
			_, msg, err := wsConn.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}
			_, err = resp.Conn.Write(msg)
			if err != nil {
				errChan <- err
				return
			}
		}
	}()

	// Keep alive / wait for connection drop
	select {
	case err := <-errChan:
		if err != io.EOF && !websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
			log.Printf("Terminal stream ended with error: %v", err)
		}
	case <-ctx.Done():
	}
}
