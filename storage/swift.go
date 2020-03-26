package storage

import (
	"github.com/ncw/swift"
	"io"
)

type Swift struct {
	container string
	conn      *swift.Connection
}

func NewSwiftStorage(container string, userName string, apiKey string, tenant string, tenantID string, authURL string) (Swift, error) {
	conn := &swift.Connection{
		AuthUrl:  authURL,
		UserName: userName,
		ApiKey:   apiKey,
		Tenant:   tenant,
		TenantId: tenantID,
	}

	if err := conn.Authenticate(); err != nil {
		return Swift{}, err
	}

	if _, _, err := conn.Container(container); err != nil {
		return Swift{}, err
	}

	s := Swift{
		container: container,
		conn:      conn,
	}
	return s, nil
}

func (s *Swift) Save(filename string, src io.Reader) error {
	_, err := s.conn.ObjectPut(s.container, filename, src, true, "", "", swift.Headers{})
	return err
}

func (s *Swift) Open(filename string) (io.ReadCloser, error) {
	file, _, err := s.conn.ObjectOpen(s.container, filename, true, nil)
	if err != nil {
		if err == swift.ObjectNotFound {
			return nil, ErrFileNotFound
		}
		return nil, err
	}
	return file, nil
}

func (s *Swift) Delete(filename string) error {
	err := s.conn.ObjectDelete(s.container, filename)
	if err != nil {
		if err == swift.ObjectNotFound {
			return ErrFileNotFound
		}
		return err
	}
	return nil
}
