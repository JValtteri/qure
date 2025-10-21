package utils

type ExampleStruct struct {
    ID               string
    Name             string;
    ShortDescription string;
    LongDescription  string;
    DtStart          Epoch;
    DtEnd            Epoch;
    StaffSlots       int;
    Staff            int;
}

var exampleJson []byte = []byte(`{
    "ID": "0",
    "name": "Test event",
    "shortDescription": "Lorem ipsum dolor sit amet, meis illud at his",
    "LongDescription": "Lorem ipsum dolor sit amet, meis illud at his, ornatus facilisis ocurreret sit ut. Duo sale tractatos at, facilisi accusamus at per. Pro magna probatus senserit te, his sumo dico lucilius at. Usu ut iisque theophrastus definitiones, tollit latine aliquid an vim, eu cum partem voluptua. Eum aliquip qualisque interpretaris ut, per affert legere dissentiunt ad. Scaevola facilisi expetendis an nam, lucilius convenire dignissim per ei.",
    "DtStart":          1735675270,
    "DtEnd":            1735687830,
    "StaffSlots":       5,
    "Staff":            1
}`)

var badJson []byte = []byte(`
{
    "spam": "eggs
}`)
