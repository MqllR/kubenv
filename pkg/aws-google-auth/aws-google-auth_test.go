package awsgoogleauth

import (
	"testing"

	utilexec "k8s.io/utils/exec"
	fakeexec "k8s.io/utils/exec/testing"
)

var testCase = []struct {
	idp      string
	sp       string
	duration int
}{
	{"A708090", "12345", 0},
	{"BH124TA23", "7890", 30000},
}

func TestNewAWSGoogleAuth(t *testing.T) {

	for i := range testCase {
		ga := NewAWSGoogleAuth(testCase[i].idp, testCase[i].sp)

		if ga.IDP != testCase[i].idp {
			t.Errorf("Expected IDP %s but found %s", testCase[i].idp, ga.IDP)
		}

		if ga.SP != testCase[i].sp {
			t.Errorf("Expected SP %s but found %s", testCase[i].sp, ga.SP)
		}
	}
}

func TestSetDefaults(t *testing.T) {
	for i := range testCase {
		ga := NewAWSGoogleAuth(testCase[i].idp, testCase[i].sp)

		if testCase[i].duration != 0 {
			ga.Duration = testCase[i].duration
		}

		ga.SetDefaults()

		if testCase[i].duration == 0 && ga.Duration != DefaultDuration {
			t.Errorf("Expected default Duration %d but found %d", DefaultDuration, ga.Duration)
		}

		if testCase[i].duration != 0 && testCase[i].duration != ga.Duration {
			t.Errorf("Expected explicit Duration %d but found %d", testCase[i].duration, ga.Duration)
		}
	}
}

func TestGetVersion(t *testing.T) {
	testOutput := []struct {
		vstring string
		Expect  string
		Err     bool
	}{
		{"aws-google-auth 0.0.34", "0.0.34", false},
		{"aws-google-auth 1.2.3", "1.2.3", false},
		{"random string", "", true},
	}

	for i := range testOutput {
		fcmd := fakeexec.FakeCmd{
			CombinedOutputScript: []fakeexec.FakeAction{
				func() ([]byte, []byte, error) { return []byte(testOutput[i].vstring), nil, nil },
			},
		}

		fexec := fakeexec.FakeExec{
			CommandScript: []fakeexec.FakeCommandAction{
				func(cmd string, args ...string) utilexec.Cmd { return fakeexec.InitFakeCmd(&fcmd, cmd, args...) },
			},
		}

		runner := New(&fexec)
		version, err := runner.GetVersion()

		if (err != nil) != testOutput[i].Err {
			t.Errorf("Expected error: %v, Got error: %v", testOutput[i].Err, err)
		}

		if testOutput[i].Expect != version {
			t.Errorf("Expected version %s but got %s", testOutput[i].Expect, version)
		}
	}
}
