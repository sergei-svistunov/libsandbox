package sandbox_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	sandbox "github.com/sergei-svistunov/libsandbox"
)

func TestSandbox(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	sbox := sandbox.New(tmpDir)

	sbox.AddEnv("TESTENV=1234")
	sbox.AddEnv("TESTENV2=1234")

	sbox.AddFile("/usr/bin/echo", "/testbin/echo", true)
	sbox.AddFile("/usr/bin/env", "/testbin/env", true)
	sbox.SetCGroup("testCg").SetNoNewNet(true).SetCpuSet("1,2").SetMemLimit(100 * 1024 * 1024)

	//sbox.SaveUsageStat(usageFile.Name())

	sbox.ExecDir("/testbin")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := sbox.CommandContext(ctx, "./echo", `"TEST"`)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(out))
		t.Fatal(err)
	}

	if string(out) != `"TEST"`+"\n" {
		t.Fatalf("command output: %s", string(out))
	}

	cmd = sbox.CommandContext(ctx, "./env")
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Println(string(out))
		t.Fatal(err)
	}

	if string(out) != "TESTENV=1234\nTESTENV2=1234\n" {
		t.Fatalf("command output: %s", string(out))
	}
}
