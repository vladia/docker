package main

import (
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestPsListContainers(t *testing.T) {
	runCmd := exec.Command(dockerBinary, "run", "-d", "busybox", "top")
	out, _, err := runCommandWithOutput(runCmd)
	errorOut(err, t, out)
	firstID := stripTrailingCharacters(out)

	runCmd = exec.Command(dockerBinary, "run", "-d", "busybox", "top")
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)
	secondID := stripTrailingCharacters(out)

	// not long running
	runCmd = exec.Command(dockerBinary, "run", "-d", "busybox", "true")
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)
	thirdID := stripTrailingCharacters(out)

	runCmd = exec.Command(dockerBinary, "run", "-d", "busybox", "top")
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)
	fourthID := stripTrailingCharacters(out)

	// make sure third one is not running
	runCmd = exec.Command(dockerBinary, "wait", thirdID)
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)

	// all
	runCmd = exec.Command(dockerBinary, "ps", "-a")
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)

	if !assertContainerList(out, []string{fourthID, thirdID, secondID, firstID}) {
		t.Error("Container list is not in the correct order")
	}

	// running
	runCmd = exec.Command(dockerBinary, "ps")
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)

	if !assertContainerList(out, []string{fourthID, secondID, firstID}) {
		t.Error("Container list is not in the correct order")
	}

	// from here all flag '-a' is ignored

	// limit
	runCmd = exec.Command(dockerBinary, "ps", "-n=2", "-a")
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)
	expected := []string{fourthID, thirdID}

	if !assertContainerList(out, expected) {
		t.Error("Container list is not in the correct order")
	}

	runCmd = exec.Command(dockerBinary, "ps", "-n=2")
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)

	if !assertContainerList(out, expected) {
		t.Error("Container list is not in the correct order")
	}

	// since
	runCmd = exec.Command(dockerBinary, "ps", "--since", firstID, "-a")
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)
	expected = []string{fourthID, thirdID, secondID}

	if !assertContainerList(out, expected) {
		t.Error("Container list is not in the correct order")
	}

	runCmd = exec.Command(dockerBinary, "ps", "--since", firstID)
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)

	if !assertContainerList(out, expected) {
		t.Error("Container list is not in the correct order")
	}

	// before
	runCmd = exec.Command(dockerBinary, "ps", "--before", thirdID, "-a")
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)
	expected = []string{secondID, firstID}

	if !assertContainerList(out, expected) {
		t.Error("Container list is not in the correct order")
	}

	runCmd = exec.Command(dockerBinary, "ps", "--before", thirdID)
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)

	if !assertContainerList(out, expected) {
		t.Error("Container list is not in the correct order")
	}

	// since & before
	runCmd = exec.Command(dockerBinary, "ps", "--since", firstID, "--before", fourthID, "-a")
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)
	expected = []string{thirdID, secondID}

	if !assertContainerList(out, expected) {
		t.Error("Container list is not in the correct order")
	}

	runCmd = exec.Command(dockerBinary, "ps", "--since", firstID, "--before", fourthID)
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)
	if !assertContainerList(out, expected) {
		t.Error("Container list is not in the correct order")
	}

	// since & limit
	runCmd = exec.Command(dockerBinary, "ps", "--since", firstID, "-n=2", "-a")
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)
	expected = []string{fourthID, thirdID}

	if !assertContainerList(out, expected) {
		t.Error("Container list is not in the correct order")
	}

	runCmd = exec.Command(dockerBinary, "ps", "--since", firstID, "-n=2")
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)

	if !assertContainerList(out, expected) {
		t.Error("Container list is not in the correct order")
	}

	// before & limit
	runCmd = exec.Command(dockerBinary, "ps", "--before", fourthID, "-n=1", "-a")
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)
	expected = []string{thirdID}

	if !assertContainerList(out, expected) {
		t.Error("Container list is not in the correct order")
	}

	runCmd = exec.Command(dockerBinary, "ps", "--before", fourthID, "-n=1")
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)

	if !assertContainerList(out, expected) {
		t.Error("Container list is not in the correct order")
	}

	// since & before & limit
	runCmd = exec.Command(dockerBinary, "ps", "--since", firstID, "--before", fourthID, "-n=1", "-a")
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)
	expected = []string{thirdID}

	if !assertContainerList(out, expected) {
		t.Error("Container list is not in the correct order")
	}

	runCmd = exec.Command(dockerBinary, "ps", "--since", firstID, "--before", fourthID, "-n=1")
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)

	if !assertContainerList(out, expected) {
		t.Error("Container list is not in the correct order")
	}

	deleteAllContainers()

	logDone("ps - test ps options")
}

func assertContainerList(out string, expected []string) bool {
	lines := strings.Split(strings.Trim(out, "\n "), "\n")
	if len(lines)-1 != len(expected) {
		return false
	}

	containerIDIndex := strings.Index(lines[0], "CONTAINER ID")
	for i := 0; i < len(expected); i++ {
		foundID := lines[i+1][containerIDIndex : containerIDIndex+12]
		if foundID != expected[i][:12] {
			return false
		}
	}

	return true
}

func TestPsListContainersSize(t *testing.T) {
	name := "test_size"
	runCmd := exec.Command(dockerBinary, "run", "--name", name, "busybox", "sh", "-c", "echo 1 > test")
	out, _, err := runCommandWithOutput(runCmd)
	errorOut(err, t, out)
	id, err := getIDByName(name)
	if err != nil {
		t.Fatal(err)
	}

	runCmd = exec.Command(dockerBinary, "ps", "-s", "-n=1")
	wait := make(chan struct{})
	go func() {
		out, _, err = runCommandWithOutput(runCmd)
		close(wait)
	}()
	select {
	case <-wait:
	case <-time.After(3 * time.Second):
		t.Fatalf("Calling \"docker ps -s\" timed out!")
	}
	errorOut(err, t, out)
	lines := strings.Split(strings.Trim(out, "\n "), "\n")
	sizeIndex := strings.Index(lines[0], "SIZE")
	idIndex := strings.Index(lines[0], "CONTAINER ID")
	foundID := lines[1][idIndex : idIndex+12]
	if foundID != id[:12] {
		t.Fatalf("Expected id %s, got %s", id[:12], foundID)
	}
	expectedSize := "2 B"
	foundSize := lines[1][sizeIndex:]
	if foundSize != expectedSize {
		t.Fatalf("Expected size %q, got %q", expectedSize, foundSize)
	}

	deleteAllContainers()
	logDone("ps - test ps size")
}

func TestPsListContainersFilterStatus(t *testing.T) {
	// FIXME: this should test paused, but it makes things hang and its wonky
	// this is because paused containers can't be controlled by signals

	// start exited container
	runCmd := exec.Command(dockerBinary, "run", "-d", "busybox")
	out, _, err := runCommandWithOutput(runCmd)
	errorOut(err, t, out)
	firstID := stripTrailingCharacters(out)

	// make sure the exited cintainer is not running
	runCmd = exec.Command(dockerBinary, "wait", firstID)
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)

	// start running container
	runCmd = exec.Command(dockerBinary, "run", "-d", "busybox", "sh", "-c", "sleep 360")
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)
	secondID := stripTrailingCharacters(out)

	// filter containers by exited
	runCmd = exec.Command(dockerBinary, "ps", "-a", "-q", "--filter=status=exited")
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)
	containerOut := strings.TrimSpace(out)
	if containerOut != firstID[:12] {
		t.Fatalf("Expected id %s, got %s for exited filter, output: %q", firstID[:12], containerOut, out)
	}

	runCmd = exec.Command(dockerBinary, "ps", "-a", "-q", "--filter=status=running")
	out, _, err = runCommandWithOutput(runCmd)
	errorOut(err, t, out)
	containerOut = strings.TrimSpace(out)
	if containerOut != secondID[:12] {
		t.Fatalf("Expected id %s, got %s for running filter, output: %q", secondID[:12], containerOut, out)
	}

	deleteAllContainers()

	logDone("ps - test ps filter status")
}
