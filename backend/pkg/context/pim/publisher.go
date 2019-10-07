package pim

import (
	public "github.com/cod1ng-earth/event-web-store/backend/pkg/context/pim/public"
)

func (p publisher) publishProduct(offset int64, fact *Product) error {
	i := &public.Product{
		Id:            fact.Id,
		Price:         fact.Price,
		Name:          fact.Name,
		Description:   fact.Description,
		Longtext:      fact.Longtext,
		Category:      fact.Category,
		SmallImageURL: fact.SmallImageURL,
		LargeImageURL: fact.LargeImageURL,
		Disabled:      fact.Disabled,
		Tax:           fact.Tax,
	}
	_, _, err := p.logProduct(offset, i)
	return err
}
