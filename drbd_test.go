package main

import "testing"

var (
	testCase1 = []byte(` 1: cs:Established ro:Secondary/Primary ds:UpToDate/UpToDate C r-----
	ns:0 nr:119598720 dw:161539164 dr:0 al:2125 bm:0 lo:0 pe:[0;0] ua:0 ap:[0;0] ep:1 wo:2 oos:0
	resync: used:0/61 hits:25476 misses:922 starving:0 locked:0 changed:592
	act_log: used:0/1237 hits:13283342 misses:15109 starving:0 locked:0 changed:12057
	blocked on activity log: 0/0/0`)
	testCase2 = []byte(` 0: cs:Established ro:Secondary/Secondary ds:UpToDate/UpToDate C r-----
	ns:0 nr:0 dw:158780236 dr:0 al:1950 bm:0 lo:0 pe:[0;0] ua:0 ap:[0;0] ep:1 wo:2 oos:0
	resync: used:0/61 hits:2 misses:285 starving:0 locked:0 changed:285
	act_log: used:0/1237 hits:13371654 misses:15913 starving:0 locked:0 changed:11885
	blocked on activity log: 0/0/0`)
	testCase3 = []byte(` 1: cs:Established ro:Primary/Secondary ds:UpToDate/UpToDate C r-----
	ns:119876764 nr:0 dw:154278296 dr:14701504 al:2070 bm:0 lo:0 pe:[0;0] ua:0 ap:[0;0] ep:1 wo:2 oos:0
	resync: used:0/61 hits:25566 misses:719 starving:0 locked:0 changed:485
	act_log: used:0/1237 hits:13318630 misses:41690 starving:0 locked:0 changed:12063
	blocked on activity log: 0/0/0`)
	testCase4 = []byte(` 1: cs:SyncSource ro:Primary/Secondary ds:UpToDate/Inconsistent C s-----
	ns:119876764 nr:0 dw:154278296 dr:14701504 al:2070 bm:0 lo:0 pe:[0;0] ua:0 ap:[0;0] ep:1 wo:2 oos:0
	resync: used:0/61 hits:25566 misses:719 starving:0 locked:0 changed:485
	act_log: used:0/1237 hits:13318630 misses:41690 starving:0 locked:0 changed:12063
	blocked on activity log: 0/0/0`)
)

func TestDrbdParse1(t *testing.T) {
	output := drbdConnection{}
	err := parseProcDRBD(testCase1, &output)
	if err != nil {
		t.Fatalf("failed to parse legit output %s", err.Error())
	}

	if output.ConnectionStatus != "Established" {
		t.Fatalf("mis-parsed testCase1 Expected 'Established' got %s", output.ConnectionStatus)
	}
	if output.MyDiskStatus != "UpToDate" {
		t.Fatalf("mis-parsed testCase1 Expected 'UpToDate' got %s", output.MyDiskStatus)
	}
	if output.RemoteRole != "Primary" {
		t.Fatalf("mis-parsed testCase1 Expected 'Primary' got %s", output.RemoteRole)
	}
}

func TestDrbdParse2(t *testing.T) {
	output := drbdConnection{}
	err := parseProcDRBD(testCase2, &output)
	if err != nil {
		t.Fatalf("failed to parse legit output %s", err.Error())
	}

	if output.ConnectionStatus != "Established" {
		t.Fatalf("mis-parsed testCase2 Expected 'Established' got %s", output.ConnectionStatus)
	}
	if output.MyDiskStatus != "UpToDate" {
		t.Fatalf("mis-parsed testCase2 Expected 'UpToDate' got %s", output.MyDiskStatus)
	}
	if output.RemoteRole != "Secondary" {
		t.Fatalf("mis-parsed testCase2 Expected 'Secondary' got %s", output.RemoteRole)
	}
}

func TestDrbdParse3(t *testing.T) {
	output := drbdConnection{}
	err := parseProcDRBD(testCase3, &output)
	if err != nil {
		t.Fatalf("failed to parse legit output %s", err.Error())
	}

	if output.ConnectionStatus != "Established" {
		t.Fatalf("mis-parsed testCase3 Expected 'Established' got %s", output.ConnectionStatus)
	}
	if output.MyDiskStatus != "UpToDate" {
		t.Fatalf("mis-parsed testCase3 Expected 'UpToDate' got %s", output.MyDiskStatus)
	}
	if output.RemoteRole != "Secondary" {
		t.Fatalf("mis-parsed testCase3 Expected 'Secondary' got %s", output.RemoteRole)
	}
	if output.MyRole != "Primary" {
		t.Fatalf("mis-parsed testCase3 Expected 'Primary' got %s", output.RemoteRole)
	}
}

func TestDrbdParse4(t *testing.T) {
	output := drbdConnection{}
	err := parseProcDRBD(testCase4, &output)
	if err != nil {
		t.Fatalf("failed to parse legit output %s", err.Error())
	}

	if output.ConnectionStatus != "SyncSource" {
		t.Fatalf("mis-parsed testCase4 Expected 'SyncSource' got %s", output.ConnectionStatus)
	}
	if output.MyDiskStatus != "UpToDate" {
		t.Fatalf("mis-parsed testCase4 Expected 'UpToDate' got %s", output.MyDiskStatus)
	}
	if output.RemoteRole != "Secondary" {
		t.Fatalf("mis-parsed testCase4 Expected 'Secondary' got %s", output.RemoteRole)
	}
	if output.MyRole != "Primary" {
		t.Fatalf("mis-parsed testCase4 Expected 'Primary' got %s", output.RemoteRole)
	}

	if !output.Suspended {
		t.Fatal("mis-parsed testCase4 Expected to be suspended, was not")
	}
}
