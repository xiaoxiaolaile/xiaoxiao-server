package core

func black(filter *Filter, str string) bool {
	if filter != nil {
		if filter.BlackMode {
			if Contains(filter.Items, str) {
				return true
			}
		} else {
			if !Contains(filter.Items, str) {
				return true
			}
		}
	}
	return false
}

func Contains(strs []string, str string) bool {
	for _, o := range strs {
		if str == o {
			return true
		}
	}
	return false
}
