package main

const MessageFile = "./resources/messages"
const ShoppingFile = "./resources/shopping"
const CalendarFile = "./resources/calendar"
const IndexFile = "src/static/templates/index.html"
const MainLog = "./log/sis50.log"

const DayCount = 7
const UserCount = 4

// functions below here are stand-ins for constant arrays as these are not
// a thing in go
func getValidPages() []string {
	return []string{
		"messages", 
		"help",
	}
}

func getStateList() []string {
	return []string{
		"present", 
		"absent",
		"cooking",
		"uncertain",
		"maybe-cooking",
		"cant-cook",
	}
}

func getPersonList() []string {
	return []string{
		"rick",
		"youri", 
		"robert",
		"milan",
	}
}

func getDayList() []string {
	return []string{
		"ma",
		"di",
		"wo",
		"do",
		"vr",
		"za",
		"zo",
	}
}

// Actions for Shopping
const (
	ADD = iota
	OPEN
	REMOVE
	EDIT
	REARRANGE
)

var shoppingActionMap = map[string]int{
	"add": ADD,
	"open-shopping": OPEN,
	"remove": REMOVE,
	"edit": EDIT,
	"rearrange": REARRANGE,
}
