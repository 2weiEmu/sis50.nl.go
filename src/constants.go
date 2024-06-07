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
func getValidPages() []string {
	return []string{
		"messages", 
		"help",
		"login",
		"profile",
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
