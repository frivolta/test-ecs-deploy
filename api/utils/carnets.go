package api

import (
	db "birdie/db/sqlc"
)

type KidCarnet = map[int64][]db.Carnet
type KidCarnetInfo = map[int64]CarnetInfo
type CarnetInfo struct {
	KidName     string
	KidSurname  string
	TotalBought int32
	TotalUsed   int32
	TotalLeft   int32
	TotalDebit  int32
}

// CUtil util for Carnet data structure
type CUtil struct {
	// carnet list
	Cs []db.Carnet
	// kid carnet structure
	Kc KidCarnet
	// kid carnet info structure
	Kci KidCarnetInfo
}

func NewCUtil(Cs []db.Carnet) *CUtil {
	return &CUtil{
		Cs: Cs,
	}
}

// ToKidCarnet I want all the carnets to be grouped by kidId in a map
func (cu *CUtil) ToKidCarnet() {
	/* {kidId.(Int64): []Carnet} */
	kCsMap := make(map[int64][]db.Carnet)
	// c carnet
	for _, c := range cu.Cs {
		if _, ok := kCsMap[c.KidID]; ok {
			kCsMap[c.KidID] = append(kCsMap[c.KidID], c)
		} else {
			kCsMap[c.KidID] = []db.Carnet{{ID: c.ID, Date: c.Date, Quantity: c.Quantity, KidID: c.KidID}}
		}
	}
	cu.Kc = kCsMap
}

func (cu *CUtil) ToKidCarnetInfo(kidNotes []db.KidNote, kids []db.Kid) {
	Kci := make(map[int64]CarnetInfo)

	//-->1 for every kid ad a record on the Kci map
	// since kids are in range 10/20 than it will be faster make calculations accessing the map
	for _, k := range kids {
		Kci[k.ID] = CarnetInfo{
			KidName:     k.Name,
			KidSurname:  k.Surname,
			TotalBought: 0,
			TotalUsed:   0,
			TotalLeft:   0,
			TotalDebit:  0,
		}
	}

	//-->2 for every kid in the kidCarnet structure
	for id, _ := range cu.Kc {
		// initially we have 0 carnet bought
		var tb int32
		tb = 0
		// loop over the carnets array and get the total
		for _, Cs := range cu.Kc[id] {
			tb = tb + Cs.Quantity
		}
		// append the total to the info
		Kci[id] = CarnetInfo{
			KidName:     Kci[id].KidName,
			KidSurname:  Kci[id].KidSurname,
			TotalBought: tb,
			TotalUsed:   0,
			TotalLeft:   0,
			TotalDebit:  0,
		}
	}

	//-->3 For every kid note get the total meals and update the structure
	for _, kn := range kidNotes {
		// if kid has meal subtract from original info structure
		if kn.HasMeal {
			Kci[kn.KidID] = CarnetInfo{
				KidName:     Kci[kn.KidID].KidName,
				KidSurname:  Kci[kn.KidID].KidSurname,
				TotalBought: Kci[kn.KidID].TotalBought,
				TotalUsed:   Kci[kn.KidID].TotalUsed + 1,
				TotalLeft:   0,
				TotalDebit:  0,
			}
		}
	}

	//-->4 For every info check the total bought / used and fill the info for left, debit
	for i, Kc := range Kci {
		// if you are in credit
		if Kc.TotalBought-Kc.TotalUsed > 0 {
			Kc.TotalLeft = Kc.TotalBought - Kc.TotalUsed
			Kc.TotalDebit = 0
			// else you are in debit
		} else {
			Kc.TotalDebit = Kc.TotalUsed - Kc.TotalBought
			Kc.TotalLeft = 0
		}
		Kci[i] = Kc
	}
	cu.Kci = Kci
}
