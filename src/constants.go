package src

const MessageFile = "./resources/messages"
const ShoppingFile = "./resources/shopping"
const CalendarFile = "./resources/calendar"
const IndexFile = "src/static/templates/index.html"
const MainLog = "./log/sis50.log"

const DayCount = 7
const UserCount = 4

// functions below here are stand-ins for constant arrays as these are not
// a thing in go

// TODO: i should probably find a better way to do this, and this will change
// with user permissions as well of course but who knows for now this is fine
func GetValidPages() []string {
	return []string{
		"messages", 
		"help",
		"login",
		"profile",
	}
}

func GetStateList() []string {
	return []string{
		"present", 
		"absent",
		"cooking",
		"uncertain",
		"maybe-cooking",
		"cant-cook",
	}
}

func GetPersonList() []string {
	return []string{
		"rick",
		"youri", 
		"robert",
		"milan",
	}
}

func GetDayList() []string {
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

var ShoppingActionMap = map[string]int{
	"add": ADD,
	"open-shopping": OPEN,
	"remove": REMOVE,
	"edit": EDIT,
	"rearrange": REARRANGE,
}
