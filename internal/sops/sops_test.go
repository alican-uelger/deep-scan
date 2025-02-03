//go:build unit

package sops

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDecryptFile(t *testing.T) {

	t.Run("Valid file", func(t *testing.T) {
		path := "test.yaml"
		expectedFileType := "yaml"

		mockContent := []byte("encrypted content")
		mockStorage := NewStorageMock(t)
		mockStorage.
			On("ReadFile", path).
			Return(mockContent, nil)

		mockDecryptedContent := []byte("decrypted content")
		expectedContent := "decrypted content"

		mockSopsClient := NewSopsAPIMock(t)
		mockSopsClient.
			On("DecryptFile", mockContent, expectedFileType).
			Return(mockDecryptedContent, nil)

		s := &Sops{
			Storage: mockStorage,
			Client:  mockSopsClient,
		}

		actualContent, err := s.DecryptFile(path)
		mockStorage.AssertNumberOfCalls(t, "ReadFile", 1)
		mockSopsClient.AssertNumberOfCalls(t, "DecryptFile", 1)
		require.NoError(t, err)
		assert.Equal(t, expectedContent, actualContent)
	})

	t.Run("Read file error", func(t *testing.T) {
		mockErr := errors.New("read error")
		expectedErr := errors.New("could not read secret file: read error")

		mockStorage := NewStorageMock(t)
		mockStorage.
			On("ReadFile", mock.Anything).
			Return([]byte{}, mockErr)
		mockSopsClient := NewSopsAPIMock(t)

		s := &Sops{
			Storage: mockStorage,
			Client:  mockSopsClient,
		}

		_, actualErr := s.DecryptFile("test")
		mockStorage.AssertNumberOfCalls(t, "ReadFile", 1)
		mockSopsClient.AssertNumberOfCalls(t, "DecryptFile", 0)
		require.Error(t, actualErr)
		assert.EqualError(t, expectedErr, actualErr.Error())
	})

	t.Run("Decrypt file error", func(t *testing.T) {
		mockErr := errors.New("decrypt error")
		expectedErr := errors.New("could not decrypt secret file: decrypt error")

		mockStorage := NewStorageMock(t)
		mockStorage.
			On("ReadFile", mock.Anything).
			Return([]byte{}, nil)
		mockSopsClient := NewSopsAPIMock(t)
		mockSopsClient.
			On("DecryptFile", mock.Anything, mock.Anything).
			Return([]byte{}, mockErr)

		s := &Sops{
			Storage: mockStorage,
			Client:  mockSopsClient,
		}

		_, actualErr := s.DecryptFile("test")
		mockStorage.AssertNumberOfCalls(t, "ReadFile", 1)
		mockSopsClient.AssertNumberOfCalls(t, "DecryptFile", 1)
		require.Error(t, actualErr)
		assert.EqualError(t, expectedErr, actualErr.Error())
	})
}
