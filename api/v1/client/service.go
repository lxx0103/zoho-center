package client

import (
	"errors"
	"zoho-center/core/database"
)

type clientService struct {
}

func NewClientService() ClientService {
	return &clientService{}
}

// ClientService represents a service for managing clients.
type ClientService interface {
	//Client Management
	GetClientByID(int64, int64) (*Client, error)
	NewClient(ClientNew, int64) (*Client, error)
	GetClientList(ClientFilter, int64) (int, *[]Client, error)
	UpdateClient(int64, ClientNew, int64) (*Client, error)
}

func (s *clientService) GetClientByID(id int64, organizationID int64) (*Client, error) {
	db := database.InitMySQL()
	query := NewClientQuery(db)
	client, err := query.GetClientByID(id, organizationID)
	return client, err
}

func (s *clientService) NewClient(info ClientNew, organizationID int64) (*Client, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewClientRepository(tx)
	exist, err := repo.CheckNameExist(info.Name, organizationID)
	if err != nil {
		return nil, err
	}
	if exist != 0 {
		msg := "客户名称重复"
		return nil, errors.New(msg)
	}
	clientID, err := repo.CreateClient(info, organizationID)
	if err != nil {
		return nil, err
	}
	client, err := repo.GetClientByID(clientID, organizationID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return client, err
}

func (s *clientService) GetClientList(filter ClientFilter, organizationID int64) (int, *[]Client, error) {
	db := database.InitMySQL()
	query := NewClientQuery(db)
	count, err := query.GetClientCount(filter, organizationID)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetClientList(filter, organizationID)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *clientService) UpdateClient(clientID int64, info ClientNew, organizationID int64) (*Client, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewClientRepository(tx)
	oldClient, err := repo.GetClientByID(clientID, organizationID)
	if err != nil {
		return nil, err
	}
	if organizationID != 0 && organizationID != oldClient.OrganizationID {
		msg := "你无权修改此客户"
		return nil, errors.New(msg)
	}
	exist, err := repo.CheckNameExist(info.Name, organizationID)
	if err != nil {
		return nil, err
	}
	if exist != 0 {
		msg := "客户名称重复"
		return nil, errors.New(msg)
	}
	_, err = repo.UpdateClient(clientID, info)
	if err != nil {
		return nil, err
	}
	client, err := repo.GetClientByID(clientID, organizationID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return client, err
}
