package libcfg_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ibraimgm/libcfg"
)

func withEnv(env map[string]string, handler func()) {
	for k, v := range env {
		if err := os.Setenv(k, v); err != nil {
			panic(err)
		}
	}

	defer func() {
		for k := range env {
			if err := os.Unsetenv(k); err != nil {
				panic(err)
			}
		}
	}()

	handler()
}

func withFileEnv(env map[string]string, handler func(string)) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}

	file, err := ioutil.TempFile(dir, "")
	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(dir)

	for k, v := range env {
		line := fmt.Sprintf("%s=%s\n", k, v)

		if _, err := file.WriteString(line); err != nil {
			panic(err)
		}
	}

	handler(file.Name())
}

func TestParseEnvArgs(t *testing.T) {
	tests := []struct {
		env     map[string]string
		abool   bool
		aint    int
		auint   uint
		astring string
	}{
		{env: map[string]string{}},
		{env: map[string]string{
			"B1": "true",
			"I1": "5",
			"U1": "9",
			"S1": "foo",
		}, abool: true, aint: 5, auint: 9, astring: "foo"},
		{env: map[string]string{
			"B2": "true",
			"I2": "6",
			"U2": "10",
			"S2": "bar",
		}, abool: true, aint: 6, auint: 10, astring: "bar"},
		{env: map[string]string{
			"B3": "true",
			"I3": "7",
			"U3": "11",
			"S3": "baz",
		}, abool: true, aint: 7, auint: 11, astring: "baz"},
		{env: map[string]string{
			"B1": "false",
			"B2": "true",
		}, abool: true},
		{env: map[string]string{
			"B1": "true",
			"B2": "false",
		}},
		{env: map[string]string{
			"B1": "true",
			"B2": "false",
			"B3": "true",
		}, abool: true},
		{env: map[string]string{
			"I1": "1",
			"I2": "2",
		}, aint: 2},
		{env: map[string]string{
			"I1": "1",
			"I2": "0",
		}},
		{env: map[string]string{
			"I1": "1",
			"I2": "0",
			"I3": "3",
		}, aint: 3},
		{env: map[string]string{
			"I1": "1",
			"I2": "2",
			"I3": "",
		}},
		{env: map[string]string{
			"U1": "1",
			"U2": "2",
		}, auint: 2},
		{env: map[string]string{
			"U1": "1",
			"U2": "0",
		}},
		{env: map[string]string{
			"U1": "1",
			"U2": "0",
			"U3": "3",
		}, auint: 3},
		{env: map[string]string{
			"U1": "1",
			"U2": "2",
			"U3": "",
		}},
		{env: map[string]string{
			"S1": "foo",
			"S2": "bar",
		}, astring: "bar"},
		{env: map[string]string{
			"S1": "1",
			"S2": "",
		}},
		{env: map[string]string{
			"S1": "foo",
			"S2": "",
			"S3": "baz",
		}, astring: "baz"},
		{env: map[string]string{
			"S1": "foo",
			"S2": "bar",
			"S3": "",
		}},
	}

	for i, test := range tests {
		cfg := libcfg.NewParser()

		abool := cfg.EnvBool(false, "B1", "B2", "B3")
		aint := cfg.EnvInt(0, "I1", "I2", "I3")
		auint := cfg.EnvUint(0, "U1", "U2", "U3")
		astring := cfg.EnvString("", "S1", "S2", "S3")

		i := i       // pin scope
		test := test // pin scope

		withEnv(test.env, func() {
			if err := cfg.RunEnv(); err != nil {
				t.Errorf("Case %d, error loading from environment: %v", i, err)
				return
			}

			if *abool != test.abool {
				t.Errorf("Case %d, wrong boolean value: expected '%v', received '%v'", i, test.abool, *abool)
			}

			if *aint != test.aint {
				t.Errorf("Case %d, wrong int value: expected '%v', received '%v'", i, test.aint, *aint)
			}

			if *auint != test.auint {
				t.Errorf("Case %d, wrong uint value: expected '%v', received '%v'", i, test.auint, *auint)
			}

			if *astring != test.astring {
				t.Errorf("Case %d, wrong string value: expected '%v', received '%v'", i, test.astring, *astring)
			}
		})
	}
}

func TestParseEnvDefault(t *testing.T) {
	tests := []struct {
		env     map[string]string
		abool   bool
		aint    int
		auint   uint
		astring string
	}{
		{env: map[string]string{}, abool: true, aint: 8, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"B1": "true",
			"I1": "5",
			"U1": "9",
			"S1": "foo",
		}, abool: true, aint: 5, auint: 9, astring: "foo"},
		{env: map[string]string{
			"B2": "true",
			"I2": "6",
			"U2": "10",
			"S2": "bar",
		}, abool: true, aint: 6, auint: 10, astring: "bar"},
		{env: map[string]string{
			"B3": "true",
			"I3": "7",
			"U3": "11",
			"S3": "baz",
		}, abool: true, aint: 7, auint: 11, astring: "baz"},
		{env: map[string]string{
			"B1": "false",
			"B2": "true",
		}, abool: true, aint: 8, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"B1": "true",
			"B2": "false",
		}, abool: false, aint: 8, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"B1": "true",
			"B2": "false",
			"B3": "true",
		}, abool: true, aint: 8, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"I1": "1",
			"I2": "2",
		}, abool: true, aint: 2, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"I1": "1",
			"I2": "0",
		}, abool: true, aint: 0, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"I1": "1",
			"I2": "0",
			"I3": "3",
		}, abool: true, aint: 3, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"I1": "1",
			"I2": "2",
			"I3": "",
		}, abool: true, aint: 8, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"U1": "1",
			"U2": "2",
		}, abool: true, aint: 8, auint: 2, astring: "xyz"},
		{env: map[string]string{
			"U1": "1",
			"U2": "0",
		}, abool: true, aint: 8, auint: 0, astring: "xyz"},
		{env: map[string]string{
			"U1": "1",
			"U2": "0",
			"U3": "3",
		}, abool: true, aint: 8, auint: 3, astring: "xyz"},
		{env: map[string]string{
			"U1": "1",
			"U2": "2",
			"U3": "",
		}, abool: true, aint: 8, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"S1": "foo",
			"S2": "bar",
		}, abool: true, aint: 8, auint: 16, astring: "bar"},
		{env: map[string]string{
			"S1": "1",
			"S2": "",
		}, abool: true, aint: 8, auint: 16, astring: ""},
		{env: map[string]string{
			"S1": "foo",
			"S2": "",
			"S3": "baz",
		}, abool: true, aint: 8, auint: 16, astring: "baz"},
		{env: map[string]string{
			"S1": "foo",
			"S2": "bar",
			"S3": "",
		}, abool: true, aint: 8, auint: 16, astring: ""},
	}

	for i, test := range tests {
		cfg := libcfg.NewParser()

		abool := cfg.EnvBool(true, "B1", "B2", "B3")
		aint := cfg.EnvInt(8, "I1", "I2", "I3")
		auint := cfg.EnvUint(16, "U1", "U2", "U3")
		astring := cfg.EnvString("xyz", "S1", "S2", "S3")

		i := i       // pin scope
		test := test // pin scope

		withEnv(test.env, func() {
			if err := cfg.RunEnv(); err != nil {
				t.Errorf("Case %d, error loading from environment: %v", i, err)
				return
			}

			if *abool != test.abool {
				t.Errorf("Case %d, wrong boolean value: expected '%v', received '%v'", i, test.abool, *abool)
			}

			if *aint != test.aint {
				t.Errorf("Case %d, wrong int value: expected '%v', received '%v'", i, test.aint, *aint)
			}

			if *auint != test.auint {
				t.Errorf("Case %d, wrong uint value: expected '%v', received '%v'", i, test.auint, *auint)
			}

			if *astring != test.astring {
				t.Errorf("Case %d, wrong string value: expected '%v', received '%v'", i, test.astring, *astring)
			}
		})
	}
}

func TestParseEnvFileArgs(t *testing.T) {
	tests := []struct {
		env     map[string]string
		abool   bool
		aint    int
		auint   uint
		astring string
	}{
		{env: map[string]string{}},
		{env: map[string]string{
			"B1": "true",
			"I1": "5",
			"U1": "9",
			"S1": "foo",
		}, abool: true, aint: 5, auint: 9, astring: "foo"},
		{env: map[string]string{
			"B2": "true",
			"I2": "6",
			"U2": "10",
			"S2": "bar",
		}, abool: true, aint: 6, auint: 10, astring: "bar"},
		{env: map[string]string{
			"B3": "true",
			"I3": "7",
			"U3": "11",
			"S3": "baz",
		}, abool: true, aint: 7, auint: 11, astring: "baz"},
		{env: map[string]string{
			"B1": "false",
			"B2": "true",
		}, abool: true},
		{env: map[string]string{
			"B1": "true",
			"B2": "false",
		}},
		{env: map[string]string{
			"B1": "true",
			"B2": "false",
			"B3": "true",
		}, abool: true},
		{env: map[string]string{
			"I1": "1",
			"I2": "2",
		}, aint: 2},
		{env: map[string]string{
			"I1": "1",
			"I2": "0",
		}},
		{env: map[string]string{
			"I1": "1",
			"I2": "0",
			"I3": "3",
		}, aint: 3},
		{env: map[string]string{
			"I1": "1",
			"I2": "2",
			"I3": "",
		}},
		{env: map[string]string{
			"U1": "1",
			"U2": "2",
		}, auint: 2},
		{env: map[string]string{
			"U1": "1",
			"U2": "0",
		}},
		{env: map[string]string{
			"U1": "1",
			"U2": "0",
			"U3": "3",
		}, auint: 3},
		{env: map[string]string{
			"U1": "1",
			"U2": "2",
			"U3": "",
		}},
		{env: map[string]string{
			"S1": "foo",
			"S2": "bar",
		}, astring: "bar"},
		{env: map[string]string{
			"S1": "1",
			"S2": "",
		}},
		{env: map[string]string{
			"S1": "foo",
			"S2": "",
			"S3": "baz",
		}, astring: "baz"},
		{env: map[string]string{
			"S1": "foo",
			"S2": "bar",
			"S3": "",
		}},
	}

	for i, test := range tests {
		cfg := libcfg.NewParser()

		abool := cfg.EnvBool(false, "B1", "B2", "B3")
		aint := cfg.EnvInt(0, "I1", "I2", "I3")
		auint := cfg.EnvUint(0, "U1", "U2", "U3")
		astring := cfg.EnvString("", "S1", "S2", "S3")

		i := i       // pin scope
		test := test // pin scope

		withFileEnv(test.env, func(filename string) {
			if err := cfg.UseFile(filename); err != nil {
				t.Errorf("Case %d, error loading file: %v", i, err)
				return
			}

			if err := cfg.RunEnv(); err != nil {
				t.Errorf("Case %d, error loading from environment: %v", i, err)
				return
			}

			if *abool != test.abool {
				t.Errorf("Case %d, wrong boolean value: expected '%v', received '%v'", i, test.abool, *abool)
			}

			if *aint != test.aint {
				t.Errorf("Case %d, wrong int value: expected '%v', received '%v'", i, test.aint, *aint)
			}

			if *auint != test.auint {
				t.Errorf("Case %d, wrong uint value: expected '%v', received '%v'", i, test.auint, *auint)
			}

			if *astring != test.astring {
				t.Errorf("Case %d, wrong string value: expected '%v', received '%v'", i, test.astring, *astring)
			}
		})
	}
}

func TestParseEnvFileDefault(t *testing.T) {
	tests := []struct {
		env     map[string]string
		abool   bool
		aint    int
		auint   uint
		astring string
	}{
		{env: map[string]string{}, abool: true, aint: 8, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"B1": "true",
			"I1": "5",
			"U1": "9",
			"S1": "foo",
		}, abool: true, aint: 5, auint: 9, astring: "foo"},
		{env: map[string]string{
			"B2": "true",
			"I2": "6",
			"U2": "10",
			"S2": "bar",
		}, abool: true, aint: 6, auint: 10, astring: "bar"},
		{env: map[string]string{
			"B3": "true",
			"I3": "7",
			"U3": "11",
			"S3": "baz",
		}, abool: true, aint: 7, auint: 11, astring: "baz"},
		{env: map[string]string{
			"B1": "false",
			"B2": "true",
		}, abool: true, aint: 8, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"B1": "true",
			"B2": "false",
		}, abool: false, aint: 8, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"B1": "true",
			"B2": "false",
			"B3": "true",
		}, abool: true, aint: 8, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"I1": "1",
			"I2": "2",
		}, abool: true, aint: 2, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"I1": "1",
			"I2": "0",
		}, abool: true, aint: 0, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"I1": "1",
			"I2": "0",
			"I3": "3",
		}, abool: true, aint: 3, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"I1": "1",
			"I2": "2",
			"I3": "",
		}, abool: true, aint: 8, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"U1": "1",
			"U2": "2",
		}, abool: true, aint: 8, auint: 2, astring: "xyz"},
		{env: map[string]string{
			"U1": "1",
			"U2": "0",
		}, abool: true, aint: 8, auint: 0, astring: "xyz"},
		{env: map[string]string{
			"U1": "1",
			"U2": "0",
			"U3": "3",
		}, abool: true, aint: 8, auint: 3, astring: "xyz"},
		{env: map[string]string{
			"U1": "1",
			"U2": "2",
			"U3": "",
		}, abool: true, aint: 8, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"S1": "foo",
			"S2": "bar",
		}, abool: true, aint: 8, auint: 16, astring: "bar"},
		{env: map[string]string{
			"S1": "1",
			"S2": "",
		}, abool: true, aint: 8, auint: 16, astring: ""},
		{env: map[string]string{
			"S1": "foo",
			"S2": "",
			"S3": "baz",
		}, abool: true, aint: 8, auint: 16, astring: "baz"},
		{env: map[string]string{
			"S1": "foo",
			"S2": "bar",
			"S3": "",
		}, abool: true, aint: 8, auint: 16, astring: ""},
	}

	for i, test := range tests {
		cfg := libcfg.NewParser()

		abool := cfg.EnvBool(true, "B1", "B2", "B3")
		aint := cfg.EnvInt(8, "I1", "I2", "I3")
		auint := cfg.EnvUint(16, "U1", "U2", "U3")
		astring := cfg.EnvString("xyz", "S1", "S2", "S3")

		i := i       // pin scope
		test := test // pin scope

		withFileEnv(test.env, func(filename string) {
			if err := cfg.UseFile(filename); err != nil {
				t.Errorf("Case %d, error loading file: %v", i, err)
				return
			}

			if err := cfg.RunEnv(); err != nil {
				t.Errorf("Case %d, error loading from environment: %v", i, err)
				return
			}

			if *abool != test.abool {
				t.Errorf("Case %d, wrong boolean value: expected '%v', received '%v'", i, test.abool, *abool)
			}

			if *aint != test.aint {
				t.Errorf("Case %d, wrong int value: expected '%v', received '%v'", i, test.aint, *aint)
			}

			if *auint != test.auint {
				t.Errorf("Case %d, wrong uint value: expected '%v', received '%v'", i, test.auint, *auint)
			}

			if *astring != test.astring {
				t.Errorf("Case %d, wrong string value: expected '%v', received '%v'", i, test.astring, *astring)
			}
		})
	}
}