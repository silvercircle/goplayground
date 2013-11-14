package lib

type User struct {
	ID    int
	NAME  string
	EMAIL string
}

func (d *Data) DoServe() {
	if d.TheDB.Ping() == nil {
		user := make([]User, 0)
		index := 0
		result := map[string]interface{}{}

		rows, err := d.TheDB.NamedQueryMap("SELECT id, name, email FROM user WHERE id = :id",
			map[string]interface{}{"id": 1})
		if err != nil {
			d.DBError(err)
		} else {
			defer rows.Close()
			for rows.Next() {
				user = append(user, User{})
				rows.StructScan(&user[index])
				index++
			}
			result["user"] = user
			rows.Close()
			//cnt, _ := LoadTemplate("test/test")
			d.Templates = append(d.Templates, "test/test")
		}
		//cnt, _ := LoadTemplate("content")
		d.Templates = append(d.Templates, "content")
		d.Context["D"] = result
		d.Context["content"] = "The content"
	}
	d.Context["httpheaders"] = map[string]string{"X-Served-By": "Golang"}
}

// main dispatcher function. Looks at the route and calls the appropriate
// action.
func (d *Data) Dispatch() {
	d.DoServe()
}
