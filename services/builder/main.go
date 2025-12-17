package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// BuildExecutor defines an interface for building images
type BuildExecutor interface {
	Build(ctx context.Context, gitURL, branch, imageName string) error
}

// LocalDockerExecutor builds images using local Docker CLI
type LocalDockerExecutor struct{}

// Build clones the repo, detects a Dockerfile, builds and pushes the image using Docker CLI
func (e *LocalDockerExecutor) Build(ctx context.Context, gitURL, branch, imageName string) error {
	tempDir, err := cloneRepo(gitURL, branch)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	if !hasDockerfile(tempDir) {
		log.Printf("No Dockerfile found in %s, skipping build\n", gitURL)
		return nil
	}

	if err := dockerBuild(tempDir, imageName); err != nil {
		return err
	}

	if err := dockerPush(imageName); err != nil {
		return err
	}

	return nil
}

// KubernetesKanikoExecutor placeholder for production
type KubernetesKanikoExecutor struct{}

func (e *KubernetesKanikoExecutor) Build(ctx context.Context, gitURL, branch, imageName string) error {
	log.Println("Kaniko build requested for", imageName)
	return nil
}

// cloneRepo clones a Git repo into a temp directory
func cloneRepo(repoURL, branch string) (string, error) {
	tempDir, err := os.MkdirTemp("", "repo-*")
	if err != nil {
		return "", err
	}

	args := []string{"clone", "--depth", "1"}
	if branch != "" {
		args = append(args, "-b", branch)
	}
	args = append(args, repoURL, tempDir)

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git clone failed: %w", err)
	}

	return tempDir, nil
}

// hasDockerfile checks if a Dockerfile exists
func hasDockerfile(repoPath string) bool {
	dockerfilePath := filepath.Join(repoPath, "Dockerfile")
	if _, err := os.Stat(dockerfilePath); err == nil {
		return true
	}
	return false
}

// dockerBuild runs `docker build` on the directory
func dockerBuild(contextDir, imageName string) error {
	cmd := exec.Command("docker", "build", "-t", imageName, contextDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("Building image %s...\n", imageName)
	return cmd.Run()
}

// dockerPush runs `docker push` for the image
func dockerPush(imageName string) error {
	cmd := exec.Command("docker", "push", imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("Pushing image %s...\n", imageName)
	return cmd.Run()
}

func main() {
	ctx := context.Background()

	// Configuration
	repoURL := "https://github.com/youruser/yourrepo.git"
	branch := "main"
	imageName := "localhost:5000/yourrepo:latest"
	runningLocally := true

	var executor BuildExecutor
	if runningLocally {
		executor = &LocalDockerExecutor{}
	} else {
		executor = &KubernetesKanikoExecutor{}
	}

	if err := executor.Build(ctx, repoURL, branch, imageName); err != nil {
		log.Fatalf("Build failed: %v", err)
	}

	log.Println("Build finished successfully!")
}
