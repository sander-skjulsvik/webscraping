package db

// Finn
type Realestate struct {
	Title, Address, URL, DateTime string
	ID, Price                     int
	Info                          map[string]string
	Updates                       map[string]Realestate // datetime string
	//Active                        bool
}

func (left Realestate) RightUpdates(right Realestate) (Realestate, bool) {
	// If diff keep Left data
	updates := Realestate{}
	isUpdated := false
	if left.Title != right.Title {
		updates.Title = right.Title
		isUpdated = true
	}
	if left.Address != right.Address {
		updates.Address = right.Address
		isUpdated = true
	}
	if left.URL != right.URL {
		updates.URL = right.URL
		isUpdated = true
	}
	//if left.DateTime != right.DateTime {
	//	updates.DateTime = right.DateTime
	//	isUpdated = true
	//}
	if left.ID != right.ID {
		updates.ID = right.ID
		isUpdated = true
	}
	if left.Price != right.Price {
		updates.Price = right.Price
		isUpdated = true
	}
	for rightKey, rightVal := range right.Info {
		leftVal, ok := left.Info[rightKey]
		if !ok || leftVal != rightVal {
			updates.Info[rightKey] = rightVal
			isUpdated = true
		}
	}
	return updates, isUpdated
}
func (left Realestate) LeftUpdates(right Realestate) (Realestate, bool) {
	return right.RightUpdates(left)
}

//func getAllKeys(c *mongo.Collection) {
//
//}
