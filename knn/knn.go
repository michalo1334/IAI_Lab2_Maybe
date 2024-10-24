package knn

import (
	"math/rand"
	"sort"
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

func (k *Knn) GetThreeMostSimilarBeers(targetBeer *Beer) []*Beer {
	type beerDistance struct {
		beer     *Beer
		distance float64
	}

	distances := []beerDistance{}

	// Compute distances between the target beer and all other beers
	for i := range k.beers.beers {
		beer := &k.beers.beers[i]
		if beer.Id == targetBeer.Id {
			continue // Skip the target beer itself
		}
		// Calculate the distance using the existing Beer.Distance method
		dist := targetBeer.Distance(beer)
		distances = append(distances, beerDistance{beer: beer, distance: dist})
	}

	// Sort the beers by distance (ascending order)
	sort.Slice(distances, func(i, j int) bool {
		return distances[i].distance < distances[j].distance
	})

	// Select the top three most similar beers
	similarBeers := []*Beer{
		distances[0].beer,
		distances[1].beer,
		distances[2].beer,
	}

	// Populate StyleName for each similar beer
	for i, beer := range similarBeers {
		// Create a copy to avoid modifying shared data
		beerCopy := *beer
		beerCopy.StyleName = k.GetStyleName(beerCopy.Style)
		similarBeers[i] = &beerCopy
	}

	return similarBeers
}
