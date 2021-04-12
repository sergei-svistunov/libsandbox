package sandbox

import (
	"context"
	"log"
	"os/exec"
	"strconv"
)

var SandboxElf = "/usr/bin/sandbox"

type Sandbox struct {
	path      string
	files     []file
	mountDirs []mountDir
	env       []string
	noNewNet  bool
	cgroup    string
	memLimit  uint64
}

type file struct {
	src      string
	dst      string
	withLibs bool
}

type mountDir struct {
	src string
	dst string
}

func New(path string) *Sandbox {
	return &Sandbox{
		path: path,
	}
}

func (s *Sandbox) AddFile(src, dst string, withLibs bool) *Sandbox {
	s.files = append(s.files, file{
		src:      src,
		dst:      dst,
		withLibs: withLibs,
	})

	return s
}

func (s *Sandbox) MountDir(src, dst string) *Sandbox {
	s.mountDirs = append(s.mountDirs, mountDir{
		src: src,
		dst: dst,
	})

	return s
}

func (s *Sandbox) AddEnv(value string) *Sandbox {
	s.env = append(s.env, value)

	return s
}

func (s *Sandbox) SetNoNewNet(v bool) *Sandbox {
	s.noNewNet = v

	return s
}

func (s *Sandbox) SetCGroup(name string) *Sandbox {
	s.cgroup = name

	return s
}

func (s *Sandbox) SetMemLimit(limit uint64) *Sandbox {
	s.memLimit = limit

	return s
}

func (s *Sandbox) Command(path string, args ...string) *exec.Cmd {
	return s.CommandContext(nil, path, args...)
}

func (s *Sandbox) CommandContext(ctx context.Context, path string, args ...string) *exec.Cmd {
	execArgs := []string{s.path}

	for _, f := range s.files {
		if f.withLibs {
			execArgs = append(execArgs, "--add_elf_file")
		} else {
			execArgs = append(execArgs, "--add_file")
		}

		execArgs = append(execArgs, f.src, f.dst)
	}

	for _, d := range s.mountDirs {
		execArgs = append(execArgs, "--mount_dir", d.src, d.dst)
	}

	for _, e := range s.env {
		execArgs = append(execArgs, "--env", e)
	}

	if s.noNewNet {
		execArgs = append(execArgs, "--no_new_net")
	}

	if s.cgroup != "" {
		execArgs = append(execArgs, "--cgroup", s.cgroup)
	}

	if s.memLimit != 0 {
		execArgs = append(execArgs, "--mem_limit", strconv.FormatUint(s.memLimit, 10))
	}

	execArgs = append(execArgs, "--", path)
	execArgs = append(execArgs, args...)

	log.Print(SandboxElf, execArgs)

	if ctx == nil {
		return exec.Command(SandboxElf, execArgs...)
	}

	return exec.CommandContext(ctx, SandboxElf, execArgs...)
}
