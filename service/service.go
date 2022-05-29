package service

import (
	"main/models"
	"main/repository"
)

// *Adapter* //
type customerService struct{
	custRepo repository.CustomerRepository
}

func NewCustomerService(custRepo repository.CustomerRepository)customerService{
	return customerService{custRepo: custRepo}
}

func (r customerService)GetProducts()([]*models.Product,error){
	return r.custRepo.FetchAll()	
}

func (r customerService)GetProduct(id int)(*models.Product, error){
	return r.custRepo.FetchById(id)
}

func (r customerService)Create(product *models.Product)error{
	return r.custRepo.Create(product)
}

func(r customerService)Update(product *models.Product)error{
	return r.custRepo.Update(product)
}

func(r customerService)Delete(id int)error{
	return r.custRepo.Delete(id)
}