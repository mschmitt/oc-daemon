package trafpol

import (
	"reflect"
	"testing"

	"github.com/T-Systems-MMS/oc-daemon/internal/cpd"
	"github.com/T-Systems-MMS/oc-daemon/internal/devmon"
	"github.com/vishvananda/netlink"
)

// TestTrafPolHandleDeviceUpdate tests handleDeviceUpdate of TrafPol
func TestTrafPolHandleDeviceUpdate(t *testing.T) {
	allowedHosts := []string{"example.com"}
	tp := NewTrafPol(allowedHosts)

	// test adding
	update := &devmon.Update{
		Add: true,
	}
	tp.handleDeviceUpdate(update)

	// test removing
	update.Add = false
	tp.handleDeviceUpdate(update)
}

// TestTrafPolHandleDNSUpdate tests handleDNSUpdate of TrafPol
func TestTrafPolHandleDNSUpdate(t *testing.T) {
	allowedHosts := []string{"example.com"}
	tp := NewTrafPol(allowedHosts)

	tp.allowHosts.Start()
	defer tp.allowHosts.Stop()
	tp.cpd.Start()
	defer tp.cpd.Stop()

	tp.handleDNSUpdate()
}

// TestTrafPolHandleCPDReport tests handleCPDReport of TrafPol
func TestTrafPolHandleCPDReport(t *testing.T) {
	allowedHosts := []string{"example.com"}
	tp := NewTrafPol(allowedHosts)

	tp.allowHosts.Start()
	defer tp.allowHosts.Stop()

	want := []string{}
	got := []string{}
	runNft = func(s string) {
		got = append(got, s)
	}

	// test not detected
	report := &cpd.Report{}
	tp.handleCPDReport(report)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	// test detected
	want = []string{
		"add element inet oc-daemon-filter allowports { 80, 443 }",
	}
	got = []string{}
	report.Detected = true
	tp.handleCPDReport(report)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	// test not detected any more
	want = []string{
		"delete element inet oc-daemon-filter allowports { 80, 443 }",
	}
	got = []string{}
	report.Detected = false
	tp.handleCPDReport(report)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestTrafPolStartStop tests Start and Stop of TrafPol
func TestTrafPolStartStop(t *testing.T) {
	allowedHosts := []string{"example.com"}
	tp := NewTrafPol(allowedHosts)

	// set dummy low level function for devmon
	devmon.RegisterLinkUpdates = func(*devmon.DevMon) chan netlink.LinkUpdate {
		return nil
	}

	tp.Start()
	tp.Stop()
}

// TestNewTrafPol tests NewTrafPol
func TestNewTrafPol(t *testing.T) {
	allowedHosts := []string{"example.com"}
	tp := NewTrafPol(allowedHosts)
	if tp.devmon == nil ||
		tp.dnsmon == nil ||
		tp.cpd == nil ||
		tp.allowDevs == nil ||
		tp.allowHosts == nil ||
		tp.loopDone == nil ||
		tp.done == nil {

		t.Errorf("got nil, want != nil")
	}
}

// TestCleanup tests Cleanup
func TestCleanup(t *testing.T) {
	want := []string{
		"delete table inet oc-daemon-filter",
	}
	got := []string{}
	runCleanupNft = func(s string) {
		got = append(got, s)
	}
	Cleanup()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
