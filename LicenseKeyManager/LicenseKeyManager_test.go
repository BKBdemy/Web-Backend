package LicenseKeyManager_test

import (
	"EntitlementServer/DatabaseAbstraction"
	"EntitlementServer/LicenseKeyManager"
	"EntitlementServer/mocks"
	"testing"
)

func TestLicenseService_AddLicenseKey(t *testing.T) {
	tests := []struct {
		name          string
		keyid         string
		encryptionkey string
		productid     int
		wantErr       bool
	}{
		{
			name:          "Pass",
			keyid:         "test",
			encryptionkey: "enc1",
			productid:     1,
		},
		{
			name:          "Empty KeyID",
			keyid:         "",
			encryptionkey: "tswedtgstg",
			productid:     1,
			wantErr:       true,
		},
		{
			name:          "Empty EncryptionKey",
			keyid:         "test",
			encryptionkey: "",
			productid:     1,
			wantErr:       true,
		},
		{
			name:          "Invalid ProductID",
			keyid:         "test",
			encryptionkey: "tswedtgstg",
			productid:     -1,
			wantErr:       true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockDB := &mocks.DBOrm{}
			mockDB.On("AddLicenseKey", DatabaseAbstraction.LicenseKey{
				KeyID:         test.keyid,
				EncryptionKey: test.encryptionkey,
				ProductID:     test.productid,
			}).Return(1, nil)
			licenseSvc := LicenseKeyManager.LicenseService{DB: mockDB}
			_, err := licenseSvc.AddLicenseKey(test.keyid, test.encryptionkey, test.productid)
			if err != nil {
				if !test.wantErr {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestLicenseService_GetProductLicenseKeys(t *testing.T) {
	mockDB := &mocks.DBOrm{}
	mockDB.On("GetProductLicenseKeys", 1).Return(
		[]DatabaseAbstraction.LicenseKey{
			{
				KeyID:         "test",
				EncryptionKey: "enc1",
				ProductID:     1,
			},
			{
				KeyID:         "test2",
				EncryptionKey: "enc2",
				ProductID:     1,
			},
		}, nil)

	licenseSvc := LicenseKeyManager.LicenseService{DB: mockDB}
	keys, err := licenseSvc.GetProductLicenseKeys(1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %v", len(keys))
	}

	if keys[0].KeyID != "test" {
		t.Errorf("Expected keyID test, got %v", keys[0].KeyID)
	}

	if keys[0].EncryptionKey.GetHex() != "enc1" {
		t.Errorf("Expected encryptionKey enc1, got %v", keys[0].EncryptionKey)
	}

	if keys[0].ProductID != 1 {
		t.Errorf("Expected productID 1, got %v", keys[0].ProductID)
	}

	if keys[1].KeyID != "test2" {
		t.Errorf("Expected keyID test2, got %v", keys[1].KeyID)
	}

	if keys[1].EncryptionKey.GetHex() != "enc2" {
		t.Errorf("Expected encryptionKey enc2, got %v", keys[1].EncryptionKey)
	}
}

func TestLicenseService_GetLicenseKey(t *testing.T) {
	mockDB := &mocks.DBOrm{}
	mockDB.On("GetLicenseKeyByKeyID", "test").Return(
		DatabaseAbstraction.LicenseKey{
			KeyID:         "test",
			EncryptionKey: "enc1",
			ProductID:     1,
		}, nil)

	licenseSvc := LicenseKeyManager.LicenseService{DB: mockDB}
	key, err := licenseSvc.GetLicenseKey("test")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if key.KeyID != "test" {
		t.Errorf("Expected keyID test, got %v", key.KeyID)
	}

	if key.EncryptionKey.GetHex() != "enc1" {
		t.Errorf("Expected encryptionKey enc1, got %v", key.EncryptionKey)
	}

	if key.ProductID != 1 {
		t.Errorf("Expected productID 1, got %v", key.ProductID)
	}
}
