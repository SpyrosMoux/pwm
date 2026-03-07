package secrets

func (svc *Service) RemoveSecret(secret string) error {
	exists, err := svc.Store.CheckFileExists(secret)
	if err != nil {
		return err
	}

	if !exists {
		return err
	}

	err = svc.Store.RemoveFile(secret)
	if err != nil {
		return err
	}

	return nil
}
