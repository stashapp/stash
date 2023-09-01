package models

import "testing"

func TestFingerprints_Equals(t *testing.T) {
	var (
		value1 = 1
		value2 = "2"
		value3 = 1.23

		fingerprint1 = Fingerprint{
			Type:        FingerprintTypeMD5,
			Fingerprint: value1,
		}
		fingerprint2 = Fingerprint{
			Type:        FingerprintTypeOshash,
			Fingerprint: value2,
		}
		fingerprint3 = Fingerprint{
			Type:        FingerprintTypePhash,
			Fingerprint: value3,
		}
	)

	tests := []struct {
		name  string
		f     Fingerprints
		other Fingerprints
		want  bool
	}{
		{
			"identical",
			Fingerprints{
				fingerprint1,
				fingerprint2,
			},
			Fingerprints{
				fingerprint1,
				fingerprint2,
			},
			true,
		},
		{
			"different order",
			Fingerprints{
				fingerprint1,
				fingerprint2,
			},
			Fingerprints{
				fingerprint2,
				fingerprint1,
			},
			true,
		},
		{
			"different length",
			Fingerprints{
				fingerprint1,
				fingerprint2,
			},
			Fingerprints{
				fingerprint1,
			},
			false,
		},
		{
			"different",
			Fingerprints{
				fingerprint1,
				fingerprint2,
			},
			Fingerprints{
				fingerprint1,
				fingerprint3,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Equals(tt.other); got != tt.want {
				t.Errorf("Fingerprints.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFingerprints_ContentsChanged(t *testing.T) {
	var (
		value1 = 1
		value2 = "2"
		value3 = 1.23

		fingerprint1 = Fingerprint{
			Type:        FingerprintTypeMD5,
			Fingerprint: value1,
		}
		fingerprint2 = Fingerprint{
			Type:        FingerprintTypeOshash,
			Fingerprint: value2,
		}
		fingerprint3 = Fingerprint{
			Type:        FingerprintTypeMD5,
			Fingerprint: value3,
		}
	)

	tests := []struct {
		name  string
		f     Fingerprints
		other Fingerprints
		want  bool
	}{
		{
			"identical",
			Fingerprints{
				fingerprint1,
				fingerprint2,
			},
			Fingerprints{
				fingerprint1,
				fingerprint2,
			},
			false,
		},
		{
			"has new",
			Fingerprints{
				fingerprint1,
				fingerprint2,
			},
			Fingerprints{
				fingerprint1,
			},
			true,
		},
		{
			"has different value",
			Fingerprints{
				fingerprint3,
				fingerprint2,
			},
			Fingerprints{
				fingerprint1,
				fingerprint2,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.ContentsChanged(tt.other); got != tt.want {
				t.Errorf("Fingerprints.ContentsChanged() = %v, want %v", got, tt.want)
			}
		})
	}
}
