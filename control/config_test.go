package control

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestServiceBlock(t *testing.T) {
	RelayState.RedisClient.FlushAll().Result()

	app := configCmdInit()

	app.SetArgs([]string{"enable", "service-block"})
	app.Execute()
	RelayState.Load()
	if !RelayState.RelayConfig.BlockService {
		t.Fatalf("Not Enabled Blocking feature for service-type actor")
	}

	app.SetArgs([]string{"enable", "-d", "service-block"})
	app.Execute()
	RelayState.Load()
	if RelayState.RelayConfig.BlockService {
		t.Fatalf("Not Disabled Blocking feature for service-type actor")
	}
}

func TestManuallyAccept(t *testing.T) {
	RelayState.RedisClient.FlushAll().Result()

	app := configCmdInit()

	app.SetArgs([]string{"enable", "manually-accept"})
	app.Execute()
	RelayState.Load()
	if !RelayState.RelayConfig.ManuallyAccept {
		t.Fatalf("Not Enabled Manually accept follow-request feature")
	}

	app.SetArgs([]string{"enable", "-d", "manually-accept"})
	app.Execute()
	RelayState.Load()
	if RelayState.RelayConfig.ManuallyAccept {
		t.Fatalf("Not Disabled Manually accept follow-request feature")
	}
}

func TestCreateAsAnnounce(t *testing.T) {
	RelayState.RedisClient.FlushAll().Result()

	app := configCmdInit()

	app.SetArgs([]string{"enable", "create-as-announce"})
	app.Execute()
	RelayState.Load()
	if !RelayState.RelayConfig.CreateAsAnnounce {
		t.Fatalf("Enable announce activity instead of relay create activity")
	}

	app.SetArgs([]string{"enable", "-d", "create-as-announce"})
	app.Execute()
	RelayState.Load()
	if RelayState.RelayConfig.CreateAsAnnounce {
		t.Fatalf("Enable announce activity instead of relay create activity")
	}
}

func TestInvalidConfig(t *testing.T) {
	RelayState.RedisClient.FlushAll().Result()

	app := configCmdInit()
	buffer := new(bytes.Buffer)
	app.SetOut(buffer)

	app.SetArgs([]string{"enable", "hoge"})
	app.Execute()

	output := buffer.String()
	if strings.Split(output, "\n")[0] != "Invalid config given" {
		t.Fatalf("Invalid Response.")
	}
}

func TestListConfig(t *testing.T) {
	RelayState.RedisClient.FlushAll().Result()

	app := configCmdInit()
	buffer := new(bytes.Buffer)
	app.SetOut(buffer)

	app.SetArgs([]string{"list"})
	app.Execute()

	output := buffer.String()
	for _, row := range strings.Split(output, "\n") {
		switch strings.Split(row, ":")[0] {
		case "Blocking for service-type actor ":
			if strings.Split(row, ":")[1] == "  true" {
				t.Fatalf("Invalid Response.")
			}
		case "Manually accept follow-request ":
			if strings.Split(row, ":")[1] == "  true" {
				t.Fatalf("Invalid Response.")
			}
		case "Announce activity instead of relay create activity ":
			if strings.Split(row, ":")[1] == "  true" {
				t.Fatalf("Invalid Response.")
			}
		}
	}
}

func TestExportConfig(t *testing.T) {
	RelayState.RedisClient.FlushAll().Result()

	app := configCmdInit()
	buffer := new(bytes.Buffer)
	app.SetOut(buffer)

	app.SetArgs([]string{"export"})
	app.Execute()

	file, err := os.Open("../misc/test/blankConfig.json")
	if err != nil {
		t.Fatalf("Test resource fetch error.")
	}
	jsonData, _ := io.ReadAll(file)
	output := buffer.String()
	if strings.Split(output, "\n")[0] != string(jsonData) {
		t.Fatalf("Invalid Response.")
	}
}

func TestImportConfig(t *testing.T) {
	RelayState.RedisClient.FlushAll().Result()

	app := configCmdInit()

	app.SetArgs([]string{"import", "--json", "../misc/test/exampleConfig.json"})
	app.Execute()
	RelayState.Load()

	buffer := new(bytes.Buffer)
	app.SetOut(buffer)

	app.SetArgs([]string{"export"})
	app.Execute()

	file, err := os.Open("../misc/test/exampleConfig.json")
	if err != nil {
		t.Fatalf("Test resource fetch error.")
	}
	jsonData, _ := io.ReadAll(file)
	output := buffer.String()
	if strings.Split(output, "\n")[0] != string(jsonData) {
		t.Fatalf("Invalid Response.")
	}
}
