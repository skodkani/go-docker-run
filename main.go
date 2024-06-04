package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

// go run main.go run <cmd> <args>
func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("help")
	}
}

func run() {
	fmt.Printf("Running %v \n", os.Args[2:])

	// Re-run the code so that host name can be set in the namespace.
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	// Wiring the stdin, stdout, stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Running the process in a separate hostname, process id and mount namespace. Also unsharing the mount
	// to not clutter the mount command on the host.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	check_error(cmd.Run())
}

func child() {
	fmt.Printf("Running %v \n", os.Args[2:])

	// Create a separate control group
	cg()

	// Run the Command
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// setting the hostname
	check_error(syscall.Sethostname([]byte("container")))
	// Mounting the rootfs
	check_error(syscall.Chroot("./rootfs"))
	check_error(os.Chdir("/"))
	// Mounting the proc
	check_error(syscall.Mount("proc", "proc", "proc", 0, ""))

	check_error(cmd.Run())

	check_error(syscall.Unmount("proc", 0))
}

func cg() {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	os.Mkdir(filepath.Join(pids, "sharad"), 0755)
	// Setting the max process which can be run in the container
	check_error(ioutil.WriteFile(filepath.Join(pids, "sharad/pids.max"), []byte("20"), 0700))
	// Removes the new cgroup in place after the container exits
	check_error(ioutil.WriteFile(filepath.Join(pids, "sharad/notify_on_release"), []byte("1"), 0700))
	// Moving the process to the cgroup
	check_error(ioutil.WriteFile(filepath.Join(pids, "sharad/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func check_error(err error) {
	if err != nil {
		panic(err)
        }
}
