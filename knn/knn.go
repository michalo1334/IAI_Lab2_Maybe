package knn

import (
	"math/rand"
)

type Knn struct {
	beers *Beers
}

func Initialize() *Knn {
	k := Knn{}
	k.beers = LoadBeers("files/beers.csv")
	return &k
}

func (k *Knn) GetStyles() *Styles {
	return &k.beers.styles
}

func (k *Knn) GetRandomBeer() *Beer {
	index := rand.Intn(len(k.beers.beers))
	b := k.beers.beers[index]
	b.StyleName = k.GetStyleName(k.beers.beers[index].Style)
	return &b
}

func (k *Knn) Get10RandomBeers() []*Beer {
	beers := []*Beer{}
	for i := 0; i < 10; i++ {
		index := rand.Intn(len(k.beers.beers))
		b := k.beers.beers[index]
		b.StyleName = k.GetStyleName(k.beers.beers[index].Style)
		beers = append(beers, &b)
	}
	return beers
}

func (k *Knn) GetBeerByID(id int) *Beer {
	for b := range k.beers.beers {
		if k.beers.beers[b].Id == id {
			k.beers.beers[b].StyleName = k.GetStyleName(k.beers.beers[b].Style)
			return &k.beers.beers[b]
		}
	}
	return nil
}

func (k *Knn) GetRecommendation() []*Beer {
	return k.beers.Recomendation()
}

func (k *Knn) GetStyleName(id int) string {
	return k.beers.styles.GetStyleName(id)
}
