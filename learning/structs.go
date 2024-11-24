package learning

///// creation
/// type Student struct {
/// 	Name string
/// 	Age  int
/// }

///// declaration
/// var a Student    // a == Student{"", 0}
/// a.Name = "Alice" // a == Student{"Alice", 0}

///// instantiation 1 - "new()" returns pointer to newly created struct
/// var pa *Student   // pa == nil
/// pa = new(Student) // pa == &Student{"", 0}
/// pa.Name = "Alice" // pa == &Student{"Alice", 0}

///// instantion 2 - struct literal
/// b := Student{ // b == Student{"Bob", 0}
/// 	Name: "Bob",
/// }
///
/// pb := &Student{ // pb == &Student{"Bob", 8}
/// 	Name: "Bob",
/// 	Age:  8,
/// }
///
/// c := Student{"Cecilia", 5} // c == Student{"Cecilia", 5}
/// d := Student{}             // d == Student{"", 0}

/*
- An element list that contains keys does not need to have an element for each struct field. Omitted fields get the zero value for that field.
- An element list that does not contain any keys must list an element for each struct field in the order in which the fields are declared.
- A literal may omit the element list; such a literal evaluates to the zero value for its type.
*/
